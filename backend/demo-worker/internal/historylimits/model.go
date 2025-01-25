package historylimits

type HistoryLimitsState struct {
	ActivityExecutionCount  int `json:"activity_execution_count"`
	ActivityExecutionIndex  int `json:"activity_execution_index"`
	ActivityResultSizeBytes int `json:"activity_result_size_bytes"`
}
