package context_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/pfcp/pfcpType"
	smf_context "github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/pkg/factory"
)

var mockIPv4NodeID = &pfcpType.NodeID{
	NodeIdType: pfcpType.NodeIdTypeIpv4Address,
	IP:         net.ParseIP("127.0.0.1"),
}

var mockIfaces = []*factory.InterfaceUpfInfoItem{
	{
		InterfaceType:    "N3",
		Endpoints:        []string{"127.0.0.1"},
		NetworkInstances: []string{"internet"},
	},
}

func convertPDUSessTypeToString(pduType uint8) string {
	switch pduType {
	case nasMessage.PDUSessionTypeIPv4:
		return "PDU Session Type IPv4"
	case nasMessage.PDUSessionTypeIPv6:
		return "PDU Session Type IPv6"
	case nasMessage.PDUSessionTypeIPv4IPv6:
		return "PDU Session Type IPv4 IPv6"
	case nasMessage.PDUSessionTypeUnstructured:
		return "PDU Session Type Unstructured"
	case nasMessage.PDUSessionTypeEthernet:
		return "PDU Session Type Ethernet"
	}

	return "Unkwown PDU Session Type"
}

func TestIP(t *testing.T) {
	testCases := []struct {
		input               *smf_context.UPFInterfaceInfo
		inputPDUSessionType uint8
		paramStr            string
		resultStr           string
		expectedIP          string
		expectedError       error
	}{
		{
			input: &smf_context.UPFInterfaceInfo{
				NetworkInstances:      []string{""},
				IPv4EndPointAddresses: []net.IP{net.ParseIP("8.8.8.8")},
				IPv6EndPointAddresses: []net.IP{net.ParseIP("2001:4860:4860::8888")},
				EndpointFQDN:          "www.google.com",
			},
			inputPDUSessionType: nasMessage.PDUSessionTypeIPv4,
			paramStr:            "select " + convertPDUSessTypeToString(nasMessage.PDUSessionTypeIPv4),
			expectedIP:          "8.8.8.8",
			expectedError:       nil,
		},
		{
			input: &smf_context.UPFInterfaceInfo{
				NetworkInstances:      []string{""},
				IPv4EndPointAddresses: []net.IP{net.ParseIP("8.8.8.8")},
				IPv6EndPointAddresses: []net.IP{net.ParseIP("2001:4860:4860::8888")},
				EndpointFQDN:          "www.google.com",
			},
			inputPDUSessionType: nasMessage.PDUSessionTypeIPv6,
			paramStr:            "select " + convertPDUSessTypeToString(nasMessage.PDUSessionTypeIPv6),
			expectedIP:          "2001:4860:4860::8888",
			expectedError:       nil,
		},
	}

	Convey("Given UPFInterfaceInfo and select PDU Session type, should return correct IP", t, func() {
		for i, testcase := range testCases {
			upfInterfaceInfo := testcase.input
			infoStr := fmt.Sprintf("testcase[%d] UPF Interface Info: %+v", i, upfInterfaceInfo)

			Convey(infoStr, func() {
				Convey(testcase.paramStr, func() {
					ip, err := upfInterfaceInfo.IP(testcase.inputPDUSessionType)
					testcase.resultStr = "IP addr should be " + testcase.expectedIP

					Convey(testcase.resultStr, func() {
						So(ip.String(), ShouldEqual, testcase.expectedIP)
						So(err, ShouldEqual, testcase.expectedError)
					})
				})
			})
		}
	})
}

func TestAddDataPath(t *testing.T) {
	// AddDataPath is simple, should only have one case
	testCases := []struct {
		tunnel        *smf_context.UPTunnel
		addedDataPath *smf_context.DataPath
		resultStr     string
		expectedExist bool
	}{
		{
			tunnel:        smf_context.NewUPTunnel(),
			addedDataPath: smf_context.NewDataPath(),
			resultStr:     "Datapath should exist",
			expectedExist: true,
		},
	}

	Convey("AddDataPath should indeed add datapath", t, func() {
		for i, testcase := range testCases {
			upTunnel := testcase.tunnel
			infoStr := fmt.Sprintf("testcase[%d]: Add Datapath", i)

			Convey(infoStr, func() {
				upTunnel.AddDataPath(testcase.addedDataPath)

				Convey(testcase.resultStr, func() {
					var exist bool
					for _, datapath := range upTunnel.DataPathPool {
						if datapath == testcase.addedDataPath {
							exist = true
						}
					}
					So(exist, ShouldEqual, testcase.expectedExist)
				})
			})
		}
	})
}

func TestAddPDR(t *testing.T) {
	testCases := []struct {
		upf           *smf_context.UPF
		resultStr     string
		expectedError error
	}{
		{
			upf:           smf_context.NewUPF(mockIPv4NodeID, mockIfaces),
			resultStr:     "AddPDR should success",
			expectedError: nil,
		},
		{
			upf:           smf_context.NewUPF(mockIPv4NodeID, mockIfaces),
			resultStr:     "AddPDR should fail",
			expectedError: fmt.Errorf("UPF[127.0.0.1] not associated with SMF"),
		},
	}

	testCases[0].upf.AssociationContext = context.Background()

	Convey("AddPDR should indeed add PDR and report error appropiately", t, func() {
		for i, testcase := range testCases {
			upf := testcase.upf
			infoStr := fmt.Sprintf("testcase[%d]: ", i)

			Convey(infoStr, func() {
				_, err := upf.AddPDR()

				Convey(testcase.resultStr, func() {
					if testcase.expectedError == nil {
						So(err, ShouldBeNil)
					} else {
						So(err, ShouldNotBeNil)
						if err != nil {
							So(err.Error(), ShouldEqual, testcase.expectedError.Error())
						}
					}
				})
			})
		}
	})
}

func TestAddFAR(t *testing.T) {
	testCases := []struct {
		upf           *smf_context.UPF
		resultStr     string
		expectedError error
	}{
		{
			upf:           smf_context.NewUPF(mockIPv4NodeID, mockIfaces),
			resultStr:     "AddFAR should success",
			expectedError: nil,
		},
		{
			upf:           smf_context.NewUPF(mockIPv4NodeID, mockIfaces),
			resultStr:     "AddFAR should fail",
			expectedError: fmt.Errorf("UPF[127.0.0.1] not associated with SMF"),
		},
	}

	testCases[0].upf.AssociationContext = context.Background()

	Convey("AddFAR should indeed add FAR and report error appropiately", t, func() {
		for i, testcase := range testCases {
			upf := testcase.upf
			infoStr := fmt.Sprintf("testcase[%d]: ", i)

			Convey(infoStr, func() {
				_, err := upf.AddFAR()

				Convey(testcase.resultStr, func() {
					if testcase.expectedError == nil {
						So(err, ShouldBeNil)
					} else {
						So(err, ShouldNotBeNil)
						if err != nil {
							So(err.Error(), ShouldEqual, testcase.expectedError.Error())
						}
					}
				})
			})
		}
	})
}

func TestAddQER(t *testing.T) {
	testCases := []struct {
		upf           *smf_context.UPF
		resultStr     string
		expectedError error
	}{
		{
			upf:           smf_context.NewUPF(mockIPv4NodeID, mockIfaces),
			resultStr:     "AddQER should success",
			expectedError: nil,
		},
		{
			upf:           smf_context.NewUPF(mockIPv4NodeID, mockIfaces),
			resultStr:     "AddQER should fail",
			expectedError: fmt.Errorf("UPF[127.0.0.1] not associated with SMF"),
		},
	}

	testCases[0].upf.AssociationContext = context.Background()

	Convey("AddQER should indeed add QER and report error appropiately", t, func() {
		for i, testcase := range testCases {
			upf := testcase.upf
			infoStr := fmt.Sprintf("testcase[%d]: ", i)

			Convey(infoStr, func() {
				_, err := upf.AddQER()

				Convey(testcase.resultStr, func() {
					if testcase.expectedError == nil {
						So(err, ShouldBeNil)
					} else {
						So(err, ShouldNotBeNil)
						if err != nil {
							So(err.Error(), ShouldEqual, testcase.expectedError.Error())
						}
					}
				})
			})
		}
	})
}

func TestAddBAR(t *testing.T) {
	testCases := []struct {
		upf           *smf_context.UPF
		resultStr     string
		expectedError error
	}{
		{
			upf:           smf_context.NewUPF(mockIPv4NodeID, mockIfaces),
			resultStr:     "AddBAR should success",
			expectedError: nil,
		},
		{
			upf:           smf_context.NewUPF(mockIPv4NodeID, mockIfaces),
			resultStr:     "AddBAR should fail",
			expectedError: fmt.Errorf("UPF[127.0.0.1] not associated with SMF"),
		},
	}

	testCases[0].upf.AssociationContext = context.Background()

	Convey("AddBAR should indeed add BAR and report error appropiately", t, func() {
		for i, testcase := range testCases {
			upf := testcase.upf
			infoStr := fmt.Sprintf("testcase[%d]: ", i)

			Convey(infoStr, func() {
				_, err := upf.AddBAR()

				Convey(testcase.resultStr, func() {
					if testcase.expectedError == nil {
						So(err, ShouldBeNil)
					} else {
						So(err, ShouldNotBeNil)
						if err != nil {
							So(err.Error(), ShouldEqual, testcase.expectedError.Error())
						}
					}
				})
			})
		}
	})
}
