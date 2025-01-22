package todolist

import "go.temporal.io/sdk/workflow"

func TodoListWorkflow(ctx workflow.Context, state TodoList) (*TodoList, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("start workflow", "state", state)

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

	selector.AddReceive(updateTodoItemChannel, func(c workflow.ReceiveChannel, _ bool) {
		var signal UpdateTodoItemSignal
		c.Receive(ctx, &signal)

		state.UpdateItem(signal.ID, signal.IsCompleted)
	})

	for {
		selector.Select(ctx)

		if state.IsAllCompleted() {
			break
		}
	}

	logger.Info("stop workflow", "state", state)

	return &state, nil
}
