package historylimits

import "context"

func HistoryLimitsActivity(ctx context.Context, payloadSizeBytes int) ([]byte, error) {
	return make([]byte, payloadSizeBytes), nil
}
