// Package payloads declares the Inertia page props and request bodies the
// frontend consumes. andurel sync payloads emits every exported struct in this
// package to resources/js/types/payloads.ts.
//
// Keep fields to standard-library types, time, and the stdlib uuid package
// the scaffold imports ("uuid"):
// payloads are data, not a place for application packages.
package payloads

import "time"

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PasswordResetRequestPayload struct {
	Email string `json:"email"`
}

type PasswordResetPayload struct {
	Token           string `json:"resetPasswordToken"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type PasswordResetPageProps struct {
	Token string `json:"token"`
}

type ContentTokenPageProps struct {
	ExpiresAt *time.Time `json:"expiresAt"`
	CSRFToken string     `json:"csrfToken"`
}

type IssuedContentToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type SharedAuthProps struct {
	Auth ViewerAuth `json:"auth"`
}

type ViewerAuth struct {
	IsAuthenticated bool `json:"isAuthenticated"`
}

type ErrorMessageProps struct {
	Message string `json:"message"`
}
