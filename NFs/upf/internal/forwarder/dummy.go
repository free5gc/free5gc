package forwarder

import (
	"sync"

	"github.com/free5gc/go-upf/internal/logger"
	"github.com/free5gc/go-upf/internal/report"
	"github.com/free5gc/go-upf/pkg/factory"
	"github.com/wmnsk/go-pfcp/ie"
)

type Dummy struct {
}

func NewDummy(wg *sync.WaitGroup, ifList []factory.IfInfo) (*Dummy, error) {
	logger.MainLog.Infof("Dummy Forwarder initialized with interfaces: %v", ifList)
	return &Dummy{}, nil
}

func (d *Dummy) Close() {

}

func (d *Dummy) CreatePDR(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) UpdatePDR(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) RemovePDR(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) CreateFAR(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) UpdateFAR(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) RemoveFAR(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) CreateQER(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) UpdateQER(uint64, *ie.IE) error {
	return nil
}
func (d *Dummy) RemoveQER(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) CreateURR(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) UpdateURR(uint64, *ie.IE) ([]report.USAReport, error) {
	return []report.USAReport{}, nil
}

func (d *Dummy) RemoveURR(uint64, *ie.IE) ([]report.USAReport, error) {
	return []report.USAReport{}, nil
}

func (d *Dummy) QueryURR(uint64, uint32) ([]report.USAReport, error) {
	return []report.USAReport{}, nil
}

func (d *Dummy) CreateBAR(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) UpdateBAR(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) RemoveBAR(uint64, *ie.IE) error {
	return nil
}

func (d *Dummy) HandleReport(report.Handler) {

}
