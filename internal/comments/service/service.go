package service

import (
	"github.com/Khvan-Group/blog-service/internal/comments/model"
	"github.com/Khvan-Group/common-library/errors"
)

type CommentService interface {
	Create(input model.CommentCreate) *errors.CustomError
	FindAll(blogId int) []model.CommentView
	Delete(id int) *errors.CustomError
}

type Comments struct {
	Service CommentService
}

func New(s CommentService) *Comments {
	return &Comments{
		Service: s,
	}
}
