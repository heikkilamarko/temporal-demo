package todolist

type TodoList struct {
	Items []*TodoItem `json:"items"`
}

type TodoItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsCompleted bool   `json:"is_completed"`
}

type UpdateTodoItemSignal struct {
	ID          string `json:"id"`
	IsCompleted bool   `json:"is_completed"`
}

func (l *TodoList) IsAllCompleted() bool {
	for _, item := range l.Items {
		if !item.IsCompleted {
			return false
		}
	}
	return true
}

func (l *TodoList) GetItemByID(id string) *TodoItem {
	for _, item := range l.Items {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func (l *TodoList) UpdateItem(id string, isCompleted bool) {
	for _, item := range l.Items {
		if item.ID == id {
			item.IsCompleted = isCompleted
			return
		}
	}
}
