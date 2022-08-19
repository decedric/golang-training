package main

import (
	"context"
	"fmt"
	"github.com/pborman/uuid"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
	"time"
)

func main() {

	dispatcher := buildDispatcher()
	domainClient := buildDomainClient(buildServiceClient(dispatcher))
	domainClient.Describe(context.Background(), Domain)
	service := buildServiceClient(dispatcher)
	startWorker(buildLogger(), &service)
	workflowOptions := client.StartWorkflowOptions{
		ID:                              WorkflowName + "_" + uuid.New(),
		TaskList:                        TaskListName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	workflowClient, _ := buildCadenceClient(buildServiceClient(dispatcher))
	we, err := workflowClient.StartWorkflow(context.Background(), workflowOptions, startFibonacciWorkflow, "100")

	if err != nil {
		fmt.Println("Failed to create workflow")
		fmt.Println(zap.Error(err))
		panic("Failed to create workflow.")
	} else {
		fmt.Println("Started Workflow", zap.String("WorkflowID", we.ID))
		fmt.Println(zap.String("RunID", we.RunID))
	}
	select {}
}
