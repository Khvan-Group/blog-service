package blogs

import (
	"github.com/dkhvan-dev/alabs_project/blog-service/internal/users/model"
	"time"
)

type Blog struct {
	Id        int         `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	CreatedBy model.User  `json:"created_by"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	Status    Status      `json:"status"`
	Category  string      `json:"category"`
	Likes     int         `json:"likes"`
	Favorites int         `json:"favorites"`
	Watches   int         `json:"watches"`
	UpdatedAt *time.Time  `json:"updated_at"`
	UpdatedBy *model.User `json:"updated_by"`
	DeletedAt *time.Time  `json:"deleted_at"`
	DeletedBy *model.User `json:"deleted_by"`
}

// DTOs
type BlogCreate struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type BlogUpdate struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type BlogView struct {
	Id        int             `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	CreatedBy *model.UserView `json:"created_by"`
	Title     string          `json:"title"`
	Content   string          `json:"content"`
	Status    Status          `json:"status"`
	Category  string          `json:"category"`
	Likes     int             `json:"likes"`
	Favorites int             `json:"favorites"`
	Watches   int             `json:"watches"`
	UpdatedAt *time.Time      `json:"updated_at"`
	UpdatedBy *model.UserView `json:"updated_by"`
}

// enums

type Status string
type Category string

const (
	// Statuses
	DRAFT     Status = "DRAFT"
	IN_REVIEW Status = "IN_REVIEW"
	ACTIVATED Status = "ACTIVATED"
	REJECTED  Status = "REJECTED"

	// Categories
	IT         Category = "IT"
	NEWS       Category = "NEWS"
	MANAGEMENT Category = "MANAGEMENT"
	BUSINESS   Category = "BUSINESS"
	GAMES      Category = "GAMES"
	TRAVEL     Category = "TRAVEL"
)

func (s Status) IsValid() bool {
	switch s {
	case DRAFT, IN_REVIEW, ACTIVATED, REJECTED:
		return true
	}

	return false
}

func (c Category) IsValid() bool {
	switch c {
	case IT, NEWS, MANAGEMENT, BUSINESS, GAMES, TRAVEL:
		return true
	}

	return false
}

// mapper

func (b *Blog) ToView() *BlogView {
	return &BlogView{
		Id:        b.Id,
		CreatedAt: b.CreatedAt,
		CreatedBy: b.CreatedBy.ToView(),
		Title:     b.Title,
		Content:   b.Content,
		Status:    b.Status,
		Category:  b.Category,
		Likes:     b.Likes,
		Favorites: b.Favorites,
		Watches:   b.Watches,
		UpdatedAt: b.UpdatedAt,
		UpdatedBy: b.UpdatedBy.ToView(),
	}
}

func ToViewList(list []Blog) []BlogView {
	var response []BlogView
	for _, l := range list {
		response = append(response, *l.ToView())
	}

	return response
}
