package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Register(app *gin.Engine) {
	app.Use(otelgin.Middleware("backend-vercel"))
}
