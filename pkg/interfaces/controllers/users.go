package controllers

import (
	"log"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserServiceInterface
}

// type UserController struct {
// 	userService *services.UserService
// }

func NewUserController(newUserService services.UserServiceInterface) *UserController {
	return &UserController{userService: newUserService}
}

func (u *UserController) SyncUseWrithIDProvider(ctx *gin.Context) {
	currentUser := ctx.MustGet("user").(models.Users)
	created, exist := u.userService.SynUser(&currentUser)
	response := gin.H{
		"error":   "",
		"status":  http.StatusOK,
		"exist":   exist,
		"created": created,
	}
	if !created && !exist {
		log.Printf("WARN | User not created and does not exist: %s", currentUser.Sub)
		response["error"] = models.UserNameCannotCreate
		response["status"] = http.StatusInternalServerError
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (u *UserController) GetUserByStub(_ *gin.Context) {

}
