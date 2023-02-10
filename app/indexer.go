package app

import (
	"fmt"
	"os"
	"path/filepath"

	dbconfig "github.com/forbole/juno/v4/database/config"
	"github.com/spf13/viper"
)

type IndexerConfigFromToml struct {
	Name               string `mapstructure:"name"`
	Host               string `mapstructure:"host"`
	Port               int64  `mapstructure:"port"`
	User               string `mapstructure:"user"`
	Password           string `mapstructure:"password"`
	SSLMode            string `mapstructure:"ssl_mode,omitempty"`
	Schema             string `mapstructure:"schema,omitempty"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
	PartitionSize      int64  `mapstructure:"partition_size"`
	PartitionBatchSize int64  `mapstructure:"partition_batch"`
}

func LoadIndexerConfig(homePath string) (indexerConfig dbconfig.Config, err error) {
	cfgToml := IndexerConfigFromToml{}
	v := viper.New()
	configPath := filepath.Join(homePath, "config")

	configFilePath := filepath.Join(configPath, "indexer.toml")
	if _, err = os.Stat(configFilePath); os.IsNotExist(err) {
		return
	}

	v.AddConfigPath(configPath)
	v.SetConfigName("indexer")
	v.SetConfigType("toml")
	err = v.ReadInConfig()
	if err != nil {
		fmt.Printf("err=%v", err)
		return
	}

	err = v.Unmarshal(&cfgToml)
	indexerConfig = dbconfig.Config{
		URL:                fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s&schema=%s", cfgToml.User, cfgToml.Password, cfgToml.Host, cfgToml.Port, cfgToml.Name, cfgToml.SSLMode, cfgToml.Schema),
		MaxOpenConnections: cfgToml.MaxOpenConnections,
		MaxIdleConnections: cfgToml.MaxIdleConnections,
		PartitionSize:      cfgToml.PartitionSize,
		PartitionBatchSize: cfgToml.PartitionBatchSize,
	}

	return
}
