package hasura

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/forbole/juno/v4/node/remote"
	"github.com/spf13/viper"
)

// Config contains the configuration about the actions module
type Config struct {
	Host string          `yaml:"host"`
	Port uint            `yaml:"port"`
	Node *remote.Details `yaml:"node,omitempty"`
}

// NewConfig returns a new Config instance
func NewConfig(host string, port uint, remoteDetails *remote.Details) *Config {
	return &Config{
		Host: host,
		Port: port,
		Node: remoteDetails,
	}
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Host: "127.0.0.1",
		Port: 3000,
		Node: nil,
	}
}

func LoadHasuraConfig(homepath string) (config *Config, err error) {
	cfgToml := Config{}
	v := viper.New()
	configPath := filepath.Join(homepath, "config")

	configFilePath := filepath.Join(configPath, "hasura.yaml")
	if _, err = os.Stat(configFilePath); os.IsNotExist(err) {
		return
	}

	v.AddConfigPath(configPath)
	v.SetConfigName("hasura")
	v.SetConfigType("yaml")
	err = v.ReadInConfig()
	if err != nil {
		fmt.Printf("err=%v", err)
		return
	}

	err = v.Unmarshal(&cfgToml)
	config = &Config{
		Host: cfgToml.Host,
		Port: cfgToml.Port,
		Node: cfgToml.Node,
	}
	return
}
