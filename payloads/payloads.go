// Package payloads declares the Inertia page props and request bodies the
// frontend consumes. andurel sync payloads emits every exported struct in this
// package to resources/js/types/payloads.ts.
//
// Keep fields to standard-library types, time, and the stdlib uuid package
// the scaffold imports ("uuid"):
// payloads are data, not a place for application packages.
package payloads

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
