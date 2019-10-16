//go:binary-only-package

package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

var log *logrus.Logger
var NgapLog *logrus.Entry

func init() {}

func SetLogLevel(level logrus.Level) {}

func SetReportCaller(bool bool) {}
