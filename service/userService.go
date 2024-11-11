package service

import (
	"context"
	"github.com/Darkhackit/events/domain"
	"github.com/Darkhackit/events/dto"
)

type UserService interface {
	CreateUser(ctx context.Context, user dto.UserRequest) (*dto.UserResponse, error)
	GetUsers(ctx context.Context) ([]dto.UserResponse, error)
	LoginUser(ctx context.Context, request dto.LoginRequest) (*dto.UserResponse, error)
}

type DefaultUserService struct {
	repo domain.UserRepository
}

func (d *DefaultUserService) LoginUser(ctx context.Context, request dto.LoginRequest) (*dto.UserResponse, error) {
	token, err := d.repo.Login(ctx, dto.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (d *DefaultUserService) GetUsers(ctx context.Context) ([]dto.UserResponse, error) {
	u, err := d.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]dto.UserResponse, len(u))
	for i, user := range u {
		users[i] = dto.UserResponse{
			Email:    user.Email,
			Username: user.Username,
		}
	}
	return users, nil
}

func (d *DefaultUserService) CreateUser(ctx context.Context, user dto.UserRequest) (*dto.UserResponse, error) {
	duser := domain.User{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}

	u, err := d.repo.CreateUser(ctx, duser)
	if err != nil {
		return nil, err
	}

	dtUser := dto.UserResponse{
		Email:    u.Email,
		Username: u.Username,
	}

	return &dtUser, nil
}

func NewUserService(repo domain.UserRepository) *DefaultUserService {
	return &DefaultUserService{repo: repo}
}
