package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type FibonacciConfig struct {
	HostPort       string
	Domain         string
	TaskListName   string
	ClientName     string
	CadenceService string
	WorkflowName   string
}

func (cfg *FibonacciConfig) SetupConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file %w", err))
	}
	err = viper.Unmarshal(cfg)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %w", err))
	}
}
