package api

import (
	"encoding/json"
	blogs "github.com/dkhvan-dev/alabs_project/blog-service/internal/blogs/model"
	"github.com/dkhvan-dev/alabs_project/common-libraries/errors"
	"github.com/dkhvan-dev/alabs_project/common-libraries/utils"
	"github.com/gorilla/context"
	"io"
	"net/http"
)

const (
	APPLICATION_JSON = "application/json"
	CONTENT_TYPE     = "Content-type"
)

func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	currentUserLogin := utils.ToString(context.Get(r, "login"))

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	input := blogs.BlogCreate{}

	err = json.Unmarshal(data, &input)
	if err != nil {
		panic(err)
	}

	err = a.blogs.Service.Create(input, currentUserLogin)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
