package router

import (
	// "net/http"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/nchalla5/react-go-app/middleware"
)

func InitializeRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", middleware.Login).Methods("POST")
	//r.HandleFunc("/api/products", ProductsHandler).Methods("GET")
	// Add more routes as needed
	r.HandleFunc("/signup", middleware.Signup).Methods("POST")
	r.HandleFunc("/create-product", middleware.CreateProductHandler).Methods("POST")
	r.HandleFunc("/products", middleware.ListProductsHandler).Methods("GET")
	r.HandleFunc("/product/{id}", middleware.GetProductHandler).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(middleware.NotFoundHandler)

	return r
}
