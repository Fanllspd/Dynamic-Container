package app

import (
	"k3s-client/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChallengeRequest struct {
	// Username    string `json:"username"`
	TemplateID uint `json:"template_id"`
}

func CreateContainerHandler(c *gin.Context) {
	var req ChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fmt.Println(req.TemplateID)

	if err, url := services.ContainerServices.CreateContainer(req.TemplateID, c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"url": url})
	}
}

func GetContainersHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err, response := services.ContainerServices.GetUserContainersInfo(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}
