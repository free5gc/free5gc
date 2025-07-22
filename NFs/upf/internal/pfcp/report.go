package pfcp

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-pfcp/ie"
	"github.com/wmnsk/go-pfcp/message"

	"github.com/free5gc/go-upf/internal/report"
	"github.com/free5gc/go-upf/pkg/factory"
)

func (s *PfcpServer) ServeReport(sr *report.SessReport) {
	s.log.Debugf("ServeReport: SEID(%#x)", sr.SEID)
	sess, err := s.lnode.Sess(sr.SEID)
	if err != nil {
		s.log.Errorln(err)
		return
	}

	addr := fmt.Sprintf("%s:%d", sess.rnode.ID, factory.UpfPfcpDefaultPort)
	laddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return
	}

	var usars []report.USAReport
	for _, rpt := range sr.Reports {
		switch r := rpt.(type) {
		case report.DLDReport:
			s.log.Debugf("ServeReport: SEID(%#x), type(%s)", sr.SEID, r.Type())
			if r.Action&report.APPLY_ACT_BUFF != 0 && len(r.BufPkt) > 0 {
				sess.Push(r.PDRID, r.BufPkt)
			}
			if r.Action&report.APPLY_ACT_NOCP == 0 {
				return
			}
			err := s.serveDLDReport(laddr, sr.SEID, r.PDRID)
			if err != nil {
				s.log.Errorln(err)
			}
		case report.USAReport:
			s.log.Debugf("ServeReport: SEID(%#x), type(%s)", sr.SEID, r.Type())
			usars = append(usars, r)
		default:
			s.log.Warnf("Unsupported Report: SEID(%#x), type(%d)", sr.SEID, rpt.Type())
		}
	}

	if len(usars) > 0 {
		err := s.serveUSAReport(laddr, sr.SEID, usars)
		if err != nil {
			s.log.Errorln(err)
		}
	}
}

func (s *PfcpServer) serveDLDReport(addr net.Addr, lSeid uint64, pdrid uint16) error {
	s.log.Infoln("serveDLDReport")

	sess, err := s.lnode.Sess(lSeid)
	if err != nil {
		return errors.Wrap(err, "serveDLDReport")
	}

	req := message.NewSessionReportRequest(
		0,
		0,
		sess.RemoteID,
		0,
		0,
		ie.NewReportType(0, 0, 0, 1),
		ie.NewDownlinkDataReport(
			ie.NewPDRID(pdrid),
			/*
				ie.NewDownlinkDataServiceInformation(
					true,
					true,
					ppi,
					qfi,
				),
			*/
		),
	)

	err = s.sendReqTo(req, addr)
	return errors.Wrap(err, "serveDLDReport")
}

func (s *PfcpServer) serveUSAReport(addr net.Addr, lSeid uint64, usars []report.USAReport) error {
	s.log.Infoln("serveUSAReport")

	sess, err := s.lnode.Sess(lSeid)
	if err != nil {
		return errors.Wrap(err, "serveUSAReport")
	}

	req := message.NewSessionReportRequest(
		0,
		0,
		sess.RemoteID,
		0,
		0,
		ie.NewReportType(0, 0, 1, 0),
	)
	for _, r := range usars {
		urrInfo, ok := sess.URRIDs[r.URRID]
		if !ok {
			sess.log.Warnf("serveUSAReport: URRInfo[%#x] not found", r.URRID)
			continue
		}
		r.URSEQN = sess.URRSeq(r.URRID)
		req.UsageReport = append(req.UsageReport,
			ie.NewUsageReportWithinSessionReportRequest(
				r.IEsWithinSessReportReq(
					urrInfo.MeasureMethod, urrInfo.MeasureInformation)...,
			))
	}

	err = s.sendReqTo(req, addr)
	return errors.Wrap(err, "serveUSAReport")
}
