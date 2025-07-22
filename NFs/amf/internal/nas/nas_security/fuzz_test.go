//go:build go1.18
// +build go1.18

package nas_security_test

import (
	"fmt"
	"reflect"
	"testing"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/internal/nas/nas_security"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/openapi/models"
)

func FuzzNASSecurity(f *testing.F) {
	f.Fuzz(func(t *testing.T, d []byte) {
		// No security
		ue := newFuzzTestAmfUe()
		//nolint:errcheck // fuzzing code
		nas_security.Decode(ue, models.AccessType__3_GPP_ACCESS, d, true)

		// With security (NIA0/NEA0)
		ue = newFuzzTestAmfUe()
		ue.SecurityContextAvailable = true
		ue.IntegrityAlg = security.AlgIntegrity128NIA0
		ue.CipheringAlg = security.AlgCiphering128NEA0
		msg0, integrityProtected0, err0 := nas_security.Decode(ue, models.AccessType__3_GPP_ACCESS, d, true)

		if len(d) >= 7 {
			ue = newFuzzTestAmfUe()
			ue.SecurityContextAvailable = true
			ue.IntegrityAlg = security.AlgIntegrity128NIA2
			ue.CipheringAlg = security.AlgCiphering128NEA2
			if err := security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, uint32(d[6]),
				security.AccessType3GPP, security.DirectionUplink, d[7:]); err == nil {
				if mac32, errNASMacCalculate := security.NASMacCalculate(ue.IntegrityAlg, ue.KnasInt, uint32(d[6]),
					security.AccessType3GPP, security.DirectionUplink, d[6:]); errNASMacCalculate == nil {
					copy(d[2:6], mac32)
					msg2, integrityProtected2, err2 := nas_security.Decode(ue, models.AccessType__3_GPP_ACCESS, d, true)
					if err0 == nil && integrityProtected0 &&
						(d[1]&0x0f == nas.SecurityHeaderTypeIntegrityProtectedAndCiphered ||
							d[1]&0x0f == nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext) {
						if err2 != nil {
							panic(fmt.Sprintf("err mismatch: %s", err2))
						}
						if !integrityProtected2 {
							panic("integrityProtected mismatch")
						}
						if !reflect.DeepEqual(msg0.GmmMessage, msg2.GmmMessage) {
							panic("msg mismatch")
						}
					}
				}
			}
		}
	})
}

func newFuzzTestAmfUe() *amf_context.AmfUe {
	ue := new(amf_context.AmfUe)
	ue.RanUe = make(map[models.AccessType]*amf_context.RanUe)
	ue.RanUe[models.AccessType__3_GPP_ACCESS] = new(amf_context.RanUe)
	ue.RanUe[models.AccessType__3_GPP_ACCESS].AmfUe = ue
	ue.RanUe[models.AccessType__3_GPP_ACCESS].Ran = new(amf_context.AmfRan)
	ue.RanUe[models.AccessType__3_GPP_ACCESS].Ran.AnType = models.AccessType__3_GPP_ACCESS
	ue.NASLog = logger.NasLog
	return ue
}
