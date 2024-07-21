package api

import (
	"encoding/json"
	comments "github.com/Khvan-Group/blog-service/internal/comments/models"
	"github.com/Khvan-Group/blog-service/internal/common/utils"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

const (
	APPLICATION_JSON = "application/json"
	CONTENT_TYPE     = "Content-type"
)

// Create
// @Summary Создание комментария
// @ID create-comment
// @Accept  json
// @Produce  json
// @Param input body comments.CommentCreate true "Информация о создаваемом комментарии"
// @Success 200
// @Failure 400 {object} errors.CustomError
// @Router /comments [post]
// @Security ApiKeyAuth
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	var input comments.CommentCreate

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		panic(err)
	}

	input.CreatedBy = utils.GetJwtUser(r).Login
	if createErr := a.comments.Service.Create(input); createErr != nil {
		errors.HandleError(w, createErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// FindAll
// @Summary Получение списка комментариев блога
// @ID find-all-comments
// @Accept  json
// @Produce  json
// @Param blogId query int true "ID блога"
// @Success 200 {array} comments.CommentView
// @Failure 400 {object} errors.CustomError
// @Router /comments [get]
func (a *API) FindAll(w http.ResponseWriter, r *http.Request) {
	blogId, err := strconv.Atoi(r.URL.Query().Get("blogId"))
	if err != nil {
		panic(err)
	}

	response := a.comments.Service.FindAll(blogId)
	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

// Delete
// @Summary Удаление комментария
// @ID delete-comment
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 400 {object} errors.CustomError
// @Router /comments/{id} [delete]
// @Security ApiKeyAuth
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	deleteErr := a.comments.Service.Delete(id)
	if deleteErr != nil {
		errors.HandleError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}
