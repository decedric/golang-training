package main

import (
	"context"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var HostPort = "127.0.0.1:7933"
var Domain = "test-domain"
var TaskListName = "SimpleWorker"
var ClientName = "cadence-client"
var CadenceService = "cadence-frontend"

func SetupCadence() client.Client {
	dispatcher := buildDispatcher()
	domainClient := buildDomainClient(buildServiceClient(dispatcher))
	domainClient.Describe(context.Background(), Domain)
	service := buildServiceClient(dispatcher)
	startWorker(buildLogger(), &service)
	workflowClient, _ := buildCadenceClient(buildServiceClient(dispatcher))
	return workflowClient

}

func buildCadenceClient(service workflowserviceclient.Interface) (client.Client, error) {
	return client.NewClient(
		service,
		Domain,
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

func buildDispatcher() *yarpc.Dispatcher {
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientName))
	if err != nil {
		panic("Failed to setup tchannel")
	}
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: ClientName,
		Outbounds: yarpc.Outbounds{
			CadenceService: {Unary: ch.NewSingleOutbound(HostPort)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		panic("Failed to start dispatcher")
	}
	return dispatcher
}

func buildServiceClient(dispatcher *yarpc.Dispatcher) workflowserviceclient.Interface {
	return workflowserviceclient.New(dispatcher.ClientConfig(CadenceService))
}

func startWorker(logger *zap.Logger, service *workflowserviceclient.Interface) {
	workerOptions := worker.Options{
		Logger: logger,
	}

	worker := worker.New(
		*service,
		Domain,
		TaskListName,
		workerOptions)

	err := worker.Start()
	if err != nil {
		panic("Failed to start worker")
	}
	logger.Info("Started Worker.", zap.String("worker", TaskListName))
}
