package di

import (
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/services"
	"minireipaz/pkg/infra/brokerclient"
	"minireipaz/pkg/infra/httpclient"
	"minireipaz/pkg/infra/redisclient"
	"minireipaz/pkg/interfaces/controllers"
)

func InitDependencies() (*controllers.WorkflowController, *services.AuthService, *controllers.UserController) {
	configZitadel := config.NewZitaldelEnvConfig()
	kafkaConfig := config.NewKafkaEnvConfig()

	// init autentication
	authContext := controllers.NewAuthContext(configZitadel)
	authContext.GetAuthController()
	authService := authContext.GetAuthService()

	userRedisClient := redisclient.NewRedisClient()
	// userHttpClient := httpclient.NewUserClientHTTP()
	userHTTPClient := &httpclient.HttpClientImpl{}
	userBrokerClient := brokerclient.NewBrokerClient(kafkaConfig)

	repoUserRedis := redisclient.NewUserRedisRepository(userRedisClient)
	repoUserHTTP := httpclient.NewUserClientHTTP(userHTTPClient) /// .UserHTTPRepository(userHTTPClient)
  repoUserBroker := brokerclient.NewUserKafkaRepository(userBrokerClient)

	userService := services.NewUserService(repoUserHTTP, repoUserRedis, repoUserBroker)
	userController := controllers.NewUserController(userService)

	workflowRedisClient := redisclient.NewRedisClient()
	repo := redisclient.NewWorkflowRepository(workflowRedisClient)
	idService := services.NewUUIDService()
	workflowService := services.NewWorkflowService(repo, idService)
	workflowController := controllers.NewWorkflowController(workflowService, authService)

	return workflowController, authService, userController

}
