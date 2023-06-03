package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/igorsfreitas/devbook-api/src/commons/auth"
	response "github.com/igorsfreitas/devbook-api/src/commons/responses"
	"github.com/igorsfreitas/devbook-api/src/db"
	"github.com/igorsfreitas/devbook-api/src/models"
	"github.com/igorsfreitas/devbook-api/src/repositories"
)

// CreatePost creates a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post
	if err = json.Unmarshal(requestBody, &post); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	post.AuthorID = userID

	if err = post.Prepare(); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewPostRepository(db)
	post.ID, err = repository.Create(post)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, post)
}

// GetPosts returns all posts
func GetPosts(w http.ResponseWriter, r *http.Request) {

}

// GetPost returns a single post
func GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID, err := strconv.ParseUint(params["postId"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := db.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewPostRepository(db)
	post, err := repository.GetPost(postID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, post)
}

// UpdatePost updates a post
func UpdatePost(w http.ResponseWriter, r *http.Request) {

}

// DeletePost deletes a post
func DeletePost(w http.ResponseWriter, r *http.Request) {

}
