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
	ClientHTTP HTTPClient
	userID     string
	privateKey []byte
	keyID      string
}

const (
	TwoDays = 48 * time.Hour
)

func NewZitadelClient(apiURL, userID, privateKey, keyID string) *ZitadelClient {
	return &ZitadelClient{
		apiURL:     apiURL,
		ClientHTTP: &Impl{}, // &http.Client{Timeout: 10 * time.Second},
		userID:     userID,
		privateKey: []byte(privateKey),
		keyID:      keyID,
	}
}

func (z *ZitadelClient) SetHTTPClient(client HTTPClient) {
	z.ClientHTTP = client
}

func (z *ZitadelClient) GetAccessToken(jwt string) (string, time.Duration, error) {
	data := fmt.Sprintf("grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer&scope=openid&assertion=%s", jwt)
	req, err := http.NewRequest("POST", z.apiURL+"/oauth/v2/token", bytes.NewBufferString(data))
	if err != nil {
		return "", TwoDays, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := z.ClientHTTP.Do(req)
	if err != nil {
		return "", TwoDays, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", TwoDays, fmt.Errorf("ERROR | failed to get access token: %d", resp.StatusCode)
	}

	var result tokenrepo.Token
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", TwoDays, fmt.Errorf("ERROR | cannot get decode token: %v", err)
	}

	return result.AccessToken, result.ExpiresIn, nil
}
