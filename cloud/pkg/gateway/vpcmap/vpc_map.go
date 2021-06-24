package vpcmap

import (
	"fmt"
	"net"
)

var VPC_MAP map[string][]SubnetInfo

type SubnetInfo struct {
	IpNet *net.IPNet
	GatewayIP string
}

func GetGatewayIP(vpcId string, ip string) (string, error) {
	subnetInfoList, ok := VPC_MAP[vpcId]
	if !ok {
		return "", fmt.Errorf("vpc %v not found.", vpcId)
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
