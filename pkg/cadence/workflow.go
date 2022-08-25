package cadence

import (
	"context"
	"math/big"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func init() {
	workflow.Register(StartFibonacciWorkflow)
	activity.Register(startFibonacciActivity)
}

var TotallyPersistentStorage = make(map[string]*big.Int)

func GetResult(id string) *big.Int {
	return TotallyPersistentStorage[id]
}

func StartFibonacciWorkflow(ctx workflow.Context, number uint, id string) error {
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
	var result *big.Int
	err = workflow.ExecuteActivity(ctx, startFibonacciActivity, number).Get(ctx, &result)
	if err != nil {
		currentState = "activity failed"
		logger.Error("Activity failed.", zap.Error(err))
		return err
	}

	TotallyPersistentStorage[id] = result
	currentState = "activity completed"
	logger.Info("Workflow completed.")
	return nil
}

func startFibonacciActivity(ctx context.Context, number uint) (*big.Int, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("fibonacci activity started")
	fib := fibonacci(number)
	// This sleep can be uncommented to create some time buffer s.t. the polling service can be tested
	//time.Sleep(15 * time.Second)
	return fib, nil
}

func fibonacci(n uint) *big.Int {
	last, current := big.NewInt(0), big.NewInt(1)
	if n <= 0 {
		return big.NewInt(0)
	}
	var i uint
	for i = 1; i < n; i++ {
		last.Add(current, last)
		last, current = current, last
	}
	return current
}
