package services

import (
	"context"
	"errors"
	"fmt"

	"mbvlabs/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailNotVerified   = errors.New("email not verified")
)

type LoginData struct {
	Email    string
	Password string
}

func (i Identity) AuthenticateUser(
	ctx context.Context,
	data LoginData,
) (models.UserEntity, error) {
	user, err := models.User.FindByEmail(ctx, i.db.Executor(), data.Email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return models.UserEntity{}, ErrInvalidCredentials
		}

		return models.UserEntity{}, fmt.Errorf("find user by email: %w", err)
	}

	validPassword, err := user.ValidPassword(data.Password, i.pepper)
	if err != nil {
		return models.UserEntity{}, fmt.Errorf("validate password: %w", err)
	}

	if !validPassword {
		return models.UserEntity{}, ErrInvalidCredentials
	}

	if !user.HasValidatedEmail() {
		return models.UserEntity{}, ErrEmailNotVerified
	}

	return user, nil
}
