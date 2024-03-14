package router

import (
	// "net/http"

	"github.com/gorilla/mux"
	"github.com/nchalla5/react-go-app/middleware"
)

func InitializeRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", middleware.Login).Methods("POST")
	//r.HandleFunc("/api/products", ProductsHandler).Methods("GET")
	// Add more routes as needed
	return r
}
