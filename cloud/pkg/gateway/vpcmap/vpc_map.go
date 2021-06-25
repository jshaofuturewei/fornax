package vpcmap

import (
	"fmt"
	"net"
	"os"

	"github.com/kubeedge/kubeedge/pkg/util"
)

var VPC_MAP map[string][]SubnetInfo

var LOCAL_IP string

func init() {
	VPC_MAP = map[string][]SubnetInfo{}

	hostnameOverride, _ := os.Hostname()
	LOCAL_IP, _ = util.GetLocalIP(hostnameOverride)
}

type SubnetInfo struct {
	IpNet *net.IPNet
	GatewayIP string
}

func GetGatewayIP(vpcId string, ip string) (string, error) {
	if len(vpcId) == 0 {
		return "", fmt.Errorf("vpcid is empty.")
	}

	subnetInfoList, ok := VPC_MAP[vpcId]
	if !ok {
		return "", fmt.Errorf("vpc %v not found. %v", vpcId, VPC_MAP)
	}

	IP := net.ParseIP(ip)
	if IP == nil {
		return "", fmt.Errorf("invalid IP address %v", ip)
	}
	for _, subnetInfo := range subnetInfoList {
		if subnetInfo.IpNet.Contains(IP) {
			return subnetInfo.GatewayIP, nil
		}
	}

	return "", fmt.Errorf("no match subnet for %v/%v", vpcId, ip)
}

func IsLocalGateway(ipAddr string) bool {
	return ipAddr == LOCAL_IP || ipAddr == "127.0.0.1" || ipAddr == "0.0.0.0"
}