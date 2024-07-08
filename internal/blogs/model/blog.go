package blogs

import (
	"github.com/Khvan-Group/blog-service/internal/users/model"
	"time"
)

type Blog struct {
	Id        int             `json:"id" db:"id"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	CreatedBy *model.UserView `json:"created_by" db:"created_by"`
	Title     string          `json:"title" db:"title"`
	Content   string          `json:"content" db:"content"`
	Status    string          `json:"status" db:"status"`
	Category  string          `json:"category" db:"category"`
	Likes     int             `json:"likes" db:"likes"`
	Favorites int             `json:"favorites" db:"favorites"`
	Watches   int             `json:"watches" db:"watches"`
	UpdatedAt *time.Time      `json:"updated_at" db:"updated_at"`
	UpdatedBy *model.UserView `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time      `json:"deleted_at" db:"deleted_at"`
	DeletedBy *model.User     `json:"deleted_by" db:"deleted_by"`
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
	Id        int             `json:"id" db:"id"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	CreatedBy *model.UserView `json:"created_by" db:"created_by"`
	Title     string          `json:"title" db:"title"`
	Content   string          `json:"content" db:"content"`
	Status    string          `json:"status" db:"status"`
	Category  string          `json:"category" db:"category"`
	Likes     int             `json:"likes" db:"likes"`
	Favorites int             `json:"favorites" db:"favorites"`
	Watches   int             `json:"watches" db:"watches"`
	UpdatedAt *time.Time      `json:"updated_at" db:"updated_at"`
	UpdatedBy *model.UserView `json:"updated_by" db:"updated_by"`
}

type BlogSearch struct {
	Page        int
	Size        int
	SortBy      []string
	Title       *string `json:"title"`
	Status      *string `json:"status"`
	Category    *string `json:"category"`
	CurrentUser model.JwtUser
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

// mapper

func (b *Blog) ToView() *BlogView {
	return &BlogView{
		Id:        b.Id,
		CreatedAt: b.CreatedAt,
		CreatedBy: b.CreatedBy,
		Title:     b.Title,
		Content:   b.Content,
		Status:    b.Status,
		Category:  b.Category,
		Likes:     b.Likes,
		Favorites: b.Favorites,
		Watches:   b.Watches,
		UpdatedAt: b.UpdatedAt,
		UpdatedBy: b.UpdatedBy,
	}
}

func ToViewList(list []Blog) []BlogView {
	var response []BlogView
	for _, l := range list {
		response = append(response, *l.ToView())
	}

	return response
}
