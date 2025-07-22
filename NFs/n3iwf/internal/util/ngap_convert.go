package util

import (
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/free5gc/aper"
	"github.com/free5gc/n3iwf/internal/logger"
	"github.com/free5gc/n3iwf/pkg/factory"
	"github.com/free5gc/ngap/ngapType"
)

func PlmnIdToNgap(plmnId factory.PLMNID) (ngapPlmnId ngapType.PLMNIdentity) {
	var hexString string
	mcc := strings.Split(plmnId.Mcc, "")
	mnc := strings.Split(plmnId.Mnc, "")
	if len(plmnId.Mnc) == 2 {
		hexString = mcc[1] + mcc[0] + "f" + mcc[2] + mnc[1] + mnc[0]
	} else {
		hexString = mcc[1] + mcc[0] + mnc[0] + mcc[2] + mnc[2] + mnc[1]
	}
	var err error
	ngapPlmnId.Value, err = hex.DecodeString(hexString)
	if err != nil {
		logger.UtilLog.Errorf("DecodeString error: %+v", err)
	}
	return
}

func N3iwfIdToNgap(n3iwfId uint16) (ngapN3iwfId *aper.BitString) {
	ngapN3iwfId = new(aper.BitString)
	ngapN3iwfId.Bytes = make([]byte, 2)
	binary.BigEndian.PutUint16(ngapN3iwfId.Bytes, n3iwfId)
	ngapN3iwfId.BitLength = 16
	return
}
