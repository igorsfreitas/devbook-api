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
		AuthRequired: true,
	},
	{
		URI:          "/user/{userId}",
		Method:       http.MethodPut,
		Function:     controllers.UpdateUser,
		AuthRequired: true,
	},
	{
		URI:          "/user/{userId}",
		Method:       http.MethodDelete,
		Function:     controllers.DeleteUser,
		AuthRequired: true,
	},
	{
		URI:          "/user/{userId}/follow",
		Method:       http.MethodPost,
		Function:     controllers.FollowUser,
		AuthRequired: true,
	},
	{
		URI:          "/user/{userId}/unfollow",
		Method:       http.MethodPost,
		Function:     controllers.UnfollowUser,
		AuthRequired: true,
	},
	{
		URI:          "/user/{userId}/followers",
		Method:       http.MethodGet,
		Function:     controllers.GetFollowers,
		AuthRequired: true,
	},
	{
		URI:          "/user/{userId}/following",
		Method:       http.MethodGet,
		Function:     controllers.GetFollowing,
		AuthRequired: true,
	},
	{
		URI:          "/user/{userId}/update-password",
		Method:       http.MethodPost,
		Function:     controllers.UpdatePassword,
		AuthRequired: true,
	},
}
