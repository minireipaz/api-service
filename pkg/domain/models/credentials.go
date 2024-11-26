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
	CredNameCannotGenerate    = "error checking Credential name existence"
	CredNameNotGenerate       = "cannot create new Credential"
	CredDirectorySaveRequired = "Directory to save is required"
	CredDirectoryInvalid      = "Directory to save must be alphanumeric with max length of 255"
	CredDateInvalid           = "Invalid date"
	CredRateLimitUpdate       = 10 * time.Second
	UserTokenExpired          = "token expired"
)

type RequestExchangeCredential struct {
	ID         string         `json:"id,omitempty" `
	Sub        string         `json:"sub,omitempty"`
	Name       string         `json:"name,omitempty"`
	Type       string         `json:"type,omitempty"`
	WorkflowID string         `json:"workflowid,omitempty"`
	NodeID     string         `json:"nodeid,omitempty"`
	CreatedAt  *CustomTime    `json:"created_at,omitempty"`
	UpdatedAt  *CustomTime    `json:"updated_at,omitempty"`
	LastUsedAt *CustomTime    `json:"last_used_at,omitempty"`
	RevokedAt  *CustomTime    `json:"revoked_at,omitempty"`
	Version    uint32         `json:"version,omitempty"`
	IsActive   bool           `json:"is_active,omitempty"`
	ExpiresAt  *CustomTime    `json:"expires_at,omitempty"`
	Data       DataCredential `json:"data,omitempty"`
}

type RequestCreateCredential struct {
	ID         string         `json:"id,omitempty"`
	Sub        string         `json:"sub,omitempty"`
	Name       string         `json:"name,omitempty" `
	Type       string         `json:"type,omitempty" `
	Timestamp  int64          `json:"timestamp,omitempty"`
	WorkflowID string         `json:"workflowid,omitempty"`
	NodeID     string         `json:"nodeid,omitempty"`
	Data       DataCredential `json:"data" binding:"required"`
}

type DataCredential struct {
	ID           string   `json:"id,omitempty"`
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret" `
	RedirectURL  string   `json:"redirectURL" `
	OAuthURL     string   `json:"oauthurl,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
	State        string   `json:"state,omitempty"`
	Code         string   `json:"code"`
	Token        string   `json:"token,omitempty"`
	TokenRefresh string   `json:"tokenrefresh,omitempty"`
}

type ResponseCreateCredential struct {
	Data   string `json:"data"`
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type ResponseGetCredential struct {
	Status      int                          `json:"status"`
	Error       string                       `json:"error"`
	Credentials *[]RequestExchangeCredential `json:"credentials"`
}

type CredentialPayload struct {
	RequestExchangeCredential
	Data string `json:"data,omitempty"`
}

type InfoCredentials struct {
	Meta                   []Meta                       `json:"meta,omitempty"`
	Data                   *[]RequestExchangeCredential `json:"data,omitempty"`
	Rows                   *int64                       `json:"rows,omitempty"`
	RowsBeforeLimitAtLeast *int64                       `json:"rows_before_limit_at_least,omitempty"`
	Statistics             *Statistics                  `json:"statistics,omitempty"`
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
