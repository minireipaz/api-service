package tests

import (
	"errors"
	"minireipaz/mocks"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/interfaces/controllers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestActionsController_CreateActionsNotion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	action798 := "action_789"

	tests := []struct {
		name               string
		setupMocks         func(*mocks.ActionsService, *mocks.AuthService)
		setupContext       func(*gin.Context) // Nuevo campo para configurar el contexto
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Success - Notion Action Created",
			setupMocks: func(actionsService *mocks.ActionsService, authService *mocks.AuthService) {
				token := "valid_oauth_token"
				authService.On("GetActionUserAccessToken").Return(&token, nil)

				actionsService.On(
					"CreateActionsNotion",
					mock.MatchedBy(func(action models.RequestGoogleAction) bool {
						return action.Type == "notionoauth" &&
							action.WorkflowID == "wf_123" &&
							action.CredentialID == "cred_456"
					}),
					mock.MatchedBy(func(token *string) bool {
						return *token == "valid_oauth_token"
					}),
				).Return(true, true, &action798)
			},
			setupContext: func(ctx *gin.Context) {
				ctx.Set(models.ActionNotionKey, models.RequestGoogleAction{
					Type:         "notionoauth",
					WorkflowID:   "wf_123",
					CredentialID: "cred_456",
					Sub:          "user_901",
					Pollmode:     models.NopollNode,
					Testmode:     true,
				})
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"error":"","data":"action_789","status":200}`,
		},
		{
			name: "Error - Invalid Token",
			setupMocks: func(actionsService *mocks.ActionsService, authService *mocks.AuthService) {
				authService.On("GetActionUserAccessToken").Return(nil, errors.New("invalid token"))
			},
			setupContext: func(ctx *gin.Context) {
				ctx.Set(models.ActionNotionKey, models.RequestGoogleAction{
					Type: models.NotionToken,
				})
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Failed to authenticate: invalid token", "status":400}`,
		},
		{
			name: "Error - Invalid Action Type",
			setupMocks: func(actionsService *mocks.ActionsService, authService *mocks.AuthService) {
				authService.On("GetActionUserAccessToken").Return(nil, errors.New("invalid token"))
			},
			setupContext: func(ctx *gin.Context) {
				ctx.Set(models.ActionNotionKey, models.RequestGoogleAction{
					Type: "invalid_type",
				})
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Failed to authenticate: invalid token", "status":400}`,
		},
		{
			name: "Error - Service Internal Failure",
			setupMocks: func(actionsService *mocks.ActionsService, authService *mocks.AuthService) {
				token := "valid_token"
				authService.On("GetActionUserAccessToken").Return(&token, nil)

				actionsService.On(
					"CreateActionsNotion",
					mock.AnythingOfType("models.RequestGoogleAction"),
					mock.AnythingOfType("*string"),
				).Return(false, false, nil)
			},
			setupContext: func(ctx *gin.Context) {
				ctx.Set(models.ActionNotionKey, models.RequestGoogleAction{
					Type:         "notionoauth",
					WorkflowID:   "wf_123",
					CredentialID: "cred_456",
				})
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"cannot create new workflow", "status":500}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockActionsService := mocks.NewActionsService(t)
			mockAuthService := mocks.NewAuthService(t)
			tt.setupMocks(mockActionsService, mockAuthService)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			tt.setupContext(ctx) // Aplica la configuración específica del contexto

			controller := controllers.NewActionsController(mockActionsService, mockAuthService)

			controller.CreateActionsNotion(ctx)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.JSONEq(t, tt.expectedResponse, w.Body.String())

			mockActionsService.AssertExpectations(t)
			mockAuthService.AssertExpectations(t)
			// mockActionsService := mocks.NewActionsService(t)
			// mockAuthService := mocks.NewAuthService(t)
			// mockActionsRedis := mocks.NewActionsRedisRepoInterface(t)
			// mockActionsBroker := mocks.NewActionsBrokerRepository(t)
			// mockActionsRepo := mocks.NewActionsHTTPRepository(t)

			// mockActionsRepo.On("SendAction", mock.Anything, mock.Anything).Return(true)

			// actionsService := services.NewActionsService(
			// 	mockActionsRedis,
			// 	mockActionsBroker,
			// 	mockActionsRepo,
			// )

			// // Inyecta el servicio en el controlador (si tu controlador lo requiere)
			// controller := controllers.NewActionsController(actionsService, mockAuthService)

			// // Configuración del contexto Gin
			// w := httptest.NewRecorder()
			// ctx, _ := gin.CreateTestContext(w)
			// ctx.Set(models.ActionNotionKey, models.RequestGoogleAction{})

			// // Ejecuta el método del controlador
			// controller.CreateActionsNotion(ctx)

			// // Verificaciones
			// assert.Equal(t, tt.expectedStatusCode, w.Code)
			// assert.JSONEq(t, tt.expectedResponse, w.Body.String())

			// // Asegura que los mocks fueron llamados correctamente
			// mockActionsService.AssertExpectations(t)
			// mockAuthService.AssertExpectations(t)
		})
	}
}
