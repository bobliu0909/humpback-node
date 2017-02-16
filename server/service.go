package server

import "github.com/bobliu0909/humpback-node/etc"
import "github.com/humpback/discovery"
import "github.com/humpback/gounits/json"
import "github.com/humpback/gounits/rand"
import "github.com/humpback/humpback-center/cluster/types"

import (
	"flag"
	"log"
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
	registOpts, err := types.NewClusterRegistOptions(service.Configuration.API.Host)
	if err != nil {
		return err
	}

	buf, err := json.EnCodeObjectToBuffer(registOpts)
	if err != nil {
		return err
	}

	log.Printf("[#service#] regist to cluster id:%s, addr %s\n", service.Key, registOpts.Addr)
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
