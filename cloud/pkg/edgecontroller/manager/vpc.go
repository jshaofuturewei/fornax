package manager

import (
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/kubeedge/kubeedge/cloud/pkg/edgecontroller/config"
)

// VpcManager manage all events of rule by SharedInformer
type VpcManager struct {
	events chan watch.Event
}

// Events return the channel save events from watch secret change
func (rem *VpcManager) Events() chan watch.Event {
	return rem.events
}

// NewVpcManager create VpcManager by SharedIndexInformer
func NewVpcManager(si cache.SharedIndexInformer) (*VpcManager, error) {
	events := make(chan watch.Event, config.Config.Buffer.VpcsEvent)
	rh := NewCommonResourceEventHandler(events)
	si.AddEventHandler(rh)

	return &VpcManager{events: events}, nil
}
