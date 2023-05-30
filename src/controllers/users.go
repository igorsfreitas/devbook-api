package controllers

import (
	"net/http"
)

// CreateUser cria um usuário
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Criando usuário"))
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
