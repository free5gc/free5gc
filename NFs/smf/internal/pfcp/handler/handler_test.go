package handler_test

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/free5gc/pfcp"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/pfcp/pfcpUdp"
	"github.com/free5gc/smf/internal/logger"
	"github.com/free5gc/smf/internal/pfcp/handler"
)

type LogCapture struct {
	buffer bytes.Buffer
}

func (lc *LogCapture) Write(p []byte) (n int, err error) {
	return lc.buffer.Write(p)
}

func (lc *LogCapture) String() string {
	return lc.buffer.String()
}

// func TestHandlePfcpHeartbeatRequest(t *testing.T) {
// }

func TestHandlePfcpManagementRequest(t *testing.T) {
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}Z)(.*)`)
	Convey("Test log message", t, func() {
		remoteAddr := &net.UDPAddr{}
		testPfcpReq := &pfcp.Message{}
		msg := pfcpUdp.NewMessage(remoteAddr, testPfcpReq)
		logCapture := &LogCapture{}
		logger.Log.SetOutput(io.MultiWriter(logCapture, logrus.StandardLogger().Out))
		handler.HandlePfcpPfdManagementRequest(msg)
		capturedLogs := re.FindStringSubmatch(logCapture.String())

		logCaptureExp := &LogCapture{}
		logger.Log.SetOutput(io.MultiWriter(logCaptureExp))
		logger.PfcpLog.Warnf("PFCP PFD Management Request handling is not implemented")
		capturedLogsExp := re.FindStringSubmatch(logCaptureExp.String())
		So(capturedLogs[2], ShouldEqual, capturedLogsExp[2])
	})
}

func TestHandlePfcpAssociationSetupRequest(t *testing.T) {
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}Z)(.*)`)
	re2 := regexp.MustCompile(`(.*)\n(.*)`)
	Convey("Test if NodeID is Nil", t, func() {
		remoteAddr := &net.UDPAddr{
			IP:   net.ParseIP("192.168.1.1"),
			Port: 12345,
		}

		testPfcpReq := &pfcp.Message{
			Header: pfcp.Header{
				Version:         1,
				MP:              0,
				S:               0,
				MessageType:     pfcp.PFCP_ASSOCIATION_SETUP_REQUEST,
				MessageLength:   9,
				SEID:            0,
				SequenceNumber:  1,
				MessagePriority: 0,
			},
			Body: pfcp.PFCPAssociationSetupRequest{
				NodeID: nil,
			},
		}

		logCapture := &LogCapture{}
		logger.Log.SetOutput(io.MultiWriter(logCapture, logrus.StandardLogger().Out))

		msg := pfcpUdp.NewMessage(remoteAddr, testPfcpReq)
		handler.HandlePfcpAssociationSetupRequest(msg)
		capturedLogs := re.FindStringSubmatch(logCapture.String())

		logCaptureExp := &LogCapture{}
		logger.Log.SetOutput(io.MultiWriter(logCaptureExp))
		logger.PfcpLog.Errorln("pfcp association needs NodeID")
		logger.PfcpLog.Infof("Handle PFCP Association Setup Request with NodeID")
		ExpLogs := re.FindStringSubmatch(logCaptureExp.String())
		fmt.Println(ExpLogs)
		if len(capturedLogs) <= 2 || len(ExpLogs) <= 2 {
			t.Errorf("The extracted log is not as expected.")
		}
		So(capturedLogs[2], ShouldEqual, ExpLogs[2])
	})
	Convey("Test if NodeID is NotNil, upf is Nil", t, func() {
		remoteAddr := &net.UDPAddr{
			IP:   net.ParseIP("192.168.1.1"),
			Port: 12345,
		}

		testPfcpReq := &pfcp.Message{
			Header: pfcp.Header{
				Version:         1,
				MP:              0,
				S:               0,
				MessageType:     pfcp.PFCP_ASSOCIATION_SETUP_REQUEST,
				MessageLength:   9,
				SEID:            0,
				SequenceNumber:  1,
				MessagePriority: 0,
			},
			Body: pfcp.PFCPAssociationSetupRequest{
				NodeID: &pfcpType.NodeID{
					NodeIdType: pfcpType.NodeIdTypeIpv4Address,
					IP:         net.ParseIP("192.168.1.1").To4(),
				},
			},
		}

		logCapture := &LogCapture{}
		logger.Log.SetOutput(io.MultiWriter(logCapture, logrus.StandardLogger().Out))

		msg := pfcpUdp.NewMessage(remoteAddr, testPfcpReq)
		handler.HandlePfcpAssociationSetupRequest(msg)
		capturedLogs := re.FindStringSubmatch(re2.FindStringSubmatch(logCapture.String())[1])

		logCaptureExp := &LogCapture{}
		logger.Log.SetOutput(io.MultiWriter(logCaptureExp))
		logger.PfcpLog.Infof("Handle PFCP Association Setup Request with NodeID[%s]", "192.168.1.1")
		logger.PfcpLog.Errorf("can't find UPF[%s]", "192.168.1.1")
		ExpLogs := re.FindStringSubmatch(re2.FindStringSubmatch(logCaptureExp.String())[1])
		if len(capturedLogs) <= 2 || len(ExpLogs) <= 2 {
			t.Errorf("The extracted log is not as expected.")
		}
		So(capturedLogs[2], ShouldEqual, ExpLogs[2])
		capturedLogs = re.FindStringSubmatch(re2.FindStringSubmatch(logCapture.String())[2])
		ExpLogs = re.FindStringSubmatch(re2.FindStringSubmatch(logCaptureExp.String())[2])
		So(capturedLogs[2], ShouldEqual, ExpLogs[2])
	})
}

func TestHandlePfcpAssociationUpdateRequest(t *testing.T) {
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{9}Z)(.*)`)
	Convey("Test logger message", t, func() {
		remoteAddr := &net.UDPAddr{}
		testPfcpReq := &pfcp.Message{}
		msg := pfcpUdp.NewMessage(remoteAddr, testPfcpReq)
		logCapture := &LogCapture{}
		logger.Log.SetOutput(io.MultiWriter(logCapture, logrus.StandardLogger().Out))
		handler.HandlePfcpAssociationUpdateRequest(msg)
		capturedLogs := re.FindStringSubmatch(logCapture.String())

		logCaptureExp := &LogCapture{}
		logger.Log.SetOutput(io.MultiWriter(logCaptureExp))
		logger.PfcpLog.Warnf("PFCP Association Update Request handling is not implemented")
		capturedLogsExp := re.FindStringSubmatch(logCaptureExp.String())
		So(capturedLogs[2], ShouldEqual, capturedLogsExp[2])
	})
}

// func TestHandlePfcpAssociationReleaseRequest(t *testing.T) {
// }
