package models

import (
	blogs "github.com/Khvan-Group/blog-service/internal/blogs/models"
	"github.com/Khvan-Group/blog-service/internal/clients"
	"github.com/Khvan-Group/blog-service/internal/common/models"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/go-resty/resty/v2"
	"time"
)

type Comment struct {
	Id        int         `json:"id" db:"id"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	CreatedBy models.User `json:"created_by" db:"created_by"`
	Blog      blogs.Blog  `json:"blog" db:"blog"`
	Comment   string      `json:"comment" db:"comment"`
}

// DTOs
type CommentCreate struct {
	BlogId    int    `json:"blog_id" db:"blog_id"`
	CreatedBy string `db:"created_by"`
	Comment   string `json:"comment" db:"comment"`
}

type CommentView struct {
	Id        int             `json:"id" db:"id"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	CreatedBy models.UserView `json:"created_by" db:"created_by"`
	Comment   string          `json:"comment" db:"comment"`
}

func (c *CommentView) FillUserInfo(client *resty.Client) *errors.CustomError {
	createdBy, err := clients.GetUserByLogin(c.CreatedBy.Login, client)
	if err != nil {
		return err
	}

	c.CreatedBy = *createdBy
	return nil
}
