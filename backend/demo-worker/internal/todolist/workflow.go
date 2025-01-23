package todolist

import (
	"errors"

	"go.temporal.io/sdk/workflow"
)

func TodoListWorkflow(ctx workflow.Context, state TodoList) (*TodoList, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("start workflow")

	err := workflow.SetQueryHandler(ctx, "get_todo_list", func() (TodoList, error) {
		return state, nil
	})
	if err != nil {
		return nil, err
	}

	err = workflow.SetQueryHandler(ctx, "get_todo_item", func(id string) (*TodoItem, error) {
		return state.GetItemByID(id), nil
	})
	if err != nil {
		return nil, err
	}

	updateTodoItemChannel := workflow.GetSignalChannel(ctx, "update_todo_item")

	selector := workflow.NewSelector(ctx)

	selector.AddReceive(ctx.Done(), func(_ workflow.ReceiveChannel, _ bool) {
		logger.Info("workflow context done")
	})

	selector.AddReceive(updateTodoItemChannel, func(c workflow.ReceiveChannel, _ bool) {
		var signal UpdateTodoItemSignal
		c.Receive(ctx, &signal)

		state.UpdateItem(signal.ID, signal.IsCompleted)
	})

	for {
		if err := ctx.Err(); errors.Is(err, workflow.ErrCanceled) {
			logger.Info("workflow canceled")
			return nil, err
		}

		if state.IsAllCompleted() {
			logger.Info("all items are completed")
			break
		}

		selector.Select(ctx)
	}

	logger.Info("stop workflow")

	return &state, nil
}
