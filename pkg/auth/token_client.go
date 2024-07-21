package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TokenClient struct {
	apiURL string
}

func NewTokenClient(apiURL string) *TokenClient {
	return &TokenClient{
		apiURL: apiURL,
	}
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (t *TokenClient) GetAccessToken(jwt string) (TokenResponse, error) {
	data := fmt.Sprintf("grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer&scope=openid&assertion=%s", jwt)
	req, err := http.NewRequest("POST", t.apiURL, bytes.NewBufferString(data))
	if err != nil {
		return TokenResponse{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TokenResponse{}, fmt.Errorf("failed to get access token: %s", resp.Status)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return TokenResponse{}, err
	}

	return tokenResponse, nil
}
