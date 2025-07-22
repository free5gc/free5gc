package nas_security

import (
	"encoding/binary"
)

// API for N3IWF
func EncapNasMsgToEnvelope(nasPDU []byte) []byte {
	// According to TS 24.502 8.2.4,
	// in order to transport a NAS message over the non-3GPP access between the UE and the N3IWF,
	// the NAS message shall be framed in a NAS message envelope as defined in subclause 9.4.
	// According to TS 24.502 9.4,
	// a NAS message envelope = Length | NAS Message
	nasEnv := make([]byte, 2)
	binary.BigEndian.PutUint16(nasEnv, uint16(len(nasPDU)))
	nasEnv = append(nasEnv, nasPDU...)
	return nasEnv
}
