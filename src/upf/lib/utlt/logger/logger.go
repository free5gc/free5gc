package main

import (
	"C"
	"github.com/sirupsen/logrus"
	"runtime"
)

var log *logrus.Logger
var UpfUtilLog *logrus.Entry

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
			// orgFilename, _ := os.Getwd()
			// log.Traceln("orgFilename", orgFilename)
			// repopath := fmt.Sprintf("%s", orgFilename)
			// repopath = strings.Replace(repopath, "/bin", "", 1)
			return "", ""
		},
	}

	UpfUtilLog = log.WithFields(logrus.Fields{"UPF": "Util"})
}

func SetLogLevel(level logrus.Level) {
	UpfUtilLog.Infoln("Set log level:", level)
	log.SetLevel(level)
}

//func SetReportCaller(reportCaller bool) {
//	UpfUtilLog.Infoln("Report caller:", reportCaller)
//	log.SetReportCaller(reportCaller)
//}

//export UpfUtilLog_SetLogLevel
func UpfUtilLog_SetLogLevel(levelString string) bool {
	level, err := logrus.ParseLevel(levelString)
	if err == nil {
		SetLogLevel(level)
		return true
	} else {
		UpfUtilLog.Errorln("Error: invalid log level: ", levelString)
		return false
	}
}

//export UpfUtilLog_Panicln
func UpfUtilLog_Panicln(comment string) {
	UpfUtilLog.Panicln(comment)
}

//export UpfUtilLog_Fatalln
func UpfUtilLog_Fatalln(comment string) {
	UpfUtilLog.Fatalln(comment)
}

//export UpfUtilLog_Errorln
func UpfUtilLog_Errorln(comment string) {
	UpfUtilLog.Errorln(comment)
}

//export UpfUtilLog_Warningln
func UpfUtilLog_Warningln(comment string) {
	UpfUtilLog.Warningln(comment)
}

//export UpfUtilLog_Infoln
func UpfUtilLog_Infoln(comment string) {
	UpfUtilLog.Infoln(comment)
}

//export UpfUtilLog_Debugln
func UpfUtilLog_Debugln(comment string) {
	UpfUtilLog.Debugln(comment)
}

//export UpfUtilLog_Traceln
func UpfUtilLog_Traceln(comment string) {
	UpfUtilLog.Traceln(comment)
}

func main() {}
