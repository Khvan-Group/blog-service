package service

import (
	blogs "github.com/Khvan-Group/blog-service/internal/blogs/models"
	blogStore "github.com/Khvan-Group/blog-service/internal/blogs/store"
	"github.com/Khvan-Group/blog-service/internal/common/models"
	"github.com/Khvan-Group/blog-service/internal/db"
	"github.com/Khvan-Group/common-library/errors"
	commonModels "github.com/Khvan-Group/common-library/models"
)

type BlogService interface {
	Create(input blogs.BlogCreate, currentUser models.JwtUser) *errors.CustomError
	Update(id int, input blogs.BlogUpdate, currentUser models.JwtUser) *errors.CustomError
	Send(id int, currentUser models.JwtUser) *errors.CustomError
	FindAll(input blogs.BlogSearch) commonModels.Page
	FindById(id int, currentUser models.JwtUser) (*blogs.BlogView, *errors.CustomError)
	Delete(id int, currentUser models.JwtUser) *errors.CustomError
	DeleteAllByUsername(username string) *errors.CustomError
	LikeOrFavorite(id int, currentUser models.JwtUser, action string)
	Confirm(id int, status string, currentUser models.JwtUser) *errors.CustomError
}

type Blogs struct {
	Service BlogService
}

func New() *Blogs {
	return &Blogs{
		Service: blogStore.New(db.DB),
	}
}
