package domain

import (
	"context"
)

type User struct {
	Username string
	Email    string
	Password string
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
}
