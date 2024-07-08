package api

import (
	"encoding/json"
	blogs "github.com/Khvan-Group/blog-service/internal/blogs/model"
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
	input := blogs.BlogCreate{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		panic(err)
	}

	if createErr := a.blogs.Service.Create(input, getJwtUser(r)); createErr != nil {
		errors.HandleError(w, createErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	input := blogs.BlogUpdate{}
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		panic(err)
	}

	if createErr := a.blogs.Service.Update(id, input, getJwtUser(r)); createErr != nil {
		errors.HandleError(w, createErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) Send(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		panic(err)
	}

	if sendErr := a.blogs.Service.Send(id, getJwtUser(r)); sendErr != nil {
		errors.HandleError(w, sendErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) FindAll(w http.ResponseWriter, r *http.Request) {
	var input blogs.BlogSearch
	page, err := strconv.Atoi(mux.Vars(r)["page"])
	if err != nil {
		page = 0
	}

	size, err := strconv.Atoi(mux.Vars(r)["size"])
	if err != nil {
		size = 50
	}

	title, ok := mux.Vars(r)["title"]
	if ok {
		input.Title = &title
	}

	sortFields := r.URL.Query()["sortFields"]
	status, existsStatus := mux.Vars(r)["status"]
	category, existsCategory := mux.Vars(r)["category"]

	if existsStatus {
		if !blogs.IsValidStatus(status) {
			errors.HandleError(w, errors.NewBadRequest("Неверный переданный статус."))
			return
		}

		input.Status = &status
	}

	if existsCategory {
		if !blogs.IsValidCategory(category) {
			errors.HandleError(w, errors.NewBadRequest("Неверная переданная категория."))
			return
		}

		input.Category = &category
	}

	input.Page = page
	input.Size = size
	input.SortBy = sortFields
	input.CurrentUser = getJwtUser(r)
	response := a.blogs.Service.FindAll(input)
	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

func (a *API) FindById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		panic(err)
	}

	response, findErr := a.blogs.Service.FindById(id, getJwtUser(r))
	if findErr != nil {
		errors.HandleError(w, findErr)
		return
	}

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

	deleteErr := a.blogs.Service.Delete(id, getJwtUser(r))
	if deleteErr != nil {
		errors.HandleError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) LikeOrFavorite(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		panic(err)
	}

	action := r.PathValue("action")
	if action != "likes" && action != "favorites" {
		errors.HandleError(w, errors.NewNotFound(""))
		return
	}

	a.blogs.Service.LikeOrFavorite(id, getJwtUser(r), action)
	w.WriteHeader(http.StatusOK)
}

func (a *API) Confirm(w http.ResponseWriter, r *http.Request) {
	status := mux.Vars(r)["status"]
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		panic(err)
	}

	confirmErr := a.blogs.Service.Confirm(id, status, getJwtUser(r))
	if confirmErr != nil {
		errors.HandleError(w, confirmErr)
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
