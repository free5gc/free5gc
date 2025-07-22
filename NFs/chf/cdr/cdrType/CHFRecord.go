package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const (
	CHFRecordPresentNothing int = iota /* No components present */
	CHFRecordPresentChargingFunctionRecord
)

// TS 32.298 5.1.5.0 CHF record (CHF-CDR)
type CHFRecord struct {
	Present int /* Choice Type */
	// For CHF CDR parameters, see 3GPP TS 32.298 5.1.5.1
	ChargingFunctionRecord *ChargingRecord `ber:"tagNum:200"`
}
