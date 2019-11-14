package n3iwf_context

var n3iwfContext = N3IWFContext{}
var ranUeNgapIdGenerator int64 = 0

type N3IWFContext struct {
	UePool map[int64]*N3IWFUe // RanUeNgapID as key
}

func init() {
	N3IWFSelf().UePool = make(map[int64]*N3IWFUe)
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
