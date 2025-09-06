package business

// Global metric information
const (
	SUBSYSTEM_NAME    = "bsf_business"
	BINDING_METRICS   = "binding"
	DISCOVERY_METRICS = "discovery"
)

// Collectors information
const (
	// PCF Bindings
	PCF_BINDING_GAUGE_NAME         = "pcf_binding_current_count"
	PCF_BINDING_GAUGE_DESC         = "Current number of PCF bindings stored in the BSF"
	PCF_BINDING_EVENT_COUNTER_NAME = "pcf_binding_events_total"
	PCF_BINDING_EVENT_COUNTER_DESC = "Count of PCF binding events (creation, updates, deletions)"

	// PCF Discovery
	PCF_DISCOVERY_COUNTER_NAME = "pcf_discovery_total"
	PCF_DISCOVERY_COUNTER_DESC = "Count of PCF discovery requests processed by BSF"

	// Binding Duration
	BINDING_DURATION_HISTOGRAM_NAME = "binding_duration_seconds"
	BINDING_DURATION_HISTOGRAM_DESC = "Histogram of binding duration in seconds"
)

// Label names
const (
	// Binding Type Labels
	BINDING_TYPE_LABEL   = "binding_type"
	BINDING_EVENT_LABEL  = "event"
	BINDING_RESULT_LABEL = "result"

	// Discovery Labels
	DISCOVERY_TYPE_LABEL   = "discovery_type"
	DISCOVERY_RESULT_LABEL = "result"
)

// Metrics Values
const (
	// Binding Types
	PCF_BINDING_TYPE_VALUE     = "pcf_binding"
	PCF_UE_BINDING_TYPE_VALUE  = "pcf_ue_binding"
	PCF_MBS_BINDING_TYPE_VALUE = "pcf_mbs_binding"

	// Binding Events
	BINDING_EVENT_CREATE_VALUE = "create"
	BINDING_EVENT_UPDATE_VALUE = "update"
	BINDING_EVENT_DELETE_VALUE = "delete"
	BINDING_EVENT_QUERY_VALUE  = "query"

	// Discovery Types
	DISCOVERY_TYPE_PCF_VALUE = "pcf"

	// Result Values
	RESULT_SUCCESS_VALUE = "success"
	RESULT_FAILURE_VALUE = "failure"
)

var bindingMetricsEnabled bool

func IsBindingMetricsEnabled() bool {
	return bindingMetricsEnabled
}

func EnableBindingMetrics() {
	bindingMetricsEnabled = true
}

var discoveryMetricsEnabled bool

func IsDiscoveryMetricsEnabled() bool {
	return discoveryMetricsEnabled
}

func EnableDiscoveryMetrics() {
	discoveryMetricsEnabled = true
}
