package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nchalla5/react-go-app/models"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Welcome to the Login Page"))
	var details models.DetailsStruct
	err := json.NewDecoder(r.Body).Decode(&details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Signup attempt with Email: %s, Name: %s, Password: %s", details.Email, details.Name, details.Password)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
