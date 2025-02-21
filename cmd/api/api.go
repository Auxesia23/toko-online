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
	Product repository.ProductRepository
	Image repository.ImageRepository
	Category repository.CategoryRepository
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
			r.Get("/profile",app.GetUserHandler)
			r.Put("/profile", app.UpdateUserHandler)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.RegisterHanlder)
			r.Post("/login", app.LoginHandler)
		})

		r.Route("/product",func(r chi.Router){
			r.Get("/",app.GetProductsListHandler)
			r.With(SuperUserAuth).Post("/create", app.CreateProductHandler)
			r.Get("/{id}", app.GetSingleProductHandler)
			r.With(SuperUserAuth).Put("/{id}",app.UpdateProductHandler)
			r.With(SuperUserAuth).Delete("/{id}", app.DeleteProductHandler)
		})

		r.Route("/category",func(r chi.Router) {
			r.Post("/create", app.CreateCategoryHanlder)
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
