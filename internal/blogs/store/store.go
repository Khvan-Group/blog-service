package store

import (
	"database/sql"
	"fmt"
	blogs "github.com/dkhvan-dev/alabs_project/blog-service/internal/blogs/model"
	"github.com/dkhvan-dev/alabs_project/common-libraries/errors"
)

type BlogStore struct {
	db *sql.DB
}

func New(db *sql.DB) *BlogStore {
	return &BlogStore{
		db: db,
	}
}

func (s *BlogStore) Create(input blogs.BlogCreate, currentUserLogin string) *errors.Error {
	if len(input.Title) == 0 || len(input.Content) == 0 {
		return errors.NewBadRequest("Пустой заголовок или содержание.")
	}

	_, err := s.db.Exec("insert into t_blogs(created_by, title, content, category) values (?, ?, ?, ?)", currentUserLogin, input.Title, input.Content, input.Category)
	if err != nil {
		return errors.NewInternal("Внутренняя ошибка: Ошибка создания блога.")
	}

	return nil
}

func (s *BlogStore) Update(id int, input blogs.BlogUpdate, currentUserLogin string) *errors.Error {
	if err := validateUpdating(id, input, s.db); err != nil {
		return err
	}

	_, err := s.db.Exec("update t_blogs set title = ?, content = ?, status = 'IN_REVIEW', category = ?, updated_at = now(), updated_by = ? where id = ?", input.Title, input.Content, input.Category, currentUserLogin, id)
	if err != nil {
		return errors.NewInternal("Внутренняя ошибка: Ошибка обновления блога.")
	}

	return nil
}

func (s *BlogStore) FindAll(page, size int) []blogs.BlogView {
	var entityList []blogs.Blog
	rows, err := s.db.Query("select * from t_blogs limit ? offset ?", size, page)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&entityList); err != nil {
			panic(err)
		}
	}

	return blogs.ToViewList(entityList)
}

func (s *BlogStore) FindById(id int) (*blogs.BlogView, *errors.Error) {
	var blog blogs.Blog
	row := s.db.QueryRow("select * from t_blogs where id = ?", id)

	if err := row.Scan(&blog); err != nil {
		return nil, errors.NewBadRequest(fmt.Sprintf("Блог с ID: %s не найден.", id))
	}

	return blog.ToView(), nil
}

func (s *BlogStore) Delete(id int) *errors.Error {
	var isExists bool
	row := s.db.QueryRow("select exists(select 1 from t_blogs where id = ?)", id)
	row.Scan(&isExists)

	if !isExists {
		return errors.NewBadRequest(fmt.Sprintf("Блог с ID: %s не найден.", id))
	}

	s.db.Exec("update t_blogs set is_deleted = true where id = ?", id)

	return nil
}

func validateUpdating(id int, input blogs.BlogUpdate, db *sql.DB) *errors.Error {
	var blog blogs.Blog
	row := db.QueryRow("select * from t_blogs where id = ? and is_deleted is false", id)

	if err := row.Scan(&blog); err != nil {
		return errors.NewBadRequest(fmt.Sprintf("Блог с ID: %s не найден.", id))
	}

	if blog.Status == blogs.IN_REVIEW {
		return errors.NewBadRequest("Нельзя обновить блог, который находится на рассмотрении модератора.")
	}

	if len(input.Title) == 0 || len(input.Content) == 0 {
		return errors.NewBadRequest("Пустой заголовок или содержание.")
	}

	return nil
}
