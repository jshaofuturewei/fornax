package gateway

import (
	"bytes"
	"fmt"
	"encoding/json"
	"net/http"

	"k8s.io/klog/v2"

	"github.com/kubeedge/kubeedge/cloud/pkg/gateway/config"
)

func startGatewayServer() {
	klog.Infof("Starting Gateway http server")

	http.HandleFunc("/intra-cluster", ServeHTTPIntraCluster)

	http.HandleFunc("/inter-cluster", ServeHTTPInterCluster)

	http.ListenAndServe(fmt.Sprintf(":%v",config.Config.Gateway.Port), nil)
}

func ServeHTTPIntraCluster(w http.ResponseWriter, req *http.Request) {
	klog.V(5).Infof("Received a request inside the cluster %v, req")

	vpcId, ipAddr, err := parseRequest(req)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error in parse request: %v", err))
		return
	}

	if IsLocalCIDR(vpcId, ipAddr) {
		fmt.Fprintf(w, "Not a request to another cluster. Dropped")
		return
	}

	addr, port, err := getDestGateway(vpcId, ipAddr)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error in getting dest GW: %v", err))
		return
	}

	gw_url := fmt.Sprintf("http://%v/%v/inter-cluster", addr, port)
	postBody, _ := json.Marshal(map[string]string{
		"name":  "Toby",
		"email": "Toby@example.com",
	 })
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(gw_url, "application/json", responseBody)
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error in posting request: %v", err))
		return

	}

	fmt.Fprintf(w, "OK")
}

func ServeHTTPInterCluster(w http.ResponseWriter, req *http.Request) {
	klog.V(4).Infof("\nreceived a request from another cluster %#v\n", req)
	fmt.Fprintf(w, "OK")
}

func parseRequest(req *http.Request) (string, string, error) {
	return "fake-vpc", "fake-addr", nil
}

func getDestGateway(vpc, ipAddr string) (string, uint16, error) {
	return "172.31.9.41", 10005, nil
}

func IsLocalCIDR(vpcId, ipAddr string) bool {
	return true
}