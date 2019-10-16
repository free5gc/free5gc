package ausf_consumer_test

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/src/ausf/ausf_consumer"
	"free5gc/src/ausf/ausf_context"
	"testing"
	"time"
)

func TestRegisterNFInstance(t *testing.T) {

	nrfInit()

	time.Sleep(200 * time.Millisecond)
	MongoDBLibrary.RestfulAPIDeleteMany("NfProfile", bson.M{})

	time.Sleep(100 * time.Millisecond)

	ausf_context.TestInit()
	self := ausf_context.GetSelf()
	nfprofile, err := ausf_consumer.BuildNFInstance(self)
	if err != nil {
		t.Error(err.Error())
	}

	uri, err1 := ausf_consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, nfprofile)
	if err1 != nil {
		t.Error(err1.Error())
	} else {
		fmt.Println("uri: ", uri)
	}
}
