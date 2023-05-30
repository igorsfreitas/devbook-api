package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	response "github.com/igorsfreitas/devbook-api/src/commons/responses"
	"github.com/igorsfreitas/devbook-api/src/db"
	"github.com/igorsfreitas/devbook-api/src/models"
	"github.com/igorsfreitas/devbook-api/src/repositories"
)

// CreateUser cria um usuário
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(requestBody, &user); err != nil {
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
	userID, err := repository.Create(user)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	user.ID = userID
	response.JSON(w, http.StatusCreated, user)
}

// FindUsers busca todos os usuários
func FindUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Buscando todos os usuários"))
}

// FindUser busca um usuário
func FindUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Buscando um usuário"))
}

// UpdateUser atualiza um usuário
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Atualizando usuário"))
}

// DeleteUser deleta um usuário
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Deletando usuário"))
}
