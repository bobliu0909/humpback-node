package server

import "github.com/humpback/discovery"
import "github.com/humpback/gounits/json"
import "github.com/humpback/gounits/rand"
import "github.com/humpback/humpback-center/cluster"
import "github.com/bobliu0909/humpback-node/etc"

import (
	"flag"
	"fmt"
	"gutils/network"
	"log"
	"net"
	"time"
)

type NodeService struct {
	Key           string
	Configuration *etc.Configuration
	discovery     *discovery.Discovery
	stopCh        chan struct{}
}

// NewNodeService exported
func NewNodeService() (*NodeService, error) {

	var conf string
	flag.StringVar(&conf, "f", "etc/config.yaml", "humpback node configuration file.")
	flag.Parse()
	configuration, err := etc.NewConfiguration(conf)
	if err != nil {
		return nil, err
	}

	key, err := rand.UUIDFile("./humpback-node.key")
	if err != nil {
		return nil, err
	}

	heartbeat, err := time.ParseDuration(configuration.Discovery.Heartbeat)
	if err != nil {
		return nil, err
	}

	ttl, err := time.ParseDuration(configuration.Discovery.TTL)
	if err != nil {
		return nil, err
	}

	configopts := map[string]string{"kv.path": configuration.Discovery.Cluster}
	d, err := discovery.New(configuration.Discovery.URIs, heartbeat, ttl, configopts)
	if err != nil {
		return nil, err
	}

	return &NodeService{
		Key:           key,
		Configuration: configuration,
		discovery:     d,
		stopCh:        make(chan struct{}),
	}, nil
}

func (service *NodeService) Startup() error {

	log.Printf("[#service#] service start...\n")
	addr, err := getServiceAddr(service.Configuration.API.Host)
	if err != nil {
		return err
	}

	regOpts := &cluster.RegistClusterOptions{Addr: addr}
	buf, err := json.EnCodeObjectToBuffer(regOpts)
	if err != nil {
		return err
	}

	service.discovery.Register(service.Key, buf, service.stopCh, func(key string, err error) {
		log.Printf("[#service#] discovery regist %s error:%s\n", key, err.Error())
	})
	return nil
}

func (service *NodeService) Stop() error {

	log.Printf("[#service#] service closed.\n")
	close(service.stopCh)
	return nil
}

func getServiceAddr(host string) (string, error) {

	ip, port, err := net.SplitHostPort(host)
	if err != nil {
		return "", err
	}

	addrs := []string{}
	if len(ip) == 0 {
		nets, err := network.GetLocalNetAddrs()
		if err != nil {
			return "", err
		}
		for _, net := range nets {
			addrs = append(addrs, net.IP+":"+port)
		}
	} else {
		addrs = append(addrs, ip+":"+port)
	}

	selectip, err := autoSelectIp(addrs)
	if err != nil {
		return "", err
	}
	return selectip, nil
}

func autoSelectIp(addrs []string) (string, error) {

	fmt.Println(addrs)
	for _, addr := range addrs {
		if _, err := net.ResolveIPAddr("ip", addr); err != nil {
			log.Printf("[#service#] auto select ip error:%s\n", err.Error())
			continue
		}
		return addr, nil
	}
	return "", fmt.Errorf("local select ipaddr invalid.")
}
