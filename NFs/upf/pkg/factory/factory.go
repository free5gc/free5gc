package factory

import (
	"net"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/free5gc/go-upf/internal/logger"
)

// TODO: Support configuration update from REST api
func InitConfigFactory(f string, cfg *Config) error {
	if f == "" {
		// Use default config path
		f = UpfDefaultConfigPath
	}

	if content, err := os.ReadFile(f); err != nil {
		return errors.Errorf("[Factory] %+v", err)
	} else {
		logger.CfgLog.Infof("Read config from [%s]", f)
		if yamlErr := yaml.Unmarshal(content, cfg); yamlErr != nil {
			return errors.Errorf("[Factory] %+v", yamlErr)
		}
	}

	return nil
}

func ReadConfig(cfgPath string) (*Config, error) {
	cfg := &Config{}
	err := InitConfigFactory(cfgPath, cfg)
	if err != nil {
		return nil, errors.Errorf("ReadConfig [%s] Error: %+v", cfgPath, err)
	}

	govalidator.TagMap["cidr"] = govalidator.Validator(func(str string) bool {
		return govalidator.IsCIDR(str)
	})
	_, err = govalidator.ValidateStruct(cfg)
	if err != nil {
		logger.CfgLog.Errorf("[-- PLEASE REFER TO SAMPLE CONFIG FILE COMMENTS --]")
		return nil, err
	}

	_, err = net.ResolveIPAddr("ip4", cfg.Pfcp.NodeID)
	if err != nil {
		return nil, errors.Errorf("cfg.Pfcp.NodeID[%s] can't be resolved", cfg.Pfcp.NodeID)
	}

	cfg.Print()
	return cfg, nil
}
