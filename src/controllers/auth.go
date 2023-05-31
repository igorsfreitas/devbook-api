package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/igorsfreitas/devbook-api/src/commons/auth"
	"github.com/igorsfreitas/devbook-api/src/commons/encrypt"
	response "github.com/igorsfreitas/devbook-api/src/commons/responses"
	"github.com/igorsfreitas/devbook-api/src/db"
	"github.com/igorsfreitas/devbook-api/src/models"
	"github.com/igorsfreitas/devbook-api/src/repositories"
)

// Login is a function to login
func Login(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)
	userOnDatabase, err := repository.FindByEmail(user.Email)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = encrypt.CheckPasswordHash(user.Password, userOnDatabase.Password); err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	token, err := auth.GenerateJWT(userOnDatabase.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, token)
}
