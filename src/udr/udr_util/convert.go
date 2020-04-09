package udr_util

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"free5gc/lib/openapi/models"
	"strconv"
)

func MapToByte(data map[string]interface{}) (ret []byte) {
	ret, _ = json.Marshal(data)
	return
}

func MapArrayToByte(data []map[string]interface{}) (ret []byte) {
	ret, _ = json.Marshal(data)
	return
}

func ToBsonM(data interface{}) bson.M {
	tmp, _ := json.Marshal(data)
	var putData = bson.M{}
	_ = json.Unmarshal(tmp, &putData)
	return putData
}

func SnssaiHexToModels(hexString string) (*models.Snssai, error) {
	sst, err := strconv.ParseInt(hexString[:2], 16, 32)
	if err != nil {
		return nil, err
	}
	sNssai := &models.Snssai{
		Sst: int32(sst),
		Sd:  hexString[2:],
	}
	return sNssai, nil
}

func SnssaiModelsToHex(snssai models.Snssai) string {
	sst := fmt.Sprintf("%02x", snssai.Sst)
	return sst + snssai.Sd
}
