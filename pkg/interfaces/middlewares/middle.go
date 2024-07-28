package middlewares

import (
	"minireipaz/pkg/domain/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Register(app *gin.Engine, authService *services.AuthService) {
	app.Use(otelgin.Middleware("backend-vercel"))
	app.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"http://localhost:3020", "http://localhost:3010"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	app.Use(AuthMiddleware(authService))
}
