package utils

import (
	"github.com/Khvan-Group/blog-service/internal/common/models"
	"github.com/Khvan-Group/common-library/utils"
	"github.com/gorilla/context"
	"net/http"
)

func GetJwtUser(r *http.Request) models.JwtUser {
	return models.JwtUser{
		Login: utils.ToString(context.Get(r, "login")),
		Role:  utils.ToString(context.Get(r, "role")),
	}
}
