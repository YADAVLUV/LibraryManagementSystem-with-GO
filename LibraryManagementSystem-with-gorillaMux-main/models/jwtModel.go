package models

import "github.com/dgrijalva/jwt-go"

var JwtKey = []byte("aofsijfwokfjsadklfj")

type Claims struct {
	Username string
	UserType string
	jwt.StandardClaims
}
