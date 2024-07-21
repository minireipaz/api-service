package di

import (
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/services"
	"minireipaz/pkg/infra/redisclient"
	"minireipaz/pkg/interfaces/controllers"
)

func InitDependencies() ( *controllers.WorkflowController,  *services.AuthService) {
  // init autentication
	config := config.NewZitaldelEnvConfig()
	authContext := controllers.NewAuthContext(config)
	authContext.GetAuthController()
	authService := authContext.GetAuthService()

	workflowRedisClient := redisclient.NewRedisClient()
	repo := redisclient.NewWorkflowRepository(workflowRedisClient)
	idService := services.NewUUIDService()
	workflowService := services.NewWorkflowService(repo, idService)
	workflowController := controllers.NewWorkflowController(workflowService, authService)

  return workflowController, authService

}
