package service

import (
	"github.com/danikarik/mux"
)

const (
	appleSerialsRoute    string = "/devices/{deviceLibraryIdentifier}/registrations/{passTypeIdentifier}"
	appleLogRoute        string = "/log"
	appleRegisterRoute   string = "/devices/{deviceLibraryIdentifier}/registrations/{passTypeIdentifier}/{serialNumber}"
	appleUnregisterRoute string = "/devices/{deviceLibraryIdentifier}/registrations/{passTypeIdentifier}/{serialNumber}"
	appleLatestRoute     string = "/passes/{passTypeIdentifier}/{serialNumber}"
)

func (s *Service) routerOptions(r *mux.Router) {
	r.Wrapper = mux.NewDefaultWrapper(s.errorHandler)
	r.NotFoundHandler = s.notFoundHandler
	r.MethodNotAllowedHandler = s.methodNotAllowedHandler
}

func (s *Service) withRouter() *Service {
	r := mux.NewRouter(s.routerOptions)

	r.Use(requestIDMiddleware)
	r.Use(realIPMiddleware)
	r.Use(loggerMiddleware(s.logger))
	r.Use(recovererMiddleware(s.logger))

	r.HandleFunc("/health", s.healthHandler).Methods("GET")

	apple := r.PathPrefix("/v1").Subrouter()
	{
		public := apple.NewRoute().Subrouter()
		public.HandleFunc(appleSerialsRoute, s.serialNumbers).Methods("GET")
		public.HandleFunc(appleLogRoute, s.errorLogs).Methods("POST")

		protected := apple.NewRoute().Subrouter()
		protected.Use(applePassMiddleware)
		protected.HandleFunc(appleRegisterRoute, s.registerDevice).Methods("POST")
		protected.HandleFunc(appleUnregisterRoute, s.unregisterDevice).Methods("DELETE")
		protected.HandleFunc(appleLatestRoute, s.latestPass).Methods("GET")
	}

	api := r.NewRoute().Subrouter()
	{
		public := api.NewRoute().Subrouter()
		public.HandleFunc("/", s.okHandler).Methods("GET")
		public.HandleFunc("/login", s.loginHandler).Methods("POST")
		public.HandleFunc("/register", s.registerHandler).Methods("POST")
		public.HandleFunc("/recover", s.recoverHandler).Methods("POST")
		public.HandleFunc("/verify", s.verifyHandler).Methods("GET").Queries(
			"type", "{type}",
			"token", "{token}",
			"redirect_url", "{redirect_url:http.+}",
		)

		protected := api.NewRoute().Subrouter()
		protected.Use(s.authMiddleware, s.csrfMiddleware)
		protected.HandleFunc("/logout", s.logoutHandler).Methods("DELETE")
		protected.HandleFunc("/account", s.accountHandler).Methods("GET")
	}

	s.handler = s.corsMiddleware(r)
	return s
}
