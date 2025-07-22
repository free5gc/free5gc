package pfcp_test

import (
	"testing"
)

// var testAddr *net.UDPAddr

// Adjust waiting time in millisecond if PFCP packets are not captured
// var testWaitingTime int = 500

func init() {
	// smfContext := context.GetSelf()

	// smfContext.CPNodeID.NodeIdType = 0
	// smfContext.CPNodeID.NodeIdValue = net.ParseIP("127.0.0.2").To4()

	// udp.Run()

	// testAddr = &net.UDPAddr{
	// 	IP:   net.ParseIP("127.0.0.2"),
	// 	Port: pfcpUdp.PFCP_PORT,
	// }
}

func TestReliablePFCPResponseDelivery(t *testing.T) {
	// conn, err := net.DialUDP("udp", nil, testAddr)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// pfcpMsg, err := message.BuildPfcpAssociationReleaseRequest()

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// msg := pfcp.Message{
	// 	Header: pfcp.Header{
	// 		Version:        pfcp.PfcpVersion,
	// 		MP:             0,
	// 		S:              pfcp.SEID_NOT_PRESENT,
	// 		MessageType:    pfcp.PFCP_ASSOCIATION_RELEASE_REQUEST,
	// 		SequenceNumber: 1,
	// 	},
	// 	Body: pfcpMsg,
	// }

	// buf, err := msg.Marshal()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// _, err = conn.Write(buf)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// time.Sleep(2 * time.Second)

	// _, err = conn.Write(buf)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// time.Sleep(2 * time.Second)

	// _, err = conn.Write(buf)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// time.Sleep(20 * time.Second)
}
