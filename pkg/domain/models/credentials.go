package models

import (
	"encoding/json"
	"strings"
	"time"
)

// maybe repeated data
const (
	CredNameRequired          = "Credentials name is required"
	CredNameInvalid           = "Credentials name must be alphanumeric with max length of 255"
	CredNameExist             = "Credentials name already exists for this user"
	CredNameNotExist          = "Credentials name not exists for this user"
	CredNameCannotGenerate    = "error checking credential name existence"
	CredNameNotGenerate       = "cannot create new credential"
	CredDirectorySaveRequired = "Directory to save is required"
	CredDirectoryInvalid      = "Directory to save must be alphanumeric with max length of 255"
	CredDateInvalid           = "Invalid date"
	CredRateLimitUpdate       = 10 * time.Second
	UserTokenExpired          = "token expired"
)

type RequestExchangeCredential struct {
	UpdatedAt            *CustomTime    `json:"updated_at,omitempty"`
	CreatedAt            *CustomTime    `json:"created_at,omitempty"`
	RevokedAt            *CustomTime    `json:"revoked_at,omitempty"`
	LastUsedAt           *CustomTime    `json:"last_used_at,omitempty"`
	ExpiresAt            *CustomTime    `json:"expires_at,omitempty"`
	Sub                  string         `json:"sub,omitempty"`
	NodeID               string         `json:"nodeid,omitempty"`
	RequestID            string         `json:"requestid,omitempty"`
	WorkflowID           string         `json:"workflowid,omitempty"`
	ID                   string         `json:"id,omitempty" `
	Type                 string         `json:"type,omitempty"`
	Name                 string         `json:"name,omitempty"`
	Data                 DataCredential `json:"data,omitempty"`
	Version              uint32         `json:"version,omitempty"`
	CredentialCreatedNew bool           `json:"credentialnew"`
	IsActive             bool           `json:"is_active"`
}

type RequestCreateCredential struct {
	ID                   string         `json:"id,omitempty"`
	Sub                  string         `json:"sub,omitempty"`
	Name                 string         `json:"name,omitempty" `
	Type                 string         `json:"type,omitempty" `
	WorkflowID           string         `json:"workflowid,omitempty"`
	NodeID               string         `json:"nodeid,omitempty"`
	RequestID            string         `json:"requestid,omitempty"`
	Data                 DataCredential `json:"data" binding:"required"`
	CredentialCreatedNew bool           `json:"credentialnew"`
}

type DataCredential struct {
	ID           string   `json:"id,omitempty"`
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret" `
	RedirectURL  string   `json:"redirectURL" `
	OAuthURL     string   `json:"oauthurl,omitempty"`
	State        string   `json:"state,omitempty"`
	Code         string   `json:"code"`
	CodeVerifier string   `json:"codeverifier"`
	Token        string   `json:"token,omitempty"`
	TokenRefresh string   `json:"tokenrefresh,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
}

type ResponseCreateCredential struct {
	Data   string `json:"data"`
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type ResponseGetCredential struct {
	Credentials *[]RequestExchangeCredential `json:"credentials"`
	Error       string                       `json:"error"`
	Status      int                          `json:"status"`
}

type CredentialPayload struct {
	Data string `json:"data,omitempty"`
	RequestExchangeCredential
}

type InfoCredentials struct {
	Data                   *[]RequestExchangeCredential `json:"data,omitempty"`
	Rows                   *int64                       `json:"rows,omitempty"`
	RowsBeforeLimitAtLeast *int64                       `json:"rows_before_limit_at_least,omitempty"`
	Statistics             *Statistics                  `json:"statistics,omitempty"`
	Meta                   []Meta                       `json:"meta,omitempty"`
}

func (dc *RequestExchangeCredential) UnmarshalJSON(data []byte) error {
	type Alias RequestExchangeCredential
	aux := &struct {
		*Alias
		Data json.RawMessage `json:"data"`
	}{
		Alias: (*Alias)(dc),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	cleanedData := cleanEscapedJSON(string(aux.Data))
	var dataCredential DataCredential
	if err := json.Unmarshal([]byte(cleanedData), &dataCredential); err != nil {
		return err
	}

	dc.Data = dataCredential

	return nil
}

func cleanEscapedJSON(escapedJSON string) string {
	cleaned := strings.Trim(escapedJSON, "\"")
	cleaned = strings.ReplaceAll(cleaned, "\\\"", "\"")
	cleaned = strings.ReplaceAll(cleaned, "\\\\", "\\")
	return cleaned
}
