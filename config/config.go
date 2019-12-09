package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

const (
	defaultClientConfigFile = "config/client.yaml"
	defaultServerConfigFile = "config/server.yaml"
	doPanic                 = true
)

// ClientConfig client
type ClientConfig struct {
	Version string `yaml:"version"`
	Server  struct {
		Host    string        `yaml:"host"`
		Port    int           `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"server"`
	Wal struct {
		Datadir string `yaml:"datadir"`
	} `yaml:"wal"`
}

// ServerConfig server
type ServerConfig struct {
	Server struct {
		Host    string        `yaml:"host"`
		Port    int           `yaml:"port"`
		DB      string        `yaml:"dbname"`
		Timeout time.Duration `yaml:"timeout"`
	} `yaml:"server"`
	Wal struct {
		Datadir string `yaml:"datadir"`
	} `yaml:"wal"`
}

func check(err error, methodSign string) {
	msg := fmt.Sprintf("Failed while running method %s, Error %v", methodSign, err)
	if !doPanic {
		log.Printf(msg)
		return
	}
	if err != nil {
		log.Fatalf(msg)
	}
}

func loadServerConfig() (cfg *ServerConfig) {
	configFile, err := ioutil.ReadFile(defaultClientConfigFile)
	check(err, "loadClientConfig")
	err = yaml.Unmarshal(configFile, &cfg)
	check(err, "loadServerConfig")
	return cfg
}

func loadClientConfig() (cfg *ClientConfig) {
	configFile, err := ioutil.ReadFile(defaultClientConfigFile)
	check(err, "loadClientConfig")
	err = yaml.Unmarshal(configFile, &cfg)
	check(err, "loadClientConfig")
	return cfg
}

// Config function loads and returns the config based on `configName` parameter
func Config(configName string) (cfg interface{}) {
	switch strings.ToLower(configName) {
	case "client":
		cfg = loadClientConfig()
	case "server":
		cfg = loadServerConfig()
	default:
		log.Fatalf("Invalid config")
	}
	return cfg
}
