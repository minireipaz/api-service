package controllers

import (
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
	if !created && !exist {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   models.UserNameCannotCreate,
			"status":  http.StatusInternalServerError,
			"exist":   false,
			"created": false,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   "",
		"status":  http.StatusOK,
		"exist":   true,
		"created": false,
	})
}

func (u *UserController) GetUserByStub(_ *gin.Context) {

}
