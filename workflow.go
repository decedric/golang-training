package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func init() {
	workflow.Register(startFibonacciWorkflow)
	activity.Register(startFibonacciActivity)
}

var WorkflowName string = "fibonacci"

var TotallyPersistentStorage = make(map[string]int)

func startFibonacciWorkflow(ctx workflow.Context, name string, id string) error {
	currentState := "started"
	err := workflow.SetQueryHandler(ctx, "current_state", func() (string, error) {
		return currentState, nil
	})
	if err != nil {
		currentState = "failed to register query handler"
		return err
	}
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("fibonacci workflow started")
	var result int
	err = workflow.ExecuteActivity(ctx, startFibonacciActivity, name).Get(ctx, &result)
	if err != nil {
		currentState = "activity failed"
		logger.Error("Activity failed.", zap.Error(err))
		return err
	}

	TotallyPersistentStorage[id] = result
	currentState = "activity completed"
	logger.Info("Workflow completed.", zap.String("Result", strconv.Itoa(result)))
	return nil
}

func startFibonacciActivity(ctx context.Context, number string) (int, error) {
	n, err := convertAndCheck(number)
	if err != nil {
		panic("please provide a valid number")
	}
	logger := activity.GetLogger(ctx)
	logger.Info("fibonacci activity started")
	fib := fibonacci(n)
	time.Sleep(15 * time.Second)
	return fib, nil
}

func fibonacci(n int) int {
	var last, current int = 0, 1
	if n <= 0 {
		return 0
	}
	for i := 1; i < n; i++ {
		temp := current
		current += last
		last = temp
	}
	return current
}

func convertAndCheck(number string) (int, error) {
	n, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(err)
	}
	return n, err
}
