package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var AppLog *logrus.Entry
var InitLog *logrus.Entry
var Handlelog *logrus.Entry
var HttpLog *logrus.Entry
var UeauLog *logrus.Entry
var UecmLog *logrus.Entry
var SdmLog *logrus.Entry
var PpLop *logrus.Entry
var EeLog *logrus.Entry
var UtilLog *logrus.Entry

func init() {
	log = logrus.New()
	log.SetReportCaller(true)

	log.Formatter = &logrus.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			orgFilename, _ := os.Getwd()
			log.Traceln("orgFilename", orgFilename)
			repopath := orgFilename
			repopath = strings.Replace(repopath, "/bin", "", 1)
			filename := strings.Replace(f.File, repopath, "", -1)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	AppLog = log.WithFields(logrus.Fields{"UDM": "app"})
	InitLog = log.WithFields(logrus.Fields{"UDM": "init"})
	Handlelog = log.WithFields(logrus.Fields{"UDM": "Handler"})
	HttpLog = log.WithFields(logrus.Fields{"UDM": "HTTP"})
	UeauLog = log.WithFields(logrus.Fields{"UDM": "UEAU"})
	UecmLog = log.WithFields(logrus.Fields{"UDM": "UECM"})
	SdmLog = log.WithFields(logrus.Fields{"UDM": "SDM"})
	PpLop = log.WithFields(logrus.Fields{"UDM": "PP"})
	EeLog = log.WithFields(logrus.Fields{"UDM": "EE"})
	UtilLog = log.WithFields(logrus.Fields{"UDM": "Util"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
