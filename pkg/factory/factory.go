/*
 * NRF Configuration Factory
 */

package factory

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/free5gc/nrf/internal/logger"
)

var NrfConfig Config

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	if content, err := ioutil.ReadFile(f); err != nil {
		return err
	} else {
		NrfConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &NrfConfig); yamlErr != nil {
			return yamlErr
		}
	}

	return nil
}

func CheckConfigVersion() error {
	currentVersion := NrfConfig.GetVersion()

	if currentVersion != NrfExpectedConfigVersion {
		return fmt.Errorf("config version is [%s], but expected is [%s].",
			currentVersion, NrfExpectedConfigVersion)
	}

	logger.CfgLog.Infof("config version [%s]", currentVersion)

	return nil
}
