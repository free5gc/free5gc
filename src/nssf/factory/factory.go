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
var Configured bool

func init() {
	Configured = false
}

func checkErr(err error) {
	if err != nil {
		err = fmt.Errorf("[Configuration] %s", err.Error())
		logger.FactoryLog.Fatal(err)
	}
}

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) {
	content, err := ioutil.ReadFile(f)
	checkErr(err)

	NssfConfig = Config{}

	err = yaml.Unmarshal([]byte(content), &NssfConfig)
	checkErr(err)

	Configured = true

	logger.FactoryLog.Infof("Successfully initialize configuration %s", f)
}
