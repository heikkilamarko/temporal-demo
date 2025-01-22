package counter

import (
	"github.com/mitchellh/mapstructure"
	"go.temporal.io/sdk/workflow"
)

func CounterWorkflow(ctx workflow.Context, state Counter) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("start counter", "counter", state)

	err := workflow.SetQueryHandler(ctx, "get_counter", func() (Counter, error) {
		return state, nil
	})
	if err != nil {
		return err
	}

	incrementCounterChannel := workflow.GetSignalChannel(ctx, "increment_counter")

	for {
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

		selector.Select(ctx)

		if state.IsReady() {
			break
		}
	}

	logger.Info("stop counter", "counter", state)

	return nil
}
