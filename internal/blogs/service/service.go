package service

import (
	blogs "github.com/dkhvan-dev/alabs_project/blog-service/internal/blogs/model"
	"github.com/dkhvan-dev/alabs_project/common-libraries/errors"
)

type BlogService interface {
	Create(input blogs.BlogCreate, currentUserLogin string) *errors.Error
	Update(id int, input blogs.BlogUpdate, currentUserLogin string) *errors.Error
	FindAll(page, size int) []blogs.BlogView
	FindById(id int) (*blogs.BlogView, *errors.Error)
	Delete(id int) *errors.Error
}

type Blogs struct {
	Service BlogService
}

func New(s BlogService) *Blogs {
	return &Blogs{
		Service: s,
	}
}
