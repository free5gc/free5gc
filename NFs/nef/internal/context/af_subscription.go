package context

import (
	"github.com/free5gc/openapi/models_nef"
	"github.com/sirupsen/logrus"
)

type AfSubscription struct {
	SubID        string
	TiSub        *models_nef.TrafficInfluSub
	AppSessID    string // use in single UE case
	InfluID      string // use in multiple UE case
	NotifCorreID string
	Log          *logrus.Entry
}

func (s *AfSubscription) PatchTiSubData(tiSubPatch *models_nef.TrafficInfluSubPatch) {
	s.TiSub.AppReloInd = tiSubPatch.AppReloInd
	s.TiSub.TrafficFilters = tiSubPatch.TrafficFilters
	s.TiSub.EthTrafficFilters = tiSubPatch.EthTrafficFilters
	s.TiSub.TrafficRoutes = tiSubPatch.TrafficRoutes
	s.TiSub.TfcCorrInd = tiSubPatch.TfcCorrInd
	s.TiSub.TempValidities = tiSubPatch.TempValidities
	s.TiSub.ValidGeoZoneIds = tiSubPatch.ValidGeoZoneIds
	s.TiSub.AfAckInd = tiSubPatch.AfAckInd
	s.TiSub.AddrPreserInd = tiSubPatch.AddrPreserInd
}
