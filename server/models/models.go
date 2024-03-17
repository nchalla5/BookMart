package models

import (
	"github.com/golang-jwt/jwt"
)

type CredsStruct struct {
	EmailOrPhone string `json:"emailOrPhone"`
	Password     string `json:"password"`
}

type DetailsStruct struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type Product struct {
	ProductID   string `json:"productId"` // Add this line
	Product     string `json:"product"`
	Image       string `json:"image,omitempty"`
	Title       string `json:"title"`
	Cost        string `json:"cost"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Status      string `json:"status,omitempty"`
	Seller      string `json:"seller,omitempty"`
	Buyer       string `json:"buyer,omitempty"`
}

type ProductsApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type Claims struct {
	MemberID string `json:"member_id"`
	jwt.StandardClaims
}
