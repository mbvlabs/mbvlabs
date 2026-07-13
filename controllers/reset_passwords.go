package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"mbvlabs/internal/hypermedia"
	"mbvlabs/internal/inertia"
	"mbvlabs/router"
	"mbvlabs/router/cookies"
	"mbvlabs/router/routes"
	"mbvlabs/services"
	"mbvlabs/views"

	"github.com/labstack/echo/v5"
)

type ResetPasswords struct {
	identity services.Identity
}

func NewResetPasswords(
	identity services.Identity,
) ResetPasswords {
	return ResetPasswords{identity}
}

func (rp ResetPasswords) RegisterRoutes(r *router.Router) error {
	errs := []error{}

	_, err := r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.PasswordNew.Path(),
		Name:    routes.PasswordNew.Name(),
		Handler: rp.New,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.PasswordCreate.Path(),
		Name:    routes.PasswordCreate.Name(),
		Handler: rp.Create,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.PasswordEdit.Path(),
		Name:    routes.PasswordEdit.Name(),
		Handler: rp.Edit,
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodPut,
		Path:    routes.PasswordUpdate.Path(),
		Name:    routes.PasswordUpdate.Name(),
		Handler: rp.Update,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (rp ResetPasswords) New(etx *echo.Context) error {
	setPrivateSEOHeaders(etx)
	return inertia.Page(etx, "Auth/ForgotPassword", inertia.Props{})
}

func (rp ResetPasswords) Create(etx *echo.Context) error {
	var payload struct {
		Email string `json:"email"`
	}

	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse password reset request payload",
			"error",
			err,
		)

		return hypermedia.RenderPage(etx, views.BadRequest())
	}

	if err := rp.identity.RequestResetPassword(
		etx.Request().Context(),
		services.RequestResetPasswordData{
			Email: payload.Email,
		},
	); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to request password reset",
			"error",
			err,
		)
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			"Failed to send password reset code",
		); flashErr != nil {
			return hypermedia.RenderPage(etx, views.InternalError())
		}

		return inertia.Redirect(etx, routes.PasswordNew.URL())
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"If an account exists with that email, you will receive password reset instructions.",
	); flashErr != nil {
		return hypermedia.RenderPage(etx, views.InternalError())
	}

	return inertia.Redirect(etx, routes.SessionNew.URL())
}

func (rp ResetPasswords) Edit(etx *echo.Context) error {
	setPrivateSEOHeaders(etx)
	etx.Response().Header().Set("Referrer-Policy", "strict-origin")

	token := etx.Param("token")
	if token == "" {
		if flashErr := cookies.AddFlash(
			etx,
			cookies.FlashError,
			"Invalid or missing reset token",
		); flashErr != nil {
			return hypermedia.RenderPage(etx, views.InternalError())
		}
		return inertia.Redirect(etx, routes.PasswordNew.URL())
	}

	return inertia.Page(etx, "Auth/ResetPassword", inertia.Props{"token": token})
}

func (rp ResetPasswords) Update(etx *echo.Context) error {
	var payload struct {
		Token           string `json:"resetPasswordToken"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err := etx.Bind(&payload); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"could not parse password reset payload",
			"error",
			err,
		)
		return hypermedia.RenderPage(etx, views.BadRequest())
	}

	if err := rp.identity.ResetPassword(
		etx.Request().Context(),
		services.ResetPasswordData{
			Token:           payload.Token,
			Password:        payload.Password,
			ConfirmPassword: payload.ConfirmPassword,
		},
	); err != nil {
		slog.ErrorContext(
			etx.Request().Context(),
			"failed to reset password",
			"error",
			err,
		)

		var errorMsg string
		switch {
		case errors.Is(err, services.ErrInvalidResetCode):
			errorMsg = "Invalid reset code"
		case errors.Is(err, services.ErrExpiredResetCode):
			errorMsg = "Reset code has expired"
		case errors.Is(err, services.ErrPasswordMismatch):
			errorMsg = "Passwords do not match"
		case errors.Is(err, services.ErrPasswordTooShort):
			errorMsg = "Password must be at least 8 characters"
		default:
			errorMsg = "Failed to reset password"
		}

		if flashErr := cookies.AddFlash(etx, cookies.FlashError, errorMsg); flashErr != nil {
			return hypermedia.RenderPage(etx, views.InternalError())
		}
		return inertia.Redirect(etx, routes.PasswordEdit.URL(payload.Token))
	}

	if flashErr := cookies.AddFlash(
		etx,
		cookies.FlashSuccess,
		"Password reset successfully! Please log in.",
	); flashErr != nil {
		return hypermedia.RenderPage(etx, views.InternalError())
	}

	return inertia.Redirect(etx, routes.SessionNew.URL())
}
