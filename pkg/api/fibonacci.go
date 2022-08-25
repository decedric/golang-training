package api

import (
	"context"
	"fmt"
	"net/http"
	"owt/fibonacci/pkg/cadence"
	"owt/fibonacci/pkg/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
)

type ControllerData struct {
	cfg            *config.FibonacciConfig
	workflowClient *client.Client
}

func setupFibonacciRoutes(r *gin.Engine, workflowClient *client.Client, cfg *config.FibonacciConfig) {
	cd := ControllerData{cfg, workflowClient}
	endpoint := r.Group("/fibonacci")
	endpoint.POST("/:number", cd.startWorkflow)
	endpoint.GET("/polling/:id", cd.pollWorkflow)
	endpoint.GET("/:id", cd.getResult)
}

type Number struct {
	Value int `uri:"number" binding:"required,gte=0"`
}

func (cd *ControllerData) startWorkflow(c *gin.Context) {
	var number Number
	if err := c.ShouldBindUri(&number); err != nil {
		c.String(http.StatusBadRequest, "Please provide positive number")
		return
	}

	workflowId := cd.cfg.WorkflowName + "_" + uuid.New()
	workflowOptions := client.StartWorkflowOptions{
		ID:                              workflowId,
		TaskList:                        cd.cfg.TaskListName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	_, err := (*cd.workflowClient).StartWorkflow(context.Background(), workflowOptions, cadence.StartFibonacciWorkflow, uint(number.Value), workflowId)
	if err != nil {
		fmt.Println("Failed to create workflow", zap.Error(err))
		c.String(http.StatusInternalServerError, "Failed to create workflow.")
	}
	c.JSON(http.StatusCreated, gin.H{
		"address": fmt.Sprintf("fibonacci/polling/%s", workflowId),
		"id":      workflowId,
	})
}

type WorkflowId struct {
	Id string `uri:"id" buiding:"required"`
}

func (cd *ControllerData) pollWorkflow(c *gin.Context) {
	var workflowId WorkflowId
	if err := c.ShouldBindUri(&workflowId); err != nil {
		c.String(http.StatusBadRequest, "Please provide positive number")
		return
	}
	response, _ := (*cd.workflowClient).QueryWorkflow(context.Background(), workflowId.Id, "", "current_state")
	var status string
	response.Get(&status)
	if status == "activity completed" {
		c.JSON(http.StatusFound, gin.H{
			"address": fmt.Sprintf("fibonacci/%s", workflowId.Id),
			"id":      workflowId.Id,
			"status":  status,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"address": fmt.Sprintf("fibonacci/polling/%s", workflowId.Id),
			"id":      workflowId.Id,
			"status":  status,
		})
	}
}

func (cd *ControllerData) getResult(c *gin.Context) {
	var workflowId WorkflowId
	if err := c.ShouldBindUri(&workflowId); err != nil {
		c.String(http.StatusBadRequest, "Please provide positive number")
		return
	}
	result := cadence.GetResult(workflowId.Id)
	c.JSON(http.StatusOK, gin.H{
		"result": result.Text(10),
	})
}
