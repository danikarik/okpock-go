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

		auth := public.NewRoute().Subrouter()
		auth.HandleFunc("/login", s.loginHandler).Methods("POST")
		auth.HandleFunc("/logout", s.logoutHandler).Methods("DELETE")
		auth.HandleFunc("/register", s.registerHandler).Methods("POST")
		auth.HandleFunc("/recover", s.recoverHandler).Methods("POST")
		auth.HandleFunc("/reset", s.resetHandler).Methods("POST")
		auth.HandleFunc("/verify", s.verifyHandler).Methods("GET").Queries(verifyQueries...)
		auth.HandleFunc("/check/email", s.checkEmailHandler).Methods("POST")
		auth.HandleFunc("/check/username", s.checkUsernameHandler).Methods("POST")

		protected := api.NewRoute().Subrouter()
		protected.Use(s.authMiddleware, s.csrfMiddleware)

		user := protected.NewRoute().Subrouter()
		user.HandleFunc("/invite", s.inviteHandler).Methods("POST")

		account := protected.PathPrefix("/account").Subrouter()
		account.HandleFunc("/info", s.accountInfoHandler).Methods("GET")
		account.HandleFunc("/email", s.emailChangeHandler).Methods("PUT")
		account.HandleFunc("/username", s.usernameChangeHandler).Methods("PUT")
		account.HandleFunc("/password", s.passwordChangeHandler).Methods("PUT")
		account.HandleFunc("/metadata", s.metaDataChangeHandler).Methods("PUT")

		projects := protected.PathPrefix("/projects").Subrouter()
		projects.HandleFunc("/check", s.checkProjectHandler).Methods("POST")
		// TODO: SaveNewProject
		projects.HandleFunc("/", s.okHandler).Methods("POST")
		// TODO: LoadProjects
		projects.HandleFunc("/", s.okHandler).Methods("GET")
		// TODO: LoadProject
		projects.HandleFunc("/{id}", s.okHandler).Methods("GET")
		// TODO: UpdateProject
		projects.HandleFunc("/{id}", s.okHandler).Methods("PUT")
		// TODO: SetBackgroundImage, SetFooterImage, SetIconImage, SetStripImage
		projects.HandleFunc("/{id}/upload", s.okHandler).Methods("PUT")
	}

	s.handler = s.corsMiddleware(r)
	return s
}
