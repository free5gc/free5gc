package perio

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/free5gc/go-upf/internal/logger"
	"github.com/free5gc/go-upf/internal/report"
)

const (
	EVENT_CHANNEL_LEN = 512
)

type EventType uint8

const (
	TYPE_PERIO_ADD EventType = iota + 1
	TYPE_PERIO_DEL
	TYPE_PERIO_TIMEOUT
	TYPE_SERVER_CLOSE
)

func (t EventType) String() string {
	s := []string{
		"", "TYPE_PERIO_ADD", "TYPE_PERIO_DEL",
		"TYPE_PERIO_TIMEOUT", "TYPE_SERVER_CLOSE",
	}
	return s[t]
}

type Event struct {
	eType  EventType
	lSeid  uint64
	urrid  uint32
	period time.Duration
}

type PERIOGroup struct {
	urrids map[uint64]map[uint32]struct{}
	period time.Duration
	ticker *time.Ticker
	stopCh chan struct{}
}

func (pg *PERIOGroup) newTicker(wg *sync.WaitGroup, evtCh chan Event) error {
	if pg.ticker != nil {
		return errors.Errorf("ticker not nil")
	}
	logger.PerioLog.Infof("new ticker [%+v]", pg.period)

	pg.ticker = time.NewTicker(pg.period)
	pg.stopCh = make(chan struct{})

	wg.Add(1)
	go func(ticker *time.Ticker, period time.Duration, evtCh chan Event) {
		defer func() {
			ticker.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-ticker.C:
				logger.PerioLog.Debugf("ticker[%v] timeout", period)
				// If the UPF had terminating, the evtCh would be nil
				if evtCh != nil {
					evtCh <- Event{
						eType:  TYPE_PERIO_TIMEOUT,
						period: period,
					}
				}
			case <-pg.stopCh:
				logger.PerioLog.Infof("ticker[%v] Stopped", period)
				return
			}
		}
	}(pg.ticker, pg.period, evtCh)

	return nil
}

func (pg *PERIOGroup) stopTicker() {
	logger.PerioLog.Debugf("stopTicker: [%+v]", pg.period)
	pg.stopCh <- struct{}{}
	close(pg.stopCh)
}

type Server struct {
	evtCh     chan Event
	perioList map[time.Duration]*PERIOGroup // key: period

	handler  report.Handler
	queryURR func(map[uint64][]uint32) (map[uint64][]report.USAReport, error)
}

func OpenServer(wg *sync.WaitGroup) (*Server, error) {
	s := &Server{
		evtCh:     make(chan Event, EVENT_CHANNEL_LEN),
		perioList: make(map[time.Duration]*PERIOGroup),
	}

	wg.Add(1)
	go s.Serve(wg)

	return s, nil
}

func (s *Server) Close() {
	s.evtCh <- Event{eType: TYPE_SERVER_CLOSE}
}

func (s *Server) Handle(
	handler report.Handler,
	queryURR func(map[uint64][]uint32) (map[uint64][]report.USAReport, error),
) {
	s.handler = handler
	s.queryURR = queryURR
}

func (s *Server) Serve(wg *sync.WaitGroup) {
	logger.PerioLog.Infof("perio server started")
	defer func() {
		logger.PerioLog.Infof("perio server stopped")
		close(s.evtCh)
		wg.Done()
	}()

	for e := range s.evtCh {
		logger.PerioLog.Infof("recv event[%s][%+v]", e.eType, e)
		switch e.eType {
		case TYPE_PERIO_ADD:
			perioGroup, ok := s.perioList[e.period]
			if !ok {
				// New ticker if no this period ticker found
				perioGroup = &PERIOGroup{
					urrids: make(map[uint64]map[uint32]struct{}),
					period: e.period,
				}
				err := perioGroup.newTicker(wg, s.evtCh)
				if err != nil {
					logger.PerioLog.Errorln(err)
					continue
				}
				s.perioList[e.period] = perioGroup
			}

			urrids := perioGroup.urrids[e.lSeid]
			if urrids == nil {
				perioGroup.urrids[e.lSeid] = make(map[uint32]struct{})
				perioGroup.urrids[e.lSeid][e.urrid] = struct{}{}
			} else {
				_, ok := perioGroup.urrids[e.lSeid][e.urrid]
				if !ok {
					perioGroup.urrids[e.lSeid][e.urrid] = struct{}{}
				}
			}
		case TYPE_PERIO_DEL:
			for period, perioGroup := range s.perioList {
				_, ok := perioGroup.urrids[e.lSeid][e.urrid]
				if ok {
					// Stop ticker if no more PERIO URR
					delete(perioGroup.urrids[e.lSeid], e.urrid)
					if len(perioGroup.urrids[e.lSeid]) == 0 {
						delete(perioGroup.urrids, e.lSeid)
						if len(perioGroup.urrids) == 0 {
							// If no urr for the ticker, this ticker could be stop and delete
							perioGroup.stopTicker()
							delete(s.perioList, period)
						}
					}
					break
				}
			}
		case TYPE_PERIO_TIMEOUT:
			var lSeidUrridsMap map[uint64][]uint32

			perioGroup, ok := s.perioList[e.period]
			if !ok {
				logger.PerioLog.Warnf("no periodGroup found for period[%v]", e.period)
				break
			}

			lSeidUrridsMap = make(map[uint64][]uint32)
			for lSeid, urrIds := range perioGroup.urrids {
				for urrId := range urrIds {
					lSeidUrridsMap[lSeid] = append(lSeidUrridsMap[lSeid], urrId)
				}
			}

			seidUsars, err := s.queryURR(lSeidUrridsMap)
			if err != nil {
				logger.PerioLog.Warnf("get Multiple USAReports error: %v", err)
				break
			}
			if len(seidUsars) == 0 {
				logger.PerioLog.Warnf("no PERIO USAReport")
				break
			}

			for seid, usars := range seidUsars {
				var rpts []report.Report

				for i := range usars {
					usars[i].USARTrigger.Flags |= report.USAR_TRIG_PERIO
					rpts = append(rpts, usars[i])
				}

				s.handler.NotifySessReport(
					report.SessReport{
						SEID:    seid,
						Reports: rpts,
					})
			}
		case TYPE_SERVER_CLOSE:
			for period, perioGroup := range s.perioList {
				perioGroup.stopTicker()
				delete(s.perioList, period)
			}
			return
		}
	}
}

func (s *Server) AddPeriodReportTimer(lSeid uint64, urrid uint32, period time.Duration) {
	s.evtCh <- Event{
		eType:  TYPE_PERIO_ADD,
		lSeid:  lSeid,
		urrid:  urrid,
		period: period,
	}
}

func (s *Server) DelPeriodReportTimer(lSeid uint64, urrid uint32) {
	s.evtCh <- Event{
		eType: TYPE_PERIO_DEL,
		lSeid: lSeid,
		urrid: urrid,
	}
}
