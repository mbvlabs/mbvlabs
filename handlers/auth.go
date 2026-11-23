package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"mbvlabs/middleware"
	"mbvlabs/payloads"
	"mbvlabs/routes"
	"mbvlabs/services"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/cache"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
	"github.com/mbvlabs/andurel/pkg/routing"
	"github.com/mbvlabs/andurel/pkg/telemetry"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type AuthHandler struct {
	identity services.Identity
	renderer *inertia.Renderer
	cache    cache.Store
	hosts    *routing.Hosts
}

// NewAuthHandler takes the host registry so a sign-in can send the visitor back
// to a page on another host (see middleware.ConsumeReturnTo).
func NewAuthHandler(
	identity services.Identity,
	renderer *inertia.Renderer,
	store cache.Store,
	hosts *routing.Hosts,
) AuthHandler {
	return AuthHandler{identity: identity, renderer: renderer, cache: store, hosts: hosts}
}

const (
	authAttemptLimit  = 5
	authAttemptWindow = 10 * time.Minute
)

// Register binds each route on the host it declares with rt.Route, so the
// route value, not a copied path, decides the domain it is served on.
func (h AuthHandler) Register(rt *anduin.Router) {
	rt.Route(routes.SessionNew).Get(h.SessionNew)
	rt.Route(routes.SessionCreate).Post(
		h.SessionCreate,
		h.throttle(h.tooManyAttempts("Auth/Login", "email")),
	)
	rt.Route(routes.SessionDestroy).Delete(h.SessionDestroy)
	// The POST route lets a plain HTML form sign out without JavaScript; the
	// DELETE route stays for the Inertia link.
	rt.Route(routes.SessionDestroy).Post(h.SessionDestroy)

	rt.Route(routes.PasswordNew).Get(h.PasswordNew)
	rt.Route(routes.PasswordCreate).Post(
		h.PasswordCreate,
		h.throttle(h.tooManyAttempts("Auth/ResetPasswordRequest", "email")),
	)
	rt.Route(routes.PasswordEdit).Get(h.PasswordEdit)
	rt.Route(routes.PasswordUpdate).Put(
		h.PasswordUpdate,
		h.throttle(h.tooManyPasswordResets),
	)
}

// throttle rate-limits an auth endpoint per client IP; onLimited re-renders the
// submitting form so Inertia shows the error inline.
func (h AuthHandler) throttle(onLimited anduin.Handler) anduin.Middleware {
	return anduin.RateLimit(h.cache, authAttemptLimit, authAttemptWindow, onLimited)
}

func (h AuthHandler) tooManyAttempts(component, field string) anduin.Handler {
	return func(ctx context.Context) error {
		return h.renderTooManyAttempts(ctx, component, field, inertia.Props{})
	}
}

func (h AuthHandler) tooManyPasswordResets(ctx context.Context) error {
	var payload payloads.PasswordResetPayload
	_ = anduin.DecodeJSON(ctx, &payload)

	return h.renderTooManyAttempts(
		ctx,
		"Auth/ResetPassword",
		"resetPasswordToken",
		inertia.FromStruct(payloads.PasswordResetPageProps{Token: payload.Token}),
	)
}

func (h AuthHandler) renderTooManyAttempts(
	ctx context.Context,
	component string,
	field string,
	props inertia.Props,
) error {
	return h.renderer.Page(ctx, component, props).
		ValidationErrors(map[string]string{
			field: "Too many attempts, try again later.",
		}).
		Status(http.StatusTooManyRequests).
		Render()
}

func (h AuthHandler) SessionNew(ctx context.Context) error {
	setPrivateSEO(ctx)
	return h.renderer.Page(ctx, "Auth/Login", inertia.Props{}).Render()
}

func setPrivateSEO(ctx context.Context) {
	anduin.SetHeader(ctx, "X-Robots-Tag", "noindex, nofollow")
}

func (h AuthHandler) SessionCreate(ctx context.Context) error {
	ctx, span := telemetry.From(anduin.Request(ctx), "sessions.create")
	defer span.End()

	var payload payloads.LoginPayload
	if err := anduin.DecodeJSON(ctx, &payload); err != nil {
		return anduin.BadRequest(err)
	}

	entity, err := h.identity.AuthenticateUser(
		ctx,
		services.LoginData{
			Email:    payload.Email,
			Password: payload.Password,
		},
	)
	if err != nil {
		message := "Failed to log in"
		switch {
		case errors.Is(err, services.ErrEmailNotVerified):
			message = "Please verify your email before logging in"
		case errors.Is(err, services.ErrInvalidCredentials):
			message = "Invalid email or password"
		}
		return h.authFlash(ctx, kiks.FlashError, message, routes.SessionNew.URL())
	}

	if err := kiks.Set(ctx, &middleware.Session{
		UserID:          entity.ID.String(),
		SessionVersion:  entity.SessionVersion,
		IsAuthenticated: true,
	}); err != nil {
		return err
	}
	if err := kiks.AddFlash(ctx, kiks.FlashSuccess, "Successfully logged in!"); err != nil {
		return err
	}

	return h.signedInRedirect(ctx)
}

// signedInRedirect sends a freshly signed-in visitor to the page they asked for
// before sign-in, or to the home page. The return-target cookie is cleared
// either way. It uses Location rather than a plain redirect because sign-in is
// an Inertia request: Location makes the client visit the target itself, on
// this host or another one.
func (h AuthHandler) signedInRedirect(ctx context.Context) error {
	location, err := middleware.ConsumeReturnTo(ctx, h.hosts, routes.HomePage.URL())
	if err != nil {
		return err
	}

	return h.renderer.Location(ctx, location)
}

func (h AuthHandler) SessionDestroy(ctx context.Context) error {
	ctx, span := telemetry.From(anduin.Request(ctx), "sessions.destroy")
	defer span.End()

	if err := kiks.Destroy[*middleware.Session](ctx); err != nil {
		return err
	}
	if err := kiks.Destroy[middleware.ReturnTo](ctx); err != nil {
		return err
	}
	return h.authFlash(ctx, kiks.FlashSuccess, "Successfully logged out!", routes.SessionNew.URL())
}

func (h AuthHandler) authFlash(ctx context.Context, level kiks.FlashType, message, location string) error {
	if err := kiks.AddFlash(ctx, level, message); err != nil {
		return err
	}
	return h.renderer.Redirect(ctx, location, http.StatusSeeOther)
}

func (h AuthHandler) PasswordNew(ctx context.Context) error {
	setPrivateSEO(ctx)
	return h.renderer.Page(ctx, "Auth/ResetPasswordRequest", inertia.Props{}).Render()
}

func (h AuthHandler) PasswordCreate(ctx context.Context) error {
	ctx, span := telemetry.From(anduin.Request(ctx), "reset_passwords.create")
	defer span.End()

	var payload payloads.PasswordResetRequestPayload
	if err := anduin.DecodeJSON(ctx, &payload); err != nil {
		return anduin.BadRequest(err)
	}

	if err := h.identity.RequestResetPassword(
		ctx,
		services.RequestResetPasswordData{
			Email: payload.Email,
		},
	); err != nil {
		telemetry.Error(
			ctx,
			"failed to request password reset",
			"error",
			err,
		)

		return h.authFlash(ctx, kiks.FlashError, "Failed to send password reset code", routes.PasswordNew.URL())
	}

	return h.authFlash(ctx, kiks.FlashSuccess, "If an account exists with that email, you will receive password reset instructions.", routes.SessionNew.URL())
}

func (h AuthHandler) PasswordEdit(ctx context.Context) error {
	setPrivateSEO(ctx)
	anduin.SetHeader(ctx, "Referrer-Policy", "strict-origin")

	token := anduin.PathValue(ctx, "token")
	if token == "" {
		return h.renderer.Redirect(ctx, routes.PasswordNew.URL(), http.StatusSeeOther)
	}

	return h.renderer.Page(
		ctx,
		"Auth/ResetPassword",
		inertia.FromStruct(payloads.PasswordResetPageProps{Token: token}),
	).Render()
}

func (h AuthHandler) PasswordUpdate(ctx context.Context) error {
	ctx, span := telemetry.From(anduin.Request(ctx), "reset_passwords.update")
	defer span.End()

	var payload payloads.PasswordResetPayload
	if err := anduin.DecodeJSON(ctx, &payload); err != nil {
		return anduin.BadRequest(err)
	}

	if err := h.identity.ResetPassword(
		ctx,
		services.ResetPasswordData{
			Token:           payload.Token,
			Password:        payload.Password,
			ConfirmPassword: payload.ConfirmPassword,
		},
	); err != nil {
		message := "Failed to reset password"
		_, isValidation := validation.As(err)
		switch {
		case errors.Is(err, services.ErrInvalidResetCode):
			message = "Invalid reset code"
		case errors.Is(err, services.ErrExpiredResetCode):
			message = "Reset code has expired"
		case isValidation && payload.Password != payload.ConfirmPassword:
			message = "Passwords do not match"
		case isValidation:
			message = "Password must be at least 8 characters"
		}
		location := routes.PasswordNew.URL()
		if payload.Token != "" {
			location = routes.PasswordEdit.URL(payload.Token)
		}
		return h.authFlash(ctx, kiks.FlashError, message, location)
	}

	return h.authFlash(ctx, kiks.FlashSuccess, "Password reset successfully! Please log in.", routes.SessionNew.URL())
}
