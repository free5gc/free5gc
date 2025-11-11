package business

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/free5gc/util/metrics/utils"
)

// pcfDiscoveryCounter Counter for PCF discovery requests processed by BSF
// labeled by discovery type and result
var pcfDiscoveryCounter *prometheus.CounterVec

func GetDiscoveryHandlerMetrics(namespace string) []prometheus.Collector {
	var collectors []prometheus.Collector

	pcfDiscoveryCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      PCF_DISCOVERY_COUNTER_NAME,
			Help:      PCF_DISCOVERY_COUNTER_DESC,
		},
		[]string{DISCOVERY_TYPE_LABEL, DISCOVERY_RESULT_LABEL},
	)

	collectors = append(collectors, pcfDiscoveryCounter)

	return collectors
}

// IncrPCFDiscoveryCounter increments the PCF discovery counter
func IncrPCFDiscoveryCounter(discoveryType string, result string) {
	if utils.IsBusinessMetricsEnabled() && IsDiscoveryMetricsEnabled() {
		pcfDiscoveryCounter.With(prometheus.Labels{
			DISCOVERY_TYPE_LABEL:   discoveryType,
			DISCOVERY_RESULT_LABEL: result,
		}).Inc()
	}
}
