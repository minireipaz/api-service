package routes

import (
	"minireipaz/pkg/common"
	"minireipaz/pkg/interfaces/controllers"
	"minireipaz/pkg/interfaces/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(app *gin.Engine, workflowController *controllers.WorkflowController, userController *controllers.UserController, dashboardController *controllers.DashboardController, authController *controllers.AuthController, credentialController *controllers.CredentialController) {
	app.NoRoute(ErrRouter)

	// Routes in groups
	api := app.Group("/api")
	{
		api.GET("/ping", common.Ping)

		workflows := api.Group("/workflows")
		{
			workflows.GET("/:iduser/workflow/:idworkflow/:usertoken", authController.VerifyUserTokenForMiddleware, middlewares.ValidateOnGetWorkflow(), workflowController.GetWorkflow)
			workflows.GET("/:iduser/:usertoken", authController.VerifyUserTokenForMiddleware, middlewares.ValidateOnGetWorkflow(), workflowController.GetAllWorkflows)
			workflows.POST("", middlewares.ValidateOnCreateWorkflow(), workflowController.CreateWorkflow)
			workflows.PUT("/:id", middlewares.ValidateOnUpdateWorkflow(), workflowController.UpdateWorkflow)
		}

		users := api.Group("/users")
		{
			users.POST("", middlewares.ValidateUser(), userController.SyncUseWrithIDProvider)
			users.GET("/:stub", userController.GetUserByStub)
		}

		credentials := api.Group("/credentials")
		{
			// credentials.POST("", middlewares.ValidateUser(), userController.SyncUseWrithIDProvider)
			credentials.GET("/:iduser/:usertoken", authController.VerifyUserTokenForMiddleware, credentialController.GetAllCredentials)
		}

		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/:iduser", middlewares.ValidateUserAuth(), dashboardController.GetUserDashboardByID)
		}

		auth := api.Group("/auth")
		{
			auth.GET("/verify/:usertoken", authController.VerifyUserToken)
		}

		credentialsGoogle := api.Group("/google")
		{
			credentialsGoogle.POST("/credential", middlewares.ValidateOnCreateCredential(), credentialController.CreateCredential)
			credentialsGoogle.POST("/exchange", middlewares.ValidateOnExchangeCredential(), credentialController.ExchangeGoogleCode)
		}
	}
}

func ErrRouter(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "Page not found",
	})
}
