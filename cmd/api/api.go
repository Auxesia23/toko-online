package main

import (
	"log"
	"net/http"

	"github.com/Auxesia23/toko-online/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type application struct {
	Config   config
	User     repository.UserRepository
	Product  repository.ProductRepository
	Image    repository.ImageRepository
	Category repository.CategoryRepository
	Cart     repository.CartRepository
	Order    repository.OrderRepository
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
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, 
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/", app.HandleHealthCheck)

		r.Route("/user", func(r chi.Router) {
			r.Use(UserAuth)
			r.Get("/profile", app.GetUserHandler)
			r.Put("/profile", app.UpdateUserHandler)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.RegisterHanlder)
			r.Post("/login", app.LoginHandler)
			r.Get("/google/callback", app.GoogleCallbackHandler)
		})

		r.Route("/product", func(r chi.Router) {
			r.Get("/", app.GetProductsListHandler)
			r.With(SuperUserAuth).Post("/create", app.CreateProductHandler)
			r.Get("/{id}", app.GetSingleProductHandler)
			r.With(SuperUserAuth).Put("/{id}", app.UpdateProductHandler)
			r.With(SuperUserAuth).Delete("/{id}", app.DeleteProductHandler)
		})

		r.Route("/category", func(r chi.Router) {
			r.With(SuperUserAuth).Post("/create", app.CreateCategoryHanlder)
			r.With(SuperUserAuth).Delete("/{id}", app.DeleteCategoryHandler)
			r.With(SuperUserAuth).Put("/{id}", app.UpdateCategoryHanlder)
			r.Get("/", app.GetCategoryListHandler)
			r.Get("/{id}", app.GetCategoryHandler)
		})

		r.Route("/cart", func(r chi.Router) {
			r.Use(UserAuth)
			r.Post("/create", app.CreatecartHandler)
			r.Get("/", app.GetCartsHandler)
			r.Delete("/{id}", app.DeleteCartHandler)
			r.Post("/increase/{id}", app.IncrementCartHandler)
			r.Post("/decrease/{id}", app.DecrementsCartHandler)
		})

		r.Route("/order", func(r chi.Router) {
			r.Use(UserAuth)
			r.Post("/preview", app.PreviewOrderHandler)
			r.Post("/create", app.CreateOrderHanlder)
			r.Get("/", app.GetListOrderhanlder)
			r.Get("/{id}", app.GetOrderHandler)
			r.Post("/{id}/create-payment", app.CreatePaymentHandler)
		})

		r.Route("/webhook", func(r chi.Router) {
			r.Post("/midtrans", app.MidtransWebhookHandler)
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
