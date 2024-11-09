package respository

import (
	"context"
	db "github.com/Darkhackit/events/db/sqlc"
	"github.com/Darkhackit/events/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryDB struct {
	q *db.Queries
}

func (us *UserRepositoryDB) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
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

	u, err := us.q.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	duser := &domain.User{
		Username: u.Username.String,
		Email:    u.Email.String,
		Password: u.Password.String,
	}

	return duser, nil

}

func NewUserRepositoryDB(q *db.Queries) UserRepositoryDB {
	return UserRepositoryDB{q: q}
}
