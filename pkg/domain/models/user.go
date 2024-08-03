package models

import (
	"fmt"
	"time"
)

type Users struct {
	ID                       string     `json:"id,omitempty"`
	AccessToken              string     `json:"access_token,omitempty"`
	Stub                     string     `json:"stub"`
	Status                   UserStatus `json:"status,omitempty"`
	ResetPasswordToken       *string    `json:"resetPasswordToken,omitempty"`
	ResetPasswordTokenSentAt *time.Time `json:"resetPasswordTokenSentAt,omitempty"`
	InvitationToken          *string    `json:"invitationToken,omitempty"`
	InvitationTokenSentAt    *time.Time `json:"invitationTokenSentAt,omitempty"`
	TrialExpiryDate          time.Time  `json:"trialExpiryDate,omitempty"`
	RoleID                   string     `json:"roleId,omitempty"`
	DeletedAt                *time.Time `json:"deleted_at,omitempty"`
	CreatedAt                time.Time  `json:"created_at,omitempty"`
	UpdatedAt                time.Time  `json:"updated_at,omitempty"`
}

type UserStatus uint8

const (
	StatusActive  UserStatus = 1
	StatusInvited UserStatus = 2
	StatusPending UserStatus = 3
	StatusBlocked UserStatus = 4

	UserStubExist        = "Stub already exists"
	UserNameExist        = "username already exists"
	UserCannotGenerate   = "error checking Stub existence"
	UserNameCannotCreate = "error checking username existence"
	UsertNameNotGenerate = "cannot create new user"
)

func (s UserStatus) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusInvited:
		return "invited"
	case StatusPending:
		return "pending"
	case StatusBlocked:
		return "blocked"
	default:
		return "unknown"
	}
}

func UserStatusFromUint8(v uint8) (UserStatus, error) {
	if v >= 1 && v <= 4 {
		return UserStatus(v), nil
	}
	return 0, fmt.Errorf("invalid user status value: %d", v)
}
