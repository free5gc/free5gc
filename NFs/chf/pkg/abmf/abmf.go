// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter server example. This is by no means a complete server.
//
// If you'd like to test diameter over SSL, generate SSL certificates:
//   go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
//
// And start the server with `-cert_file cert.pem -key_file key.pem`.
//
// By default this server runs in a single OS thread. If you want to
// make it run on more, set the GOMAXPROCS=n environment variable.
// See Go's FAQ for details: http://golang.org/doc/faq#Why_no_multi_CPU

package abmf

import (
	"bytes"
	"context"
	"math"
	_ "net/http/pprof"
	"strconv"
	"sync"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/sm"
	"go.mongodb.org/mongo-driver/bson"

	charging_datatype "github.com/free5gc/chf/ccs_diameter/datatype"
	charging_dict "github.com/free5gc/chf/ccs_diameter/dict"
	"github.com/free5gc/chf/internal/logger"
	"github.com/free5gc/chf/pkg/factory"
	"github.com/free5gc/util/mongoapi"
)

const chargingDatasColl = "policyData.ues.chargingData"

func OpenServer(ctx context.Context, wg *sync.WaitGroup) {
	// Load our custom dictionary on top of the default one, which
	// always have the Base Protocol (RFC6733) and Credit Control
	// Application (RFC4006).
	logger.AcctLog.Infof("Open Account Balance Management Server")

	mongodb := factory.ChfConfig.Configuration.Mongodb
	// Connect to MongoDB
	if err := mongoapi.SetMongoDB(mongodb.Name, mongodb.Url); err != nil {
		logger.InitLog.Errorf("InitpcfContext err: %+v", err)
		return
	}

	err := dict.Default.Load(bytes.NewReader([]byte(charging_dict.AbmfDictionary)))
	if err != nil {
		logger.RatingLog.Error(err)
	}
	settings := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity("server"),
		OriginRealm:      datatype.DiameterIdentity("go-diameter"),
		VendorID:         13,
		ProductName:      "go-diameter",
		FirmwareRevision: 1,
	}

	// Create the state machine (mux) and set its message handlers.
	mux := sm.New(settings)
	mux.Handle("CCR", handleCCR())
	mux.HandleFunc("ALL", handleALL) // Catch all.

	// Print error reports.
	go printErrors(mux.ErrorReports())
	go func() {
		defer func() {
			logger.AcctLog.Infoln("ABMF server stopped")
			wg.Done()
		}()
		<-ctx.Done()
	}()
	abmfDiameter := factory.ChfConfig.Configuration.AbmfDiameter
	addr := abmfDiameter.HostIPv4 + ":" + strconv.Itoa(abmfDiameter.Port)
	go func() {
		errListen := diam.ListenAndServeTLS(addr, abmfDiameter.Tls.Pem, abmfDiameter.Tls.Key, mux, nil)
		if errListen != nil {
			logger.AcctLog.Errorf("ABMF server fail to listen: %V", errListen)
		}
	}()
}

func printErrors(ec <-chan *diam.ErrorReport) {
	for err := range ec {
		logger.AcctLog.Errorf("Diam Error Report: %v", err)
	}
}

// Commenting the unused function
// func listen(addr, cert, key string, handler diam.Handler) error {
// 	// Start listening for connections.
// 	if len(cert) > 0 && len(key) > 0 {
// 		logger.AcctLog.Infof("Starting secure diameter server on %s", addr)
// 		return diam.ListenAndServeTLS(addr, cert, key, handler, nil)
// 	}

// 	logger.AcctLog.Infof("Starting diameter server on %s", addr)
// 	return diam.ListenAndServe(addr, handler, nil)
// }

// Diameter Credit-Control-Request
func handleCCR() diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var ccr charging_datatype.AccountDebitRequest
		var cca charging_datatype.AccountDebitResponse
		var subscriberId string
		var creditControl *charging_datatype.MultipleServicesCreditControl

		if err := m.Unmarshal(&ccr); err != nil {
			logger.AcctLog.Errorf("Failed to parse message from %s: %s\n%s",
				c.RemoteAddr(), err, m)
			return
		}

		switch ccr.SubscriptionId.SubscriptionIdType {
		case charging_datatype.END_USER_IMSI:
			subscriberId = "imsi-" + string(ccr.SubscriptionId.SubscriptionIdData)
		}

		mscc := ccr.MultipleServicesCreditControl
		rg := mscc.RatingGroup

		// Retrieve quota into mongoDB
		filter := bson.M{"ueId": subscriberId, "ratingGroup": rg}
		queryStrength := 2
		chargingInterface, err := mongoapi.RestfulAPIGetOne(chargingDatasColl, filter, queryStrength)

		if chargingInterface == nil {
			logger.AcctLog.Errorf("chargingInterface is nil, err: %v", err)
			return
		}

		quotaStr := chargingInterface["quota"].(string)
		quota, err := strconv.ParseInt(quotaStr, 10, 64)
		if err != nil {
			logger.AcctLog.Errorf("srtconv ParseInt error: %+v", err)
			return
		}

		switch ccr.RequestedAction {
		case charging_datatype.CHECK_BALANCE:
			logger.AcctLog.Errorf("CHECK_BALANCE not supported")
		case charging_datatype.PRICE_ENQUIRY:
			logger.AcctLog.Errorf("Should use rating function for PRICE_ENQUIRY")
		case charging_datatype.REFUND_ACCOUNT:
			logger.AcctLog.Infof("Refund Account")
			refundQuota := int64(mscc.RequestedServiceUnit.CCTotalOctets)
			quota += refundQuota
		case charging_datatype.DIRECT_DEBITING:
			switch ccr.CcRequestType {
			case charging_datatype.INITIAL_REQUEST, charging_datatype.UPDATE_REQUEST:
				var finalUnitIndication *charging_datatype.FinalUnitIndication
				requestQuota := int64(mscc.RequestedServiceUnit.CCTotalOctets)
				if requestQuota > quota {
					finalUnitIndication = &charging_datatype.FinalUnitIndication{
						FinalUnitAction: charging_datatype.TERMINATE,
					}

					requestQuota = quota
				}

				creditControl = &charging_datatype.MultipleServicesCreditControl{
					RatingGroup: rg,
					GrantedServiceUnit: &charging_datatype.GrantedServiceUnit{
						CCTotalOctets: datatype.Unsigned64(requestQuota),
					},
					FinalUnitIndication: finalUnitIndication,
				}

				quota -= requestQuota
			case charging_datatype.TERMINATION_REQUEST:
				usedQuota := int64(mscc.UsedServiceUnit.CCTotalOctets)
				quota -= usedQuota
			}

			// Convert quota into value digits and exponential expression
			quotaStr = strconv.FormatInt(quota, 10)
			quotaInt, err1 := strconv.ParseInt(quotaStr, 10, 64)
			if err1 != nil {
				logger.AcctLog.Errorf("srtconv ParseInt error: %+v", err1)
				return
			}

			quotaLen := len(quotaStr)
			quotaExp := quotaLen - 1
			quotaVal := quotaInt / int64(math.Pow10(quotaExp))

			cca = charging_datatype.AccountDebitResponse{
				SessionId:       ccr.SessionId,
				OriginHost:      ccr.DestinationHost,
				OriginRealm:     ccr.DestinationRealm,
				CcRequestType:   ccr.CcRequestType,
				CcRequestNumber: ccr.CcRequestNumber,
				EventTimestamp:  datatype.Time(time.Now()),
				RemainingBalance: &charging_datatype.RemainingBalance{
					UnitValue: &charging_datatype.UnitValue{
						ValueDigits: datatype.Integer64(quotaVal),
						Exponent:    datatype.Integer32(quotaExp),
					},
				},
				MultipleServicesCreditControl: creditControl,
			}
		}

		logger.AcctLog.Infof("UE [%s], Rating group [%d], quota [%d]", subscriberId, rg, quota)

		chargingBsonM := make(bson.M)
		chargingBsonM["quota"] = strconv.FormatInt(quota, 10)
		logger.AcctLog.Warnln("quota:", quota)
		if _, err1 := mongoapi.RestfulAPIPutOne(chargingDatasColl, filter, chargingBsonM); err1 != nil {
			logger.AcctLog.Errorf("RestfulAPIPutOne err: %+v", err1)
		}

		a := m.Answer(diam.Success)

		err = a.Marshal(&cca)
		if err != nil {
			logger.AcctLog.Errorf("Marshal CCA Err: %+v:", err)
		}

		_, err = a.WriteTo(c)
		if err != nil {
			logger.AcctLog.Errorf("Failed to write message to %s: %s\n%s\n",
				c.RemoteAddr(), err, a)
			return
		}
	}
}

func handleALL(c diam.Conn, m *diam.Message) {
	logger.AcctLog.Warnf("Received unexpected message from %s:\n%s", c.RemoteAddr(), m)
}
