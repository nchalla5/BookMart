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
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods("GET")
	r.HandleFunc("/signup", middleware.Signup).Methods("POST")
	r.HandleFunc("/create-product", middleware.CreateProductHandler).Methods("POST")
	r.HandleFunc("/products", middleware.ListProductsHandler).Methods("GET")
	r.HandleFunc("/product/{id}", middleware.GetProductHandler).Methods("GET")
	r.HandleFunc("/product/{id}/purchase", middleware.PurchaseProductHandler).Methods("POST")
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	r.NotFoundHandler = http.HandlerFunc(middleware.NotFoundHandler)

	return r
}
