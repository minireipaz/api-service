package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minireipaz/pkg/infra/tokenrepo"
	"net/http"
	"time"
)

type ZitadelClient struct {
	apiURL     string
	client     *http.Client
	userID     string
	privateKey []byte
	keyID      string
}

const (
	twoDays = 48 * time.Hour
)

func NewZitadelClient(apiURL, userID, privateKey, keyID string) *ZitadelClient {
	return &ZitadelClient{
		apiURL:     apiURL,
		client:     &http.Client{Timeout: 10 * time.Second},
		userID:     userID,
		privateKey: []byte(privateKey),
		keyID:      keyID,
	}
}

func (z *ZitadelClient) GetAccessToken(jwt string) (string, time.Duration, error) {
	data := fmt.Sprintf("grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer&scope=openid&assertion=%s", jwt)
	req, err := http.NewRequest("POST", z.apiURL+"/oauth/v2/token", bytes.NewBufferString(data))
	if err != nil {
		return "", twoDays, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := z.client.Do(req)
	if err != nil {
		return "", twoDays, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", twoDays, fmt.Errorf("ERROR | failed to get access token: %s", resp.Status)
	}

	var result tokenrepo.Token
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", twoDays, fmt.Errorf("ERROR | cannot get decode token: %v", err)
	}

	return result.AccessToken, result.ExpiresIn, nil
}
