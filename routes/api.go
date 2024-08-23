package routes

import (
	"k3s-client/app/controllers/app"
	"k3s-client/app/middleware"
	"k3s-client/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetApiGroupRoutes(router *gin.RouterGroup) {

	// PING TEST
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// TODO
	authRouter := router.Group("").Use(middleware.JWTAuth(services.AppGuardName))
	{
		// 创建容器
		authRouter.POST("/containers", app.CreateContainerHandler)
		// 获取指定容器的详细信息
		authRouter.GET("/containers/:id", app.GetContainersHandler)
	}
	// router.GET("/containers/:id", func(c *gin.Context) {
	// 	id := c.Param("id")
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code": 200,
	// 		"msg":  "success",
	// 		"data": id,
	// 	})
	// })

	// TODO
	// 获取容器列表
	router.GET("/containers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": []string{"container1", "container2"},
		})
	})

	// TODO
	// 删除指定容器
	router.DELETE("/containers/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": id,
		})
	})
	// TODO
	// 延时容器
	router.PATCH("/containers/:id/extend", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": id,
		})
	})
}
