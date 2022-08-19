package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func setupRouter(workflowClient client.Client) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/fibonacci/:number", func(c *gin.Context) {
		number := c.Params.ByName("number")
		workflowId := WorkflowName + "_" + uuid.New()
		workflowOptions := client.StartWorkflowOptions{
			ID:                              workflowId,
			TaskList:                        TaskListName,
			ExecutionStartToCloseTimeout:    time.Minute,
			DecisionTaskStartToCloseTimeout: time.Minute,
		}

		we, err := workflowClient.StartWorkflow(context.Background(), workflowOptions, startFibonacciWorkflow, number)

		if err != nil {
			fmt.Println("Failed to create workflow")
			fmt.Println(zap.Error(err))
			panic("Failed to create workflow.")
		} else {
			fmt.Println("Started Workflow", zap.String("WorkflowID", we.ID))
			fmt.Println(zap.String("RunID", we.RunID))
		}
		c.JSON(http.StatusCreated, gin.H{
			"address": fmt.Sprintf("fibonacci/polling/%s", workflowId),
		})
	})

	r.GET("/fibonacci/polling/:id", func(c *gin.Context) {

	})

	r.GET("/fibonacci/:id", func(c *gin.Context) {

	})

	return r
}

func main() {

	dispatcher := buildDispatcher()
	domainClient := buildDomainClient(buildServiceClient(dispatcher))
	domainClient.Describe(context.Background(), Domain)
	service := buildServiceClient(dispatcher)
	startWorker(buildLogger(), &service)
	workflowClient, _ := buildCadenceClient(buildServiceClient(dispatcher))

	r := setupRouter(workflowClient)
	r.Run(":8080")
	select {}
}
