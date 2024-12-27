package models

import "time"

const (
	TwoDays = 48 * time.Hour
	OneDay  = 24 * time.Hour
)

type VerifyTokenUser struct {
	Iat           *int64   `json:"iat,omitempty"`
	Active        *bool    `json:"active,omitempty"`
	AuthTime      *int64   `json:"auth_time,omitempty"`
	ClientID      *string  `json:"client_id,omitempty"`
	Exp           *int64   `json:"exp,omitempty"`
	EmailVerified *bool    `json:"email_verified,omitempty"`
	Iss           *string  `json:"iss,omitempty"`
	Jti           *string  `json:"jti,omitempty"`
	Nbf           *int64   `json:"nbf,omitempty"`
	Scope         *string  `json:"scope,omitempty"`
	Sub           *string  `json:"sub,omitempty"`
	TokenType     *string  `json:"token_type,omitempty"`
	Aud           []string `json:"aud,omitempty"`
}
