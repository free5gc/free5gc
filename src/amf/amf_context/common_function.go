package amf_context

import (
	"encoding/binary"
	"github.com/mohae/deepcopy"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/logger"
	"reflect"
)

func CompareUserLocation(loc1 models.UserLocation, loc2 models.UserLocation) bool {
	if loc1.EutraLocation != nil && loc2.EutraLocation != nil {
		eutraloc1 := deepcopy.Copy(*loc1.EutraLocation).(models.EutraLocation)
		eutraloc2 := deepcopy.Copy(*loc2.EutraLocation).(models.EutraLocation)
		eutraloc1.UeLocationTimestamp = nil
		eutraloc2.UeLocationTimestamp = nil
		return reflect.DeepEqual(eutraloc1, eutraloc2)
	}
	if loc1.N3gaLocation != nil && loc2.N3gaLocation != nil {
		return reflect.DeepEqual(loc1, loc2)
	}
	if loc1.NrLocation != nil && loc2.NrLocation != nil {
		nrloc1 := deepcopy.Copy(*loc1.NrLocation).(models.NrLocation)
		nrloc2 := deepcopy.Copy(*loc2.NrLocation).(models.NrLocation)
		nrloc1.UeLocationTimestamp = nil
		nrloc2.UeLocationTimestamp = nil
		return reflect.DeepEqual(nrloc1, nrloc2)
	}

	return false
}

func InTaiList(servedTai models.Tai, taiList []models.Tai) bool {
	for _, tai := range taiList {
		if reflect.DeepEqual(tai, servedTai) {
			return true
		}
	}
	return false
}

func TacInAreas(targetTac string, areas []models.Area) bool {
	for _, area := range areas {
		for _, tac := range area.Tacs {
			if targetTac == tac {
				return true
			}
		}
	}
	return false
}

func GetSecurityCount(overflow uint16, sqn uint8) []uint8 {
	var r = make([]byte, 4)
	binary.BigEndian.PutUint16(r, overflow)
	r[3] = sqn
	r[2] = r[1]
	r[1] = r[0]
	r[0] = 0x00
	return r
}
func AttachSourceUeTargetUe(sourceUe, targetUe *RanUe) {
	if sourceUe == nil {
		logger.ContextLog.Error("Source Ue is Nil")
		return
	}
	if targetUe == nil {
		logger.ContextLog.Error("Target Ue is Nil")
		return
	}
	amfUe := sourceUe.AmfUe
	if amfUe == nil {
		logger.ContextLog.Error("AmfUe is Nil")
		return
	}
	targetUe.AmfUe = amfUe
	targetUe.SourceUe = sourceUe
	sourceUe.TargetUe = targetUe
}

func DetachSourceUeTargetUe(ranUe *RanUe) {
	if ranUe == nil {
		logger.ContextLog.Error("ranUe is Nil")
		return
	}
	if ranUe.TargetUe != nil {
		targetUe := ranUe.TargetUe

		ranUe.TargetUe = nil
		targetUe.SourceUe = nil

	} else if ranUe.SourceUe != nil {
		source := ranUe.SourceUe

		ranUe.SourceUe = nil
		source.TargetUe = nil
	}

}
