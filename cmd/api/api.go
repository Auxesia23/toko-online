package main

import (
	"log"
	"net/http"

	"github.com/Auxesia23/toko-online/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	Config   config
	User repository.UserRepository
}

type config struct {
	addr string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/", app.HandleHealthCheck)

		r.Route("/user", func(r chi.Router){
			r.Use(JWTAuthMiddleware)
			r.Get("/me",app.GetUser)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.RegisterHanlder)
			r.Post("/login", app.LoginHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:    app.Config.addr,
		Handler: mux,
	}

	log.Println("Server running on port" + app.Config.addr)

	return srv.ListenAndServe()
}
