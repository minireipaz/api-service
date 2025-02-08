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
	ApiURL     string
	ClientHTTP HTTPClient
	UserID     string
	PrivateKey []byte
	KeyID      string
	ClientID   string
	ProjectID  string
}

func NewZitadelClient(apiURL, userID, privateKey, keyID, projectID, clientID string) *ZitadelClient {
	return &ZitadelClient{
		ApiURL:     apiURL,
		ClientHTTP: NewClientImpl(models.TimeoutRequest), // &http.Client{Timeout: 10 * time.Second},
		UserID:     userID,
		PrivateKey: []byte(privateKey),
		KeyID:      keyID,
		ProjectID:  projectID,
		ClientID:   clientID,
	}
}

func (z *ZitadelClient) SetHTTPClient(client HTTPClient) {
	z.ClientHTTP = client
}

func (z *ZitadelClient) GenerateServiceUserAccessToken(jwt string) (*string, time.Duration, error) {
	if jwt == "" {
		return nil, models.OneDay, fmt.Errorf("ERROR | token empty")
	}
	data := fmt.Sprintf(`grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer&scope='openid profile urn:zitadel:iam:org:project:id:%s:aud'&assertion=%s`, z.ProjectID, jwt)
	req, err := http.NewRequest("POST", z.ApiURL+"/oauth/v2/token", bytes.NewBufferString(data))
	if err != nil {
		return nil, models.OneDay, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := z.ClientHTTP.Do(req)
	if err != nil {
		return nil, models.OneDay, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, models.OneDay, fmt.Errorf("ERROR | failed to get access token response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result tokenrepo.Token
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, models.OneDay, fmt.Errorf("ERROR | cannot get decode token: %v", err)
	}

	return result.AccessToken, result.ExpiresIn, nil
}

func (z *ZitadelClient) ValidateUserToken(userToken, jwtToken string) (bool, int64, error) {
	data := fmt.Sprintf("client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer&client_assertion=%s&token=%s", jwtToken, userToken)
	req, err := http.NewRequest("POST", z.ApiURL+"/oauth/v2/introspect", strings.NewReader(data))
	if err != nil {
		return false, 0, fmt.Errorf("ERROR | cannot generate sol: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := z.ClientHTTP.Do(req)
	if err != nil {
		return false, 0, fmt.Errorf("ERROR | cannot send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, 0, fmt.Errorf("ERROR | response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.VerifyTokenUser
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, 0, fmt.Errorf("ERROR | cannot decode token: %v body msg: %s", err, string(bodyBytes))
	}
	if result.Active == nil || result.Exp == nil {
		// nil error because its expired not fault
		return false, 0, nil
	}
	return *result.Active, *result.Exp, nil
}

func (z *ZitadelClient) ValidateServiceUserAccessToken(serviceUserToken, jwtToken *string) (bool, error) {
	data := fmt.Sprintf("client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer&client_assertion=%s&token=%s", *jwtToken, *serviceUserToken)
	req, err := http.NewRequest("POST", z.ApiURL+"/oauth/v2/introspect", strings.NewReader(data))
	if err != nil {
		return false, fmt.Errorf("ERROR | cannot create HTTP newRequest: %v", err)
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

	return *result.Active, nil
}
