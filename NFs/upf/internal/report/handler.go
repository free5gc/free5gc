package report

type Handler interface {
	NotifySessReport(SessReport)
	PopBufPkt(uint64, uint16) ([]byte, bool)
}
