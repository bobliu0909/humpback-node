package etc

import "gopkg.in/yaml.v2"

import (
	"io/ioutil"
	"os"
)

type Configuration struct {

	//base options
	Version string `yaml:"version"`
	//service discovery options
	Discovery struct {
		URIs      string `yaml:"uris"`
		Cluster   string `yaml:"cluster"`
		Heartbeat string `yaml:"heartbeat"`
		TTL       string `yaml:"ttl"`
	} `yaml:"discovery"`

	//api options
	API struct {
		Host string `yaml:"host"`
	} `yaml:"api"`
}

func NewConfiguration(file string) (*Configuration, error) {

	fd, err := os.OpenFile(file, os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}

	defer fd.Close()
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	configuration := &Configuration{}
	if err := yaml.Unmarshal([]byte(data), configuration); err != nil {
		return nil, err
	}
	return configuration, nil
}
