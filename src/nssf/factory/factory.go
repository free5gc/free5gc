/*
 * NSSF Configuration Factory
 */

package factory

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"free5gc/src/nssf/logger"
)

var NssfConfig Config

func checkErr(err error) {
	if err != nil {
		err = fmt.Errorf("[Configuration] %s", err.Error())
		logger.Factory.Fatal(err)
	}
}

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) {
	content, err := ioutil.ReadFile(f)
	checkErr(err)

	NssfConfig = Config{}

	err = yaml.Unmarshal([]byte(content), &NssfConfig)
	checkErr(err)

	logger.Factory.Infof("Successfully initialize configuration %s", f)
}
