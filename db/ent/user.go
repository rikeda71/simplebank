package db

import (
	"context"

	"github.com/s14t284/simplebank/ent"
)

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func (store *SQLStore) CreateUser(ctx context.Context, arg CreateUserParams) (*ent.User, error) {
	user, err := store.entClient.User.Create().
		SetID(arg.Username).
		SetHashedPassword(arg.HashedPassword).
		SetFullName(arg.FullName).
		SetEmail(arg.Email).
		Save(ctx)

	return user, err
}

func (store *SQLStore) GetUser(ctx context.Context, username string) (*ent.User, error) {
	return store.entClient.User.Get(ctx, username)
}
