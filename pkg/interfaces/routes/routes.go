package routes

import (
	"minireipaz/pkg/common"
	"minireipaz/pkg/interfaces/controllers"
	"minireipaz/pkg/interfaces/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(app *gin.Engine, workflowController *controllers.WorkflowController, userController *controllers.UserController, dashboardController *controllers.DashboardController, authController *controllers.AuthController) {
	app.NoRoute(ErrRouter)

	// Routes in groups
	api := app.Group("/api")
	{
		api.GET("/ping", common.Ping)

		workflows := api.Group("/workflows")
		{
			workflows.POST("", middlewares.ValidateOnCreateWorkflow(), workflowController.CreateWorkflow)
			workflows.PUT(":id", middlewares.ValidateOnUpdateWorkflow(), workflowController.UpdateWorkflow)
			workflows.GET("/:iduser/:idworkflow", middlewares.ValidateOnGetWorkflow(), workflowController.GetWorkflow)
			workflows.GET("/:iduser", middlewares.ValidateOnGetWorkflow(), workflowController.GetAllWorkflows)
		}

		users := api.Group("/users")
		{
			users.POST("", middlewares.ValidateUser(), userController.SyncUseWrithIDProvider)
			users.GET("/:stub", userController.GetUserByStub)
		}

		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/:iduser", middlewares.ValidateUserAuth(), dashboardController.GetUserDashboardByID)
		}

		auth := api.Group("/auth")
		{
			auth.GET("/verify/:id", authController.VerifyUserToken)
		}
	}
}

func ErrRouter(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Page not found",
	})
}
