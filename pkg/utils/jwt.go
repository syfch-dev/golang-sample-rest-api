package utils

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.StandardClaims
}
