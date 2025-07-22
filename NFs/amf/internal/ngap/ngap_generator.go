//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/free5gc/aper"
	"github.com/free5gc/ngap/ngapType"
)

type messageDirection int

const (
	messageDirectionRANtoAMF messageDirection = 1
	messageDirectionAMFtoRAN messageDirection = 2
	messageDirectionBoth     messageDirection = messageDirectionRANtoAMF | messageDirectionAMFtoRAN
)

// NGAP IE definition
type IEInfo struct {
	Criticality   aper.Enumerated
	Type          string
	Presence      aper.Enumerated
	GoID          string
	GoField       string
	GoVar         string
	GoType        string
	Unimplemented bool
}

// NGAP message definition
type MsgInfo struct {
	ProcName    string
	ProcCode    string
	Criticality aper.Enumerated
	IEs         map[string]*IEInfo
	IEorder     []string
	Type        int
	GoField     string
	GoTypeVar   string
	GoMsgVar    string
	TypeDesc    string
}

var MsgTable map[string]*MsgInfo
var msgNames []string

// Convert golang name from ASN.1 name
func convGoName(name string) string {
	if name == "CoreNetworkAssistanceInformationForInactive" {
		return "CoreNetworkAssistanceInformation"
	}
	if name == "UE-associatedLogicalNG-connectionList" {
		return "UEAssociatedLogicalNGConnectionList"
	}
	return strings.Replace(name, "-", "", -1)
}

// Convert golang name for local symbol from ASN.1 name
func convGoLocalName(name string) string {
	lname := convGoName(name)
	lname = strings.ToLower(lname[:1]) + lname[1:]
	return lname
}

func main() {
	readASN1()

	fixIEs()

	generateHandler()

	generateDispatcher()
}

func readASN1() {
	reMsgName := regexp.MustCompile(`^(.*)IEs\s+NGAP-PROTOCOL-IES\s+::=\s+\{\s*$`)
	reEntry := regexp.MustCompile(`^\s*\{\s*ID\s+(\S+)\s+CRITICALITY\s+(\S+)\s+TYPE\s+(\S+)\s+PRESENCE\s+(\S+)\s*\}\s*[,|]\s*$`)
	reMsgEnd := regexp.MustCompile(`^\s*\}\s*$`)
	reElem := regexp.MustCompile(`^(.*)\s+NGAP-ELEMENTARY-PROCEDURE\s+::=\s*\{\s*$`)
	reElemEntry := regexp.MustCompile(`^\s*(INITIATING|SUCCESSFUL|UNSUCCESSFUL)\s+(MESSAGE|OUTCOME)\s+(\S+)\s*$`)
	reElemEntry2 := regexp.MustCompile(`^\s*CRITICALITY\s+(\S+)\s*$`)
	reElemEntry3 := regexp.MustCompile(`^\s*PROCEDURE\s+CODE\s+id-(\S+)\s*$`)

	MsgTable = make(map[string]*MsgInfo)

	// parse ASN.1 file
	fAsn, err := os.Open("asn1/38413-fd0.asn")
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(fAsn)
	for s.Scan() {
		// parse IE infos
		m := reMsgName.FindStringSubmatch(s.Text())
		if m != nil {
			msgName := m[1]
			msgName = strings.TrimSuffix(msgName, "-")
			// fmt.Printf("%s\n", msgName)
			for s.Scan() {
				m := reMsgEnd.FindStringSubmatch(s.Text())
				if m != nil {
					break
				}
				m = reEntry.FindStringSubmatch(s.Text())
				if m != nil {
					ie := &IEInfo{}
					if MsgTable[msgName] == nil {
						continue
					}
					ieId := m[1]
					ieType := m[3]
					// XXX No definition is exist in NGAP package currently.
					if msgName == "InitialContextSetupRequest" && ieId == "id-LocationReportingRequestType" {
						continue
					}
					if msgName == "SecondaryRATDataUsageReport" && ieId == "id-UserLocationInformation" {
						continue
					}
					if ieId == "id-OldAssociatedQosFlowList-ULendmarkerexpected" ||
						ieId == "id-CNTypeRestrictionsForEquivalent" ||
						ieId == "id-CNTypeRestrictionsForServing" ||
						ieId == "id-NewGUAMI" ||
						ieId == "id-ULForwarding" ||
						ieId == "id-ULForwardingUP-TNLInformation" ||
						ieId == "id-CNAssistedRANTuning" ||
						ieId == "id-CommonNetworkInstance" ||
						ieId == "id-NGRAN-TNLAssociationToRemoveList" ||
						ieId == "id-TNLAssociationTransportLayerAddressNGRAN" ||
						ieId == "id-EndpointIPAddressAndPort" ||
						ieId == "id-LocationReportingAdditionalInfo" ||
						ieId == "id-SourceToTarget-AMFInformationReroute" ||
						ieId == "id-AdditionalULForwardingUPTNLInformation" ||
						ieId == "id-SCTP-TLAs" ||
						ieId == "id-SelectedPLMNIdentity" {
						continue
					}
					MsgTable[msgName].IEs[ieId] = ie
					MsgTable[msgName].IEorder = append(MsgTable[msgName].IEorder, ieId)
					ie.Criticality = str2Criticality(m[2])
					ie.Type = ieType
					ie.Presence = str2Presence(m[4])
					ie.GoID = "ngapType.ProtocolIEID" + convGoName(ieId[3:])
					ie.GoField = "Value." + convGoName(ieId[3:])
					ie.GoVar = convGoLocalName(ieId[3:])
					ie.GoType = "ngapType." + convGoName(ieType)
					// fmt.Printf("%+v\n", *ie)
				}
			}
		}

		// parse message infos
		m = reElem.FindStringSubmatch(s.Text())
		if m != nil {
			procName := m[1]
			for s.Scan() {
				m := reMsgEnd.FindStringSubmatch(s.Text())
				if m != nil {
					break
				}
				m = reElemEntry.FindStringSubmatch(s.Text())
				if m != nil {
					msgName := m[3]
					// fmt.Println(msgName)
					if _, exist := MsgTable[msgName]; exist {
						panic(fmt.Sprintf("%s is exist", msgName))
					}
					MsgTable[msgName] = &MsgInfo{
						IEs:     make(map[string]*IEInfo),
						IEorder: make([]string, 0),
					}
					var goField string
					switch m[1] {
					case "INITIATING":
						MsgTable[msgName].Type = ngapType.NGAPPDUPresentInitiatingMessage
						goField = "InitiatingMessage"
						MsgTable[msgName].TypeDesc = "Initiating Message"
					case "SUCCESSFUL":
						MsgTable[msgName].Type = ngapType.NGAPPDUPresentSuccessfulOutcome
						goField = "SuccessfulOutcome"
						MsgTable[msgName].TypeDesc = "Successful Outcome"
					case "UNSUCCESSFUL":
						MsgTable[msgName].Type = ngapType.NGAPPDUPresentUnsuccessfulOutcome
						goField = "UnsuccessfulOutcome"
						MsgTable[msgName].TypeDesc = "Unsuccessful Outcome"
					}
					MsgTable[msgName].GoField = goField
					MsgTable[msgName].GoTypeVar = convGoLocalName(goField)
					MsgTable[msgName].GoMsgVar = convGoLocalName(msgName)
					MsgTable[msgName].ProcName = procName
				}
				m = reElemEntry2.FindStringSubmatch(s.Text())
				if m != nil {
					criticality := str2Criticality(m[1])
					for _, mInfo := range MsgTable {
						if mInfo.ProcName == procName {
							mInfo.Criticality = criticality
						}
					}
				}
				m = reElemEntry3.FindStringSubmatch(s.Text())
				if m != nil {
					for _, mInfo := range MsgTable {
						if mInfo.ProcName == procName {
							mInfo.ProcCode = m[1]
						}
					}
				}
			}
		}
	}
	fAsn.Close()

	msgNames = make([]string, 0, len(MsgTable))
	for msgName := range MsgTable {
		// XXX Skip private message
		if msgName == "PrivateMessage" {
			continue
		}
		msgNames = append(msgNames, msgName)
	}
	sort.Strings(msgNames)
}

func fixIEs() {
	// Not implemented IEs
	MsgTable["AMFConfigurationUpdateAcknowledge"].IEs["id-AMF-TNLAssociationSetupList"].Unimplemented = true
	MsgTable["AMFConfigurationUpdateAcknowledge"].IEs["id-AMF-TNLAssociationFailedToSetupList"].Unimplemented = true
	MsgTable["AMFConfigurationUpdateFailure"].IEs["id-TimeToWait"].Unimplemented = true
	MsgTable["HandoverRequired"].IEs["id-DirectForwardingPathAvailability"].Unimplemented = true
	MsgTable["InitialUEMessage"].IEs["id-AMFSetID"].Unimplemented = true
	MsgTable["InitialUEMessage"].IEs["id-AllowedNSSAI"].Unimplemented = true
	MsgTable["NGSetupRequest"].IEs["id-UERetentionInformation"].Unimplemented = true
	MsgTable["RANConfigurationUpdate"].IEs["id-RANNodeName"].Unimplemented = true
	MsgTable["RANConfigurationUpdate"].IEs["id-DefaultPagingDRX"].Unimplemented = true
	MsgTable["RANConfigurationUpdate"].IEs["id-GlobalRANNodeID"].Unimplemented = true
	// MsgTable["UERadioCapabilityCheckResponse"].IEs["id-AMF-UE-NGAP-ID"].Unimplemented = true
	// MsgTable["UERadioCapabilityCheckResponse"].IEs["id-RAN-UE-NGAP-ID"].Unimplemented = true
	MsgTable["UERadioCapabilityCheckResponse"].IEs["id-IMSVoiceSupportIndicator"].Unimplemented = true
	MsgTable["UplinkRANConfigurationTransfer"].IEs["id-ENDC-SONConfigurationTransferUL"].Unimplemented = true
	MsgTable["UplinkRANStatusTransfer"].IEs["id-RANStatusTransfer-TransparentContainer"].Unimplemented = true
	MsgTable["UplinkUEAssociatedNRPPaTransport"].IEs["id-NRPPa-PDU"].Unimplemented = true
}

// generate NGAP handler file
func generateHandler() {
	fOut := newOutputFile("handler_generated.go",
		"ngap",
		[]string{
			"\"github.com/free5gc/amf/internal/context\"",
			"\"github.com/free5gc/amf/internal/logger\"",
			"ngap_message \"github.com/free5gc/amf/internal/ngap/message\"",
			"\"github.com/free5gc/ngap\"",
			"\"github.com/free5gc/ngap/ngapType\"",
		})

	// generate handler functions
	for _, msgName := range msgNames {
		mInfo := MsgTable[msgName]

		// generate function header
		messageAppend := ""
		if msgName == "InitialUEMessage" {
			messageAppend = ", message *ngapType.NGAPPDU"
		}
		fmt.Fprintf(fOut, "func handler%s(ran *context.AmfRan%s, %s *ngapType.%s) {\n", msgName, messageAppend, mInfo.GoTypeVar, mInfo.GoField)

		// setup variables
		isRequest := false
		mainFuncArgDefs := []string{"ran *context.AmfRan"}
		mainFuncArgs := []string{"ran"}
		if mInfo.Type == ngapType.NGAPPDUPresentInitiatingMessage && msgName != "ErrorIndication" {
			isRequest = true
		}
		amfIdIe := mInfo.IEs["id-AMF-UE-NGAP-ID"]
		ranIdIe := mInfo.IEs["id-RAN-UE-NGAP-ID"]
		amfIdIeVar := "nil"
		if amfIdIe != nil {
			amfIdIeVar = amfIdIe.GoVar
		}
		ranIdIeVar := "nil"
		if ranIdIe != nil {
			ranIdIeVar = ranIdIe.GoVar
		}
		hasRanUe := false
		ranUeMayNil := false
		if msgName == "InitialUEMessage" {
			mainFuncArgDefs = append(mainFuncArgDefs, "message *ngapType.NGAPPDU")
			mainFuncArgs = append(mainFuncArgs, "message")
		}

		// generate IE variables
		for _, ieName := range mInfo.IEorder {
			ieInfo := mInfo.IEs[ieName]
			// fmt.Fprintf(fOut, "// %s\n", ieName)
			fmt.Fprintf(fOut, "var %s *%s\n", ieInfo.GoVar, ieInfo.GoType)
		}
		fmt.Fprintln(fOut, "")
		fmt.Fprintln(fOut, "var syntaxCause *ngapType.Cause")
		fmt.Fprintln(fOut, "var iesCriticalityDiagnostics ngapType.CriticalityDiagnosticsIEList")
		fmt.Fprintln(fOut, "abort := false")

		// generate extract IEs code
		fmt.Fprintln(fOut, "")
		fmt.Fprintf(fOut, "%s := %s.Value.%s\n", mInfo.GoMsgVar, mInfo.GoTypeVar, msgName)
		fmt.Fprintf(fOut, "if %s == nil {\n", mInfo.GoMsgVar)
		fmt.Fprintf(fOut, "ran.Log.Error(\"%s is nil\")\n", msgName)
		fmt.Fprintln(fOut, "return")
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")
		fmt.Fprintf(fOut, "ran.Log.Info(\"Handle %s\")\n", msgName)
		fmt.Fprintln(fOut, "")
		fmt.Fprintf(fOut, "for _, ie := range %s.ProtocolIEs.List {\n", mInfo.GoMsgVar)
		fmt.Fprintln(fOut, "switch ie.Id.Value {")
		for _, ieName := range mInfo.IEorder {
			ieInfo := mInfo.IEs[ieName]
			fmt.Fprintf(fOut, "case %s: // %s, %s\n", ieInfo.GoID, presence2Str(ieInfo.Presence), criticality2Str(ieInfo.Criticality))
			// supported IE

			// duplicate check code
			fmt.Fprintf(fOut, "if %s !=nil {\n", ieInfo.GoVar)
			fmt.Fprintf(fOut, "ran.Log.Error(\"Duplicate IE %s\")\n", ieInfo.Type)
			if isRequest {
				fmt.Fprint(fOut, `
syntaxCause = &ngapType.Cause{
	Present: ngapType.CausePresentProtocol,
	Protocol: &ngapType.CauseProtocol{
		Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,
	},
}
`[1:])
			}
			fmt.Fprintln(fOut, "abort = true")
			fmt.Fprintln(fOut, "break")
			fmt.Fprintf(fOut, "}\n")

			fmt.Fprintf(fOut, "%s = ie.%s\n", ieInfo.GoVar, ieInfo.GoField)
			fmt.Fprintf(fOut, "ran.Log.Trace(\"Decode IE %s\")\n", ieInfo.Type)
		}
		fmt.Fprintln(fOut, "default:")
		fmt.Fprintln(fOut, "switch ie.Criticality.Value {")
		fmt.Fprintln(fOut, "case ngapType.CriticalityPresentReject:")
		fmt.Fprintln(fOut, "ran.Log.Errorf(\"Not comprehended IE ID 0x%04x (criticality: reject)\", ie.Id.Value)")
		fmt.Fprintln(fOut, "case ngapType.CriticalityPresentIgnore:")
		fmt.Fprintln(fOut, "ran.Log.Infof(\"Not comprehended IE ID 0x%04x (criticality: ignore)\", ie.Id.Value)")
		fmt.Fprintln(fOut, "case ngapType.CriticalityPresentNotify:")
		fmt.Fprintln(fOut, "ran.Log.Warnf(\"Not comprehended IE ID 0x%04x (criticality: notify)\", ie.Id.Value)")
		if !isRequest {
			fmt.Fprintln(fOut, "item := buildCriticalityDiagnosticsIEItem(ie.Criticality.Value, ie.Id.Value, ngapType.TypeOfErrorPresentNotUnderstood)")
			fmt.Fprintln(fOut, "iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)")
		}
		fmt.Fprintln(fOut, "}")
		if isRequest {
			fmt.Fprintln(fOut, "if ie.Criticality.Value != ngapType.CriticalityPresentIgnore {")
			fmt.Fprintln(fOut, "item := buildCriticalityDiagnosticsIEItem(ie.Criticality.Value, ie.Id.Value, ngapType.TypeOfErrorPresentNotUnderstood)")
			fmt.Fprintln(fOut, "iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)")
			fmt.Fprintln(fOut, "if ie.Criticality.Value == ngapType.CriticalityPresentReject {")
			fmt.Fprintln(fOut, "abort = true")
			fmt.Fprintln(fOut, "}")
			fmt.Fprintln(fOut, "}")
		}
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "}")

		if isRequest {
			// check code to lack mandatory IEs
			fmt.Fprintln(fOut, "")
			for _, ieName := range mInfo.IEorder {
				ieInfo := mInfo.IEs[ieName]
				if ieInfo.Presence == ngapType.PresencePresentMandatory && ieInfo.Criticality == ngapType.CriticalityPresentReject {
					fmt.Fprintf(fOut, "if %s == nil {\n", ieInfo.GoVar)
					fmt.Fprintf(fOut, "ran.Log.Error(\"Missing IE %s\")\n", ieInfo.Type)
					fmt.Fprintf(fOut, "item := buildCriticalityDiagnosticsIEItem(ngapType.CriticalityPresentReject, %s, ngapType.TypeOfErrorPresentMissing)\n", ieInfo.GoID)
					fmt.Fprintln(fOut, "iesCriticalityDiagnostics.List = append(iesCriticalityDiagnostics.List, item)")
					fmt.Fprintln(fOut, "abort = true")
					fmt.Fprintln(fOut, "}")
				}
			}
		}

		// Generate Error Indication
		fmt.Fprintln(fOut, "")
		fmt.Fprintln(fOut, "if syntaxCause != nil || len(iesCriticalityDiagnostics.List) > 0 {")
		fmt.Fprintln(fOut, "ran.Log.Trace(\"Has IE error\")")
		genErrorIndicationCommon(fOut, mInfo)
		fmt.Fprintln(fOut, "var pIesCriticalityDiagnostics *ngapType.CriticalityDiagnosticsIEList")
		fmt.Fprintln(fOut, "if len(iesCriticalityDiagnostics.List) > 0 {")
		fmt.Fprintln(fOut, "pIesCriticalityDiagnostics = &iesCriticalityDiagnostics")
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "criticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, pIesCriticalityDiagnostics)")
		// Must report error by other message than ErrorIndication by these messages
		switch msgName {
		// AMF to RAN message
		// case "AMFConfigurationUpdate":
		// 	fmt.Fprintf(fOut, "ngap_message.SendAMFConfigurationUpdateFailure(ran, syntaxCause, &criticalityDiagnostics)\n")

		case "HandoverRequired":
			fmt.Fprintf(fOut, "if %s != nil && %s != nil {\n", amfIdIeVar, ranIdIeVar)
			avoidNilSyntaxCause(fOut)
			fmt.Fprintf(fOut, "rawSendHandoverPreparationFailure(ran, *%s, *%s, *syntaxCause, &criticalityDiagnostics)\n", amfIdIeVar, ranIdIeVar)
			fmt.Fprintln(fOut, "} else {")
			fmt.Fprintf(fOut, "ngap_message.SendErrorIndication(ran, %s, %s, syntaxCause, &criticalityDiagnostics)\n", amfIdIeVar, ranIdIeVar)
			fmt.Fprintln(fOut, "}")

		// AMF to RAN message
		// case "HandoverRequest":
		// 	fmt.Fprintf(fOut, "ngap_message.SendHandoverFailure(ran, syntaxCause, &criticalityDiagnostics)\n")

		// AMF to RAN message
		// case "InitialContextSetupRequest":
		// 	fmt.Fprintf(fOut, "ngap_message.SendInitialContextSetupFailure(ran, syntaxCause, &criticalityDiagnostics)\n")

		case "NGSetupRequest":
			avoidNilSyntaxCause(fOut)
			fmt.Fprintf(fOut, "rawSendNGSetupFailure(ran, *syntaxCause, nil, &criticalityDiagnostics)\n")

		// Cannot fill mandatory IEs
		// case "PathSwitchRequest":
		// 	fmt.Fprintf(fOut, "ngap_message.SendPathSwitchRequestFailure(ran, syntaxCause, &criticalityDiagnostics)\n")

		case "RANConfigurationUpdate":
			avoidNilSyntaxCause(fOut)
			fmt.Fprintf(fOut, "rawSendRANConfigurationUpdateFailure(ran, *syntaxCause, nil, &criticalityDiagnostics)\n")

		// AMF to RAN message
		// case "UEContextModificationRequest":
		// 	fmt.Fprintf(fOut, "ngap_message.SendUEContextModificationFailure(ran, syntaxCause, &criticalityDiagnostics)\n")

		default:
			fmt.Fprintf(fOut, "ngap_message.SendErrorIndication(ran, %s, %s, syntaxCause, &criticalityDiagnostics)\n", amfIdIeVar, ranIdIeVar)
		}
		fmt.Fprintln(fOut, "}")
		fmt.Fprintln(fOut, "")
		fmt.Fprintln(fOut, "if abort {")
		fmt.Fprintln(fOut, "return")
		fmt.Fprintln(fOut, "}")

		// To avoid Coverity's false positive, generate this check for Request messages too
		fmt.Fprintln(fOut, "")
		for _, ieName := range mInfo.IEorder {
			ieInfo := mInfo.IEs[ieName]
			if ieInfo.Presence == ngapType.PresencePresentMandatory {
				fmt.Fprintf(fOut, "if %s == nil {\n", ieInfo.GoVar)
				if ieInfo.Criticality == ngapType.CriticalityPresentReject {
					fmt.Fprintf(fOut, "ran.Log.Error(\"Missing IE %s\")\n", ieInfo.Type)
					fmt.Fprintln(fOut, "return")
				} else {
					fmt.Fprintf(fOut, "ran.Log.Warn(\"Missing IE %s\")\n", ieInfo.Type)
				}
				fmt.Fprintln(fOut, "}")
			}
			if ieInfo.Unimplemented {
				fmt.Fprintf(fOut, "if %s != nil {\n", ieInfo.GoVar)
				fmt.Fprintf(fOut, "ran.Log.Warn(\"IE %s is not implemented\")\n", ieInfo.Type)
				fmt.Fprintln(fOut, "}")
			}
		}

		// generate UE ID handling codes
		if amfIdIe != nil && msgName != "ErrorIndication" {
			hasRanUe = true
			firstReturnedMessage := "false"
			sendErrorIndication := "true"
			if msgName == "UEContextReleaseComplete" {
				sendErrorIndication = "false"
			}
			if msgName == "HandoverRequestAcknowledge" || msgName == "HandoverFailure" {
				firstReturnedMessage = "true"
			}
			fmt.Fprintln(fOut, "")
			fmt.Fprintf(fOut, "// AMF: %s, %s\n", presence2Str(amfIdIe.Presence), criticality2Str(amfIdIe.Criticality))
			if ranIdIe != nil {
				fmt.Fprintf(fOut, "// RAN: %s, %s\n", presence2Str(ranIdIe.Presence), criticality2Str(ranIdIe.Criticality))
			}
			fmt.Fprintln(fOut, "var ranUe *context.RanUe")
			if !(amfIdIe.Presence == ngapType.PresencePresentMandatory && amfIdIe.Criticality == ngapType.CriticalityPresentReject) {
				fmt.Fprintf(fOut, "if %s != nil {\n", amfIdIeVar)
				ranUeMayNil = true
			}
			fmt.Fprintln(fOut, "var err error")
			fmt.Fprintf(fOut, "ranUe, err = ranUeFind(ran, %s, %s, %s, %s)\n", amfIdIeVar, ranIdIeVar, firstReturnedMessage, sendErrorIndication)
			fmt.Fprintln(fOut, "if err != nil {")
			fmt.Fprintf(fOut, "ran.Log.Errorf(\"Handle %s: %%s\", err)\n", msgName)
			fmt.Fprintln(fOut, "return")
			fmt.Fprintln(fOut, "}")
			fmt.Fprintln(fOut, "if ranUe == nil {")
			fmt.Fprintf(fOut, "ran.Log.Error(\"Handle %s: No UE Context\")\n", msgName)
			fmt.Fprintln(fOut, "return")
			fmt.Fprintln(fOut, "}")
			fmt.Fprintf(fOut, "ranUe.Log.Infof(\"Handle %s (RAN UE NGAP ID: %%d)\", ranUe.RanUeNgapId)", msgName)
			if !(amfIdIe.Presence == ngapType.PresencePresentMandatory && amfIdIe.Criticality == ngapType.CriticalityPresentReject) {
				fmt.Fprintln(fOut, "}")
			}
			fmt.Fprintln(fOut, "")
		}

		if hasRanUe {
			mayNil := ""
			if ranUeMayNil {
				mayNil = " /* may be nil */"
			}
			mainFuncArgDefs = append(mainFuncArgDefs, "ranUe *context.RanUe")
			mainFuncArgs = append(mainFuncArgs, "ranUe"+mayNil)
		}
		for _, ieName := range mInfo.IEorder {
			ieInfo := mInfo.IEs[ieName]
			if !ieInfo.Unimplemented && (!hasRanUe || (ieName != "id-AMF-UE-NGAP-ID" && ieName != "id-RAN-UE-NGAP-ID") || (msgName == "HandoverRequestAcknowledge" && ieName == "id-RAN-UE-NGAP-ID")) {
				mayNil := ""
				if !(ieInfo.Presence == ngapType.PresencePresentMandatory && ieInfo.Criticality == ngapType.CriticalityPresentReject) {
					mayNil = " /* may be nil */"
				}
				mainFuncArgDefs = append(mainFuncArgDefs, fmt.Sprintf("%s *%s", ieInfo.GoVar, ieInfo.GoType))
				mainFuncArgs = append(mainFuncArgs, ieInfo.GoVar+mayNil)
			}
		}
		// Call main code of message handler
		fmt.Fprintln(fOut, "")
		fmt.Fprintf(fOut, "\t// func handle%sMain(%s) {\n", msgName, strings.Join(mainFuncArgDefs, ",\n\t//\t"))
		fmt.Fprintf(fOut, "handle%sMain(%s)\n", msgName, strings.Join(mainFuncArgs, ","))
		fmt.Fprintf(fOut, "}\n\n")

		if !isRANtoAMFMessage(msgName) ||
			msgName == "PWSCancelResponse" || // XXX not implemented
			msgName == "PWSFailureIndication" || // XXX not implemented
			msgName == "PWSRestartIndication" || // XXX not implemented
			msgName == "SecondaryRATDataUsageReport" || // XXX not implemented
			msgName == "TraceFailureIndication" || // XXX not implemented
			msgName == "WriteReplaceWarningResponse" { // XXX not implemented
			stubCause := "CauseProtocolPresentUnspecified"
			stubMessage := "not implemented"
			if isAMFtoRANMessage(msgName) {
				stubCause = "CauseProtocolPresentMessageNotCompatibleWithReceiverState"
				stubMessage = "AMF to RAN message"
			}
			// Stub message handler
			fmt.Fprintf(fOut, "func handle%sMain(%s){\n", msgName, strings.Join(mainFuncArgDefs, ", "))
			fmt.Fprintf(fOut, "ran.Log.Error(\"Handle %s: %s\")\n", msgName, stubMessage)
			if isRequest {
				genErrorIndicationCommon(fOut, mInfo)
				fmt.Fprintln(fOut, "notImplementedCriticalityDiagnostics := buildCriticalityDiagnostics(&procedureCode, &triggeringMessage, &procedureCriticality, nil)")
				fmt.Fprintln(fOut, "notImplementedCause := &ngapType.Cause{")
				fmt.Fprintln(fOut, "Present: ngapType.CausePresentProtocol,")
				fmt.Fprintln(fOut, "Protocol: &ngapType.CauseProtocol{")
				fmt.Fprintf(fOut, "Value: ngapType.%s,\n", stubCause)
				fmt.Fprintln(fOut, "},")
				fmt.Fprintln(fOut, "}")
				fmt.Fprintln(fOut, "ngap_message.SendErrorIndication(ran, nil, nil, notImplementedCause, &notImplementedCriticalityDiagnostics)")
			}
			fmt.Fprintln(fOut, "}")
			fmt.Fprintln(fOut, "")
		}
	}

	generateBuilder(fOut)
	generateSender(fOut)

	fOut.Close()
}

func avoidNilSyntaxCause(f *outputFile) {
	fmt.Fprintln(f, "if syntaxCause == nil {")
	fmt.Fprintln(f, "syntaxCause = &ngapType.Cause{")
	fmt.Fprintln(f, "Present: ngapType.CausePresentProtocol,")
	fmt.Fprintln(f, "Protocol: &ngapType.CauseProtocol{")
	fmt.Fprintln(f, "Value: ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage,")
	fmt.Fprintln(f, "},")
	fmt.Fprintln(f, "}")
	fmt.Fprintln(f, "}")
}

func generateDispatcher() {
	fOut := newOutputFile("dispatcher_generated.go",
		"ngap",
		[]string{
			"\"github.com/free5gc/amf/internal/context\"",
			"ngap_message \"github.com/free5gc/amf/internal/ngap/message\"",
			"\"github.com/free5gc/ngap/ngapType\"",
		})

	// Generate message dispatcher codes
	fmt.Fprintln(fOut, "")
	fmt.Fprintln(fOut, "func dispatchMain(ran *context.AmfRan, message *ngapType.NGAPPDU) {")
	fmt.Fprintln(fOut, "switch message.Present {")
	for _, present := range []string{
		"InitiatingMessage",
		"SuccessfulOutcome",
		"UnsuccessfulOutcome",
	} {
		fmt.Fprintf(fOut, "case ngapType.NGAPPDUPresent%s:\n", present)
		presentVar := convGoLocalName(present)
		fmt.Fprintf(fOut, "%s := message.%s\n", presentVar, present)
		fmt.Fprintf(fOut, "if %s == nil {\n", presentVar)
		fmt.Fprintf(fOut, "ran.Log.Errorln(\"%s is nil\")\n", present)
		fmt.Fprintln(fOut, "return")
		fmt.Fprintln(fOut, "}")
		fmt.Fprintf(fOut, "switch %s.ProcedureCode.Value {\n", presentVar)
		for _, msgName := range msgNames {
			mInfo := MsgTable[msgName]
			if mInfo.GoField == present {
				fmt.Fprintf(fOut, "case ngapType.ProcedureCode%s:\n", mInfo.ProcCode)
				messageAppend := ""
				if msgName == "InitialUEMessage" {
					messageAppend = ", message"
				}
				fmt.Fprintf(fOut, "handler%s(ran%s, %s)\n", msgName, messageAppend, presentVar)
			}
		}
		fmt.Fprintln(fOut, "default:")
		fmt.Fprintln(fOut, "cause := ngapType.Cause{")
		fmt.Fprintln(fOut, "Present:  ngapType.CausePresentProtocol,")
		fmt.Fprintln(fOut, "Protocol: &ngapType.CauseProtocol{},")
		fmt.Fprintln(fOut, "}")
		fmt.Fprintf(fOut, "switch %s.Criticality.Value {\n", presentVar)
		fmt.Fprintln(fOut, "case ngapType.CriticalityPresentReject:")
		fmt.Fprintf(fOut, "ran.Log.Errorf(\"Not comprehended procedure code of %s (criticality: reject, procedureCode:0x%%02x)\", %s.ProcedureCode.Value)\n", present, presentVar)
		fmt.Fprintln(fOut, "cause.Protocol.Value = ngapType.CauseProtocolPresentAbstractSyntaxErrorReject")
		fmt.Fprintln(fOut, "case ngapType.CriticalityPresentIgnore:")
		fmt.Fprintf(fOut, "ran.Log.Infof(\"Not comprehended procedure code of %s (criticality: ignore, procedureCode:0x%%02x)\", %s.ProcedureCode.Value)\n", present, presentVar)
		fmt.Fprintln(fOut, "return")
		fmt.Fprintln(fOut, "case ngapType.CriticalityPresentNotify:")
		fmt.Fprintf(fOut, "ran.Log.Warnf(\"Not comprehended procedure code of %s (criticality: notify, procedureCode:0x%%02x)\", %s.ProcedureCode.Value)\n", present, presentVar)
		fmt.Fprintln(fOut, "cause.Protocol.Value = ngapType.CauseProtocolPresentAbstractSyntaxErrorIgnoreAndNotify")
		fmt.Fprintln(fOut, "}")
		genTriggeringMessage(fOut, present)
		fmt.Fprintf(fOut, "criticalityDiagnostics := buildCriticalityDiagnostics(&%s.ProcedureCode.Value, &triggeringMessage, &%s.Criticality.Value, nil)\n", presentVar, presentVar)
		fmt.Fprintln(fOut, "ngap_message.SendErrorIndication(ran, nil, nil, &cause, &criticalityDiagnostics)")
		fmt.Fprintln(fOut, "}")
	}
	fmt.Fprintln(fOut, "}")
	fmt.Fprintln(fOut, "}")

	fOut.Close()
}

func generateBuilder(fOut io.Writer) {
	for _, msgName := range msgNames {
		mInfo := MsgTable[msgName]
		if msgName == "HandoverPreparationFailure" ||
			msgName == "NGSetupFailure" ||
			msgName == "PathSwitchRequestFailure" ||
			msgName == "RANConfigurationUpdateFailure" {
			var argDefs []string
			for _, ieName := range mInfo.IEorder {
				ieInfo := mInfo.IEs[ieName]
				argDefs = append(argDefs, fmt.Sprintf("%s *%s", ieInfo.GoVar, ieInfo.GoType))
			}
			fmt.Fprintf(fOut, "func rawBuild%s(%s) ([]byte, error) {\n", msgName, strings.Join(argDefs, ","))
			fmt.Fprintln(fOut, "var pdu ngapType.NGAPPDU")
			fmt.Fprintln(fOut, "")
			fmt.Fprintf(fOut, "pdu.Present = ngapType.NGAPPDUPresent%s\n", mInfo.GoField)
			fmt.Fprintf(fOut, "pdu.%s = new(ngapType.%s)\n", mInfo.GoField, mInfo.GoField)
			fmt.Fprintln(fOut, "")
			fmt.Fprintf(fOut, "%s := pdu.%s\n", mInfo.GoTypeVar, mInfo.GoField)
			fmt.Fprintf(fOut, "%s.ProcedureCode.Value = ngapType.ProcedureCode%s\n", mInfo.GoTypeVar, mInfo.ProcCode)
			switch mInfo.Criticality {
			case ngapType.CriticalityPresentReject:
				fmt.Fprintf(fOut, "%s.Criticality.Value = ngapType.CriticalityPresentReject\n", mInfo.GoTypeVar)
			case ngapType.CriticalityPresentIgnore:
				fmt.Fprintf(fOut, "%s.Criticality.Value = ngapType.CriticalityPresentIgnore\n", mInfo.GoTypeVar)
			case ngapType.CriticalityPresentNotify:
				fmt.Fprintf(fOut, "%s.Criticality.Value = ngapType.CriticalityPresentNotify\n", mInfo.GoTypeVar)
			}
			fmt.Fprintln(fOut, "")
			fmt.Fprintf(fOut, "%s.Value.Present = ngapType.%sPresent%s\n", mInfo.GoTypeVar, mInfo.GoField, msgName)
			fmt.Fprintf(fOut, "%s.Value.%s = new(ngapType.%s)\n", mInfo.GoTypeVar, msgName, msgName)
			fmt.Fprintln(fOut, "")
			fmt.Fprintf(fOut, "%s := %s.Value.%s\n", mInfo.GoMsgVar, mInfo.GoTypeVar, msgName)
			fmt.Fprintf(fOut, "%sIEs := &%s.ProtocolIEs\n", mInfo.GoMsgVar, mInfo.GoMsgVar)
			fmt.Fprintf(fOut, "%sIEs.List = make([]ngapType.%sIEs, 0, %d)\n", mInfo.GoMsgVar, msgName, len(mInfo.IEorder))
			for _, ieName := range mInfo.IEorder {
				ieInfo := mInfo.IEs[ieName]
				fmt.Fprintln(fOut, "")
				if ieInfo.Presence != ngapType.PresencePresentMandatory {
					fmt.Fprintf(fOut, "if %s != nil {\n", ieInfo.GoVar)
				} else {
					fmt.Fprintln(fOut, "{")
				}
				fmt.Fprintf(fOut, "ie := ngapType.%sIEs{}\n", msgName)
				fmt.Fprintf(fOut, "ie.Id.Value = %s\n", ieInfo.GoID)
				fmt.Fprintf(fOut, "ie.Criticality.Value = ngapType.CriticalityPresentIgnore\n")
				fmt.Fprintf(fOut, "ie.Value.Present = ngapType.%sIEsPresent%s\n", msgName, convGoName(ieInfo.Type))
				fmt.Fprintf(fOut, "ie.%s = %s\n", ieInfo.GoField, ieInfo.GoVar)
				fmt.Fprintf(fOut, "%sIEs.List = append(%sIEs.List, ie)\n", mInfo.GoMsgVar, mInfo.GoMsgVar)
				fmt.Fprintln(fOut, "}")
			}
			fmt.Fprintln(fOut, "")
			fmt.Fprintln(fOut, "return ngap.Encoder(pdu)")
			fmt.Fprintln(fOut, "}")
			fmt.Fprintln(fOut, "")
		}
	}
}

func generateSender(fOut io.Writer) {
	for _, msgName := range msgNames {
		mInfo := MsgTable[msgName]
		if msgName == "HandoverPreparationFailure" ||
			msgName == "NGSetupFailure" ||
			msgName == "PathSwitchRequestFailure" ||
			msgName == "RANConfigurationUpdateFailure" {
			argDefs := []string{"ran *context.AmfRan"}
			for _, ieName := range mInfo.IEorder {
				ieInfo := mInfo.IEs[ieName]
				if ieInfo.Presence == ngapType.PresencePresentMandatory {
					argDefs = append(argDefs, fmt.Sprintf("%s %s", ieInfo.GoVar, ieInfo.GoType))
				} else {
					argDefs = append(argDefs, fmt.Sprintf("%s *%s", ieInfo.GoVar, ieInfo.GoType))
				}
			}
			fmt.Fprintf(fOut, "func rawSend%s(%s) {\n", msgName, strings.Join(argDefs, ","))
			fmt.Fprintln(fOut, "if ran == nil {")
			fmt.Fprintln(fOut, "logger.NgapLog.Error(\"Ran is nil\")")
			fmt.Fprintln(fOut, "return")
			fmt.Fprintln(fOut, "}")
			fmt.Fprintln(fOut, "")
			fmt.Fprintf(fOut, "ran.Log.Info(\"Send %s\")\n", msgName)
			fmt.Fprintln(fOut, "")
			var args []string
			for _, ieName := range mInfo.IEorder {
				ieInfo := mInfo.IEs[ieName]
				if ieInfo.Presence == ngapType.PresencePresentMandatory {
					args = append(args, fmt.Sprintf("&%s", ieInfo.GoVar))
				} else {
					args = append(args, fmt.Sprintf("%s", ieInfo.GoVar))
				}
			}
			fmt.Fprintf(fOut, "pkt, err := rawBuild%s(%s)\n", msgName, strings.Join(args, ","))
			fmt.Fprintln(fOut, "if err != nil {")
			fmt.Fprintf(fOut, "ran.Log.Errorf(\"Build %s failed : %%s\", err.Error())\n", msgName)
			fmt.Fprintln(fOut, "return")
			fmt.Fprintln(fOut, "}")
			fmt.Fprintln(fOut, "")
			fmt.Fprintln(fOut, "ngap_message.SendToRan(ran, pkt)")
			fmt.Fprintln(fOut, "}")
			fmt.Fprintln(fOut, "")
		}
	}
}

type outputFile struct {
	*bytes.Buffer
	name string
}

func newOutputFile(name string, pkgname string, imports []string) *outputFile {
	o := outputFile{
		Buffer: new(bytes.Buffer),
		name:   name,
	}

	fmt.Fprintln(o, "// Code generated by ngap_generator.go, DO NOT EDIT.")
	fmt.Fprintf(o, "package %s\n", pkgname)
	fmt.Fprintln(o, "")
	fmt.Fprintf(o, "import (\n\n%s\n)\n", strings.Join(imports, "\n"))
	fmt.Fprintln(o, "")

	return &o
}

func (o *outputFile) Close() {
	// Output to file
	if false {
		os.Stdout.Write(o.Bytes())
		return
	}
	out, err := format.Source(o.Bytes())
	if err != nil {
		panic(err)
	}
	fWrite, err := os.Create(o.name)
	if err != nil {
		panic(err)
	}
	defer fWrite.Close()
	_, err = fWrite.Write(out)
	if err != nil {
		panic(err)
	}
}

func genErrorIndicationCommon(f io.Writer, mInfo *MsgInfo) {
	fmt.Fprintf(f, "procedureCode := ngapType.ProcedureCode%s\n", mInfo.ProcCode)
	genTriggeringMessage(f, mInfo.GoField)
	switch mInfo.Criticality {
	case ngapType.CriticalityPresentReject:
		fmt.Fprintln(f, "procedureCriticality := ngapType.CriticalityPresentReject")
	case ngapType.CriticalityPresentIgnore:
		fmt.Fprintln(f, "procedureCriticality := ngapType.CriticalityPresentIgnore")
	case ngapType.CriticalityPresentNotify:
		fmt.Fprintln(f, "procedureCriticality := ngapType.CriticalityPresentNotify")
	}
}

func genTriggeringMessage(f io.Writer, present string) {
	if present == "UnsuccessfulOutcome" {
		fmt.Fprintf(f, "triggeringMessage := ngapType.TriggeringMessagePresentUnsuccessfullOutcome\n")
	} else {
		fmt.Fprintf(f, "triggeringMessage := ngapType.TriggeringMessagePresent%s\n", present)
	}
}

func str2Presence(presence string) aper.Enumerated {
	switch presence {
	case "optional":
		return ngapType.PresencePresentOptional
	case "conditional":
		return ngapType.PresencePresentConditional
	case "mandatory":
		return ngapType.PresencePresentMandatory
	default:
		panic(fmt.Sprintf("Unknown Presence %s", presence))
	}
}

func presence2Str(presence aper.Enumerated) string {
	switch presence {
	case ngapType.PresencePresentOptional:
		return "optional"
	case ngapType.PresencePresentConditional:
		return "conditional"
	case ngapType.PresencePresentMandatory:
		return "mandatory"
	default:
		panic(fmt.Sprintf("Unknown Presence %d", presence))
	}
}

func str2Criticality(criticality string) aper.Enumerated {
	switch criticality {
	case "reject":
		return ngapType.CriticalityPresentReject
	case "ignore":
		return ngapType.CriticalityPresentIgnore
	case "notify":
		return ngapType.CriticalityPresentNotify
	default:
		panic(fmt.Sprintf("Unknown Criticality %s", criticality))
	}
}

func criticality2Str(criticality aper.Enumerated) string {
	switch criticality {
	case ngapType.CriticalityPresentReject:
		return "reject"
	case ngapType.CriticalityPresentIgnore:
		return "ignore"
	case ngapType.CriticalityPresentNotify:
		return "notify"
	default:
		panic(fmt.Sprintf("Unknown Criticality %d", criticality))
	}
}

func getMessageDirection(msgName string) messageDirection {
	switch msgName {
	// PDU Session Management Messages
	case "PDUSessionResourceSetupRequest":
		return messageDirectionAMFtoRAN
	case "PDUSessionResourceSetupResponse":
		return messageDirectionRANtoAMF
	case "PDUSessionResourceReleaseCommand":
		return messageDirectionAMFtoRAN
	case "PDUSessionResourceReleaseResponse":
		return messageDirectionRANtoAMF
	case "PDUSessionResourceModifyRequest":
		return messageDirectionAMFtoRAN
	case "PDUSessionResourceModifyResponse":
		return messageDirectionRANtoAMF
	case "PDUSessionResourceNotify":
		return messageDirectionRANtoAMF
	case "PDUSessionResourceModifyIndication":
		return messageDirectionRANtoAMF
	case "PDUSessionResourceModifyConfirm":
		return messageDirectionAMFtoRAN

	// UE Context Management Messages
	case "InitialContextSetupRequest":
		return messageDirectionAMFtoRAN
	case "InitialContextSetupResponse":
		return messageDirectionRANtoAMF
	case "InitialContextSetupFailure":
		return messageDirectionRANtoAMF
	case "UEContextReleaseRequest":
		return messageDirectionRANtoAMF
	case "UEContextReleaseCommand":
		return messageDirectionAMFtoRAN
	case "UEContextReleaseComplete":
		return messageDirectionRANtoAMF
	case "UEContextModificationRequest":
		return messageDirectionAMFtoRAN
	case "UEContextModificationResponse":
		return messageDirectionRANtoAMF
	case "UEContextModificationFailure":
		return messageDirectionRANtoAMF
	case "RRCInactiveTransitionReport":
		return messageDirectionRANtoAMF

	// UE Mobility Management Messages
	case "HandoverRequired":
		return messageDirectionRANtoAMF
	case "HandoverCommand":
		return messageDirectionAMFtoRAN
	case "HandoverPreparationFailure":
		return messageDirectionAMFtoRAN
	case "HandoverRequest":
		return messageDirectionAMFtoRAN
	case "HandoverRequestAcknowledge":
		return messageDirectionRANtoAMF
	case "HandoverFailure":
		return messageDirectionRANtoAMF
	case "HandoverNotify":
		return messageDirectionRANtoAMF
	case "PathSwitchRequest":
		return messageDirectionRANtoAMF
	case "PathSwitchRequestAcknowledge":
		return messageDirectionAMFtoRAN
	case "PathSwitchRequestFailure":
		return messageDirectionAMFtoRAN
	case "HandoverCancel":
		return messageDirectionRANtoAMF
	case "HandoverCancelAcknowledge":
		return messageDirectionAMFtoRAN
	case "UplinkRANStatusTransfer":
		return messageDirectionRANtoAMF
	case "DownlinkRANStatusTransfer":
		return messageDirectionAMFtoRAN

	// Paging Messages
	case "Paging":
		return messageDirectionAMFtoRAN

	// NAS Transport Messages
	case "InitialUEMessage":
		return messageDirectionRANtoAMF
	case "DownlinkNASTransport":
		return messageDirectionAMFtoRAN
	case "UplinkNASTransport":
		return messageDirectionRANtoAMF
	case "NASNonDeliveryIndication":
		return messageDirectionRANtoAMF
	case "RerouteNASRequest":
		return messageDirectionAMFtoRAN

	// Interface Management Messages
	case "NGSetupRequest":
		return messageDirectionRANtoAMF
	case "NGSetupResponse":
		return messageDirectionAMFtoRAN
	case "NGSetupFailure":
		return messageDirectionAMFtoRAN
	case "RANConfigurationUpdate":
		return messageDirectionRANtoAMF
	case "RANConfigurationUpdateAcknowledge":
		return messageDirectionAMFtoRAN
	case "RANConfigurationUpdateFailure":
		return messageDirectionAMFtoRAN
	case "AMFConfigurationUpdate":
		return messageDirectionAMFtoRAN
	case "AMFConfigurationUpdateAcknowledge":
		return messageDirectionRANtoAMF
	case "AMFConfigurationUpdateFailure":
		return messageDirectionRANtoAMF
	case "AMFStatusIndication":
		return messageDirectionAMFtoRAN
	case "NGReset":
		return messageDirectionBoth
	case "NGResetAcknowledge":
		return messageDirectionBoth
	case "ErrorIndication":
		return messageDirectionBoth
	case "OverloadStart":
		return messageDirectionAMFtoRAN
	case "OverloadStop":
		return messageDirectionAMFtoRAN

	// Configuration Transfer Messages
	case "UplinkRANConfigurationTransfer":
		return messageDirectionRANtoAMF
	case "DownlinkRANConfigurationTransfer":
		return messageDirectionAMFtoRAN

	// Warning Message Transmission Messages
	case "WriteReplaceWarningRequest":
		return messageDirectionAMFtoRAN
	case "WriteReplaceWarningResponse":
		return messageDirectionRANtoAMF
	case "PWSCancelRequest":
		return messageDirectionAMFtoRAN
	case "PWSCancelResponse":
		return messageDirectionRANtoAMF
	case "PWSRestartIndication":
		return messageDirectionRANtoAMF
	case "PWSFailureIndication":
		return messageDirectionRANtoAMF

	// NRPPa Transport Messages
	case "DownlinkUEAssociatedNRPPaTransport":
		return messageDirectionAMFtoRAN
	case "UplinkUEAssociatedNRPPaTransport":
		return messageDirectionRANtoAMF
	case "DownlinkNonUEAssociatedNRPPaTransport":
		return messageDirectionAMFtoRAN
	case "UplinkNonUEAssociatedNRPPaTransport":
		return messageDirectionRANtoAMF

	// Trace Messages
	case "TraceStart":
		return messageDirectionAMFtoRAN
	case "TraceFailureIndication":
		return messageDirectionRANtoAMF
	case "DeactivateTrace":
		return messageDirectionAMFtoRAN
	case "CellTrafficTrace":
		return messageDirectionRANtoAMF

	// Location Reporting Messages
	case "LocationReportingControl":
		return messageDirectionAMFtoRAN
	case "LocationReportingFailureIndication":
		return messageDirectionRANtoAMF
	case "LocationReport":
		return messageDirectionRANtoAMF

	// UE TNLA Binding Messages
	case "UETNLABindingReleaseRequest":
		return messageDirectionAMFtoRAN

	// UE Radio Capability Management Messages
	case "UERadioCapabilityInfoIndication":
		return messageDirectionRANtoAMF
	case "UERadioCapabilityCheckRequest":
		return messageDirectionAMFtoRAN
	case "UERadioCapabilityCheckResponse":
		return messageDirectionRANtoAMF

	// Data Usage Reporting Messages
	case "SecondaryRATDataUsageReport":
		return messageDirectionRANtoAMF

	default:
		panic(fmt.Sprintf("Unknown message %s", msgName))
	}
}

func isAMFtoRANMessage(msgName string) bool {
	return getMessageDirection(msgName)&messageDirectionAMFtoRAN != 0
}

func isRANtoAMFMessage(msgName string) bool {
	return getMessageDirection(msgName)&messageDirectionRANtoAMF != 0
}
