package metric

import (
	"fmt"
	"time"

	"github.com/ncarlier/apimon/pkg/config"
)

// Queue metric queue
var Queue = make(chan Metric)

// Metric DTO
type Metric struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	Error     string        `json:"error,omitempty"`
}

var producer *Producer

// StartMetricProducer start metric producer
func StartMetricProducer(conf config.Output) error {
	var err error
	producer, err = NewMetricProducer(conf)
	if err != nil {
		return fmt.Errorf("unable to start metric producer: %s - %s", conf.Target, err)
	}
	producer.Start()
	return nil
}

// StopMetricProducer stops metric producer
func StopMetricProducer() {
	producer.Stop()
}
