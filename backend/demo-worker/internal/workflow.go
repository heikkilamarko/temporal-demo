package internal

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func DemoWorkflow(ctx workflow.Context, data DemoData) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("DemoWorkflow started", "data", data)

	var activities *Activities

	demoCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	})

	err := workflow.ExecuteActivity(demoCtx, activities.DemoActivity).Get(demoCtx, nil)
	if err != nil {
		return err
	}

	data.Count++

	logger.Info("DemoWorkflow completed", "data", data)

	return nil
}
