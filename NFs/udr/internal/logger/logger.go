package logger

import (
	"github.com/sirupsen/logrus"

	logger_util "github.com/free5gc/util/logger"
)

var (
	Log         *logrus.Logger
	NfLog       *logrus.Entry
	MainLog     *logrus.Entry
	InitLog     *logrus.Entry
	CfgLog      *logrus.Entry
	CtxLog      *logrus.Entry
	DataRepoLog *logrus.Entry
	UtilLog     *logrus.Entry
	HttpLog     *logrus.Entry
	ConsumerLog *logrus.Entry
	GinLog      *logrus.Entry
	ProcLog     *logrus.Entry
	SBILog      *logrus.Entry
	DbLog       *logrus.Entry
)

func init() {
	fieldsOrder := []string{
		logger_util.FieldNF,
		logger_util.FieldCategory,
	}

	Log = logger_util.New(fieldsOrder)
	NfLog = Log.WithField(logger_util.FieldNF, "UDR")
	MainLog = NfLog.WithField(logger_util.FieldCategory, "Main")
	InitLog = NfLog.WithField(logger_util.FieldCategory, "Init")
	CfgLog = NfLog.WithField(logger_util.FieldCategory, "CFG")
	CtxLog = NfLog.WithField(logger_util.FieldCategory, "CTX")
	GinLog = NfLog.WithField(logger_util.FieldCategory, "GIN")
	ConsumerLog = NfLog.WithField(logger_util.FieldCategory, "Consumer")
	DataRepoLog = NfLog.WithField(logger_util.FieldCategory, "DataRepo")
	ProcLog = NfLog.WithField(logger_util.FieldCategory, "Proc")
	HttpLog = NfLog.WithField(logger_util.FieldCategory, "HTTP")
	UtilLog = NfLog.WithField(logger_util.FieldCategory, "Util")
	SBILog = NfLog.WithField(logger_util.FieldCategory, "SBI")
	DbLog = NfLog.WithField(logger_util.FieldCategory, "DB")
}
