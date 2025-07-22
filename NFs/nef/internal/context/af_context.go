package context

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/openapi/models_nef"
	"github.com/sirupsen/logrus"
)

type AfData struct {
	AfID       string
	NumSubscID uint64
	NumTransID uint64
	Subs       map[string]*AfSubscription
	PfdTrans   map[string]*AfPfdTransaction
	Mu         sync.RWMutex
	Log        *logrus.Entry
}

func (a *AfData) NewSub(numCorreID uint64, tiSub *models_nef.TrafficInfluSub) *AfSubscription {
	a.NumSubscID++
	sub := AfSubscription{
		NotifCorreID: strconv.FormatUint(numCorreID, 10),
		SubID:        strconv.FormatUint(a.NumSubscID, 10),
		TiSub:        tiSub,
		Log:          a.Log.WithField(logger.FieldSubID, fmt.Sprintf("SUB:%d", a.NumSubscID)),
	}
	sub.Log.Infoln("New subscription")
	return &sub
}

func (a *AfData) NewPfdTrans() *AfPfdTransaction {
	a.NumTransID++
	pfdTr := AfPfdTransaction{
		TransID:   strconv.FormatUint(a.NumTransID, 10),
		ExtAppIDs: make(map[string]struct{}),
		Log:       a.Log.WithField(logger.FieldPfdTransID, fmt.Sprintf("PFDT:%d", a.NumTransID)),
	}
	pfdTr.Log.Infoln("New pfd transcation")
	return &pfdTr
}

func (a *AfData) IsAppIDExisted(appID string) (string, bool) {
	for _, pfdTrans := range a.PfdTrans {
		if _, ok := pfdTrans.ExtAppIDs[appID]; ok {
			return pfdTrans.TransID, true
		}
	}
	return "", false
}
