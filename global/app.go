package global

import (
	"k3s-client/config"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type Application struct {
	ConfigViper   *viper.Viper
	Config        config.Configuration
	Log           *zap.Logger
	DB            *gorm.DB
	Client        *kubernetes.Clientset
	DynamicClient *dynamic.DynamicClient
}

var App = new(Application)
