//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type QosFlowNotifyItem struct {
	QosFlowIdentifier QosFlowIdentifier
	NotificationCause NotificationCause
	IEExtensions      *ProtocolExtensionContainerQosFlowNotifyItemExtIEs `aper:"optional"`
}
