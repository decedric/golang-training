package main

import (
	"owt/fibonacci/pkg/api"
	cad "owt/fibonacci/pkg/cadence"
	"owt/fibonacci/pkg/config"
)

func main() {
	var cfg config.FibonacciConfig
	cfg.SetupConfig()
	workflowClient := cad.SetupCadence(&cfg)
	r := api.SetupRouter(workflowClient, &cfg)
	r.Run(":8080")
}
