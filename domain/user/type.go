package user

import (
	"ecommerce/dto/request"
)

type Repository interface {
	Persist(u *User) (*User, error)

	GetUserByUserPass(email, password string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetByIds(ids []int) ([]User, error)
}

type Service interface {
	CreateUserIfNotAny(req request.CreateUserRequest) (*User, error)
	GetByIds(ids []int) ([]User, error)

	Login(email, password string) (*User, string, error)
}
