package store

import (
	"fmt"
	blogs "github.com/Khvan-Group/blog-service/internal/blogs/models"
	"github.com/Khvan-Group/blog-service/internal/common/models"
	"github.com/Khvan-Group/blog-service/internal/core/rabbitmq"
	"github.com/Khvan-Group/blog-service/internal/db"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/errors"
	commonModels "github.com/Khvan-Group/common-library/models"
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	LIKE     = "LIKE"
	FAVORITE = "FAVORITE"
)

type BlogStore struct {
	db     *sqlx.DB
	client *resty.Client
}

func New(db *sqlx.DB) *BlogStore {
	return &BlogStore{
		db:     db,
		client: resty.New(),
	}
}

type blogValidate struct {
	CreatedBy string `json:"created_by" db:"created_by"`
	Status    string `json:"status" db:"status"`
}

func (s *BlogStore) Create(input blogs.BlogCreate, currentUser models.JwtUser) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		input.CreatedBy = currentUser.Login
		var existsCategory bool
		if len(input.Title) == 0 || len(input.Content) == 0 {
			return errors.NewBadRequest("Пустой заголовок или содержание.")
		}

		tx.Get(&existsCategory, "select exists(select 1 from t_blog_categories where code = $1)", input.Category)
		if !existsCategory {
			return errors.NewBadRequest("Такой категории блогов не существует.")
		}

		query := `insert into t_blogs(created_by, title, content, category)
			  values (:created_by, :title, :content, :category)`
		_, err := tx.NamedExec(query, input)
		if err != nil {
			return errors.NewInternal("Failed to insert blog transaction")
		}

		return nil
	})
}

func (s *BlogStore) Update(id int, input blogs.BlogUpdate, currentUser models.JwtUser) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
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
		_, err := tx.NamedExec(query, input)
		if err != nil {
			return errors.NewInternal("Failed to update blog")
		}

		return nil
	})
}

func (s *BlogStore) Send(id int, currentUser models.JwtUser) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var blog blogValidate
		err := tx.Get(&blog, "select created_by, status from t_blogs where id = $1", id)

		if err != nil {
			return errors.NewBadRequest(fmt.Sprintf("Блог с ID: %d не найден.", id))
		}

		if blog.Status != blogs.DRAFT {
			return errors.NewBadRequest("На рассмотрение можно отправить блог, только в статусе 'Черновик'")
		}

		if currentUser.Role == "USER" && blog.CreatedBy != currentUser.Login {
			return errors.NewForbidden("Недостаточно прав для отправки блога на рассмотрение.")
		}

		if _, err = tx.Exec("update t_blogs set status = 'IN_REVIEW', updated_at = now(), updated_by = $2 where id = $1", id, currentUser.Login); err != nil {
			return errors.NewInternal("Failed to update blog status")
		}

		return nil
	})
}

func (s *BlogStore) FindAll(input blogs.BlogSearch) commonModels.Page {
	var blogList []blogs.BlogView
	var totalElements int
	query := buildQuery(input)

	err := s.db.Select(&blogList, query, input.Size, input.Page*input.Size)
	if err != nil || blogList == nil {
		blogList = make([]blogs.BlogView, 0)
	}

	query = "select count(*) from t_blogs "
	query += addQueryByValidate(query, input)
	err = s.db.Get(&totalElements, query)
	if err != nil {
		totalElements = len(blogList)
	}

	for i := range blogList {
		if fillError := blogList[i].FillUserInfo(s.client); fillError != nil {
			panic(fillError)
		}
	}

	return commonModels.Page{
		Result:        blogList,
		Page:          input.Page,
		Size:          input.Size,
		TotalElements: totalElements,
	}
}

func (s *BlogStore) FindById(id int, currentUser models.JwtUser) (*blogs.BlogView, *errors.CustomError) {
	var response blogs.BlogView
	query := `
		select b.id, b.created_at, b.created_by as "created_by.login", b.updated_at, 
		       coalesce(b.updated_by, '') as "updated_by.login", b.title, b.content, b.status, 
		       b.category as "category.code", bc.name as "category.name", b.likes, b.favorites 
		from t_blogs b
		inner join t_blog_categories bc on bc.code = b.category
		where b.id = $1
	`
	err := s.db.Get(&response, query, id)

	if err != nil {
		return nil, errors.NewBadRequest(fmt.Sprintf("Блог с ID: %d не найден.", id))
	}

	if currentUser.Role == "USER" && response.Status != blogs.ACTIVATED && response.CreatedBy.Login != currentUser.Login {
		return nil, errors.NewForbidden("Недостаточно прав для просмотра.")
	}

	if fillError := response.FillUserInfo(s.client); fillError != nil {
		return nil, fillError
	}

	return &response, nil
}

func (s *BlogStore) Delete(id int, currentUser models.JwtUser) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var createdBy string
		err := tx.Get(&createdBy, "select created_by from t_blogs where id = $1", id)

		if err != nil {
			return errors.NewBadRequest(fmt.Sprintf("Блог с ID: %d не найден.", id))
		}

		if currentUser.Role == "USER" && createdBy != currentUser.Login {
			return errors.NewForbidden("Недостаточно прав для удаления блога.")
		}

		_, err = tx.Exec("delete from t_blogs where id = $1", id)
		if err != nil {
			panic(err)
		}

		return nil
	})
}

func (s *BlogStore) DeleteAllByUsername(username string) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		_, err := tx.Exec("delete from t_blogs where created_by = $1", username)
		if err != nil {
			return errors.NewInternal("Failed to delete blog transaction")
		}

		return nil
	})
}

func (s *BlogStore) LikeOrFavorite(id int, currentUser models.JwtUser, action string) {
	var exists bool
	s.db.Get(&exists, "select exists(select 1 from t_users_blogs where blog_id = $1 and user_login = $2)", id, currentUser.Login)

	blogQuery := "update t_blogs set likes = likes+1, updated_at = now(), updated_by = $1 where id = $2; "
	userBlogQuery := "update t_users_blogs set likes = true where user_login = $1 and blog_id = $2 "
	if !exists {
		userBlogQuery = "insert into t_users_blogs (user_login, blog_id, likes) values ($1, $2, true)"
	}

	if action == "favorites" {
		blogQuery = "update t_blogs set favorites = favorites+1, updated_at = now(), updated_by = $1 where id = $2; "

		if exists {
			userBlogQuery = "update t_users_blogs set favorites = true where user_login = $1 and blog_id = $2 "
		} else {
			userBlogQuery = "insert into t_users_blogs (user_login, blog_id, favorites) values ($1, $2, true);"
		}
	}

	db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		_, err := tx.Exec(blogQuery, currentUser.Login, id)
		if err != nil {
			panic(err)
		}

		_, err = tx.Exec(userBlogQuery, currentUser.Login, id)
		if err != nil {
			panic(err)
		}

		return nil
	})
}

func (s *BlogStore) Confirm(id int, status string, currentUser models.JwtUser) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
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

		_, execErr := tx.Exec("update t_blogs set status = $1, updated_at = now() where id = $2", status, id)
		if execErr != nil {
			panic(err)
		}

		if status == blogs.ACTIVATED {
			msg := models.WalletUpdate{
				Username: blog.CreatedBy.Login,
				Total:    100,
				Action:   models.WALLET_TOTAL_ADD,
			}

			if err = rabbitmq.SendToUpdateWallet(msg); err != nil {
				return err
			}
		}

		return nil
	})
}

func validateUpdating(id int, input blogs.BlogUpdate, currentUserLogin string, db *sqlx.DB) *errors.CustomError {
	var blog blogValidate
	query := `
		select b.created_by, b.status from t_blogs b
		where b.id = $1
	`
	err := db.Get(&blog, query, id)

	if err != nil {
		return errors.NewBadRequest(fmt.Sprintf("Блог с ID: %d не найден.", id))
	}

	if blog.CreatedBy != currentUserLogin {
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
	query := `
		select b.id, b.created_at, b.created_by as "created_by.login", b.updated_at, 
		       coalesce(b.updated_by, '') as "updated_by.login", b.title, b.content, b.status, 
		       b.category as "category.code", bc.name as "category.name", b.likes, b.favorites 
		from t_blogs b
		inner join t_blog_categories bc on bc.code = b.category
	`

	query = addQueryByValidate(query, input)

	if len(input.SortBy) > 0 {
		query += " order by "
		sortJoins := make([]string, 0)

		for _, field := range input.SortBy {
			sortJoins = append(sortJoins, fmt.Sprintf("%s %s ", field.SortBy, field.Direction))
		}

		query += strings.Join(sortJoins, ", ")
	}

	return query + " limit $1 offset $2"
}

func addQueryByValidate(query string, input blogs.BlogSearch) string {
	query += "where true "
	if input.CurrentUser.Role != constants.MODERATOR && input.CurrentUser.Role != constants.ADMIN {
		query += "status = 'ACTIVATED' "
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

	return query
}
