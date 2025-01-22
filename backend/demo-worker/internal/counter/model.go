package counter

import (
	"fmt"
	"time"
)

const TaskQueueName = "demo_counter"

type Counter struct {
	Name        string        `json:"name"`
	MaxValue    int           `json:"max_value"`
	Value       int           `json:"value"`
	IsExpired   bool          `json:"is_expired"`
	TTL         string        `json:"ttl"`
	TTLDuration time.Duration `json:"-"`
}

type CounterResult struct {
	Value     int  `json:"value"`
	IsExpired bool `json:"is_expired"`
}

type IncrementCounterSignal struct {
	Value int `json:"value"`
}

func (c *Counter) Validate() error {
	var err error
	c.TTLDuration, err = time.ParseDuration(c.TTL)
	if err != nil {
		return fmt.Errorf("parse ttl: %w", err)
	}
	return nil
}

func (c *Counter) IsReady() bool {
	return c.IsExpired || c.MaxValue < c.Value
}

func (c *Counter) Increment(value int) {
	c.Value += value
}

func (c *Counter) Reset() {
	c.Value = 0
}

func (c *Counter) SetExpired() {
	c.IsExpired = true
}
