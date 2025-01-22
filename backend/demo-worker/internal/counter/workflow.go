package counter

import (
	"github.com/mitchellh/mapstructure"
	"go.temporal.io/sdk/workflow"
)

func CounterWorkflow(ctx workflow.Context, state Counter) (*CounterResult, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("start counter", "counter", state)

	if err := state.Validate(); err != nil {
		return nil, err
	}

	err := workflow.SetQueryHandler(ctx, "get_counter", func() (Counter, error) {
		return state, nil
	})
	if err != nil {
		return nil, err
	}

	incrementCounterChannel := workflow.GetSignalChannel(ctx, "increment_counter")
	resetCounterChannel := workflow.GetSignalChannel(ctx, "reset_counter")

	selector := workflow.NewSelector(ctx)

	selector.AddReceive(incrementCounterChannel, func(c workflow.ReceiveChannel, _ bool) {
		var signal any
		c.Receive(ctx, &signal)

		var message IncrementCounterSignal
		err := mapstructure.Decode(signal, &message)
		if err != nil {
			logger.Error("invalid signal type %v", err)
			return
		}

		state.Increment(message.Value)
	})

	selector.AddReceive(resetCounterChannel, func(c workflow.ReceiveChannel, _ bool) {
		var signal any
		c.Receive(ctx, &signal)

		state.Reset()
	})

	selector.AddFuture(workflow.NewTimer(ctx, state.TTLDuration), func(f workflow.Future) {
		state.SetExpired()
	})

	for {
		selector.Select(ctx)

		if state.IsReady() {
			break
		}
	}

	logger.Info("stop counter", "counter", state)

	return &CounterResult{
		Value:     state.Value,
		IsExpired: state.IsExpired,
	}, nil
}
