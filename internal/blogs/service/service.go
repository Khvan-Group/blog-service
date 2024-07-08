package service

import (
	blogs "github.com/Khvan-Group/blog-service/internal/blogs/model"
	"github.com/Khvan-Group/blog-service/internal/users/model"
	"github.com/Khvan-Group/common-library/errors"
)

type BlogService interface {
	Create(input blogs.BlogCreate, currentUser model.JwtUser) *errors.CustomError
	Update(id int, input blogs.BlogUpdate, currentUser model.JwtUser) *errors.CustomError
	Send(id int, currentUser model.JwtUser) *errors.CustomError
	FindAll(input blogs.BlogSearch) []blogs.BlogView
	FindById(id int, currentUser model.JwtUser) (*blogs.BlogView, *errors.CustomError)
	Delete(id int, currentUser model.JwtUser) *errors.CustomError
	LikeOrFavorite(id int, currentUser model.JwtUser, action string)
	Confirm(id int, status string, currentUser model.JwtUser) *errors.CustomError
}

type Blogs struct {
	Service BlogService
}

func New(s BlogService) *Blogs {
	return &Blogs{
		Service: s,
	}
}
