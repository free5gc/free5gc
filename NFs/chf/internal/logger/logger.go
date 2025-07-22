package logger

import (
	golog "github.com/fclairamb/go-log"
	adapter "github.com/fclairamb/go-log/logrus"
	"github.com/sirupsen/logrus"

	logger_util "github.com/free5gc/util/logger"
)

var (
	Log                 *logrus.Logger
	NfLog               *logrus.Entry
	MainLog             *logrus.Entry
	InitLog             *logrus.Entry
	CfgLog              *logrus.Entry
	CtxLog              *logrus.Entry
	SBILog              *logrus.Entry
	ConsumerLog         *logrus.Entry
	ProcLog             *logrus.Entry
	GinLog              *logrus.Entry
	ChargingdataPostLog *logrus.Entry
	NotifyEventLog      *logrus.Entry
	RechargingLog       *logrus.Entry
	RatingLog           *logrus.Entry
	AcctLog             *logrus.Entry
	CgfLog              *logrus.Entry
	UtilLog             *logrus.Entry
	FtpServerLog        golog.Logger
)

func init() {
	fieldsOrder := []string{
		logger_util.FieldNF,
		logger_util.FieldCategory,
	}

	Log = logger_util.New(fieldsOrder)

	NfLog = Log.WithField(logger_util.FieldNF, "CHF")
	MainLog = NfLog.WithField(logger_util.FieldCategory, "Main")
	InitLog = NfLog.WithField(logger_util.FieldCategory, "Init")
	CfgLog = NfLog.WithField(logger_util.FieldCategory, "CFG")
	CtxLog = NfLog.WithField(logger_util.FieldCategory, "CTX")
	SBILog = NfLog.WithField(logger_util.FieldCategory, "SBI")
	ConsumerLog = NfLog.WithField(logger_util.FieldCategory, "Consumer")
	ProcLog = NfLog.WithField(logger_util.FieldCategory, "Proc")
	GinLog = NfLog.WithField(logger_util.FieldCategory, "GIN")
	ChargingdataPostLog = NfLog.WithField(logger_util.FieldCategory, "ChargingPost")
	NotifyEventLog = NfLog.WithField(logger_util.FieldCategory, "NotifyEvent")
	RechargingLog = NfLog.WithField(logger_util.FieldCategory, "Recharge")
	CgfLog = NfLog.WithField(logger_util.FieldCategory, "CGF")
	RatingLog = NfLog.WithField(logger_util.FieldCategory, "Rating")
	AcctLog = NfLog.WithField(logger_util.FieldCategory, "Acct")
	UtilLog = NfLog.WithField(logger_util.FieldCategory, "Util")
	FtpServerLog = adapter.NewWrap(CgfLog.Logger).With("component", "CHF", "category", "FTP")
}
