package dimodel

import (
	"minireipaz/pkg/domain/repos"
	"minireipaz/pkg/interfaces/controllers"
)

type Dependencies struct {
	WorkflowController   *controllers.WorkflowController
	AuthService          *repos.AuthService
	UserController       *controllers.UserController
	DashboardController  *controllers.DashboardController
	AuthController       *controllers.AuthController
	CredentialController *controllers.CredentialController
	ActionsController    *controllers.ActionsController
}
