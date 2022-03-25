package test

import (
	"encoding/json"

	"github.com/calee0219/fatal"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/mongoapi"
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
	if _, err := mongoapi.RestfulAPIPutOne(collName, filter, putData); err != nil {
		fatal.Fatalf("InsertAuthSubscriptionToMongoDB err: %+v", err)
	}
}

func GetAuthSubscriptionFromMongoDB(ueId string) (authSubs *models.AuthenticationSubscription) {
	collName := "subscriptionData.authenticationData.authenticationSubscription"
	filter := bson.M{"ueId": ueId}
	getData, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		fatal.Fatalf("GetAuthSubscriptionFromMongoDB err: %+v", err)
	}
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
	if err := mongoapi.RestfulAPIDeleteMany(collName, filter); err != nil {
		fatal.Fatalf("DelAuthSubscriptionToMongoDB err: %+v", err)
	}
}

func InsertAccessAndMobilitySubscriptionDataToMongoDB(
	ueId string, amData models.AccessAndMobilitySubscriptionData, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.amData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	putData := toBsonM(amData)
	putData["ueId"] = ueId
	putData["servingPlmnId"] = servingPlmnId
	if _, err := mongoapi.RestfulAPIPutOne(collName, filter, putData); err != nil {
		fatal.Fatalf("InsertAccessAndMobilitySubscriptionDataToMongoDB err: %+v", err)
	}
}

func GetAccessAndMobilitySubscriptionDataFromMongoDB(
	ueId string, servingPlmnId string) (amData *models.AccessAndMobilitySubscriptionData) {
	collName := "subscriptionData.provisionedData.amData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	getData, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		fatal.Fatalf("GetAccessAndMobilitySubscriptionDataFromMongoDB err: %+v", err)
	}
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
	if err := mongoapi.RestfulAPIDeleteMany(collName, filter); err != nil {
		fatal.Fatalf("DelAccessAndMobilitySubscriptionDataFromMongoDB err: %+v", err)
	}
}

func InsertSessionManagementSubscriptionDataToMongoDB(
	ueId string, servingPlmnId string, smDatas []models.SessionManagementSubscriptionData) {
	var putDatas = make([]interface{}, 0, len(smDatas))
	collName := "subscriptionData.provisionedData.smData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	for _, smData := range smDatas {
		putData := toBsonM(smData)
		putData["ueId"] = ueId
		putData["servingPlmnId"] = servingPlmnId
		putDatas = append(putDatas, putData)
	}
	if err := mongoapi.RestfulAPIPostMany(collName, filter, putDatas); err != nil {
		fatal.Fatalf("InsertSessionManagementSubscriptionDataToMongoDB err: %+v", err)
	}
}

func GetSessionManagementDataFromMongoDB(
	ueId string, servingPlmnId string) (amData *models.SessionManagementSubscriptionData) {
	collName := "subscriptionData.provisionedData.smData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	getData, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		fatal.Fatalf("GetSessionManagementDataFromMongoDB err: %+v", err)
	}
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
	if err := mongoapi.RestfulAPIDeleteMany(collName, filter); err != nil {
		fatal.Fatalf("DelSessionManagementSubscriptionDataFromMongoDB err: %+v", err)
	}
}

func InsertSmfSelectionSubscriptionDataToMongoDB(
	ueId string, smfSelData models.SmfSelectionSubscriptionData, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	putData := toBsonM(smfSelData)
	putData["ueId"] = ueId
	putData["servingPlmnId"] = servingPlmnId
	if _, err := mongoapi.RestfulAPIPutOne(collName, filter, putData); err != nil {
		fatal.Fatalf("InsertSmfSelectionSubscriptionDataToMongoDB err: %+v", err)
	}
}

func GetSmfSelectionSubscriptionDataFromMongoDB(
	ueId string, servingPlmnId string) (smfSelData *models.SmfSelectionSubscriptionData) {
	collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
	getData, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		fatal.Fatalf("GetSmfSelectionSubscriptionDataFromMongoDB err: %+v", err)
	}
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
	if err := mongoapi.RestfulAPIDeleteMany(collName, filter); err != nil {
		fatal.Fatalf("DelSmfSelectionSubscriptionDataFromMongoDB err: %+v", err)
	}
}

func InsertAmPolicyDataToMongoDB(ueId string, amPolicyData models.AmPolicyData) {
	collName := "policyData.ues.amData"
	filter := bson.M{"ueId": ueId}
	putData := toBsonM(amPolicyData)
	putData["ueId"] = ueId
	if _, err := mongoapi.RestfulAPIPutOne(collName, filter, putData); err != nil {
		fatal.Fatalf("InsertAmPolicyDataToMongoDB err: %+v", err)
	}
}

func GetAmPolicyDataFromMongoDB(ueId string) (amPolicyData *models.AmPolicyData) {
	collName := "policyData.ues.amData"
	filter := bson.M{"ueId": ueId}
	getData, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		fatal.Fatalf("GetAmPolicyDataFromMongoDB err: %+v", err)
	}
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
	if err := mongoapi.RestfulAPIDeleteMany(collName, filter); err != nil {
		fatal.Fatalf("DelAmPolicyDataFromMongoDB err: %+v", err)
	}
}

func InsertSmPolicyDataToMongoDB(ueId string, smPolicyData models.SmPolicyData) {
	collName := "policyData.ues.smData"
	filter := bson.M{"ueId": ueId}
	putData := toBsonM(smPolicyData)
	putData["ueId"] = ueId
	if _, err := mongoapi.RestfulAPIPutOne(collName, filter, putData); err != nil {
		fatal.Fatalf("InsertSmPolicyDataToMongoDB err: %+v", err)
	}
}

func GetSmPolicyDataFromMongoDB(ueId string) (smPolicyData *models.SmPolicyData) {
	collName := "policyData.ues.smData"
	filter := bson.M{"ueId": ueId}
	getData, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		fatal.Fatalf("GetSmPolicyDataFromMongoDB err: %+v", err)
	}
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
	if err := mongoapi.RestfulAPIDeleteMany(collName, filter); err != nil {
		fatal.Fatalf("DelSmPolicyDataFromMongoDB err: %+v", err)
	}
}
