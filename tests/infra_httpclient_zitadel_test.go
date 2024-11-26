package tests

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"minireipaz/pkg/auth"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/infra/httpclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetAccessToken(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   *http.Response
		mockError      error
		expectedToken  string
		expectedExpire time.Duration
		expectedErr    string
	}{
		{
			name: "HTTP error response",
			mockResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			},
			expectedToken:  "",
			expectedExpire: models.OneDay,
			expectedErr:    "ERROR | failed to get access token response: 500, body: Internal Server Error",
		},
		{
			name: "successful request",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"access_token":"valid-token",
					"expires_in":3600000000000
				}`)),
			},
			expectedToken:  "valid-token",
			expectedExpire: 3600 * time.Second,
			expectedErr:    "",
		},
		{
			name:           "error creating request",
			mockError:      errors.New("request creation error"),
			expectedToken:  "",
			expectedExpire: models.OneDay,
			expectedErr:    "request creation error",
		},

		{
			name: "error decoding JSON",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
			},
			expectedToken:  "",
			expectedExpire: models.OneDay,
			expectedErr:    "ERROR | cannot get decode token: invalid character 'i' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockHTTPClient)
			mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(tt.mockResponse, tt.mockError)

			configZitadel := config.NewZitaldelEnvConfig()
			client := httpclient.NewZitadelClient(
				configZitadel.GetZitadelURI(),
				configZitadel.GetZitadelServiceUserID(),
				configZitadel.GetZitadelServiceUserKeyPrivate(),
				configZitadel.GetZitadelServiceUserKeyID(),
				configZitadel.GetZitadelProjectID(),
				configZitadel.GetZitadelKeyClientID(),
			)

			jwtGenerator := auth.NewJWTGenerator(auth.JWTGeneratorConfig{
				ServiceUser: auth.ServiceUserConfig{
					UserID:     configZitadel.GetZitadelServiceUserID(),
					PrivateKey: []byte(configZitadel.GetZitadelServiceUserKeyPrivate()),
					KeyID:      configZitadel.GetZitadelServiceUserKeyID(),
				},
				BackendApp: auth.BackendAppConfig{
					AppID:      configZitadel.GetZitadelBackendID(),
					PrivateKey: []byte(configZitadel.GetZitadelBackendKeyPrivate()),
					KeyID:      configZitadel.GetZitadelBackendKeyID(),
				},
				APIURL:    configZitadel.GetZitadelURI(),
				ProjectID: configZitadel.GetZitadelProjectID(),
				ClientID:  configZitadel.GetZitadelKeyClientID(),
			})

			client.SetHTTPClient(mockClient)

			jwt, err := jwtGenerator.GenerateAppInstrospectJWT(time.Hour)
			if err != nil {
				t.Errorf("%v", err)
			}

			token, expires, err := client.GenerateServiceUserAccessToken(jwt)
			var tokenStr = ""
			if token != nil && *token != "" {
				tokenStr = *token
			}
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Equal(t, tt.expectedToken, tokenStr)
				assert.Equal(t, tt.expectedExpire, expires)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, tokenStr)
				assert.Equal(t, tt.expectedExpire, expires)
			}
		})
	}
}
