package gmm

import (
	"free5gc/lib/fsm"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/gmm/gmm_state"
)

func NewGmmFuncTable(anType models.AccessType) fsm.FuncTable {
	table := fsm.NewFuncTable()
	if anType == models.AccessType__3_GPP_ACCESS {
		table[gmm_state.DE_REGISTERED] = DeRegistered_3gpp
		table[gmm_state.REGISTERED] = Registered_3gpp
		table[gmm_state.AUTHENTICATION] = Authentication_3gpp
		table[gmm_state.SECURITY_MODE] = SecurityMode_3gpp
		table[gmm_state.INITIAL_CONTEXT_SETUP] = InitialContextSetup_3gpp
	} else {
		table[gmm_state.DE_REGISTERED] = DeRegistered_non_3gpp
		table[gmm_state.REGISTERED] = Registered_non_3gpp
		table[gmm_state.AUTHENTICATION] = Authentication_non_3gpp
		table[gmm_state.SECURITY_MODE] = SecurityMode_non_3gpp
		table[gmm_state.INITIAL_CONTEXT_SETUP] = InitialContextSetup_non_3gpp
	}

	table[gmm_state.EXCEPTION] = Exception

	return table
}

func InitAmfUeSm(ue *amf_context.AmfUe) (err error) {
	table := NewGmmFuncTable(models.AccessType__3_GPP_ACCESS)
	ue.Sm[models.AccessType__3_GPP_ACCESS], err = fsm.NewFSM(gmm_state.DE_REGISTERED, table)
	if err != nil {
		return
	}
	table = NewGmmFuncTable(models.AccessType_NON_3_GPP_ACCESS)
	ue.Sm[models.AccessType_NON_3_GPP_ACCESS], err = fsm.NewFSM(gmm_state.DE_REGISTERED, table)
	return
}
