package udm_consumer_test

import (
	"fmt"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/src/udm/udm_consumer"
	"free5gc/src/udm/udm_context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestRegisterNFInstance(t *testing.T) {

	nrfInit()

	time.Sleep(200 * time.Millisecond)
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})
	time.Sleep(200 * time.Millisecond)

	udm_context.TestInit()
	self := udm_context.UDM_Self()
	NfProfile, err := udm_consumer.BuildNFInstance(self)
	if err != nil {
		t.Error(err.Error())
	}

	uri, _, err1 := udm_consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, NfProfile)
	if err1 != nil {
		t.Error(err1.Error())
	} else {
		fmt.Println("uri: ", uri)
	}
}
