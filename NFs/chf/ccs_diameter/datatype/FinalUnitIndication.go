package datatype

import (
	diam_datatype "github.com/fiorix/go-diameter/diam/datatype"
)

type MultipleServicesCreditControl struct {
	GrantedServiceUnit   *GrantedServiceUnit      `avp:"Granted-Service-Unit"`
	RequestedServiceUnit *RequestedServiceUnit    `avp:"Requested-Service-Unit"`
	UsedServiceUnit      *UsedServiceUnit         `avp:"Used-Service-Unit"`
	TariffChangeUsage    diam_datatype.Enumerated `avp:"Tariff-Change-Usage"`
	ServiceIdentifier    diam_datatype.Unsigned32 `avp:"Service-Identifier"`
	RatingGroup          diam_datatype.Unsigned32 `avp:"Rating-Group"`
	GSUPoolReference     diam_datatype.Grouped    `avp:"G-S-U-Pool-Reference"`
	ValidityTime         diam_datatype.Unsigned32 `avp:"Validity-Time"`
	ResultCode           diam_datatype.Unsigned32 `avp:"Result-Code"`
	FinalUnitIndication  *FinalUnitIndication     `avp:"Final-Unit-Indication"`
}
