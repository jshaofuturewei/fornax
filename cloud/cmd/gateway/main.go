package main

import (
	"fmt"
	"time"

	"k8s.io/klog/v2"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func init() {
	klog.InitFlags(nil)
}
func main() {
	

	version := pcap.Version()
	fmt.Println(version)


	handle, _ := pcap.OpenLive(
		"eth0",
		int32(65535),
		false,
		-1*time.Second,
	)

	handle.SetBPFFilter("tcp and port 6081")

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		fmt.Println(packet)
	}

	defer handle.Close()

	klog.Infof("Here we are!")
}
