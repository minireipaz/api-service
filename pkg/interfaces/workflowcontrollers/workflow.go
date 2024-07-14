package workflowcontrollers

import (
	"github.com/gin-gonic/gin"
	"log"
	"minireipaz/pkg/domain/models"
	"minireipaz/pkg/domain/workflowservices"
)

func CreateWorkflow(c *gin.Context) {
	newWorkflow := c.MustGet("workflow").(models.Workflow)
	createdWorkflow := services.CreateWorkflow(newWorkflow)
	log.Printf("%v", createdWorkflow)
}
