//go:binary-only-package

package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var MongoDBLog *logrus.Entry

func init() {}

func SetLogLevel(level logrus.Level) {}

func SetReportCaller(bool bool) {}
