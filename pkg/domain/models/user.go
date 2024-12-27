package models

import (
	"fmt"
)

type SyncUserRequest struct {
	TrialExpiryDate          *CustomTime `json:"trialExpiryDate,omitempty"`
	ResetPasswordToken       *string     `json:"resetPasswordToken,omitempty"`
	ResetPasswordTokenSentAt *CustomTime `json:"resetPasswordTokenSentAt,omitempty"`
	InvitationToken          *string     `json:"invitationToken,omitempty"`
	InvitationTokenSentAt    *CustomTime `json:"invitationTokenSentAt,omitempty"`
	DeletedAt                *CustomTime `json:"deleted_at,omitempty"`
	CreatedAt                *CustomTime `json:"created_at,omitempty"`
	UpdatedAt                *CustomTime `json:"updated_at,omitempty"`
	AccessToken              string      `json:"access_token" binding:"required,max=1000"`
	Sub                      string      `json:"user_id" binding:"required,max=50"`
	ID                       string      `json:"id,omitempty"`
	Status                   UserStatus  `json:"status" binding:"omitempty,oneof=1 2 3 4"`
	RoleID                   UserRoleID  `json:"roleId" binding:"omitempty,oneof=1 2 3 4 5 6 7 8 9 10"`
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

	InvalidJSON = "Invalid JSON data"

	UserSubExist           = "Sub already exists"
	UserSubIsRequired      = "Sub is required"
	UserNameExist          = "username already exists"
	UserCannotGenerate     = "error checking Sub existence"
	UserNameCannotCreate   = "error checking username existence"
	UsertNameNotGenerate   = "cannot create new user"
	UserSubInvalid         = "Sub must be a numeric string with max length of 50"
	UserAccessTokenInvalid = "Access token must be a valid JWT with max length of 1000"
	UserSubRequired        = "Sub user is required"
	UserSubMustBe          = "Sub user must greater than 3 characters"
	UserInvalidStatus      = "Invalid status"
	UserInvalidRole        = "Invalid role ID"
	AuthInvalid            = "Authorization header is required"
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

type SyncUserResponse struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Exist   bool   `json:"exist"`
	Created bool   `json:"created"`
}

type UnauthorizedError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type InvalidRequestError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type UnsupportedMediaTypeError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type TooManyRequestsError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}
