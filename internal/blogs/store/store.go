package store

import (
	"fmt"
	blogs "github.com/Khvan-Group/blog-service/internal/blogs/model"
	"github.com/Khvan-Group/blog-service/internal/users/model"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/jmoiron/sqlx"
	"strings"
)

type BlogStore struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *BlogStore {
	return &BlogStore{
		db: db,
	}
}

func (s *BlogStore) Create(input blogs.BlogCreate, currentUser model.JwtUser) *errors.CustomError {
	input.CreatedBy = currentUser.Login
	if len(input.Title) == 0 || len(input.Content) == 0 {
		return errors.NewBadRequest("Пустой заголовок или содержание.")
	}

	query := `insert into t_blogs(created_by, title, content, category)
			  values (:created_by, :title, :content, :category)
	`
	_, err := s.db.NamedExec(query, input)
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *BlogStore) Update(id int, input blogs.BlogUpdate, currentUser model.JwtUser) *errors.CustomError {
	input.Id = id
	input.UpdatedBy = currentUser.Login
	if err := validateUpdating(id, input, currentUser.Login, s.db); err != nil {
		return err
	}

	query := `update t_blogs set
                   title = :title,
                   content = :content,
                   status = 'DRAFT',
                   category = :category,
                   updated_at = now(),
                   updated_by = :updated_by
               where id = :id
	`
	_, err := s.db.NamedExec(query, input)
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *BlogStore) Send(id int, currentUser model.JwtUser) *errors.CustomError {
	var blog blogs.Blog
	err := s.db.Get(&blog, "select * from t_blogs where id = ?", id)

	if err != nil {
		return errors.NewBadRequest(fmt.Sprintf("Блог с ID: %d не найден.", id))
	}

	if blog.Status != blogs.DRAFT {
		return errors.NewBadRequest("На рассмотрение можно отправить блог, только в статусе 'Черновик'")
	}

	if currentUser.Role == "USER" && blog.CreatedBy.Login != currentUser.Login {
		return errors.NewForbidden("Недостаточно прав для отправки блога на рассмотрение.")
	}

	if _, err = s.db.Exec("update t_blogs set status = 'IN_REVIEW' where id =?", id); err != nil {
		panic(err)
	}

	return nil
}

func (s *BlogStore) FindAll(input blogs.BlogSearch) []blogs.BlogView {
	var response []blogs.BlogView
	query := buildQuery(input)

	err := s.db.Select(&response, query, input.Size, input.Page*input.Size)
	if err != nil {
		panic(err)
	}

	return response
}

func (s *BlogStore) FindById(id int, currentUser model.JwtUser) (*blogs.BlogView, *errors.CustomError) {
	var response blogs.BlogView
	err := s.db.Get(&response, "select * from t_blogs where id = ?", id)

	if err != nil {
		return nil, errors.NewBadRequest(fmt.Sprintf("Блог с ID: %d не найден.", id))
	}

	if currentUser.Role == "USER" && response.Status != blogs.ACTIVATED && response.CreatedBy.Login != currentUser.Login {
		return nil, errors.NewForbidden("Недостаточно прав для просмотра.")
	}

	return &response, nil
}

func (s *BlogStore) Delete(id int, currentUser model.JwtUser) *errors.CustomError {
	var blog blogs.Blog
	err := s.db.Get(&blog, "select * from t_blogs where id = ?", id)

	if err != nil {
		return errors.NewBadRequest(fmt.Sprintf("Блог с ID: %d не найден.", id))
	}

	if currentUser.Role == "USER" && blog.CreatedBy.Login != currentUser.Login {
		return errors.NewForbidden("Недостаточно прав для удаления блога.")
	}

	_, err = s.db.Exec("update t_blogs set is_deleted = true, deleted_at = now() where id = ?", id)
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *BlogStore) LikeOrFavorite(id int, currentUser model.JwtUser, action string) {
	query := `
		update t_blogs set likes = likes+1;
		insert into t_users_blogs (user_login, blog_id, favorites) values (?, ?, true)
	`

	if action == "favorites" {
		query = `
			update t_blogs set favorites = favorites+1;
			insert into t_users_blogs (user_login, blog_id, likes) values (?, ?, true);
		`
	}

	_, err := s.db.Exec(query, currentUser.Login, id)
	if err != nil {
		panic(err)
	}
}

func (s *BlogStore) Confirm(id int, status string, currentUser model.JwtUser) *errors.CustomError {
	if status != blogs.ACTIVATED && status != blogs.REJECTED {
		return errors.NewBadRequest("Неверный переданный статус.")
	}

	blog, err := s.FindById(id, currentUser)
	if err != nil {
		return err
	}

	if blog.Status != blogs.IN_REVIEW {
		return errors.NewBadRequest("Подтвердить блог можно только находящийся в статусе 'На рассмотрении'.")
	}

	_, execErr := s.db.Exec("update t_blogs set status = ? where id = ?", status)
	if execErr != nil {
		panic(err)
	}

	return nil
}

func validateUpdating(id int, input blogs.BlogUpdate, currentUserLogin string, db *sqlx.DB) *errors.CustomError {
	var blog blogs.Blog
	err := db.Get(&blog, "select * from t_blogs where id = ? and is_deleted is false", id)

	if err != nil {
		return errors.NewBadRequest(fmt.Sprintf("Блог с ID: %d не найден.", id))
	}

	if blog.CreatedBy.Login != currentUserLogin {
		return errors.NewForbidden("Вы не являетесь создателем данного блога.")
	}

	if blog.Status == blogs.IN_REVIEW {
		return errors.NewBadRequest("Нельзя обновить блог, который находится на рассмотрении модератора.")
	}

	if len(input.Title) == 0 || len(input.Content) == 0 {
		return errors.NewBadRequest("Пустой заголовок или содержание.")
	}

	return nil
}

func buildQuery(input blogs.BlogSearch) string {
	query := "select * from t_blogs "

	if input.CurrentUser.Role == "USER" {
		query += "where status = 'ACTIVATED' and is_deleted is false "
	}

	if input.Title != nil && len(strings.Trim(*input.Title, " ")) != 0 {
		query += "and title like '%" + *input.Title + "%' "
	}

	if input.Status != nil {
		query += fmt.Sprintf("and status = '%s' ", *input.Status)
	}

	if input.Category != nil {
		query += fmt.Sprintf("and category = '%s' ", *input.Category)
	}

	if len(input.SortBy) > 0 {
		sortQuery := " order by " + strings.Join(input.SortBy, ", ")
		query += sortQuery
	}

	return query + " limit $1 offset $2"
}
