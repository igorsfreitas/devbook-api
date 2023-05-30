package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/igorsfreitas/devbook-api/src/db"
	"github.com/igorsfreitas/devbook-api/src/models"
	"github.com/igorsfreitas/devbook-api/src/repositories"
)

// CreateUser cria um usuário
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var usuario models.User
	if err = json.Unmarshal(requestBody, &usuario); err != nil {
		log.Fatal(err)
	}

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	repository := repositories.NewUserRepository(db)
	userID, err := repository.Create(usuario)
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(fmt.Sprintf("Usuário inserido com sucesso! ID: %d", userID)))
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
