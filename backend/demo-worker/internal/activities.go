package internal

import (
	"context"
	"log/slog"
)

type Activities struct {
	Message string
}

func NewActivities(message string) *Activities {
	return &Activities{message}
}

func (a *Activities) DemoActivity(ctx context.Context) error {
	slog.Info("DemoActivity", "message", a.Message)
	return nil
}
