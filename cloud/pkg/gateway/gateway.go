package gateway

import (
	"k8s.io/klog/v2"

	"github.com/kubeedge/beehive/pkg/core"
	"github.com/kubeedge/kubeedge/cloud/pkg/common/modules"
	"github.com/kubeedge/kubeedge/cloud/pkg/gateway/config"
	configv1alpha1 "github.com/kubeedge/kubeedge/pkg/apis/componentconfig/cloudcore/v1alpha1"
)

type Gateway struct {
	enable             bool
}

func newGateway(gw *configv1alpha1.Gateway) *Gateway {
	return &Gateway{
		enable:                gw.Enable,
	}
}

func Register(gw *configv1alpha1.Gateway) {
	config.InitConfigure(gw)
	core.Register(newGateway(gw))
}

func (gw *Gateway) Name() string {
	return modules.GatewayModuleName
}

func (gw *Gateway) Group() string {
	return modules.GatewayGroupName
}

func (gw *Gateway) Enable() bool {
	return gw.enable
}

func (gw *Gateway) Start() {
	klog.Infof("Gateway started.")
}

