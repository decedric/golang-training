package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/cadence/client"
	"owt/fibonacci/pkg/config"
)

var validate *validator.Validate

func SetupRouter(workflowClient *client.Client, cfg *config.FibonacciConfig) *gin.Engine {
	validate = validator.New()
	r := gin.Default()
	r.SetTrustedProxies(nil)
	setupFibonacciRoutes(r, workflowClient, cfg)
	return r
}
