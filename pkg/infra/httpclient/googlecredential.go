package httpclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"minireipaz/pkg/domain/models"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type CredentialGoogleHTTPRepository struct {
	client HTTPClient
}

func NewGoogleCredentialRepository(httpCli HTTPClient) *CredentialGoogleHTTPRepository {
	return &CredentialGoogleHTTPRepository{
		client: httpCli,
	}
}

func (c *CredentialGoogleHTTPRepository) GenerateAuthURL(credential *models.RequestCreateCredential, credentialCreatedNew *bool) *string {
	codeVerifier := oauth2.GenerateVerifier()

	credential.Timestamp = time.Now().UTC().Unix()
	credential.Data.Scopes = []string{"https://www.googleapis.com/auth/spreadsheets.readonly"}
	credential.Data.Code = codeVerifier
	credential.Data.CodeVerifier = codeVerifier
	credential.Data.OAuthURL = google.Endpoint.AuthURL
	credential.CredentialCreatedNew = *credentialCreatedNew // in case to update workflow
	stateJSON, _ := json.Marshal(credential)
	stateToken := base64.URLEncoding.EncodeToString(stateJSON)

	var googleOauthConfig = oauth2.Config{
		RedirectURL:  credential.Data.RedirectURL,
		ClientID:     credential.Data.ClientID,
		ClientSecret: credential.Data.ClientSecret,
		Scopes:       credential.Data.Scopes,
		Endpoint:     google.Endpoint,
	}

	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
		oauth2.S256ChallengeOption(codeVerifier), // PKCE challenge
	}

	url := googleOauthConfig.AuthCodeURL(
		stateToken,
		opts...,
	)
	return &url
}

func (c *CredentialGoogleHTTPRepository) ExchangeGoogleCredential(currentCredential *models.RequestExchangeCredential) (accessToken, refreshToken *string, expire *time.Time, stateInfo *models.RequestExchangeCredential, err error) {
	stateJSON, err := base64.URLEncoding.DecodeString(currentCredential.Data.State)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("invalid state format: %v", err)
	}

	if err := json.Unmarshal(stateJSON, &stateInfo); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("invalid state data: %v", err)
	}

	if stateInfo.Data.Code == "" {
		return nil, nil, nil, nil, fmt.Errorf("ERROR | missing code verifier")
	}

	// replay attacks not implemented
	// if time.Now().UTC().Unix()-stateInfo.Timestamp > 3600 { // 1 hour
	// 	return nil, nil, nil, nil, fmt.Errorf("ERROR | state token expired")
	// }

	stateInfo.Data.Code = currentCredential.Data.Code
	var googleOauthConfig = oauth2.Config{
		RedirectURL:  stateInfo.Data.RedirectURL,
		ClientID:     stateInfo.Data.ClientID,
		ClientSecret: stateInfo.Data.ClientSecret,
		Scopes:       stateInfo.Data.Scopes,
		Endpoint:     google.Endpoint, // endpoint uri added to DATA?
	}

	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
		oauth2.VerifierOption(stateInfo.Data.CodeVerifier), // PKCE verification
	}
	// token expiry in 1hr
	token, err := googleOauthConfig.Exchange(context.Background(), currentCredential.Data.Code, opts...)
	if err != nil {
		// TODO: can be expired
		return nil, nil, nil, nil, fmt.Errorf("ERROR | cannot exchange code: %v", err)
	}
	return &token.AccessToken, &token.RefreshToken, &token.Expiry, stateInfo, nil
}
