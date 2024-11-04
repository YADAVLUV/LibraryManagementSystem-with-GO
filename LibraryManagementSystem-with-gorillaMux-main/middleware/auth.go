package middleware

import (
	"fmt"
	"learningGorillamux/models"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func ValidateUser(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("JwtToken")
		if tokenStr == "" {
			w.Write([]byte("User not logged In"))
			return
		}
		claims := &models.Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return models.JwtKey, nil
		})
		if err != nil {
			log.Println(err)
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("you are not logged in!"))
			return
		}

		f(w, r)

	}
}

func ValidateOwner(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("JwtToken")
		if tokenStr == "" {
			fmt.Println(tokenStr)
			w.Write([]byte("User not logged In"))
			return
		}
		claims := &models.Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return models.JwtKey, nil
		})
		if err != nil {
			log.Println(err)
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("you are not logged in!"))
			return
		}
		if claims.UserType != "owner" {
			w.Write([]byte("You are not a owner"))
			return
		}
		f(w, r)

	}
}

var NubmerOfRequests = 0

func TrackNumberOfRequests(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NubmerOfRequests = NubmerOfRequests + 1
		fmt.Println("Request Number : ", NubmerOfRequests)

		f.ServeHTTP(w, r)
	})
}
