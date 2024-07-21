package service

import (
	"github.com/Khvan-Group/blog-service/internal/categories/models"
	"github.com/Khvan-Group/blog-service/internal/categories/store"
	"github.com/Khvan-Group/blog-service/internal/db"
	"github.com/Khvan-Group/common-library/errors"
)

type CategoryService interface {
	Save(category models.Category)
	FindAll() []models.Category
	Delete(code string) *errors.CustomError
}

type Categories struct {
	Service CategoryService
}

func New() *Categories {
	return &Categories{
		Service: store.New(db.DB),
	}
}
