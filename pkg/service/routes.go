package service

import (
	"github.com/danikarik/mux"
)

const (
	appleSerialsRoute    string = "/devices/{deviceID}/registrations/{passTypeID}"
	appleLogRoute        string = "/log"
	appleRegisterRoute   string = "/devices/{deviceID}/registrations/{passTypeID}/{serialNumber}"
	appleUnregisterRoute string = "/devices/{deviceID}/registrations/{passTypeID}/{serialNumber}"
	appleLatestRoute     string = "/passes/{passTypeID}/{serialNumber}"
)

var verifyQueries = []string{
	"type", "{type}",
	"token", "{token}",
	"redirect_url", "{redirect_url:http.+}",
}

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
	r.HandleFunc("/version", s.versionHandler).Methods("GET")

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
		public.HandleFunc("/logout", s.logoutHandler).Methods("DELETE")
		public.HandleFunc("/register", s.registerHandler).Methods("POST")
		public.HandleFunc("/recover", s.recoverHandler).Methods("POST")
		public.HandleFunc("/reset", s.resetHandler).Methods("POST")
		public.HandleFunc("/verify", s.verifyHandler).Methods("GET").Queries(verifyQueries...)
		public.HandleFunc("/check/email", s.checkEmailHandler).Methods("POST")
		public.HandleFunc("/check/username", s.checkUsernameHandler).Methods("POST")

		protected := api.NewRoute().Subrouter()
		protected.Use(s.authMiddleware, s.csrfMiddleware)
		protected.HandleFunc("/invite", s.inviteHandler).Methods("POST")
		protected.HandleFunc("/account", s.accountHandler).Methods("GET")
		protected.HandleFunc("/account/email", s.emailChangeHandler).Methods("PUT")
		protected.HandleFunc("/account/username", s.usernameChangeHandler).Methods("PUT")
		protected.HandleFunc("/account/password", s.passwordChangeHandler).Methods("PUT")
		protected.HandleFunc("/account/metadata", s.metaDataChangeHandler).Methods("PUT")
	}

	s.handler = s.corsMiddleware(r)
	return s
}
