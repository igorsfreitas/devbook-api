package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/igorsfreitas/devbook-api/src/commons/auth"
	response "github.com/igorsfreitas/devbook-api/src/commons/responses"
)

// Logger is a middleware to log the requests
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Host)
		next(w, r)
	}
}

// Auth is a middleware to validate the token
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := auth.ValidateToken(r); err != nil {
			response.Error(w, http.StatusUnauthorized, fmt.Errorf("Token inválido"))
			return
		}
		next(w, r)
	}
}
