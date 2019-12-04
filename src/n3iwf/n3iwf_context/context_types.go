package n3iwf_context

type N3IWFNFInfo struct {
	GlobalN3IWFID   GlobalN3IWFID     `yaml:"GlobalN3IWFID"`
	RanNodeName     string            `yaml:"Name,omitempty"`
	SupportedTAList []SupportedTAItem `yaml:"SupportedTAList"`
}

type GlobalN3IWFID struct {
	PLMNID  PLMNID `yaml:"PLMNID"`
	N3IWFID uint16 `yaml:"N3IWFID"` // with length 2 bytes
}

type SupportedTAItem struct {
	TAC               string              `yaml:"TAC"`
	BroadcastPLMNList []BroadcastPLMNItem `yaml:"BroadcastPLMNList"`
}

type BroadcastPLMNItem struct {
	PLMNID              PLMNID             `yaml:"PLMNID"`
	TAISliceSupportList []SliceSupportItem `yaml:"TAISliceSupportList"`
}

type PLMNID struct {
	Mcc string `yaml:"MCC"`
	Mnc string `yaml:"MNC"`
}

type SliceSupportItem struct {
	SNSSAI SNSSAIItem `yaml:"SNSSAI"`
}

type SNSSAIItem struct {
	SST string `yaml:"SST"`
	SD  string `yaml:"SD,omitempty"`
}
