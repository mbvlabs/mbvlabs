package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"mbvlabs/config"
	"mbvlabs/database"
	"mbvlabs/internal/storage"
	"mbvlabs/models"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

type seedUserParams struct {
	fx.In

	Config config.Config
	DB     storage.Pool
}

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("load .env: %v", err)
	}

	var data seedUserData
	flag.StringVar(&data.email, "email", os.Getenv("SEED_USER_EMAIL"), "email for the user to seed")
	flag.StringVar(&data.password, "password", os.Getenv("SEED_USER_PASSWORD"), "password for the user to seed")
	flag.BoolVar(&data.admin, "admin", envBool("SEED_USER_ADMIN", true), "whether the seeded user is an admin")
	flag.BoolVar(&data.verified, "verified", envBool("SEED_USER_VERIFIED", true), "whether the seeded user's email is already verified")
	flag.Parse()

	ctx := context.Background()
	app := fx.New(
		fx.NopLogger,
		fx.Provide(func() context.Context { return ctx }),
		config.Module,
		database.Module,
		fx.Invoke(func(params seedUserParams) error {
			defer params.DB.Close()
			return seedUser(ctx, params, data)
		}),
	)

	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
}

type seedUserData struct {
	email    string
	password string
	admin    bool
	verified bool
}

func seedUser(ctx context.Context, params seedUserParams, data seedUserData) error {
	if data.email == "" {
		return errors.New("missing user email: pass -email or set SEED_USER_EMAIL")
	}
	if data.password == "" {
		return errors.New("missing user password: pass -password or set SEED_USER_PASSWORD")
	}

	tx, err := params.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := models.User.FindByEmail(ctx, tx, data.email); err == nil {
		return fmt.Errorf("user %q already exists", data.email)
	} else if !errors.Is(err, models.ErrNotFound) {
		return fmt.Errorf("check existing user: %w", err)
	}

	user, err := models.User.Create(ctx, tx, params.Config.Auth.Pepper, models.CreateUserData{
		Email: data.email,
		PasswordPair: models.PasswordPair{
			Password:        data.password,
			ConfirmPassword: data.password,
		},
	})
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	emailValidatedAt := sql.NullTime{}
	if data.verified {
		emailValidatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	user, err = models.User.Update(ctx, tx, models.UpdateUserData{
		ID:               user.ID,
		Email:            user.Email,
		EmailValidatedAt: emailValidatedAt,
		Password:         user.Password,
		IsAdmin:          data.admin,
	})
	if err != nil {
		return fmt.Errorf("update seeded user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	fmt.Printf("seeded user %s (admin=%t verified=%t)\n", user.Email, user.IsAdmin, user.HasValidatedEmail())
	return nil
}

func envBool(name string, fallback bool) bool {
	switch os.Getenv(name) {
	case "1", "true", "TRUE", "yes", "YES":
		return true
	case "0", "false", "FALSE", "no", "NO":
		return false
	default:
		return fallback
	}
}
