package app

type NetworkFunction interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)
	Start(tlsKeyLogPath string)
	Terminate()
}
