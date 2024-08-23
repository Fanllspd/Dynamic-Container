package bootstrap

import (
	"context"
	"k3s-client/global"
	"k3s-client/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	apiRouter := r.Group("/api")
	routes.SetApiGroupRoutes(apiRouter)

	return r
}

func RunServer() {
	r := setupRouter()

	srv := &http.Server{
		Addr:    "127.0.0.1:" + global.App.Config.App.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.App.Log.Error("listen: ", zap.Any("err", err))
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.App.Log.Error("Server Shutdown:", zap.Any("err", err))
	}
	log.Println("Server Exiting")

}
