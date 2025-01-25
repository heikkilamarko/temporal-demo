package historylimits

import "context"

func HistoryLimitsActivity(ctx context.Context, resultSizeBytes int) ([]byte, error) {
	return make([]byte, resultSizeBytes), nil
}
