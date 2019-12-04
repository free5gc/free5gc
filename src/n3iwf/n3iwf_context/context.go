package n3iwf_context

import (
	"github.com/sirupsen/logrus"

	"free5gc/src/n3iwf/logger"
)

var contextLog *logrus.Entry

var n3iwfContext = N3IWFContext{}
var ranUeNgapIdGenerator int64 = 0

type N3IWFContext struct {
	NFInfo  N3IWFNFInfo
	UePool  map[int64]*N3IWFUe   // RanUeNgapID as key
	AMFPool map[string]*N3IWFAMF // SCTPSessionID as key
}

func init() {
	// init log
	contextLog = logger.ContextLog

	// init context
	N3IWFSelf().UePool = make(map[int64]*N3IWFUe)
	N3IWFSelf().AMFPool = make(map[string]*N3IWFAMF)
}

// Create new N3IWF context
func N3IWFSelf() *N3IWFContext {
	return &n3iwfContext
}

func (context *N3IWFContext) RanUeNgapIDAlloc() int64 {
	ranUeNgapIdGenerator %= MaxValueOfRanUeNgapID
	ranUeNgapIdGenerator++
	for {
		if _, double := context.UePool[ranUeNgapIdGenerator]; double {
			ranUeNgapIdGenerator++
		} else {
			break
		}
	}
	return ranUeNgapIdGenerator
}

func (context *N3IWFContext) NewN3iwfUe() *N3IWFUe {
	self := N3IWFSelf()
	n3iwfUe := N3IWFUe{}
	n3iwfUe.init()
	n3iwfUe.RanUeNgapId = context.RanUeNgapIDAlloc()
	n3iwfUe.AmfUeNgapId = AmfUeNgapIdUnspecified
	self.UePool[n3iwfUe.RanUeNgapId] = &n3iwfUe
	return &n3iwfUe
}

func (context *N3IWFContext) NewN3iwfAmf(sessionID string) *N3IWFAMF {
	if amf, ok := context.AMFPool[sessionID]; ok {
		contextLog.Warn("[Context] NewN3iwfAmf(): AMF entry already exists.")
		return amf
	} else {
		amf = &N3IWFAMF{}
		context.AMFPool[sessionID] = amf
		return amf
	}
}

func (context *N3IWFContext) FindAMFBySCTPSessionID(sessionID string) *N3IWFAMF {
	amf, ok := context.AMFPool[sessionID]
	if !ok {
		contextLog.Warnf("[Context] FindAMFBySCTPSessionID(): AMF not found. SessionID: %s", sessionID)
	}
	return amf
}

func (context *N3IWFContext) FindUeByAmfUeNgapID(amfUeNgapID int64) *N3IWFUe {
	self := N3IWFSelf()

	for _, ue := range self.UePool {
		if ue.AmfUeNgapId == amfUeNgapID {
			return ue
		}
	}

	return nil
}

func (context *N3IWFContext) FindUeByRanUeNgapID(ranUeNgapID int64) *N3IWFUe {
	self := N3IWFSelf()

	for _, ue := range self.UePool {
		if ue.RanUeNgapId == ranUeNgapID {
			return ue
		}
	}

	return nil
}
