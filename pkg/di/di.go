package di

import (
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/services"
	"minireipaz/pkg/infra/brokerclient"
	"minireipaz/pkg/infra/httpclient"
	"minireipaz/pkg/infra/redisclient"
	"minireipaz/pkg/interfaces/controllers"
)

func InitDependencies() (*controllers.WorkflowController, *services.AuthService, *controllers.UserController, *controllers.DashboardController, *controllers.AuthController) {
	configZitadel := config.NewZitaldelEnvConfig()
	kafkaConfig := config.NewKafkaEnvConfig()
	clickhouseConfig := config.NewClickhouseEnvConfig()

	// init autentication
	authContext := controllers.NewAuthContext(configZitadel)
	authService := authContext.GetAuthService()
	authController := authContext.GetAuthController()

	userRedisClient := redisclient.NewRedisClient()
	// userHttpClient := httpclient.NewUserClientHTTP()
	userHTTPClient := &httpclient.ClientImpl{}
	userBrokerClient := brokerclient.NewBrokerClient(kafkaConfig)

	repoUserRedis := redisclient.NewUserRedisRepository(userRedisClient)
	repoUserHTTP := httpclient.NewUserClientHTTP(userHTTPClient)
	repoUserBroker := brokerclient.NewUserKafkaRepository(userBrokerClient)

	userService := services.NewUserService(repoUserHTTP, repoUserRedis, repoUserBroker)
	userController := controllers.NewUserController(userService)

	workflowRedisClient := redisclient.NewRedisClient()
	repo := redisclient.NewWorkflowRepository(workflowRedisClient)
	idService := services.NewUUIDService()
	workflowService := services.NewWorkflowService(repo, idService)
	workflowController := controllers.NewWorkflowController(workflowService, authService)

	dashboardHTTPClient := &httpclient.ClientImpl{}
	dashboardRepo := httpclient.NewDashboardRepository(dashboardHTTPClient, clickhouseConfig)
	dashboardService := services.NewDashboardService(dashboardRepo)
	dashboardController := controllers.NewDashboardController(dashboardService, authService)

	return workflowController, authService, userController, dashboardController, authController
}
