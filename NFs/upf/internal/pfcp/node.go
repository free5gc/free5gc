package pfcp

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wmnsk/go-pfcp/ie"

	"github.com/free5gc/go-upf/internal/forwarder"
	"github.com/free5gc/go-upf/internal/report"
	logger_util "github.com/free5gc/util/logger"
)

const (
	BUFFQ_LEN = 512
)

type PDRInfo struct {
	RelatedURRIDs map[uint32]struct{}
}

type URRInfo struct {
	removed bool
	SEQN    uint32
	report.MeasureMethod
	report.MeasureInformation
	refPdrNum uint16
}

type Sess struct {
	rnode    *RemoteNode
	LocalID  uint64
	RemoteID uint64
	PDRIDs   map[uint16]*PDRInfo    // key: PDR_ID
	FARIDs   map[uint32]struct{}    // key: FAR_ID
	QERIDs   map[uint32]struct{}    // key: QER_ID
	URRIDs   map[uint32]*URRInfo    // key: URR_ID
	BARIDs   map[uint8]struct{}     // key: BAR_ID
	q        map[uint16]chan []byte // key: PDR_ID
	qlen     int
	log      *logrus.Entry
}

func (s *Sess) Close() []report.USAReport {
	for id := range s.FARIDs {
		i := ie.NewRemoveFAR(ie.NewFARID(id))
		err := s.RemoveFAR(i)
		if err != nil {
			s.log.Errorf("Remove FAR err: %+v", err)
		}
	}
	for id := range s.QERIDs {
		i := ie.NewRemoveQER(ie.NewQERID(id))
		err := s.RemoveQER(i)
		if err != nil {
			s.log.Errorf("Remove QER err: %+v", err)
		}
	}

	var usars []report.USAReport
	for id := range s.URRIDs {
		i := ie.NewRemoveURR(ie.NewURRID(id))
		rs, err := s.RemoveURR(i)
		if err != nil {
			s.log.Errorf("Remove URR err: %+v", err)
			continue
		}
		if rs != nil {
			usars = append(usars, rs...)
		}
	}
	for id := range s.BARIDs {
		i := ie.NewRemoveBAR(ie.NewBARID(id))
		err := s.RemoveBAR(i)
		if err != nil {
			s.log.Errorf("Remove BAR err: %+v", err)
		}
	}
	for id := range s.PDRIDs {
		i := ie.NewRemovePDR(ie.NewPDRID(id))
		rs, err := s.RemovePDR(i)
		if err != nil {
			s.log.Errorf("remove PDR err: %+v", err)
		}
		if rs != nil {
			usars = append(usars, rs...)
		}
	}
	for _, q := range s.q {
		close(q)
	}
	return usars
}

func (s *Sess) CreatePDR(req *ie.IE) error {
	ies, err := req.CreatePDR()
	if err != nil {
		return err
	}

	var pdrid uint16
	urrids := make(map[uint32]struct{})
	for _, i := range ies {
		switch i.Type {
		case ie.PDRID:
			v, err1 := i.PDRID()
			if err1 != nil {
				break
			}
			pdrid = v
		case ie.URRID:
			v, err1 := i.URRID()
			if err1 != nil {
				break
			}
			urrids[v] = struct{}{}
			urrInfo, ok := s.URRIDs[v]
			if ok {
				urrInfo.refPdrNum++
			}
		}
	}

	s.PDRIDs[pdrid] = &PDRInfo{
		RelatedURRIDs: urrids,
	}

	err = s.rnode.driver.CreatePDR(s.LocalID, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sess) diassociateURR(urrid uint32) []report.USAReport {
	urrInfo, ok := s.URRIDs[urrid]
	if !ok {
		return nil
	}

	if urrInfo.refPdrNum > 0 {
		urrInfo.refPdrNum--
		if urrInfo.refPdrNum == 0 {
			// indicates usage report being reported for a URR due to dissociated from the last PDR
			usars, err := s.rnode.driver.QueryURR(s.LocalID, urrid)
			if err != nil {
				return nil
			}
			for i := range usars {
				usars[i].USARTrigger.Flags |= report.USAR_TRIG_TERMR
			}
			return usars
		}
	} else {
		s.log.Warnf("diassociateURR: wrong refPdrNum(%d)", urrInfo.refPdrNum)
	}
	return nil
}

func (s *Sess) UpdatePDR(req *ie.IE) ([]report.USAReport, error) {
	ies, err := req.UpdatePDR()
	if err != nil {
		return nil, err
	}

	var pdrid uint16
	newUrrids := make(map[uint32]struct{})
	for _, i := range ies {
		switch i.Type {
		case ie.PDRID:
			v, err1 := i.PDRID()
			if err1 != nil {
				break
			}
			pdrid = v
		case ie.URRID:
			v, err1 := i.URRID()
			if err1 != nil {
				break
			}
			newUrrids[v] = struct{}{}
		}
	}

	pdrInfo, ok := s.PDRIDs[pdrid]
	if !ok {
		return nil, errors.Errorf("UpdatePDR: PDR(%#x) not found", pdrid)
	}

	err = s.rnode.driver.UpdatePDR(s.LocalID, req)
	if err != nil {
		return nil, err
	}

	var usars []report.USAReport
	for urrid := range pdrInfo.RelatedURRIDs {
		_, ok = newUrrids[urrid]
		if !ok {
			usar := s.diassociateURR(urrid)
			if len(usar) > 0 {
				usars = append(usars, usar...)
			}
		}
	}
	pdrInfo.RelatedURRIDs = newUrrids

	return usars, err
}

func (s *Sess) RemovePDR(req *ie.IE) ([]report.USAReport, error) {
	pdrid, err := req.PDRID()
	if err != nil {
		return nil, err
	}

	pdrInfo, ok := s.PDRIDs[pdrid]
	if !ok {
		return nil, errors.Errorf("RemovePDR: PDR(%#x) not found", pdrid)
	}

	err = s.rnode.driver.RemovePDR(s.LocalID, req)
	if err != nil {
		return nil, err
	}

	var usars []report.USAReport
	for urrid := range pdrInfo.RelatedURRIDs {
		usar := s.diassociateURR(urrid)
		if len(usar) > 0 {
			usars = append(usars, usar...)
		}
	}
	delete(s.PDRIDs, pdrid)
	return usars, nil
}

func (s *Sess) CreateFAR(req *ie.IE) error {
	id, err := req.FARID()
	if err != nil {
		return err
	}
	s.FARIDs[id] = struct{}{}

	err = s.rnode.driver.CreateFAR(s.LocalID, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sess) UpdateFAR(req *ie.IE) error {
	id, err := req.FARID()
	if err != nil {
		return err
	}

	_, ok := s.FARIDs[id]
	if !ok {
		return errors.Errorf("UpdateFAR: FAR(%#x) not found", id)
	}
	return s.rnode.driver.UpdateFAR(s.LocalID, req)
}

func (s *Sess) RemoveFAR(req *ie.IE) error {
	id, err := req.FARID()
	if err != nil {
		return err
	}

	_, ok := s.FARIDs[id]
	if !ok {
		return errors.Errorf("RemoveFAR: FAR(%#x) not found", id)
	}

	err = s.rnode.driver.RemoveFAR(s.LocalID, req)
	if err != nil {
		return err
	}

	delete(s.FARIDs, id)
	return nil
}

func (s *Sess) CreateQER(req *ie.IE) error {
	id, err := req.QERID()
	if err != nil {
		return err
	}
	s.QERIDs[id] = struct{}{}

	err = s.rnode.driver.CreateQER(s.LocalID, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sess) UpdateQER(req *ie.IE) error {
	id, err := req.QERID()
	if err != nil {
		return err
	}

	_, ok := s.QERIDs[id]
	if !ok {
		return errors.Errorf("UpdateQER: QER(%#x) not found", id)
	}
	return s.rnode.driver.UpdateQER(s.LocalID, req)
}

func (s *Sess) RemoveQER(req *ie.IE) error {
	id, err := req.QERID()
	if err != nil {
		return err
	}

	_, ok := s.QERIDs[id]
	if !ok {
		return errors.Errorf("RemoveQER: QER(%#x) not found", id)
	}

	err = s.rnode.driver.RemoveQER(s.LocalID, req)
	if err != nil {
		return err
	}

	delete(s.QERIDs, id)
	return nil
}

func (s *Sess) CreateURR(req *ie.IE) error {
	id, err := req.URRID()
	if err != nil {
		return err
	}

	mInfo := &ie.IE{}
	for _, x := range req.ChildIEs {
		if x.Type == ie.MeasurementInformation {
			mInfo = x
			break
		}
	}
	s.URRIDs[id] = &URRInfo{
		MeasureMethod: report.MeasureMethod{
			DURAT: req.HasDURAT(),
			VOLUM: req.HasVOLUM(),
			EVENT: req.HasEVENT(),
		},
		MeasureInformation: report.MeasureInformation{
			MBQE: mInfo.HasMBQE(),
			INAM: mInfo.HasINAM(),
			RADI: mInfo.HasRADI(),
			ISTM: mInfo.HasISTM(),
			MNOP: mInfo.HasMNOP(),
		},
	}

	err = s.rnode.driver.CreateURR(s.LocalID, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sess) UpdateURR(req *ie.IE) ([]report.USAReport, error) {
	id, err := req.URRID()
	if err != nil {
		return nil, err
	}

	urrInfo, ok := s.URRIDs[id]
	if !ok {
		return nil, errors.Errorf("UpdateURR: URR[%#x] not found", id)
	}
	for _, x := range req.ChildIEs {
		switch x.Type {
		case ie.MeasurementMethod:
			urrInfo.DURAT = x.HasDURAT()
			urrInfo.VOLUM = x.HasVOLUM()
			urrInfo.EVENT = x.HasEVENT()
		case ie.MeasurementInformation:
			urrInfo.MBQE = x.HasMBQE()
			urrInfo.INAM = x.HasINAM()
			urrInfo.RADI = x.HasRADI()
			urrInfo.ISTM = x.HasISTM()
			urrInfo.MNOP = x.HasMNOP()
		}
	}

	usars, err := s.rnode.driver.UpdateURR(s.LocalID, req)
	if err != nil {
		return nil, err
	}
	return usars, nil
}

func (s *Sess) RemoveURR(req *ie.IE) ([]report.USAReport, error) {
	id, err := req.URRID()
	if err != nil {
		return nil, err
	}

	info, ok := s.URRIDs[id]
	if !ok {
		return nil, errors.Errorf("RemoveURR: URR[%#x] not found", id)
	}
	info.removed = true // remove URRInfo later

	usars, err := s.rnode.driver.RemoveURR(s.LocalID, req)
	if err != nil {
		return nil, err
	}

	// indicates usage report being reported for a URR due to the removal of the URR
	for i := range usars {
		usars[i].USARTrigger.Flags |= report.USAR_TRIG_TERMR
	}
	return usars, nil
}

func (s *Sess) QueryURR(req *ie.IE) ([]report.USAReport, error) {
	id, err := req.URRID()
	if err != nil {
		return nil, err
	}

	_, ok := s.URRIDs[id]
	if !ok {
		return nil, errors.Errorf("QueryURR: URR[%#x] not found", id)
	}

	usars, err := s.rnode.driver.QueryURR(s.LocalID, id)
	if err != nil {
		return nil, err
	}

	// indicates an immediate report reported on CP function demand
	for i := range usars {
		usars[i].USARTrigger.Flags |= report.USAR_TRIG_IMMER
	}
	return usars, nil
}

func (s *Sess) CreateBAR(req *ie.IE) error {
	id, err := req.BARID()
	if err != nil {
		return err
	}
	s.BARIDs[id] = struct{}{}

	err = s.rnode.driver.CreateBAR(s.LocalID, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sess) UpdateBAR(req *ie.IE) error {
	id, err := req.BARID()
	if err != nil {
		return err
	}

	_, ok := s.BARIDs[id]
	if !ok {
		return errors.Errorf("UpdateBAR: BAR(%#x) not found", id)
	}
	return s.rnode.driver.UpdateBAR(s.LocalID, req)
}

func (s *Sess) RemoveBAR(req *ie.IE) error {
	id, err := req.BARID()
	if err != nil {
		return err
	}

	_, ok := s.BARIDs[id]
	if !ok {
		return errors.Errorf("RemoveBAR: BAR(%#x) not found", id)
	}

	err = s.rnode.driver.RemoveBAR(s.LocalID, req)
	if err != nil {
		return err
	}

	delete(s.BARIDs, id)
	return nil
}

func (s *Sess) Push(pdrid uint16, p []byte) {
	pkt := make([]byte, len(p))
	copy(pkt, p)
	q, ok := s.q[pdrid]
	if !ok {
		s.q[pdrid] = make(chan []byte, s.qlen)
		q = s.q[pdrid]
	}

	select {
	case q <- pkt:
		s.log.Debugf("Push bufPkt to q[%d](len:%d)", pdrid, len(q))
	default:
		s.log.Debugf("q[%d](len:%d) is full, drop it", pdrid, len(q))
	}
}

func (s *Sess) Len(pdrid uint16) int {
	q, ok := s.q[pdrid]
	if !ok {
		return 0
	}
	return len(q)
}

func (s *Sess) Pop(pdrid uint16) ([]byte, bool) {
	q, ok := s.q[pdrid]
	if !ok {
		return nil, ok
	}
	select {
	case pkt := <-q:
		s.log.Debugf("Pop bufPkt from q[%d](len:%d)", pdrid, len(q))
		return pkt, true
	default:
		return nil, false
	}
}

func (s *Sess) URRSeq(urrid uint32) uint32 {
	info, ok := s.URRIDs[urrid]
	if !ok {
		return 0
	}
	seq := info.SEQN
	info.SEQN++
	return seq
}

type RemoteNode struct {
	ID     string
	addr   net.Addr
	local  *LocalNode
	sess   map[uint64]struct{} // key: Local SEID
	driver forwarder.Driver
	log    *logrus.Entry
}

func NewRemoteNode(
	id string,
	addr net.Addr,
	local *LocalNode,
	driver forwarder.Driver,
	log *logrus.Entry,
) *RemoteNode {
	n := new(RemoteNode)
	n.ID = id
	n.addr = addr
	n.local = local
	n.sess = make(map[uint64]struct{})
	n.driver = driver
	n.log = log
	return n
}

func (n *RemoteNode) Reset() {
	for id := range n.sess {
		n.DeleteSess(id)
	}
	n.sess = make(map[uint64]struct{})
}

func (n *RemoteNode) Sess(lSeid uint64) (*Sess, error) {
	_, ok := n.sess[lSeid]
	if !ok {
		return nil, errors.Errorf("Sess: sess not found (lSeid:%#x)", lSeid)
	}
	return n.local.Sess(lSeid)
}

func (n *RemoteNode) NewSess(rSeid uint64) *Sess {
	s := n.local.NewSess(rSeid, BUFFQ_LEN)
	n.sess[s.LocalID] = struct{}{}
	s.rnode = n
	s.log = n.log.WithFields(
		logrus.Fields{
			logger_util.FieldUserPlaneSEID:    fmt.Sprintf("%#x", s.LocalID),
			logger_util.FieldControlPlaneSEID: fmt.Sprintf("%#x", rSeid),
		})
	s.log.Infoln("New session")
	return s
}

func (n *RemoteNode) DeleteSess(lSeid uint64) []report.USAReport {
	_, ok := n.sess[lSeid]
	if !ok {
		return nil
	}
	delete(n.sess, lSeid)
	usars, err := n.local.DeleteSess(lSeid)
	if err != nil {
		n.log.Warnln(err)
		return nil
	}
	return usars
}

type LocalNode struct {
	sess []*Sess
	free []uint64
}

func (n *LocalNode) Reset() {
	for _, sess := range n.sess {
		if sess != nil {
			sess.Close()
		}
	}
	n.sess = []*Sess{}
	n.free = []uint64{}
}

func (n *LocalNode) Sess(lSeid uint64) (*Sess, error) {
	if lSeid == 0 {
		return nil, errors.New("Sess: invalid lSeid:0")
	}
	i := int(lSeid) - 1
	if i >= len(n.sess) {
		return nil, errors.Errorf("Sess: sess not found (lSeid:%#x)", lSeid)
	}
	sess := n.sess[i]
	if sess == nil {
		return nil, errors.Errorf("Sess: sess not found (lSeid:%#x)", lSeid)
	}
	return sess, nil
}

func (n *LocalNode) RemoteSess(rSeid uint64, addr net.Addr) (*Sess, error) {
	for _, s := range n.sess {
		if s.RemoteID == rSeid && s.rnode.addr.String() == addr.String() {
			return s, nil
		}
	}
	return nil, errors.Errorf("RemoteSess: invalid rSeid:%#x, addr:%s ", rSeid, addr)
}

func (n *LocalNode) NewSess(rSeid uint64, qlen int) *Sess {
	s := &Sess{
		RemoteID: rSeid,
		PDRIDs:   make(map[uint16]*PDRInfo),
		FARIDs:   make(map[uint32]struct{}),
		QERIDs:   make(map[uint32]struct{}),
		URRIDs:   make(map[uint32]*URRInfo),
		BARIDs:   make(map[uint8]struct{}),
		q:        make(map[uint16]chan []byte),
		qlen:     qlen,
	}
	last := len(n.free) - 1
	if last >= 0 {
		s.LocalID = n.free[last]
		n.free = n.free[:last]
		n.sess[s.LocalID-1] = s
	} else {
		n.sess = append(n.sess, s)
		s.LocalID = uint64(len(n.sess))
	}
	return s
}

func (n *LocalNode) DeleteSess(lSeid uint64) ([]report.USAReport, error) {
	if lSeid == 0 {
		return nil, errors.New("DeleteSess: invalid lSeid:0")
	}
	i := int(lSeid) - 1
	if i >= len(n.sess) {
		return nil, errors.Errorf("DeleteSess: sess not found (lSeid:%#x)", lSeid)
	}
	if n.sess[i] == nil {
		return nil, errors.Errorf("DeleteSess: sess not found (lSeid:%#x)", lSeid)
	}
	n.sess[i].log.Infoln("sess deleted")
	usars := n.sess[i].Close()
	n.sess[i] = nil
	n.free = append(n.free, lSeid)
	return usars, nil
}
