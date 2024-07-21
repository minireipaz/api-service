package middlewares

import (
	"minireipaz/pkg/domain/services"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Register(app *gin.Engine, authService *services.AuthService) {
	app.Use(otelgin.Middleware("backend-vercel"))
  app.Use(AuthMiddleware(authService))
}
