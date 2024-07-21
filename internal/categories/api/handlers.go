package api

import (
	"encoding/json"
	"github.com/Khvan-Group/blog-service/internal/categories/models"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

// Save
// @Summary Создание/Изменение категории
// @ID create-update-category
// @Accept  json
// @Produce  json
// @Param input body models.Category true "Информация о создаваемой/изменяемой категории"
// @Success 200
// @Router /categories [post]
func (a *API) Save(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	var input models.Category

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		panic(err)
	}

	a.categories.Service.Save(input)
	w.WriteHeader(http.StatusOK)
}

// FindAll
// @Summary Получение списка категорий
// @ID find-all-categories
// @Accept  json
// @Produce  json
// @Success 200 {array} []models.Category
// @Router /categories [get]
func (a *API) FindAll(w http.ResponseWriter, r *http.Request) {
	response := a.categories.Service.FindAll()

	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

// Delete
// @Summary Удаление категории
// @ID delete-category
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 400 {object} errors.CustomError
// @Router /categories/{code} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	err := a.categories.Service.Delete(code)
	if err != nil {
		errors.HandleError(w, err)
	}

	w.WriteHeader(http.StatusOK)
}
