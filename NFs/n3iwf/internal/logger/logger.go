package logger

import (
	"github.com/sirupsen/logrus"

	logger_util "github.com/free5gc/util/logger"
)

var (
	Log      *logrus.Logger
	N3iwfLog *logrus.Entry
	MainLog  *logrus.Entry
	InitLog  *logrus.Entry
	CfgLog   *logrus.Entry
	CtxLog   *logrus.Entry
	GinLog   *logrus.Entry
	NasLog   *logrus.Entry
	NgapLog  *logrus.Entry
	IKELog   *logrus.Entry
	GTPLog   *logrus.Entry
	NWuCPLog *logrus.Entry
	NWuUPLog *logrus.Entry
	RelayLog *logrus.Entry
	UtilLog  *logrus.Entry
	GmmLog   *logrus.Entry
)

func UpdateN3iwfLog() {
	N3iwfLog = Log.WithField(logger_util.FieldNF, "N3IWF")
	// update logs created from N3iwfLog
	MainLog = N3iwfLog.WithField(logger_util.FieldCategory, "Main")
	InitLog = N3iwfLog.WithField(logger_util.FieldCategory, "Init")
	CfgLog = N3iwfLog.WithField(logger_util.FieldCategory, "CFG")
	CtxLog = N3iwfLog.WithField(logger_util.FieldCategory, "CTX")
	GinLog = N3iwfLog.WithField(logger_util.FieldCategory, "GIN")
	NgapLog = N3iwfLog.WithField(logger_util.FieldCategory, "NGAP")
	IKELog = N3iwfLog.WithField(logger_util.FieldCategory, "IKE")
	GTPLog = N3iwfLog.WithField(logger_util.FieldCategory, "GTP")
	NWuCPLog = N3iwfLog.WithField(logger_util.FieldCategory, "NWuCP")
	NWuUPLog = N3iwfLog.WithField(logger_util.FieldCategory, "NWuUP")
	RelayLog = N3iwfLog.WithField(logger_util.FieldCategory, "Relay")
	UtilLog = N3iwfLog.WithField(logger_util.FieldCategory, "Util")
}

func init() {
	fieldsOrder := []string{
		logger_util.FieldNF,
		logger_util.FieldCategory,
	}
	Log = logger_util.New(fieldsOrder)
	UpdateN3iwfLog()
}
