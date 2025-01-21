package tests

import (
	"minireipaz/pkg/domain/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCredentialRepo struct {
	mock.Mock
}

type MockGoogleOAuthRepo struct {
	mock.Mock
}

func (m *MockGoogleOAuthRepo) GenerateAuthURL(credential *models.RequestExchangeCredential, isNew *bool) *string {
	args := m.Called(credential, isNew)
	if str, ok := args.Get(0).(*string); ok {
		return str
	}
	return nil
}

func (m *MockCredentialRepo) GetCredentialByID(sub *string, id *string) *models.ResponseGetCredential {
	args := m.Called(sub, id)
	if response, ok := args.Get(0).(*models.ResponseGetCredential); ok {
		return response
	}
	return nil
}

func TestCreateCredential(t *testing.T) {
	mockCredRepo := new(MockCredentialRepo)
	mockGoogleOAuthRepo := new(MockGoogleOAuthRepo)

	tests := []struct {
		name           string
		input          *models.RequestCreateCredential
		mockSetup      func()
		expectedError  bool
		expectedOutput *models.RequestExchangeCredential
	}{
		{
			name: "Successfully create Google Sheets credential",
			input: &models.RequestCreateCredential{
				Sub:  "test-sub",
				ID:   "test-id",
				Type: "googlesheets",
			},
			mockSetup: func() {
				mockCredRepo.On("GetCredentialByID", mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
					Return(&models.ResponseGetCredential{
						Status: 200,
						Credentials: &[]models.RequestExchangeCredential{
							{
								Type: "googlesheets",
							},
						},
					})

				redirectURL := "https://oauth.google.com/auth"
				mockGoogleOAuthRepo.On("GenerateAuthURL", mock.AnythingOfType("*models.RequestExchangeCredential"), mock.AnythingOfType("*bool")).
					Return(&redirectURL)
			},
			expectedError: false,
			expectedOutput: &models.RequestExchangeCredential{
				Type: "googlesheets",
				Data: models.DataCredential{
					RedirectURL: "https://oauth.google.com/auth",
				},
			},
		},
		{
			name: "Fail when credential not found in DB",
			input: &models.RequestCreateCredential{
				Sub: "test-sub",
				ID:  "non-existent",
			},
			mockSetup: func() {
				mockCredRepo.On("GetCredentialByID", mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
					Return(&models.ResponseGetCredential{
						Status: 404,
					})
			},
			expectedError:  true,
			expectedOutput: nil,
		},
		{
			name: "Fail with unsupported credential type",
			input: &models.RequestCreateCredential{
				Sub:  "test-sub",
				ID:   "test-id",
				Type: "unsupported",
			},
			mockSetup: func() {
				mockCredRepo.On("GetCredentialByID", mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
					Return(&models.ResponseGetCredential{
						Status: 200,
						Credentials: &[]models.RequestExchangeCredential{
							{
								Type: "unsupported",
							},
						},
					})
			},
			expectedError:  true,
			expectedOutput: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.mockSetup()

			// Create service with mocks
			service := &CredentialServiceImpl{
				credentialRepo:  mockCredRepo,
				googleOAuthRepo: mockGoogleOAuthRepo,
			}
			NewCredentialService

			// Execute test
			result, err := service.CreateCredential(tt.input)

			// Verify results
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedOutput.Type, result.Type)
				assert.Equal(t, tt.expectedOutput.Data.RedirectURL, result.Data.RedirectURL)
			}
		})
	}

	// Verify that all expected mock calls were made
	mockCredRepo.AssertExpectations(t)
	mockGoogleOAuthRepo.AssertExpectations(t)
}
