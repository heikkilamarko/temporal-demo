package counter

const TaskQueueName = "demo_counter"

type Counter struct {
	Name     string `json:"name"`
	MaxValue int    `json:"max_value"`
	Value    int    `json:"value"`
}

type IncrementCounterSignal struct {
	Value int `json:"value"`
}

func NewCounter(name string, maxValue int) Counter {
	return Counter{
		Name:     name,
		MaxValue: maxValue,
		Value:    0,
	}
}

func (c *Counter) Increment(value int) {
	c.Value += value
}

func (c *Counter) IsReady() bool {
	return c.MaxValue < c.Value
}
