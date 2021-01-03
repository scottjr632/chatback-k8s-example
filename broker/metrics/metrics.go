package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var metricsInstance *metrics
var metricsOnce sync.Once

type metrics struct {
	activeClients prometheus.Counter
	messagesSent  prometheus.Counter
}

// Metrics is an interface to the broker metrics
type Metrics interface {
	IncActiceClients()
	DecActiveClients()

	IncMessagesSent()
	DecMessagesSent()
}

// Get returns a metrics instance
func Get() Metrics {
	metricsOnce.Do(func() {
		metrics := &metrics{
			activeClients: createActiveClientCounter(),
			messagesSent:  createMessagesSentCounter(),
		}
		metricsInstance = metrics
	})

	return metricsInstance
}

func (m *metrics) IncActiceClients() {
	m.activeClients.Inc()
}

func (m *metrics) DecActiveClients() {
	m.activeClients.Desc()
}

func (m *metrics) IncMessagesSent() {
	m.messagesSent.Inc()
}

func (m *metrics) DecMessagesSent() {
	m.messagesSent.Desc()
}

func createActiveClientCounter() prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Name: "chatback_broker_active_clients",
		Help: "The total number of active clients connected",
	})
}

func createMessagesSentCounter() prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Name: "chatback_broker_messages_sent",
		Help: "The total number of messages sent",
	})
}
