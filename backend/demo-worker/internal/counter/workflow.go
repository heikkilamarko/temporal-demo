package counter

import "go.temporal.io/sdk/workflow"

func CounterWorkflow(ctx workflow.Context, state Counter) (*CounterResult, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("start workflow", "state", state)

	if err := state.Validate(); err != nil {
		return nil, err
	}

	err := workflow.SetQueryHandler(ctx, "get_counter", func() (Counter, error) {
		return state, nil
	})
	if err != nil {
		return nil, err
	}

	timerCtx, cancelTimer := workflow.WithCancel(ctx)

	incrementCounterChannel := workflow.GetSignalChannel(ctx, "increment_counter")
	resetCounterChannel := workflow.GetSignalChannel(ctx, "reset_counter")

	selector := workflow.NewSelector(ctx)

	selector.AddReceive(incrementCounterChannel, func(c workflow.ReceiveChannel, _ bool) {
		var signal IncrementCounterSignal
		c.Receive(ctx, &signal)

		state.Increment(signal.Value)
	})

	selector.AddReceive(resetCounterChannel, func(c workflow.ReceiveChannel, _ bool) {
		var signal any
		c.Receive(ctx, &signal)

		state.Reset()
	})

	selector.AddFuture(workflow.NewTimer(timerCtx, state.TTLDuration), func(_ workflow.Future) {
		state.SetExpired()
	})

	for {
		selector.Select(ctx)

		if state.IsReady() {
			cancelTimer()
			break
		}
	}

	logger.Info("stop workflow", "state", state)

	return &CounterResult{
		Value:     state.Value,
		IsExpired: state.IsExpired,
	}, nil
}
