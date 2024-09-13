package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/infra/tokenrepo"
	"net/http"
	"strings"
	"time"
)

type ZitadelClient struct {
	apiURL     string
	ClientHTTP HTTPClient
	userID     string
	privateKey []byte
	keyID      string
	clientID   string
	projectID  string
}

const (
	TwoDays = 48 * time.Hour
)

func NewZitadelClient(apiURL, userID, privateKey, keyID, projectID, clientID string) *ZitadelClient {
	return &ZitadelClient{
		apiURL:     apiURL,
		ClientHTTP: &ClientImpl{}, // &http.Client{Timeout: 10 * time.Second},
		userID:     userID,
		privateKey: []byte(privateKey),
		keyID:      keyID,
		projectID:  projectID,
		clientID:   clientID,
	}
}

func (z *ZitadelClient) SetHTTPClient(client HTTPClient) {
	z.ClientHTTP = client
}

func (z *ZitadelClient) GetServiceUserAccessToken(jwt string) (string, time.Duration, error) {
	data := fmt.Sprintf(`grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer&scope='openid profile urn:zitadel:iam:org:project:id:%s:aud'&assertion=%s`, z.projectID, jwt)
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
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", TwoDays, fmt.Errorf("ERROR | failed to get access token response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result tokenrepo.Token
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", TwoDays, fmt.Errorf("ERROR | cannot get decode token: %v", err)
	}

	return result.AccessToken, result.ExpiresIn, nil
}

func (z *ZitadelClient) ValidateUserToken(userToken, jwtToken string) (bool, error) {
	data := fmt.Sprintf("client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer&client_assertion=%s&token=%s", jwtToken, userToken)
	req, err := http.NewRequest("POST", z.apiURL+"/oauth/v2/introspect", strings.NewReader(data))
	if err != nil {
		return false, fmt.Errorf("error creando la solicitud: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := z.ClientHTTP.Do(req)
	if err != nil {
		return false, fmt.Errorf("ERROR | cannot send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("ERROR | response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.VerifyTokenUser
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("ERROR | cannot decode token: %v body msg: %s", err, string(bodyBytes))
	}

	return result.Active, nil
}
