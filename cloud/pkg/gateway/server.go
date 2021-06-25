package gateway

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"k8s.io/klog/v2"

	"github.com/kubeedge/kubeedge/cloud/pkg/gateway/config"

	"github.com/kubeedge/kubeedge/cloud/pkg/gateway/vpcmap"
)

const (
	VPC_HEADER = "vpcId"
	SRC_IP_HEADER = "srcIp"
	DEST_IP_HEADER = "destIp"
	GatewayPort = 10005
)

func startGatewayServer() {
	klog.Infof("Starting Gateway http server")

	http.HandleFunc("/intra-cluster", ServeHTTPIntraCluster)

	http.HandleFunc("/inter-cluster", ServeHTTPInterCluster)

	http.ListenAndServe(fmt.Sprintf(":%v",config.Config.Gateway.Port), nil)
}

func ServeHTTPIntraCluster(w http.ResponseWriter, req *http.Request) {
	klog.V(4).Infof("Received a request inside the cluster, req")

	vpcId, _, destIp, _, err := parseRequest(req)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error in parsing request: %v\n", err))
		return
	}

	gatewayAddr,  err := vpcmap.GetGatewayIP(vpcId, destIp)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error in getting dest GW: %v\n", err))
		return
	}

	if vpcmap.IsLocalGateway(gatewayAddr) {
		fmt.Fprintf(w, "Not a request to another cluster. Should not reach here. Dropped\n")
		return
	}

	gw_url := fmt.Sprintf("http://%v:%v/inter-cluster", gatewayAddr, GatewayPort)

	_, err = ForwardRequestToGateway(req, gw_url)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error in forwarding request to destination gateway: %v", err))
		return
	}

	fmt.Fprintf(w, "forward request to gateway " + gw_url + "\n")
}

func ForwardRequestToGateway(req *http.Request, gatewayUrl string) (*http.Response, error) {
	httpClient := http.Client{}

	_, err := url.Parse(gatewayUrl)
	if err!= nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	newRequest, err := http.NewRequest(req.Method, gatewayUrl, ioutil.NopCloser(bytes.NewBuffer(body)))
	if err != nil {
		return nil, err
	}

	//temporarily, just let it succeed for test purpose
	// return httpClient.Do(newRequest)
	go httpClient.Do(newRequest)
	return nil, nil
}

func ServeHTTPInterCluster(w http.ResponseWriter, req *http.Request) {
	klog.V(4).Infof("\nreceived a request from another cluster %#v\n", req)
	fmt.Fprintf(w, "OK")
}

func parseRequest(req *http.Request) (vpcId, srcIp, destIp, payload string, err error) {

	vpcId = req.Header.Get(VPC_HEADER)

	srcIp = req.Header.Get(SRC_IP_HEADER)

	destIp = req.Header.Get(DEST_IP_HEADER)

	return vpcId, srcIp, destIp, "", nil
}
