package model

import (
	blogs "github.com/Khvan-Group/blog-service/internal/blogs/model"
	"github.com/Khvan-Group/blog-service/internal/users/model"
	"time"
)

type Comment struct {
	Id        int        `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	CreatedBy model.User `json:"created_by" db:"created_by"`
	Blog      blogs.Blog `json:"blog" db:"blog"`
	Comment   string     `json:"comment" db:"comment"`
}

// DTOs
type CommentCreate struct {
	BlogId    int    `json:"blog_id" db:"blog_id"`
	CreatedBy string `db:"created_by"`
	Comment   string `json:"comment" db:"comment"`
}

type CommentView struct {
	Id        int            `json:"id" db:"id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	CreatedBy model.UserView `json:"created_by" db:"created_by"`
	Comment   string         `json:"comment" db:"comment"`
}
