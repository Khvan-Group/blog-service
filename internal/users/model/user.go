package model

type User struct {
	Login      string  `json:"login" db:"login"`
	Password   string  `json:"password" db:"password"`
	Email      string  `json:"email" db:"email"`
	FirstName  string  `json:"firstName" db:"first_name"`
	MiddleName *string `json:"middleName" db:"middle_name"`
	LastName   string  `json:"lastName" db:"last_name"`
	Birthdate  string  `json:"birthdate" db:"birthdate"`
	Role       Role    `json:"role" db:"role"`
}

type Role struct {
	Code string `json:"code" db:"code"`
	Name string `json:"name" db:"name"`
}

// DTOs

type UserView struct {
	Login      string  `json:"login" db:"login"`
	Email      string  `json:"email" db:"email"`
	FirstName  string  `json:"firstName" db:"first_name"`
	MiddleName *string `json:"middleName" db:"middle_name"`
	LastName   string  `json:"lastName" db:"last_name"`
	Birthdate  string  `json:"birthdate" db:"birthdate"`
	Role       Role    `json:"role" db:"role"`
}

type JwtUser struct {
	Login string
	Role  string
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
