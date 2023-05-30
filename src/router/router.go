package router

import (
	"github.com/gorilla/mux"
	"github.com/igorsfreitas/devbook-api/src/router/routes"
)

// Gerar retorna um router com as rotas configuradas
func Gerar() *mux.Router {
	r := mux.NewRouter()

	return routes.ConfigRoutes(r)
}
