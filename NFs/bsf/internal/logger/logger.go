/*
 * BSF Logger
 *
 * BSF Logging System
 * © 2025, 3GPP Organizational Partners (ARIB, ATIS, CCSA, ETSI, TSDSI, TTA, TTC).
 * All rights reserved.
 */

package logger

import (
	"github.com/sirupsen/logrus"

	logger_util "github.com/free5gc/util/logger"
)

var (
	Log      *logrus.Logger
	NfLog    *logrus.Entry
	MainLog  *logrus.Entry
	InitLog  *logrus.Entry
	CfgLog   *logrus.Entry
	SbiLog   *logrus.Entry
	CtxLog   *logrus.Entry
	GinLog   *logrus.Entry
	ProcLog  *logrus.Entry
	ConsLog  *logrus.Entry
	MongoLog *logrus.Entry
)

func init() {
	fieldsOrder := []string{
		logger_util.FieldNF,
		logger_util.FieldCategory,
	}
	Log = logger_util.New(fieldsOrder)
	NfLog = Log.WithField(logger_util.FieldNF, "BSF")
	MainLog = NfLog.WithField(logger_util.FieldCategory, "Main")
	InitLog = NfLog.WithField(logger_util.FieldCategory, "Init")
	CfgLog = NfLog.WithField(logger_util.FieldCategory, "CFG")
	SbiLog = NfLog.WithField(logger_util.FieldCategory, "SBI")
	CtxLog = NfLog.WithField(logger_util.FieldCategory, "CTX")
	GinLog = NfLog.WithField(logger_util.FieldCategory, "GIN")
	ProcLog = NfLog.WithField(logger_util.FieldCategory, "Proc")
	ConsLog = NfLog.WithField(logger_util.FieldCategory, "Cons")
	MongoLog = NfLog.WithField(logger_util.FieldCategory, "MongoDB")
}
