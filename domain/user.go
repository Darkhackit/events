package domain

import (
	"context"
	"github.com/Darkhackit/events/dto"
)

type User struct {
	Username string
	Email    string
	Password string
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
	Login(ctx context.Context, logins dto.LoginRequest) (*dto.UserResponse, error)
}
