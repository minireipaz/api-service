package controllers

import (
	"log"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/repos"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService repos.UserService
}

func NewUserController(newUserService repos.UserService) *UserController {
	return &UserController{userService: newUserService}
}

func (u *UserController) SyncUseWrithIDProvider(ctx *gin.Context) {
	currentUser := ctx.MustGet("user").(models.SyncUserRequest)
	created, exist := u.userService.SynUser(&currentUser)
	response := models.SyncUserResponse{
		Error:   "",
		Status:  http.StatusOK,
		Exist:   exist,
		Created: created,
	}
	if !created && !exist {
		log.Printf("WARN | User not created and does not exist: %s", currentUser.Sub)
		response.Error = models.UserNameCannotCreate
		response.Status = http.StatusInternalServerError
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (u *UserController) GetUserByStub(_ *gin.Context) {

}
