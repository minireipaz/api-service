package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/services"
	"minireipaz/pkg/interfaces/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var _ services.UserServiceInterface = (*MockUserService)(nil)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) SynUser(user *models.Users) (created, exist bool) {
	args := m.Called(user)
	return args.Bool(0), args.Bool(1)
}

func TestSyncUseWrithIDProvider(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(m *MockUserService)
		want      int
		user      models.Users
	}{
		{
			name: "User exists",
			mockSetup: func(m *MockUserService) {
				m.On("SynUser", mock.Anything).Return(false, true)
			},
			want: http.StatusOK,
			user: models.Users{Sub: "testUser"},
		},
		{
			name: "User created",
			mockSetup: func(m *MockUserService) {
				m.On("SynUser", mock.Anything).Return(true, false)
			},
			want: http.StatusOK,
			user: models.Users{Sub: "testUser"},
		},
		{
			name: "User creation failed",
			mockSetup: func(m *MockUserService) {
				m.On("SynUser", mock.Anything).Return(false, false)
			},
			want: http.StatusInternalServerError,
			user: models.Users{Sub: "testUser"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Config mock
			mockUserService := new(MockUserService)
			tt.mockSetup(mockUserService)
			controller := controllers.NewUserController(mockUserService)

			// Context Gin
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			// Add user inside context
			ctx.Set("user", tt.user)

			// RUN Function testing
			controller.SyncUseWrithIDProvider(ctx)

			// Want and get
			assert.Equal(t, tt.want, ctx.Writer.Status())

			// Verifyr both
			mockUserService.AssertExpectations(t)
		})
	}
}
