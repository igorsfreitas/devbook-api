package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/igorsfreitas/devbook-api/src/middlewares"
)

type Route struct {
	URI          string
	Method       string
	Function     func(w http.ResponseWriter, r *http.Request)
	AuthRequired bool
}

// ConfigRoutes insere as rotas no router
func ConfigRoutes(r *mux.Router) *mux.Router {
	routes := userRoutes
	routes = append(routes, loginRoute)

	for _, route := range routes {

		if route.AuthRequired {
			r.HandleFunc(route.URI,
				middlewares.Logger(middlewares.Auth(route.Function)),
			).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, middlewares.Logger(route.Function)).Methods(route.Method)
		}
	}

	return r
}
