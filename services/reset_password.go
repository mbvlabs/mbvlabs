package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"mbvlabs/config"
	"mbvlabs/email"
	"mbvlabs/models"
	"mbvlabs/queue/jobs"
	"mbvlabs/router/routes"
)

const userResetPassword = "user_password_reset"

var (
	ErrInvalidResetCode = errors.New("invalid reset code")
	ErrExpiredResetCode = errors.New("reset code has expired")
	ErrPasswordMismatch = errors.New("passwords do not match")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)

type RequestResetPasswordData struct {
	Email string
}

func (i Identity) RequestResetPassword(
	ctx context.Context,
	data RequestResetPasswordData,
) error {
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin password reset request transaction: %w", err)
	}

	user, err := models.User.FindByEmail(ctx, tx, data.Email)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, models.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("find password reset user: %w", err)
	}

	meta, err := json.Marshal(map[string]string{
		"email": user.Email,
	})
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("marshal reset token metadata: %w", err)
	}

	token, err := models.Token.Create(
		ctx,
		tx,
		i.pepper,
		userResetPassword,
		time.Now().Add(1*time.Hour),
		meta,
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("create password reset token: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit password reset request transaction: %w", err)
	}

	resetURL := fmt.Sprintf("%s%s", config.BaseURL, routes.PasswordEdit.URL(token))

	rpEmail := email.ResetPassword{ResetURL: resetURL}

	html, err := rpEmail.ToHTML()
	if err != nil {
		return fmt.Errorf("render reset password email html: %w", err)
	}

	text, err := rpEmail.ToText()
	if err != nil {
		return fmt.Errorf("render reset password email text: %w", err)
	}

	if _, err := i.insertOnly.Insert(ctx, jobs.SendTransactionalEmailArgs{
		Data: email.TransactionalData{
			To:       user.Email,
			From:     config.DefaultSenderSignature,
			Subject:  "Reset Your Password",
			HTMLBody: html,
			TextBody: text,
		},
	}, nil); err != nil {
		return fmt.Errorf("queue reset password email: %v", err)
	}

	return nil
}

type ResetPasswordData struct {
	Token           string
	Password        string
	ConfirmPassword string
}

func (i Identity) ResetPassword(
	ctx context.Context,
	data ResetPasswordData,
) error {
	if data.Password != data.ConfirmPassword {
		return ErrPasswordMismatch
	}

	if len(data.Password) < 8 {
		return ErrPasswordTooShort
	}

	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin password reset transaction: %w", err)
	}

	token, err := models.Token.FindByScopeAndHash(
		ctx,
		tx,
		i.pepper,
		userResetPassword,
		data.Token,
	)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, models.ErrNotFound) {
			return ErrInvalidResetCode
		}
		return fmt.Errorf("find password reset token: %w", err)
	}

	if !token.IsValid(data.Token, i.pepper) {
		_ = tx.Rollback()
		return ErrExpiredResetCode
	}

	var meta map[string]string
	if err := json.Unmarshal(token.MetaData, &meta); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("unmarshal reset token metadata: %w", err)
	}

	emailAddr, ok := meta["email"]
	if !ok {
		_ = tx.Rollback()
		return errors.New("reset token metadata missing email")
	}

	user, err := models.User.FindByEmail(ctx, tx, emailAddr)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, models.ErrNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("find reset user: %w", err)
	}

	hashedPassword, err := models.HashPassword(data.Password, i.pepper)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("hash reset password: %w", err)
	}

	_, err = models.User.Update(ctx, tx, models.UpdateUserData{
		ID:               user.ID,
		Email:            user.Email,
		EmailValidatedAt: user.EmailValidatedAt,
		Password:         []byte(hashedPassword),
		IsAdmin:          user.IsAdmin,
	})
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("update reset password: %w", err)
	}

	if err := models.Token.Destroy(ctx, tx, token.ID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("destroy reset token: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit password reset transaction: %w", err)
	}

	return nil
}
