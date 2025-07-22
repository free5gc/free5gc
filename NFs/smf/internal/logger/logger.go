package logger

import (
	"github.com/sirupsen/logrus"

	logger_util "github.com/free5gc/util/logger"
)

const (
	FieldSupi         = "supi"
	FieldPDUSessionID = "pdu_session_id"
)

var (
	Log         *logrus.Logger
	NfLog       *logrus.Entry
	MainLog     *logrus.Entry
	InitLog     *logrus.Entry
	CfgLog      *logrus.Entry
	CtxLog      *logrus.Entry
	GinLog      *logrus.Entry
	SBILog      *logrus.Entry
	ConsumerLog *logrus.Entry
	GsmLog      *logrus.Entry
	PfcpLog     *logrus.Entry
	PduSessLog  *logrus.Entry
	ChargingLog *logrus.Entry
	UtilLog     *logrus.Entry
)

func init() {
	fieldsOrder := []string{
		logger_util.FieldNF,
		logger_util.FieldCategory,
	}

	Log = logger_util.New(fieldsOrder)
	NfLog = Log.WithField(logger_util.FieldNF, "SMF")
	MainLog = NfLog.WithField(logger_util.FieldCategory, "Main")
	InitLog = NfLog.WithField(logger_util.FieldCategory, "Init")
	CfgLog = NfLog.WithField(logger_util.FieldCategory, "CFG")
	CtxLog = NfLog.WithField(logger_util.FieldCategory, "CTX")
	GinLog = NfLog.WithField(logger_util.FieldCategory, "GIN")
	SBILog = NfLog.WithField(logger_util.FieldCategory, "SBI")
	ConsumerLog = NfLog.WithField(logger_util.FieldCategory, "Consumer")
	GsmLog = NfLog.WithField(logger_util.FieldCategory, "GSM")
	PfcpLog = NfLog.WithField(logger_util.FieldCategory, "PFCP")
	PduSessLog = NfLog.WithField(logger_util.FieldCategory, "PduSess")
	ChargingLog = NfLog.WithField(logger_util.FieldCategory, "Charging")
	UtilLog = NfLog.WithField(logger_util.FieldCategory, "Util")
}
