//go:binary-only-package

package MongoDBLibrary

import (
	"context"
	"encoding/json"
	"free5gc/lib/MongoDBLibrary/logger"
	"free5gc/lib/openapi/models"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = nil
var dbName string

func SetMongoDB(setdbName string, url string) {}

func RestfulAPIGetOne(collName string, filter bson.M) map[string]interface{} {}

func RestfulAPIGetMany(collName string, filter bson.M) []map[string]interface{} {}

func RestfulAPIPutOne(collName string, filter bson.M, putData map[string]interface{}) bool {}

func RestfulAPIPutOneNotUpdate(collName string, filter bson.M, putData map[string]interface{}) bool {}

func RestfulAPIPutMany(collName string, filterArray []bson.M, putDataArray []map[string]interface{}) bool {}

func RestfulAPIDeleteOne(collName string, filter bson.M) {}

func RestfulAPIDeleteMany(collName string, filter bson.M) {}

func RestfulAPIMergePatch(collName string, filter bson.M, patchData map[string]interface{}) models.ProblemDetails {}

func RestfulAPIJSONPatch(collName string, filter bson.M, patchJSON []byte) bool {}

func RestfulAPIJSONPatchExtend(collName string, filter bson.M, patchJSON []byte, dataName string) bool {}

func RestfulAPIPost(collName string, filter bson.M, postData map[string]interface{}) bool {}
