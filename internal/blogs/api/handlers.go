package api

import (
	"encoding/json"
	blogs "github.com/Khvan-Group/blog-service/internal/blogs/models"
	"github.com/Khvan-Group/blog-service/internal/blogs/store"
	"github.com/Khvan-Group/blog-service/internal/common/utils"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/Khvan-Group/common-library/models"
	"github.com/gorilla/mux"
	"golang.org/x/tools/container/intsets"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	APPLICATION_JSON = "application/json"
	CONTENT_TYPE     = "Content-type"
)

// Create
// @Summary Создание блога
// @ID create-blog
// @Accept  json
// @Produce  json
// @Param input body blogs.BlogCreate true "Информация о создаваемом блоге"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /blogs [post]
// @Security ApiKeyAuth
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	var input blogs.BlogCreate

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		panic(err)
	}

	if createErr := a.blogs.Service.Create(input, utils.GetJwtUser(r)); createErr != nil {
		errors.HandleError(w, createErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Update
// @Summary Обновление блога
// @ID update-blog
// @Accept  json
// @Produce  json
// @Param id path int true "ID блога"
// @Param input body blogs.BlogUpdate true "Информация об обновляемом блоге"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /blogs/{id} [put]
// @Security ApiKeyAuth
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	var input blogs.BlogUpdate
	id, err := strconv.Atoi(mux.Vars(r)["id"])

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

	if createErr := a.blogs.Service.Update(id, input, utils.GetJwtUser(r)); createErr != nil {
		errors.HandleError(w, createErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Send
// @Summary Отправить блог на рассмотрение модератора
// @ID send-blog
// @Accept  json
// @Produce  json
// @Param id path int true "ID блога"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /blogs/{id}/send [post]
// @Security ApiKeyAuth
func (a *API) Send(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	if sendErr := a.blogs.Service.Send(id, utils.GetJwtUser(r)); sendErr != nil {
		errors.HandleError(w, sendErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// FindAll
// @Summary Получить список блогов
// @ID find-all-blogs
// @Accept  json
// @Produce  json
// @Param id path int true "ID блога"
// @Param page query int false "Номер страницы"
// @Param size query int false "Количество элементов"
// @Param sortFields query string false "Сортируемые поля"
// @Param title query string false "Заголовок блога"
// @Param status query string false "Статус блога"
// @Param category query string false "Категория блога"
// @Success 200 {array} blogs.BlogView
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /blogs [get]
func (a *API) FindAll(w http.ResponseWriter, r *http.Request) {
	var input blogs.BlogSearch
	queryParams := r.URL.Query()
	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil {
		page = 0
	}

	size, err := strconv.Atoi(queryParams.Get("size"))
	if err != nil {
		size = intsets.MaxInt
	}

	title := queryParams.Get("title")
	if len(title) > 0 {
		input.Title = &title
	}

	sortFields := r.URL.Query()["sortFields"]

	for _, field := range sortFields {
		parts := strings.Split(field, ":")

		if len(parts) == 2 {
			input.SortBy = append(input.SortBy, models.SortField{SortBy: parts[0], Direction: parts[1]})
		}
	}

	status := queryParams.Get("status")
	category := queryParams.Get("category")

	if len(status) > 0 {
		if !blogs.IsValidStatus(status) {
			errors.HandleError(w, errors.NewBadRequest("Неверный переданный статус."))
			return
		}

		input.Status = &status
	}

	if len(category) > 0 {
		if !blogs.IsValidCategory(category) {
			errors.HandleError(w, errors.NewBadRequest("Неверная переданная категория."))
			return
		}

		input.Category = &category
	}

	input.Page = page
	input.Size = size
	input.CurrentUser = utils.GetJwtUser(r)
	response := a.blogs.Service.FindAll(input)

	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

// FindById
// @Summary Получить блог по ID
// @ID find-by-id-blog
// @Accept  json
// @Produce  json
// @Param id path int true "ID блога"
// @Success 200 {object} blogs.BlogView
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /blogs/{id} [get]
func (a *API) FindById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	response, findErr := a.blogs.Service.FindById(id, utils.GetJwtUser(r))
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

// Delete
// @Summary Удалить блог
// @ID delete-blog
// @Accept  json
// @Produce  json
// @Param id path int true "ID блога"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /blogs/{id} [delete]
// @Security ApiKeyAuth
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	deleteErr := a.blogs.Service.Delete(id, utils.GetJwtUser(r))
	if deleteErr != nil {
		errors.HandleError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteAllByUsername
// @Summary Удалить все блоги пользователя
// @ID delete-all-blogs-by-user
// @Accept  json
// @Produce  json
// @Param username path string true "Логин пользователя"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /blogs/{username}/delete [delete]
// @Security ApiKeyAuth
func (a *API) DeleteAllByUsername(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	deleteErr := a.blogs.Service.DeleteAllByUsername(username)
	if deleteErr != nil {
		errors.HandleError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// LikeOrFavorite
// @Summary Лайкнуть или добавить блог в избранное
// @ID like-or-favorite-blog
// @Accept  json
// @Produce  json
// @Param id path int true "ID блога"
// @Param action query string true "Лайкнуть или добавить в избранные"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /blogs/{id} [post]
// @Security ApiKeyAuth
func (a *API) LikeOrFavorite(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	action := r.URL.Query().Get("action")

	if action != store.LIKE && action != store.FAVORITE {
		errors.HandleError(w, errors.NewBadRequest("Неверное действие над блогом."))
		return
	}

	a.blogs.Service.LikeOrFavorite(id, utils.GetJwtUser(r), action)
	w.WriteHeader(http.StatusOK)
}

// Confirm
// @Summary Утвердить или отклонить блог
// @ID confirm-blog
// @Accept  json
// @Produce  json
// @Param id path int true "ID блога"
// @Param status query string true "Статус блога"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /blogs/{id}/confirm [post]
// @Security ApiKeyAuth
func (a *API) Confirm(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	confirmErr := a.blogs.Service.Confirm(id, status, utils.GetJwtUser(r))
	if confirmErr != nil {
		errors.HandleError(w, confirmErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}
