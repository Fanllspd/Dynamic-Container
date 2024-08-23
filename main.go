package main

import (
	"k3s-client/app/services"
	"k3s-client/global"

	"k3s-client/bootstrap"

	"go.uber.org/zap"
)

func main() {

	bootstrap.InitializeConfig()

	global.App.Log = bootstrap.InitializeLog()
	// global.App.Log.Info("log init success!")

	global.App.DB = bootstrap.InitializeDB()

	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			db.Close()
		}
	}()

	global.App.Client, global.App.DynamicClient = bootstrap.InitKubernetesClient()
	TokenOutPut, err, _ := services.JwtService.CreateToken(services.AppGuardName, 233 /*, "admin"*/)
	if err != nil {
		global.App.Log.Error("token create failed", zap.Any("err", err))
	}
	global.App.Log.Debug(TokenOutPut.AccessToken)
	// global.App.Log.Info(token)
	bootstrap.RunServer()
}
