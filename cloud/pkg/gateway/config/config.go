package config

import (
	"sync"

	configv1alpha1 "github.com/kubeedge/kubeedge/pkg/apis/componentconfig/cloudcore/v1alpha1"
)

var Config Configure
var once sync.Once

type Configure struct {
	Gateway *configv1alpha1.Gateway
}

func InitConfigure(gw *configv1alpha1.Gateway) {
	once.Do(func() {
		Config = Configure{
			Gateway: gw,
		}
	})
}
