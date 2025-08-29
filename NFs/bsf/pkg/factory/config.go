/*
 * BSF Configuration Factory
 */

package factory

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/free5gc/bsf/internal/logger"
)

const (
	BsfDefaultTLSKeyLogPath  = "./log/bsfsslkey.log"
	BsfDefaultConfigPath     = "./config/bsfcfg.yaml"
	BsfExpectedConfigVersion = "1.5.0"
)

type Config struct {
	Info          *Info          `yaml:"info"`
	Configuration *Configuration `yaml:"configuration"`
	Logger        *Logger        `yaml:"logger"`
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type Configuration struct {
	BsfName string   `yaml:"bsfName,omitempty"`
	Sbi     *Sbi     `yaml:"sbi,omitempty"`
	NrfUri  string   `yaml:"nrfUri,omitempty"`
	MongoDB *MongoDB `yaml:"mongodb,omitempty"`
}

type Logger struct {
	BSF            *LogSetting `yaml:"BSF,omitempty"`
	OpenApi        *LogSetting `yaml:"OpenApi,omitempty"`
	MongoDBLibrary *LogSetting `yaml:"MongoDBLibrary,omitempty"`
}

type LogSetting struct {
	DebugLevel   string `yaml:"debugLevel,omitempty"`
	ReportCaller bool   `yaml:"ReportCaller,omitempty"`
}

type Sbi struct {
	Scheme       string `yaml:"scheme"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty"`
	RegisterIPv6 string `yaml:"registerIPv6,omitempty"`
	BindingIPv4  string `yaml:"bindingIPv4,omitempty"`
	BindingIPv6  string `yaml:"bindingIPv6,omitempty"`
	Port         int    `yaml:"port,omitempty"`
	Tls          *Tls   `yaml:"tls,omitempty"`
}

type Tls struct {
	Pem string `yaml:"pem,omitempty"`
	Key string `yaml:"key,omitempty"`
}

type MongoDB struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

var BsfConfig *Config

func ReadConfig(cfgPath string) (*Config, error) {
	if cfgPath == "" {
		cfgPath = BsfDefaultConfigPath
	}

	if content, err := os.ReadFile(cfgPath); err != nil {
		return nil, fmt.Errorf("[Configuration] %s: %+v", cfgPath, err)
	} else {
		logger.CfgLog.Infof("Read config from [%s]", cfgPath)
		if yamlErr := yaml.Unmarshal(content, &BsfConfig); yamlErr != nil {
			return nil, fmt.Errorf("[Configuration] %s: %+v", cfgPath, yamlErr)
		}
	}

	return BsfConfig, nil
}

func CheckConfigVersion() error {
	currentVersion := BsfConfig.Info.Version
	if currentVersion != BsfExpectedConfigVersion {
		return fmt.Errorf("BSF config version [%s] is not supported; expected config version [%s].",
			currentVersion, BsfExpectedConfigVersion)
	}

	logger.CfgLog.Infof("BSF config version [%s]", currentVersion)

	return nil
}
