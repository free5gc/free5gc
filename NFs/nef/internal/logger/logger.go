package logger

import (
	logger_util "github.com/free5gc/util/logger"
	"github.com/sirupsen/logrus"
)

var (
	Log          *logrus.Logger
	NfLog        *logrus.Entry
	MainLog      *logrus.Entry
	InitLog      *logrus.Entry
	CfgLog       *logrus.Entry
	CmiLog       *logrus.Entry
	CtxLog       *logrus.Entry
	GinLog       *logrus.Entry
	SBILog       *logrus.Entry
	ConsumerLog  *logrus.Entry
	ProcessorLog *logrus.Entry
	TrafInfluLog *logrus.Entry
	PFDManageLog *logrus.Entry
	PFDFLog      *logrus.Entry
	OamLog       *logrus.Entry
)

const (
	FieldAFID       string = "AFID"
	FieldSubID      string = "SubID"
	FieldPfdTransID string = "PfdTRID"
)

func init() {
	fieldsOrder := []string{
		logger_util.FieldNF,
		logger_util.FieldCategory,
		FieldAFID,
		FieldSubID,
		FieldPfdTransID,
	}
	Log = logger_util.New(fieldsOrder)
	NfLog = Log.WithField(logger_util.FieldNF, "NEF")
	MainLog = NfLog.WithField(logger_util.FieldCategory, "Main")
	InitLog = NfLog.WithField(logger_util.FieldCategory, "Init")
	CfgLog = NfLog.WithField(logger_util.FieldCategory, "CFG")
	CtxLog = NfLog.WithField(logger_util.FieldCategory, "CTX")
	CmiLog = NfLog.WithField(logger_util.FieldCategory, "CMI")
	GinLog = NfLog.WithField(logger_util.FieldCategory, "GIN")
	SBILog = NfLog.WithField(logger_util.FieldCategory, "SBI")
	ConsumerLog = NfLog.WithField(logger_util.FieldCategory, "Consumer")
	ProcessorLog = NfLog.WithField(logger_util.FieldCategory, "Proc")
	TrafInfluLog = NfLog.WithField(logger_util.FieldCategory, "TraffInfl")
	PFDManageLog = NfLog.WithField(logger_util.FieldCategory, "PFDMng")
	PFDFLog = NfLog.WithField(logger_util.FieldCategory, "PFDF")
	OamLog = NfLog.WithField(logger_util.FieldCategory, "OAM")
}
