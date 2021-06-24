package config

import (
	"net"
)

var VPC_MAP map[string][]SubnetInfo

type SubnetInfo {
	ipnet *net.IPNet
	gatewayIP string
}

func GetGatewayIP(vpcId string, ip string) (string, error) {
	subnetInfoList, ok := VPC_MAP[vpcId]
	if !ok {
		return "", fmt.Errorf("vpc %v not found.", vpcId)
	}

	for _, subnetInfo := range subnetInfoList {
		if subnetInfo.ipnet.Contains(ip) {
			return subnetInfo.gatewayIP
		}
	}

	return "", fmt.Errorf("no match subnet for %v/%v", vpcId, ip)
}
