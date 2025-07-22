package forwarder

import (
	"github.com/wmnsk/go-pfcp/ie"

	"github.com/free5gc/go-upf/internal/report"
)

type Empty struct{}

func (Empty) Close() {
}

func (Empty) CreatePDR(uint64, *ie.IE) error {
	return nil
}

func (Empty) UpdatePDR(uint64, *ie.IE) error {
	return nil
}

func (Empty) RemovePDR(uint64, *ie.IE) error {
	return nil
}

func (Empty) CreateFAR(uint64, *ie.IE) error {
	return nil
}

func (Empty) UpdateFAR(uint64, *ie.IE) error {
	return nil
}

func (Empty) RemoveFAR(uint64, *ie.IE) error {
	return nil
}

func (Empty) CreateQER(uint64, *ie.IE) error {
	return nil
}

func (Empty) UpdateQER(uint64, *ie.IE) error {
	return nil
}

func (Empty) RemoveQER(uint64, *ie.IE) error {
	return nil
}

func (Empty) CreateURR(uint64, *ie.IE) error {
	return nil
}

func (Empty) UpdateURR(uint64, *ie.IE) ([]report.USAReport, error) {
	return nil, nil
}

func (Empty) RemoveURR(uint64, *ie.IE) ([]report.USAReport, error) {
	return nil, nil
}

func (Empty) CreateBAR(uint64, *ie.IE) error {
	return nil
}

func (Empty) UpdateBAR(uint64, *ie.IE) error {
	return nil
}

func (Empty) RemoveBAR(uint64, *ie.IE) error {
	return nil
}

func (Empty) QueryURR(uint64, uint32) ([]report.USAReport, error) {
	return nil, nil
}

func (Empty) HandleReport(report.Handler) {
}
