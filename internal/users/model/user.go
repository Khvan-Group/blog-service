package model

type User struct {
	Login      string  `json:"login"`
	Password   string  `json:"password"`
	Email      string  `json:"email"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
	Birthdate  string  `json:"birthdate"`
	Role       Role    `json:"role"`
}

type Role struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// DTOs

type UserView struct {
	Login      string  `json:"login"`
	Email      string  `json:"email"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName"`
	Birthdate  string  `json:"birthdate"`
	Role       Role    `json:"role"`
}

// mapper

func (u *User) ToView() *UserView {
	return &UserView{
		Login:      u.Login,
		Email:      u.Email,
		FirstName:  u.FirstName,
		MiddleName: u.MiddleName,
		LastName:   u.LastName,
		Birthdate:  u.Birthdate,
		Role:       u.Role,
	}
}
