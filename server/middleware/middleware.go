package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nchalla5/react-go-app/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Welcome to the Login Page"))
	var creds models.CredsStruct
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Login attempt with Email/Phone: %s, Password: %s", creds.EmailOrPhone, creds.Password)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})

}
