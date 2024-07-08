package store

import (
	"fmt"
	"github.com/Khvan-Group/blog-service/internal/comments/model"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/jmoiron/sqlx"
)

type CommentStore struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *CommentStore {
	return &CommentStore{
		db: db,
	}
}

func (s *CommentStore) Create(input model.CommentCreate) *errors.CustomError {
	var blogExists bool
	err := s.db.Get(&blogExists, "select exists(select 1 from t_blogs where id = ? and is_deleted is false)", input.BlogId)
	if err != nil {
		return errors.NewBadRequest(fmt.Sprintf("Блог с ID: %d не найден.", input.BlogId))
	}

	_, err = s.db.NamedExec("insert into t_comments (created_by, blog_id, comment) values (:created_by, :blog_id, :comment)", input)
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *CommentStore) FindAll(blogId int) []model.CommentView {
	var comments []model.CommentView
	query := `
		select c.id, c.created_by, c.comment, created_by.username from t_comments c
        	inner join t_users created_by ON c.created_by = created_by.id
        where c.blog_id = ?
		order by c.created_at desc
	`

	err := s.db.Select(&comments, query, blogId)
	if err != nil {
		panic(err)
	}

	return comments
}

func (s *CommentStore) Delete(id int) *errors.CustomError {
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
}
