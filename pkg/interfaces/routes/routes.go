package routes

import (
	"minireipaz/pkg/common"
	"minireipaz/pkg/interfaces/controllers"
	"minireipaz/pkg/interfaces/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(app *gin.Engine, workflowController *controllers.WorkflowController) {
	app.NoRoute(ErrRouter)

	// Configuración de rutas
	api := app.Group("/api")
	{
		api.GET("/ping", common.Ping)

		workflows := api.Group("/workflows")
		{
			workflows.POST("", middlewares.ValidateWorkflow(), workflowController.CreateWorkflow)
			// workflows.GET("/:uuid", workflowController.GetWorkflow)
		}

	}
}

func ErrRouter(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Page not found",
	})
}
