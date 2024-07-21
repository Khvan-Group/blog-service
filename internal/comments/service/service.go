package service

import (
	"github.com/Khvan-Group/blog-service/internal/comments/models"
	commentStore "github.com/Khvan-Group/blog-service/internal/comments/store"
	"github.com/Khvan-Group/common-library/errors"
)

type CommentService interface {
	Create(input models.CommentCreate) *errors.CustomError
	FindAll(blogId int) []models.CommentView
	Delete(id int) *errors.CustomError
}

type Comments struct {
	Service CommentService
}

func New() *Comments {
	return &Comments{
		Service: commentStore.New(),
	}
}
