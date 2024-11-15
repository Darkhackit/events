// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (Users, error)
	GetUser(ctx context.Context, username pgtype.Text) (Users, error)
	GetUsers(ctx context.Context) ([]Users, error)
}

var _ Querier = (*Queries)(nil)
