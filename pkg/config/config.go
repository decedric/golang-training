package config

type FibonacciConfig struct {
	HostPort       string
	Domain         string
	TaskListName   string
	ClientName     string
	CadenceService string
	WorkflowName   string
}

func (cfg *FibonacciConfig) SetupConfig() {
	cfg.HostPort = "127.0.0.1:7933"
	cfg.Domain = "test-domain"
	cfg.TaskListName = "fibonacci"
	cfg.ClientName = "cadence-client"
	cfg.CadenceService = "cadence-frontend"
	cfg.WorkflowName = "fibonacci"
}
