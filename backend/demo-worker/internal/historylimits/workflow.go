package historylimits

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func HistoryLimitsWorkflow(ctx workflow.Context, state *HistoryLimitsState) (*HistoryLimitsState, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("start workflow")

	err := workflow.SetQueryHandler(ctx, "get_state", func() (*HistoryLimitsState, error) {
		return state, nil
	})
	if err != nil {
		return nil, err
	}

	for state.ActivityExecutionIndex < state.ActivityExecutionCount {
		if workflow.GetInfo(ctx).GetContinueAsNewSuggested() {
			return nil, workflow.NewContinueAsNewError(ctx, HistoryLimitsWorkflow, state)
		}

		_, err := executeActivity(ctx, state)
		if err != nil {
			return nil, err
		}

		state.ActivityExecutionIndex++
	}

	logger.Info("stop workflow")

	return state, nil
}

func executeActivity(ctx workflow.Context, state *HistoryLimitsState) ([]byte, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	})

	var result []byte
	err := workflow.ExecuteActivity(ctx, HistoryLimitsActivity, state.ActivityResultSizeBytes).Get(ctx, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
