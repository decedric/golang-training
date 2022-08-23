package cadence

import (
	"context"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"owt/fibonacci/pkg/config"
)

func SetupCadence(c *config.FibonacciConfig) *client.Client {
	dispatcher := buildDispatcher(c)
	domainClient := buildDomainClient(buildServiceClient(dispatcher, c))
	domainClient.Describe(context.Background(), c.Domain)
	service := buildServiceClient(dispatcher, c)
	startWorker(buildLogger(), &service, c)
	workflowClient, _ := buildCadenceClient(buildServiceClient(dispatcher, c), c)
	return &workflowClient

}

func buildCadenceClient(service workflowserviceclient.Interface, c *config.FibonacciConfig) (client.Client, error) {
	return client.NewClient(
		service,
		c.Domain,
		&client.Options{
			Identity: "fibonacci",
		}), nil
}

func buildDomainClient(service workflowserviceclient.Interface) client.DomainClient {
	return client.NewDomainClient(service, &client.Options{})
}

func buildLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	var err error
	logger, err := config.Build()
	if err != nil {
		panic("Failed to setup logger")
	}
	return logger
}

func buildDispatcher(c *config.FibonacciConfig) *yarpc.Dispatcher {
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(c.ClientName))
	if err != nil {
		panic("Failed to setup tchannel")
	}
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: c.ClientName,
		Outbounds: yarpc.Outbounds{
			c.CadenceService: {Unary: ch.NewSingleOutbound(c.HostPort)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		panic("Failed to start dispatcher")
	}
	return dispatcher
}

func buildServiceClient(dispatcher *yarpc.Dispatcher, c *config.FibonacciConfig) workflowserviceclient.Interface {
	return workflowserviceclient.New(dispatcher.ClientConfig(c.CadenceService))
}

func startWorker(logger *zap.Logger, service *workflowserviceclient.Interface, c *config.FibonacciConfig) {
	workerOptions := worker.Options{
		Logger: logger,
	}

	worker := worker.New(
		*service,
		c.Domain,
		c.TaskListName,
		workerOptions)

	err := worker.Start()
	if err != nil {
		panic("Failed to start worker")
	}
	logger.Info("Started Worker.", zap.String("worker", c.TaskListName))
}
