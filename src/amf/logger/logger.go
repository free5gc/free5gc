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
var ContextLog *logrus.Entry
var NgapLog *logrus.Entry
var HandlerLog *logrus.Entry
var HttpLog *logrus.Entry
var GmmLog *logrus.Entry
var MtLog *logrus.Entry
var ProducerLog *logrus.Entry
var LocationLog *logrus.Entry
var CommLog *logrus.Entry
var CallbackLog *logrus.Entry
var UtilLog *logrus.Entry
var NasLog *logrus.Entry
var ConsumerLog *logrus.Entry

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

	AppLog = log.WithFields(logrus.Fields{"AMF": "app"})
	InitLog = log.WithFields(logrus.Fields{"AMF": "init"})
	ContextLog = log.WithFields(logrus.Fields{"AMF": "Context"})
	NgapLog = log.WithFields(logrus.Fields{"AMF": "NGAP"})
	HandlerLog = log.WithFields(logrus.Fields{"AMF": "Handler"})
	HttpLog = log.WithFields(logrus.Fields{"AMF": "HTTP"})
	GmmLog = log.WithFields(logrus.Fields{"AMF": "Gmm"})
	MtLog = log.WithFields(logrus.Fields{"AMF": "MT"})
	ProducerLog = log.WithFields(logrus.Fields{"AMF": "Producer"})
	LocationLog = log.WithFields(logrus.Fields{"AMF": "LocInfo"})
	CommLog = log.WithFields(logrus.Fields{"AMF": "Comm"})
	CallbackLog = log.WithFields(logrus.Fields{"AMF": "Callback"})
	UtilLog = log.WithFields(logrus.Fields{"AMF": "Util"})
	NasLog = log.WithFields(logrus.Fields{"AMF": "NAS"})
	ConsumerLog = log.WithFields(logrus.Fields{"AMF": "Consumer"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
