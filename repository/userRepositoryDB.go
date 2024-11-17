package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	db "github.com/Darkhackit/events/db/sqlc"
	"github.com/Darkhackit/events/domain"
	"github.com/Darkhackit/events/dto"
	"github.com/Darkhackit/events/events"
	"github.com/Darkhackit/events/sessions"
	token2 "github.com/Darkhackit/events/token"
	"github.com/Darkhackit/events/worker"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserRepositoryDB struct {
	q           *pgxpool.Pool
	PasetoToken *token2.PasetoToken
	distributor worker.TaskDistributor
	redisClient *sessions.RedisClient
}

func (us *UserRepositoryDB) Login(ctx context.Context, logins dto.LoginRequest) (*dto.UserResponse, error) {
	q := db.New(us.q)
	user, err := q.GetUser(ctx, pgtype.Text{
		String: logins.Username,
		Valid:  logins.Username != "",
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("username or password is not correct")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(logins.Password))
	if err != nil {
		return nil, fmt.Errorf("username or password is not correct")
	}

	token, payload, err := us.PasetoToken.CreateToken(user.Username.String, time.Hour*3)
	if err != nil {
		return nil, err
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	err = us.redisClient.CreateSession(ctx, payload.ID.String(), string(jsonPayload), time.Hour*2)
	result, err := q.GetUserRolesPermissions(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	permissionPayload, err := json.Marshal(result.Permissions)
	if err != nil {
		return nil, err
	}
	err = us.redisClient.CreateSession(ctx, fmt.Sprintf("permissions_%s", payload.ID.String()), string(permissionPayload), time.Hour*2)
	if err != nil {
		return nil, err
	}
	rolesPayload, err := json.Marshal(result.Roles)
	if err != nil {
		return nil, err
	}
	err = us.redisClient.CreateSession(ctx, fmt.Sprintf("roles_%s", payload.ID.String()), string(rolesPayload), time.Hour*2)
	if err != nil {
		return nil, err
	}
	userR := &dto.UserResponse{
		Email:    user.Email.String,
		Token:    token,
		Username: payload.Username,
	}

	return userR, nil

}

func (us *UserRepositoryDB) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	tx, err := us.q.BeginTx(ctx, pgx.TxOptions{})
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()
	arg := db.CreateUserParams{
		Username: pgtype.Text{
			String: user.Username,
			Valid:  user.Username != "",
		},
		Email: pgtype.Text{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Password: pgtype.Text{
			String: string(hash),
			Valid:  user.Password != "",
		},
	}
	q := db.New(tx)
	u, err := q.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	dUser := domain.User{
		Username: u.Username.String,
		Email:    u.Email.String,
		Password: u.Password.String,
	}

	events.Dispatch.Dispatch(events.UserCreatedEvent{User: dUser, TaskDistributor: us.distributor})
	return &dUser, nil

}

func (us *UserRepositoryDB) GetUsers(ctx context.Context) ([]domain.User, error) {
	q := db.New(us.q)
	rows, err := q.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, len(rows))
	for i, row := range rows {
		users[i] = domain.User{
			Username: row.Username.String,
			Email:    row.Email.String,
		}
	}
	return users, nil
}

func mapRoles(roles []domain.Role) []dto.RoleResponse {
	mappedRoles := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		mappedRoles[i] = dto.RoleResponse{
			ID:   role.ID,
			Name: role.Name,
		}
	}
	return mappedRoles
}
func mapPermissions(permissions []domain.Permission) []dto.PermissionResponse {
	mappedPermissions := make([]dto.PermissionResponse, len(permissions))
	for i, permission := range permissions {
		mappedPermissions[i] = dto.PermissionResponse{
			ID:   uint(permission.ID),
			Name: permission.Name,
		}
	}
	return mappedPermissions
}

func NewUserRepositoryDB(q *pgxpool.Pool, p *token2.PasetoToken, distributor worker.TaskDistributor, redisClient *sessions.RedisClient) *UserRepositoryDB {
	return &UserRepositoryDB{q: q, PasetoToken: p, distributor: distributor, redisClient: redisClient}
}
