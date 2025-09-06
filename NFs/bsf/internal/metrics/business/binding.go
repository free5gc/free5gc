package business

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/free5gc/util/metrics/utils"
)

var (
	// pcfBindingGauge Gauge for the number of PCF bindings currently stored in BSF
	// labeled by binding type
	pcfBindingGauge *prometheus.GaugeVec
	// pcfBindingEventCounter Counter for PCF binding events (create, update, delete, query)
	// labeled by binding type, event, and result
	pcfBindingEventCounter *prometheus.CounterVec
	// bindingDuration Histogram for time spent processing binding operations in seconds
	// labeled by binding type and event
	bindingDuration *prometheus.HistogramVec
)

func GetBindingHandlerMetrics(namespace string) []prometheus.Collector {
	var collectors []prometheus.Collector

	pcfBindingGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      PCF_BINDING_GAUGE_NAME,
			Help:      PCF_BINDING_GAUGE_DESC,
		},
		[]string{BINDING_TYPE_LABEL},
	)

	// Initialize gauge with zero values for all binding types
	pcfBindingGauge.With(prometheus.Labels{BINDING_TYPE_LABEL: PCF_BINDING_TYPE_VALUE}).Set(0)
	pcfBindingGauge.With(prometheus.Labels{BINDING_TYPE_LABEL: PCF_UE_BINDING_TYPE_VALUE}).Set(0)
	pcfBindingGauge.With(prometheus.Labels{BINDING_TYPE_LABEL: PCF_MBS_BINDING_TYPE_VALUE}).Set(0)

	pcfBindingEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      PCF_BINDING_EVENT_COUNTER_NAME,
			Help:      PCF_BINDING_EVENT_COUNTER_DESC,
		},
		[]string{BINDING_TYPE_LABEL, BINDING_EVENT_LABEL, BINDING_RESULT_LABEL},
	)

	bindingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      BINDING_DURATION_HISTOGRAM_NAME,
			Help:      BINDING_DURATION_HISTOGRAM_DESC,
			Buckets: []float64{
				0.0010, // 1ms
				0.0050, // 5ms
				0.0100, // 10ms
				0.0250, // 25ms
				0.0500, // 50ms
				0.1000, // 100ms
				0.2500, // 250ms
				0.5000, // 500ms
				1.0000, // 1s
			},
		},
		[]string{BINDING_TYPE_LABEL, BINDING_EVENT_LABEL},
	)

	collectors = append(collectors, pcfBindingGauge, pcfBindingEventCounter, bindingDuration)

	return collectors
}

// IncrPCFBindingGauge increments the PCF binding gauge for the given binding type
func IncrPCFBindingGauge(bindingType string) {
	if utils.IsBusinessMetricsEnabled() && IsBindingMetricsEnabled() {
		pcfBindingGauge.With(prometheus.Labels{
			BINDING_TYPE_LABEL: bindingType,
		}).Inc()
	}
}

// DecrPCFBindingGauge decrements the PCF binding gauge for the given binding type
func DecrPCFBindingGauge(bindingType string) {
	if utils.IsBusinessMetricsEnabled() && IsBindingMetricsEnabled() {
		pcfBindingGauge.With(prometheus.Labels{
			BINDING_TYPE_LABEL: bindingType,
		}).Dec()
	}
}

// IncrPCFBindingEventCounter increments the PCF binding event counter
func IncrPCFBindingEventCounter(bindingType string, event string, result string) {
	if utils.IsBusinessMetricsEnabled() && IsBindingMetricsEnabled() {
		pcfBindingEventCounter.With(prometheus.Labels{
			BINDING_TYPE_LABEL:   bindingType,
			BINDING_EVENT_LABEL:  event,
			BINDING_RESULT_LABEL: result,
		}).Inc()
	}
}

// ObserveBindingDuration observes the duration of a binding operation
func ObserveBindingDuration(bindingType string, event string, startTime time.Time) {
	if utils.IsBusinessMetricsEnabled() && IsBindingMetricsEnabled() {
		if startTime.IsZero() {
			return
		}

		duration := time.Since(startTime).Seconds()

		bindingDuration.With(prometheus.Labels{
			BINDING_TYPE_LABEL:  bindingType,
			BINDING_EVENT_LABEL: event,
		}).Observe(duration)
	}
}
