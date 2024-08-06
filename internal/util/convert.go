package util

import (
	"encoding/json"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/free5gc/openapi/models"
)

func SnssaisToBsonM(snssais string) []bson.M {
	snssaisString := strings.Trim(snssais, "{}")
	slicesStrings := strings.Split(snssaisString, "},{")

	var bsonMArray []bson.M
	for _, slice := range slicesStrings {
		slice = "{" + slice + "}"

		snssaiStruct := &models.Snssai{}
		err := json.Unmarshal([]byte(slice), snssaiStruct)
		if err != nil {
			return nil
		}

		snssaiBsonM := bson.M{}
		snssaiBsonM["sst"] = snssaiStruct.Sst
		if snssaiStruct.Sd != "" {
			snssaiBsonM["sd"] = snssaiStruct.Sd
		}

		bsonMArray = append(bsonMArray, snssaiBsonM)
	}
	return bsonMArray
}
