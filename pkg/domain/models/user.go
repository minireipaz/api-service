package models

import (
	"fmt"
	"time"
)

type Users struct {
	ID                       string     `json:"id,omitempty"`
	AccessToken              string     `json:"access_token,omitempty"`
	Sub                      string     `json:"sub"`
	Status                   UserStatus `json:"status,omitempty"`
	RoleID                   UserRoleID `json:"roleId,omitempty"`
	ResetPasswordToken       *string    `json:"resetPasswordToken,omitempty"`
	ResetPasswordTokenSentAt *time.Time `json:"resetPasswordTokenSentAt,omitempty"`
	InvitationToken          *string    `json:"invitationToken,omitempty"`
	InvitationTokenSentAt    *time.Time `json:"invitationTokenSentAt,omitempty"`
	TrialExpiryDate          *time.Time `json:"trialExpiryDate,omitempty"`
	DeletedAt                *time.Time `json:"deleted_at,omitempty"`
	CreatedAt                *time.Time `json:"created_at,omitempty"`
	UpdatedAt                *time.Time `json:"updated_at,omitempty"`
}

type UserStatus uint8
type UserRoleID uint8

const (
	StatusActive  UserStatus = 1
	StatusInvited UserStatus = 2
	StatusPending UserStatus = 3
	StatusBlocked UserStatus = 4

	RoleAdmin     UserRoleID = 1
	RoleUser      UserRoleID = 2
	RoleModerator UserRoleID = 3
	RoleEditor    UserRoleID = 4
	RoleGuest     UserRoleID = 5
	RoleCustomer  UserRoleID = 6
	RoleSupport   UserRoleID = 7
	RoleManager   UserRoleID = 8
	RoleAnalyst   UserRoleID = 9
	RoleDeveloper UserRoleID = 10

	UserInvalidJSON = "Invalid JSON data"

	UserSubExist         = "Sub already exists"
	UserNameExist        = "username already exists"
	UserCannotGenerate   = "error checking Sub existence"
	UserNameCannotCreate = "error checking username existence"
	UsertNameNotGenerate = "cannot create new user"

	UserSubRequired = "Sub user is required"
	UserSubMustBe   = "Sub user must greater than 3 characters"
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
