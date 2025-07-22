package context

import (
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pfcp"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/smf/internal/logger"
)

func (smContext *SMContext) HandleReports(
	usageReportRequest []*pfcp.UsageReportPFCPSessionReportRequest,
	usageReportModification []*pfcp.UsageReportPFCPSessionModificationResponse,
	usageReportDeletion []*pfcp.UsageReportPFCPSessionDeletionResponse,
	nodeId pfcpType.NodeID, reportTpye models.ChfConvergedChargingTriggerType,
) {
	var usageReport UsageReport
	upf := RetrieveUPFNodeByNodeID(nodeId)
	upfId := upf.UUID()

	for _, report := range usageReportRequest {
		usageReport.UrrId = report.URRID.UrrIdValue
		usageReport.UpfId = upfId
		usageReport.TotalVolume = report.VolumeMeasurement.TotalVolume
		usageReport.UplinkVolume = report.VolumeMeasurement.UplinkVolume
		usageReport.DownlinkVolume = report.VolumeMeasurement.DownlinkVolume
		usageReport.TotalPktNum = report.VolumeMeasurement.TotalPktNum
		usageReport.UplinkPktNum = report.VolumeMeasurement.UplinkPktNum
		usageReport.DownlinkPktNum = report.VolumeMeasurement.DownlinkPktNum
		usageReport.ReportTpye = identityTriggerType(report.UsageReportTrigger)

		if reportTpye != "" {
			usageReport.ReportTpye = reportTpye
		}

		smContext.UrrReports = append(smContext.UrrReports, usageReport)
	}
	for _, report := range usageReportModification {
		usageReport.UrrId = report.URRID.UrrIdValue
		usageReport.UpfId = upfId
		usageReport.TotalVolume = report.VolumeMeasurement.TotalVolume
		usageReport.UplinkVolume = report.VolumeMeasurement.UplinkVolume
		usageReport.DownlinkVolume = report.VolumeMeasurement.DownlinkVolume
		usageReport.TotalPktNum = report.VolumeMeasurement.TotalPktNum
		usageReport.UplinkPktNum = report.VolumeMeasurement.UplinkPktNum
		usageReport.DownlinkPktNum = report.VolumeMeasurement.DownlinkPktNum
		usageReport.ReportTpye = identityTriggerType(report.UsageReportTrigger)

		if reportTpye != "" {
			usageReport.ReportTpye = reportTpye
		}

		smContext.UrrReports = append(smContext.UrrReports, usageReport)
	}
	for _, report := range usageReportDeletion {
		usageReport.UrrId = report.URRID.UrrIdValue
		usageReport.UpfId = upfId
		usageReport.TotalVolume = report.VolumeMeasurement.TotalVolume
		usageReport.UplinkVolume = report.VolumeMeasurement.UplinkVolume
		usageReport.DownlinkVolume = report.VolumeMeasurement.DownlinkVolume
		usageReport.TotalPktNum = report.VolumeMeasurement.TotalPktNum
		usageReport.UplinkPktNum = report.VolumeMeasurement.UplinkPktNum
		usageReport.DownlinkPktNum = report.VolumeMeasurement.DownlinkPktNum
		usageReport.ReportTpye = identityTriggerType(report.UsageReportTrigger)

		if reportTpye != "" {
			usageReport.ReportTpye = reportTpye
		}

		smContext.UrrReports = append(smContext.UrrReports, usageReport)
	}
}

func identityTriggerType(usarTrigger *pfcpType.UsageReportTrigger) models.ChfConvergedChargingTriggerType {
	var trigger models.ChfConvergedChargingTriggerType

	switch {
	case usarTrigger.Volth:
		trigger = models.ChfConvergedChargingTriggerType_QUOTA_THRESHOLD
	case usarTrigger.Volqu:
		trigger = models.ChfConvergedChargingTriggerType_QUOTA_EXHAUSTED
	case usarTrigger.Quvti:
		trigger = models.ChfConvergedChargingTriggerType_VALIDITY_TIME
	case usarTrigger.Start:
		trigger = models.ChfConvergedChargingTriggerType_START_OF_SERVICE_DATA_FLOW
	case usarTrigger.Immer:
		logger.PduSessLog.Trace("Reports Query by SMF, trigger should be filled later")
		return ""
	case usarTrigger.Termr:
		trigger = models.ChfConvergedChargingTriggerType_FINAL
	default:
		logger.PduSessLog.Trace("Report is not a charging trigger")
		return ""
	}

	return trigger
}
