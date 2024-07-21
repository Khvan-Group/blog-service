package clients

import (
	"encoding/json"
	"github.com/Khvan-Group/blog-service/internal/common/models"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/Khvan-Group/common-library/utils"
	"github.com/go-resty/resty/v2"
	"net/http"
)

var AUTH_SERVICE_URL = utils.GetEnv("AUTH_SERVICE_URL")

func GetUserByLogin(login string, client *resty.Client) (*models.UserView, *errors.CustomError) {
	var user models.UserView
	request := client.R()
	request.Header.Set(constants.X_IS_INTERNAL_SERVICE, "true")
	response, err := request.Get(AUTH_SERVICE_URL + "/admin/users" + login)
	if err != nil {
		panic(err)
	}

	if response.StatusCode() != http.StatusOK {
		return nil, errors.NewBadRequest("Ошибка получения пользователя по логину.")
	}

	err = json.Unmarshal(response.Body(), &user)
	if err != nil {
		panic(err)
	}

	return &user, nil
}
