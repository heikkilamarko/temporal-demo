package main

import (
	"context"
	"demo-worker/internal"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort:  client.DefaultHostPort,
		Namespace: client.DefaultNamespace,
	})
	if err != nil {
		log.Fatalln("create temporal client", err)
	}

	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:                                       "demo",
		TaskQueue:                                internal.TaskQueueName,
		WorkflowExecutionErrorWhenAlreadyStarted: true,
	}

	var data internal.DemoData

	_, err = c.ExecuteWorkflow(context.Background(), options, internal.DemoWorkflow, data)
	if err != nil {
		log.Fatalln("execute demo workflow", err)
	}
}
