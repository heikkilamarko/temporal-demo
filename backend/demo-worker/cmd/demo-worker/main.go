package main

import (
	"demo-worker/internal/counter"
	"demo-worker/internal/historylimits"
	"demo-worker/internal/todolist"
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

	w := worker.New(c, "demo_task_queue", worker.Options{})

	w.RegisterWorkflow(counter.CounterWorkflow)
	w.RegisterWorkflow(todolist.TodoListWorkflow)
	w.RegisterWorkflow(historylimits.HistoryLimitsWorkflow)
	w.RegisterActivity(historylimits.HistoryLimitsActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("run worker", err)
	}
}
