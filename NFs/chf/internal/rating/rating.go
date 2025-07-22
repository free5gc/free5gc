package rating

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/sm/smpeer"

	charging_code "github.com/free5gc/chf/ccs_diameter/code"
	charging_datatype "github.com/free5gc/chf/ccs_diameter/datatype"
	chf_context "github.com/free5gc/chf/internal/context"
	"github.com/free5gc/chf/internal/logger"
	"github.com/free5gc/chf/pkg/factory"
)

func SendServiceUsageRequest(
	ue *chf_context.ChfUe, sur *charging_datatype.ServiceUsageRequest,
) (*charging_datatype.ServiceUsageResponse, error) {
	ue.RatingMux.Handle("SUA", HandleSUA(ue.RatingChan))
	rfDiameter := factory.ChfConfig.Configuration.RfDiameter
	addr := rfDiameter.HostIPv4 + ":" + strconv.Itoa(rfDiameter.Port)
	conn, err := ue.RatingClient.DialNetworkTLS(rfDiameter.Protocol, addr, rfDiameter.Tls.Pem, rfDiameter.Tls.Key)
	if err != nil {
		return nil, err
	}

	meta, ok := smpeer.FromContext(conn.Context())
	if !ok {
		return nil, fmt.Errorf("peer metadata unavailable")
	}

	sur.DestinationRealm = datatype.DiameterIdentity(meta.OriginRealm)
	sur.DestinationHost = datatype.DiameterIdentity(meta.OriginHost)

	msg := diam.NewRequest(charging_code.ServiceUsageMessage, charging_code.Re_interface, dict.Default)
	err = msg.Marshal(sur)
	if err != nil {
		return nil, fmt.Errorf("marshal SUR Failed: %s", err)
	}

	_, err = msg.WriteTo(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to send message from %s: %s",
			conn.RemoteAddr(), err)
	}

	select {
	case m := <-ue.RatingChan:
		var sua charging_datatype.ServiceUsageResponse
		if errMarshal := m.Unmarshal(&sua); errMarshal != nil {
			return nil, fmt.Errorf("failed to parse message from %v", errMarshal)
		}
		return &sua, nil
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout: no rate answer received")
	}
}

func HandleSUA(rgChan chan *diam.Message) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		logger.RatingLog.Tracef("Received SUA from %s", c.RemoteAddr())

		rgChan <- m
	}
}
