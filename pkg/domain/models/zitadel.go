package models

import "time"

const (
	TwoDays = 48 * time.Hour
)

type VerifyTokenUser struct {
	Active        bool   `json:"active"`
	TokenType     string `json:"token_type"`
	Exp           int    `json:"exp"`
	EmailVerified bool   `json:"email_verified"`
}
