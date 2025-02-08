package tests

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"minireipaz/pkg/infra/httpclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock generado automáticamente por Mockery
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHTTPClient) DoRequest(method, url, authToken string, body interface{}) ([]byte, error) {
	args := m.Called(method, url, authToken, body)
	return args.Get(0).([]byte), args.Error(1)
}

const (
	HostZitadel = "https://dev-instance-p17owx.us1.zitadel.cloud"
	ProjectID   = "285990866588292160"
	FakeJWT     = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjI5MjAxMjAxMTc3NTA3NjM4NCIsInR5cCI6IkpXVCJ9..."
)

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
			expectedExpire: time.Hour * 24, // models.OneDay
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
			expectedExpire: time.Hour * 24, // models.OneDay
			expectedErr:    "request creation error",
		},
		{
			name: "error decoding JSON",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
			},
			expectedToken:  "",
			expectedExpire: time.Hour * 24, // models.OneDay
			expectedErr:    "ERROR | cannot get decode token: invalid character 'i' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockHTTPClient)
			mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
				if req.Method != http.MethodPost {
					return false
				}
				if req.URL.String() != fmt.Sprintf("%s/oauth/v2/token", HostZitadel) {
					return false
				}
				contentType := req.Header.Get("Content-Type")
				if contentType != "application/x-www-form-urlencoded" {
					return false
				}
				bodyBytes, err := io.ReadAll(req.Body)
				if err != nil {
					return false
				}
				body := string(bodyBytes)
				if !strings.Contains(body, "grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer") ||
					!strings.Contains(body, "scope='openid profile urn:zitadel:iam:org:project:id:") ||
					!strings.Contains(body, "assertion=") {
					return false
				}
				return true
			})).Return(tt.mockResponse, tt.mockError)

			client := httpclient.ZitadelClient{
				ClientHTTP: mockClient,
				ApiURL:     HostZitadel,
				ProjectID:  ProjectID,
			}

			jwt := FakeJWT
			token, expires, err := client.GenerateServiceUserAccessToken(jwt)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			if token != nil {
				assert.Equal(t, tt.expectedToken, *token)
			} else {
				assert.Equal(t, tt.expectedToken, "")
			}

			assert.Equal(t, tt.expectedExpire, expires)

			mockClient.AssertExpectations(t)
		})
	}
}
