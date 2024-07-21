package store

import (
	"github.com/Khvan-Group/blog-service/internal/categories/models"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/jmoiron/sqlx"
)

type CategoryStore struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *CategoryStore {
	return &CategoryStore{
		db: db,
	}
}

func (s *CategoryStore) Save(category models.Category) {
	var existsCategory bool
	s.db.Get(&existsCategory, "select exists(select 1 from t_blog_categories where code = $1)", category.Name)

	var err error
	if !existsCategory {
		_, err = s.db.NamedExec("insert into t_blog_categories values (:code, :name)", category)
	} else {
		_, err = s.db.Exec("update t_blog_categories set name = :name where code = :code", category)
	}

	if err != nil {
		panic(err)
	}
}

func (s *CategoryStore) FindAll() []models.Category {
	var response []models.Category
	err := s.db.Select(&response, "select * from t_blog_categories")
	if err != nil {
		panic(err)
	}

	return response
}

func (s *CategoryStore) Delete(code string) *errors.CustomError {
	var existsCategory bool
	s.db.Get(&existsCategory, "select exists(select 1 from t_blog_categories where code = $1)", code)

	if !existsCategory {
		return errors.NewBadRequest("Данной категории не существует.")
	}

	_, err := s.db.Exec("delete from t_blog_categories where code = $1", code)
	if err != nil {
		panic(err)
	}

	return nil
}
