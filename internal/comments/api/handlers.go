package api

import (
	"encoding/json"
	comments "github.com/Khvan-Group/blog-service/internal/comments/model"
	"github.com/Khvan-Group/blog-service/internal/users/model"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/Khvan-Group/common-library/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

const (
	APPLICATION_JSON = "application/json"
	CONTENT_TYPE     = "Content-type"
)

func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	input := comments.CommentCreate{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		panic(err)
	}

	input.CreatedBy = getJwtUser(r).Login
	if createErr := a.comments.Service.Create(input); createErr != nil {
		errors.HandleError(w, createErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) FindAll(w http.ResponseWriter, r *http.Request) {
	blogId, err := strconv.Atoi(mux.Vars(r)["blogId"])
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

func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
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

func getJwtUser(r *http.Request) model.JwtUser {
	return model.JwtUser{
		Login: utils.ToString(context.Get(r, "login")),
		Role:  utils.ToString(context.Get(r, "role")),
	}
}
