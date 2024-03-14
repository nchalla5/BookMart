package main

import (
	//"fmt"
	"github.com/gorilla/handlers"
	// "github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/nchalla5/react-go-app/router"
)

func main() {
	// r := router.Router()
	r := router.InitializeRouter()
	log.Println("Starting server on :8080")
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"}) // React app's address
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	if err := http.ListenAndServe(":8080", handlers.CORS(headersOk, originsOk, methodsOk)(r)); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
