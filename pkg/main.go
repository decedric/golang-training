package main

func main() {
	workflowClient := SetupCadence()
	r := SetupRouter(workflowClient)
	r.Run(":8080")
}
