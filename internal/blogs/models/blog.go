package blogs

import (
	categoryModel "github.com/Khvan-Group/blog-service/internal/categories/models"
	"github.com/Khvan-Group/blog-service/internal/clients"
	"github.com/Khvan-Group/blog-service/internal/common/models"
	"github.com/Khvan-Group/common-library/errors"
	commonModels "github.com/Khvan-Group/common-library/models"
	"github.com/go-resty/resty/v2"
	"time"
)

type Blog struct {
	Id        int                    `json:"id" db:"id"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	CreatedBy *models.UserView       `json:"created_by" db:"created_by"`
	Title     string                 `json:"title" db:"title"`
	Content   string                 `json:"content" db:"content"`
	Status    string                 `json:"status" db:"status"`
	Category  categoryModel.Category `json:"category" db:"category"`
	Likes     int                    `json:"likes" db:"likes"`
	Favorites int                    `json:"favorites" db:"favorites"`
	UpdatedAt *time.Time             `json:"updated_at" db:"updated_at"`
	UpdatedBy *models.UserView       `json:"updated_by" db:"updated_by"`
}

// DTOs
type BlogCreate struct {
	Title     string `json:"title" db:"title"`
	Content   string `json:"content" db:"content"`
	Category  string `json:"category" db:"category"`
	CreatedBy string `db:"created_by"`
}

type BlogUpdate struct {
	Id        int    `db:"id"`
	Title     string `json:"title" db:"title"`
	Content   string `json:"content" db:"content"`
	Category  string `json:"category" db:"category"`
	UpdatedBy string `db:"updated_by"`
}

type BlogView struct {
	Id        int                    `json:"id" db:"id"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	CreatedBy *models.UserView       `json:"created_by" db:"created_by"`
	Title     string                 `json:"title" db:"title"`
	Content   string                 `json:"content" db:"content"`
	Status    string                 `json:"status" db:"status"`
	Category  categoryModel.Category `json:"category" db:"category"`
	Likes     int                    `json:"likes" db:"likes"`
	Favorites int                    `json:"favorites" db:"favorites"`
	UpdatedAt *time.Time             `json:"updated_at" db:"updated_at"`
	UpdatedBy *models.UserView       `json:"updated_by" db:"updated_by"`
}

type BlogSearch struct {
	Page        int
	Size        int
	SortBy      []commonModels.SortField
	Title       *string `json:"title"`
	Status      *string `json:"status"`
	Category    *string `json:"category"`
	CurrentUser models.JwtUser
}

const (
	// Statuses
	DRAFT     = "DRAFT"
	IN_REVIEW = "IN_REVIEW"
	ACTIVATED = "ACTIVATED"
	REJECTED  = "REJECTED"

	// Categories
	IT         = "IT"
	NEWS       = "NEWS"
	MANAGEMENT = "MANAGEMENT"
	BUSINESS   = "BUSINESS"
	GAMES      = "GAMES"
	TRAVEL     = "TRAVEL"
)

func IsValidStatus(status string) bool {
	switch status {
	case DRAFT, IN_REVIEW, ACTIVATED, REJECTED:
		return true
	}

	return false
}

func IsValidCategory(category string) bool {
	switch category {
	case IT, NEWS, MANAGEMENT, BUSINESS, GAMES, TRAVEL:
		return true
	}

	return false
}

func IsValidCategoryList(list []string) bool {
	for _, c := range list {
		if !IsValidCategory(c) {
			return false
		}
	}

	return true
}

func (b *BlogView) FillUserInfo(client *resty.Client) *errors.CustomError {
	createdBy, err := clients.GetUserByLogin(b.CreatedBy.Login, client)
	if err != nil {
		return err
	}

	if len(b.UpdatedBy.Login) > 0 {
		updatedBy, err := clients.GetUserByLogin(b.UpdatedBy.Login, client)
		if err != nil {
			return err
		}

		b.UpdatedBy = updatedBy
	}

	b.CreatedBy = createdBy
	return nil
}
