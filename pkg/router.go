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

func SetupRouter(workflowClient client.Client) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/fibonacci/:number", func(c *gin.Context) {
		number := c.Params.ByName("number")
		if number[0] == '-' {
			c.String(http.StatusBadRequest, "Please provide a positive number.")
		}
		workflowId := WorkflowName + "_" + uuid.New()
		workflowOptions := client.StartWorkflowOptions{
			ID:                              workflowId,
			TaskList:                        TaskListName,
			ExecutionStartToCloseTimeout:    time.Minute,
			DecisionTaskStartToCloseTimeout: time.Minute,
		}

		we, err := workflowClient.StartWorkflow(context.Background(), workflowOptions, startFibonacciWorkflow, number, workflowId)

		if err != nil {
			fmt.Println("Failed to create workflow", zap.Error(err))
			c.String(http.StatusInternalServerError, "Failed to create workflow.")
		} else {
			fmt.Println("Started Workflow", zap.String("WorkflowID", we.ID))
		}
		c.JSON(http.StatusCreated, gin.H{
			"address": fmt.Sprintf("fibonacci/polling/%s", workflowId),
			"id":      workflowId,
		})
	})

	r.GET("/fibonacci/polling/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		response, _ := workflowClient.QueryWorkflow(context.Background(), id, "", "current_state")
		var status string
		response.Get(&status)
		if status == "activity completed" {
			c.JSON(http.StatusFound, gin.H{
				"address": fmt.Sprintf("fibonacci/%s", id),
				"id":      id,
				"status":  status,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"address": fmt.Sprintf("fibonacci/polling/%s", id),
				"id":      id,
				"status":  status,
			})
		}
	})

	r.GET("/fibonacci/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		result := TotallyPersistentStorage[id]
		c.JSON(http.StatusOK, gin.H{
			"result": result.Text(10),
		})
	})

	return r
}
