package test

import (
	"encoding/json"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/openapi/models"

	"github.com/calee0219/fatal"
	"go.mongodb.org/mongo-driver/bson"
)

func toBsonM(data interface{}) bson.M {
	tmp, err := json.Marshal(data)
	if err != nil {
		fatal.Fatalf("Marshal error in toBsonM: %+v", err)
	}
	var putData = bson.M{}
	err = json.Unmarshal(tmp, &putData)
	if err != nil {
		fatal.Fatalf("Unmarshal error in toBsonM: %+v", err)
	}
	return putData
}

func InsertAuthSubscriptionToMongoDB(ueId string, authSubs models.AuthenticationSubscription) {
	collName := "subscriptionData.authenticationData.authenticationSubscription"
	filter := bson.M{"ueId": ueId}
	putData := toBsonM(authSubs)
	putData["ueId"] = ueId
	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func GetAuthSubscriptionFromMongoDB(ueId string) (authSubs *models.AuthenticationSubscription) {
	collName := "subscriptionData.authenticationData.authenticationSubscription"
	filter := bson.M{"ueId": ueId}
	getData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
	if getData == nil {
		return
	}
	tmp, err := json.Marshal(getData)
	if err != nil {
		return
	}
	authSubs = new(models.AuthenticationSubscription)
	err = json.Unmarshal(tmp, authSubs)
	if err != nil {
		fatal.Fatalf("Unmarshal error in GetAuthSubscriptionFromMongoDB: %+v", err)
	}
	return
}

func DelAuthSubscriptionToMongoDB(ueId string) {
	collName := "subscriptionData.authenticationData.authenticationSubscription"
	filter := bson.M{"ueId": ueId}
	MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
}

func InsertAccessAndMobilitySubscriptionDataToMongoDB(
	ueId string, amData models.AccessAndMobilitySubscriptionData, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.amData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	putData := toBsonM(amData)
	putData["ueId"] = ueId
	putData["servingPlmnId"] = servingPlmnId
	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func GetAccessAndMobilitySubscriptionDataFromMongoDB(
	ueId string, servingPlmnId string) (amData *models.AccessAndMobilitySubscriptionData) {
	collName := "subscriptionData.provisionedData.amData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	getData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
	if getData == nil {
		return
	}
	tmp, err := json.Marshal(getData)
	if err != nil {
		return
	}
	amData = new(models.AccessAndMobilitySubscriptionData)
	err = json.Unmarshal(tmp, amData)
	if err != nil {
		fatal.Fatalf("Unmarshal error in GetAccessAndMobilitySubscriptionDataFromMongoDB: %+v", err)
	}
	return
}

func DelAccessAndMobilitySubscriptionDataFromMongoDB(ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.amData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
}

func InsertSessionManagementSubscriptionDataToMongoDB(
	ueId string, servingPlmnId string, smData models.SessionManagementSubscriptionData) {
	collName := "subscriptionData.provisionedData.smData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	putData := toBsonM(smData)
	putData["ueId"] = ueId
	putData["servingPlmnId"] = servingPlmnId
	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func GetSessionManagementDataFromMongoDB(
	ueId string, servingPlmnId string) (amData *models.SessionManagementSubscriptionData) {
	collName := "subscriptionData.provisionedData.smData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	getData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
	if getData == nil {
		return
	}
	tmp, err := json.Marshal(getData)
	if err != nil {
		return
	}
	amData = new(models.SessionManagementSubscriptionData)
	err = json.Unmarshal(tmp, amData)
	if err != nil {
		fatal.Fatalf("Unmarshal error in GetSessionManagementDataFromMongoDB: %+v", err)
	}
	return
}

func DelSessionManagementSubscriptionDataFromMongoDB(ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.smData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
}

func InsertSmfSelectionSubscriptionDataToMongoDB(
	ueId string, smfSelData models.SmfSelectionSubscriptionData, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	putData := toBsonM(smfSelData)
	putData["ueId"] = ueId
	putData["servingPlmnId"] = servingPlmnId
	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func GetSmfSelectionSubscriptionDataFromMongoDB(
	ueId string, servingPlmnId string) (smfSelData *models.SmfSelectionSubscriptionData) {
	collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	getData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
	if getData == nil {
		return
	}
	tmp, err := json.Marshal(getData)
	if err != nil {
		return
	}
	smfSelData = new(models.SmfSelectionSubscriptionData)
	err = json.Unmarshal(tmp, smfSelData)
	if err != nil {
		fatal.Fatalf("Unmarshal error in GetSmfSelectionSubscriptionDataFromMongoDB: %+v", err)
	}
	return
}

func DelSmfSelectionSubscriptionDataFromMongoDB(ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
}

func InsertAmPolicyDataToMongoDB(ueId string, amPolicyData models.AmPolicyData) {
	collName := "policyData.ues.amData"
	filter := bson.M{"ueId": ueId}
	putData := toBsonM(amPolicyData)
	putData["ueId"] = ueId
	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func GetAmPolicyDataFromMongoDB(ueId string) (amPolicyData *models.AmPolicyData) {
	collName := "policyData.ues.amData"
	filter := bson.M{"ueId": ueId}
	getData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
	if getData == nil {
		return
	}
	tmp, err := json.Marshal(getData)
	if err != nil {
		return
	}
	amPolicyData = new(models.AmPolicyData)
	err = json.Unmarshal(tmp, amPolicyData)
	if err != nil {
		fatal.Fatalf("Unmarshal error in GetAmPolicyDataFromMongoDB: %+v", err)
	}
	return
}

func DelAmPolicyDataFromMongoDB(ueId string) {
	collName := "policyData.ues.amData"
	filter := bson.M{"ueId": ueId}
	MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
}

func InsertSmPolicyDataToMongoDB(ueId string, smPolicyData models.SmPolicyData) {
	collName := "policyData.ues.smData"
	filter := bson.M{"ueId": ueId}
	putData := toBsonM(smPolicyData)
	putData["ueId"] = ueId
	MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func GetSmPolicyDataFromMongoDB(ueId string) (smPolicyData *models.SmPolicyData) {
	collName := "policyData.ues.smData"
	filter := bson.M{"ueId": ueId}
	getData := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
	if getData == nil {
		return
	}
	tmp, err := json.Marshal(getData)
	if err != nil {
		return
	}
	smPolicyData = new(models.SmPolicyData)
	err = json.Unmarshal(tmp, smPolicyData)
	if err != nil {
		fatal.Fatalf("Unmarshal error in GetSmPolicyDataFromMongoDB: %+v", err)
	}
	return
}

func DelSmPolicyDataFromMongoDB(ueId string) {
	collName := "policyData.ues.smData"
	filter := bson.M{"ueId": ueId}
	MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
}
