package store

import (
	"fmt"
	"github.com/Khvan-Group/blog-service/internal/comments/models"
	"github.com/Khvan-Group/blog-service/internal/db"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
)

type CommentStore struct {
	db     *sqlx.DB
	client *resty.Client
}

func New() *CommentStore {
	return &CommentStore{
		db:     db.DB,
		client: resty.New(),
	}
}

func (s *CommentStore) Create(input models.CommentCreate) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var blogExists bool
		err := s.db.Get(&blogExists, "select exists(select 1 from t_blogs where id = $1)", input.BlogId)
		if err != nil {
			return errors.NewBadRequest(fmt.Sprintf("Блог с ID: %d не найден.", input.BlogId))
		}

		_, err = s.db.NamedExec("insert into t_comments (created_by, blog_id, comment) values (:created_by, :blog_id, :comment)", input)
		if err != nil {
			panic(err)
		}

		return nil
	})
}

func (s *CommentStore) FindAll(blogId int) []models.CommentView {
	var comments []models.CommentView
	query := `
		select c.id, c.created_at, c.comment, c.created_by as "created_by.login"
		from t_comments c
        where c.blog_id = $1
		order by c.created_at desc
	`

	err := s.db.Select(&comments, query, blogId)
	if err != nil {
		panic(err)
	}

	for i := range comments {
		err := comments[i].FillUserInfo(s.client)
		if err != nil {
			panic(err)
		}
	}

	return comments
}

func (s *CommentStore) Delete(id int) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var commentExists bool
		err := s.db.Get(&commentExists, "select exists(select 1 from t_comments where id = ?)", id)
		if err != nil {
			return errors.NewBadRequest(fmt.Sprintf("Комментарий с ID: %d не найден.", id))
		}

		_, err = s.db.Exec("delete from t_comments where id = ?", id)
		if err != nil {
			panic(err)
		}

		return nil
	})
}
