package routes

import (
	"net/http"

	"github.com/igorsfreitas/devbook-api/src/controllers"
)

var loginRoute = Route{
	URI:          "/login",
	Method:       http.MethodPost,
	Function:     controllers.Login,
	AuthRequired: false,
}
