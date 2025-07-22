package logger

import (
	"github.com/sirupsen/logrus"

	logger_util "github.com/free5gc/util/logger"
)

var (
	Log      *logrus.Logger
	NfLog    *logrus.Entry
	MainLog  *logrus.Entry
	CfgLog   *logrus.Entry
	PfcpLog  *logrus.Entry
	BuffLog  *logrus.Entry
	PerioLog *logrus.Entry
	FwderLog *logrus.Entry
)

func init() {
	fieldsOrder := []string{
		logger_util.FieldNF,
		logger_util.FieldCategory,
		logger_util.FieldListenAddr,
		logger_util.FieldPFCPTxTransaction,
		logger_util.FieldPFCPRxTransaction,
		logger_util.FieldControlPlaneNodeID,
		logger_util.FieldControlPlaneSEID,
		logger_util.FieldUserPlaneSEID,
	}
	Log = logger_util.New(fieldsOrder)
	NfLog = Log.WithField(logger_util.FieldNF, "UPF")
	MainLog = NfLog.WithField(logger_util.FieldCategory, "Main")
	CfgLog = NfLog.WithField(logger_util.FieldCategory, "CFG")
	PfcpLog = NfLog.WithField(logger_util.FieldCategory, "PFCP")
	BuffLog = NfLog.WithField(logger_util.FieldCategory, "BUFF")
	PerioLog = NfLog.WithField(logger_util.FieldCategory, "Perio")
	FwderLog = NfLog.WithField(logger_util.FieldCategory, "FWD")
}
