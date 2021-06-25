package gateway

import (
	"github.com/kubeedge/beehive/pkg/core"
	"github.com/kubeedge/kubeedge/cloud/pkg/common/modules"
	"github.com/kubeedge/kubeedge/cloud/pkg/gateway/config"
	configv1alpha1 "github.com/kubeedge/kubeedge/pkg/apis/componentconfig/cloudcore/v1alpha1"
)

type Gateway struct {
	enable             bool
	address            string
	port               uint32
}

func newGateway(c *configv1alpha1.Gateway) *Gateway {
	gw := Gateway{
		enable:                c.Enable,
		address:               c.Address,
		port:                  c.Port,
	}

	return &gw
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
	
	go startGatewayServer()	
}