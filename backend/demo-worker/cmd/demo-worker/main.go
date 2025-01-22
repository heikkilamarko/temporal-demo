package main

import (
	"demo-worker/internal"
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort:  os.Getenv("APP_TEMPORAL_HOSTPORT"),
		Namespace: client.DefaultNamespace,
	})
	if err != nil {
		log.Fatalln("create temporal client", err)
	}

	defer c.Close()

	w := worker.New(c, internal.TaskQueueName, worker.Options{})

	activities := internal.NewActivities("hello world")

	w.RegisterWorkflow(internal.DemoWorkflow)
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("run worker", err)
	}
}
