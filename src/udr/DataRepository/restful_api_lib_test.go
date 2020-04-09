package DataRepository_test

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"free5gc/lib/MongoDBLibrary"
)

var Client *mongo.Client = nil

func SetMongoDB(setdbName string, url string) {
	MongoDBLibrary.SetMongoDB(setdbName, url)
}

func RestfulAPIGetOne(collName string, filter bson.M) map[string]interface{} {
	return MongoDBLibrary.RestfulAPIGetOne(collName, filter)
}

func RestfulAPIGetMany(collName string, filter bson.M) []map[string]interface{} {
	return MongoDBLibrary.RestfulAPIGetMany(collName, filter)
}

func RestfulAPIPutOne(collName string, filter bson.M, putData map[string]interface{}) bool {
	return MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData)
}

func RestfulAPIPutOneNotUpdate(collName string, filter bson.M, putData map[string]interface{}) bool {
	return MongoDBLibrary.RestfulAPIPutOneNotUpdate(collName, filter, putData)
}

func RestfulAPIPutMany(collName string, filterArray []bson.M, putDataArray []map[string]interface{}) bool {
	return MongoDBLibrary.RestfulAPIPutMany(collName, filterArray, putDataArray)
}

func RestfulAPIDeleteOne(collName string, filter bson.M) {
	MongoDBLibrary.RestfulAPIDeleteOne(collName, filter)
}

func RestfulAPIDeleteMany(collName string, filter bson.M) {
	MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
}

func RestfulAPIMergePatch(collName string, filter bson.M, patchData map[string]interface{}) bool {
	return MongoDBLibrary.RestfulAPIMergePatch(collName, filter, patchData)
}

func RestfulAPIJSONPatch(collName string, filter bson.M, patchJSON []byte) bool {
	return MongoDBLibrary.RestfulAPIJSONPatch(collName, filter, patchJSON)
}
func RestfulAPIJSONPatchExtend(collName string, filter bson.M, patchJSON []byte, dataName string) bool {
	return MongoDBLibrary.RestfulAPIJSONPatchExtend(collName, filter, patchJSON, dataName)
}

func RestfulAPIPost(collName string, filter bson.M, postData map[string]interface{}) bool {
	return MongoDBLibrary.RestfulAPIPost(collName, filter, postData)
}

func toBsonM(data interface{}) bson.M {
	tmp, _ := json.Marshal(data)
	var putData = bson.M{}
	json.Unmarshal(tmp, &putData)
	return putData
}
