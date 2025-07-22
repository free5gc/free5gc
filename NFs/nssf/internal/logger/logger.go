package logger

import (
	"github.com/sirupsen/logrus"

	logger_util "github.com/free5gc/util/logger"
)

var (
	Log           *logrus.Logger
	NfLog         *logrus.Entry
	MainLog       *logrus.Entry
	InitLog       *logrus.Entry
	CfgLog        *logrus.Entry
	CtxLog        *logrus.Entry
	GinLog        *logrus.Entry
	SBILog        *logrus.Entry
	ConsumerLog   *logrus.Entry
	ProcLog       *logrus.Entry
	NsselLog      *logrus.Entry
	NssaiavailLog *logrus.Entry
	UtilLog       *logrus.Entry
)

func init() {
	fieldsOrder := []string{
		logger_util.FieldNF,
		logger_util.FieldCategory,
	}
	Log = logger_util.New(fieldsOrder)
	NfLog = Log.WithField(logger_util.FieldNF, "NSSF")
	MainLog = NfLog.WithField(logger_util.FieldCategory, "Main")
	InitLog = NfLog.WithField(logger_util.FieldCategory, "Init")
	CfgLog = NfLog.WithField(logger_util.FieldCategory, "CFG")
	CtxLog = NfLog.WithField(logger_util.FieldCategory, "CTX")
	GinLog = NfLog.WithField(logger_util.FieldCategory, "GIN")
	SBILog = NfLog.WithField(logger_util.FieldCategory, "SBI")
	ConsumerLog = NfLog.WithField(logger_util.FieldCategory, "Consumer")
	ProcLog = NfLog.WithField(logger_util.FieldCategory, "Proc")
	NsselLog = NfLog.WithField(logger_util.FieldCategory, "NsSel")
	NssaiavailLog = NfLog.WithField(logger_util.FieldCategory, "NssaiAvail")
	UtilLog = NfLog.WithField(logger_util.FieldCategory, "Util")
}
