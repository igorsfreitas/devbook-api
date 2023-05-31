package routes

import (
	"net/http"

	"github.com/igorsfreitas/devbook-api/src/controllers"
)

var userRoutes = []Route{
	{
		URI:          "/user",
		Method:       http.MethodPost,
		Function:     controllers.CreateUser,
		AuthRequired: false,
	},
	{
		URI:          "/users",
		Method:       http.MethodGet,
		Function:     controllers.FindUsers,
		AuthRequired: true,
	},
	{
		URI:          "/user/{userId}",
		Method:       http.MethodGet,
		Function:     controllers.FindUser,
		AuthRequired: false,
	},
	{
		URI:          "/user/{userId}",
		Method:       http.MethodPut,
		Function:     controllers.UpdateUser,
		AuthRequired: false,
	},
	{
		URI:          "/user/{userId}",
		Method:       http.MethodDelete,
		Function:     controllers.DeleteUser,
		AuthRequired: false,
	},
}
