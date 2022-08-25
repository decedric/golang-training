package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/cadence/client"
	"owt/fibonacci/pkg/config"
)

func SetupRouter(workflowClient *client.Client, cfg *config.FibonacciConfig) *gin.Engine {
	r := gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		panic("could not disable proxies")
	}
	setupFibonacciRoutes(r, workflowClient, cfg)
	return r
}
