/*
 * BSF Configuration Factory
 */

package factory

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/free5gc/bsf/internal/logger"
)

const (
	BsfDefaultTLSKeyLogPath  = "./log/bsfsslkey.log"
	BsfDefaultConfigPath     = "./config/bsfcfg.yaml"
	BsfExpectedConfigVersion = "1.0.0"

	// Metrics defaults
	BsfMetricsDefaultEnabled   = false
	BsfMetricsDefaultScheme    = "https"
	BsfMetricsDefaultPort      = 9091
	BsfMetricsDefaultNamespace = "free5gc"
	BsfMetricsDefaultKeyPath   = "./cert/bsf.key"
	BsfMetricsDefaultPemPath   = "./cert/bsf.pem"
)

type Config struct {
	Info          *Info          `yaml:"info"`
	Configuration *Configuration `yaml:"configuration"`
	Logger        *Logger        `yaml:"logger"`
	sync.RWMutex
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
	Metrics *Metrics `yaml:"metrics,omitempty"`
}

type Logger struct {
	Enable       bool   `yaml:"enable" valid:"type(bool)"`
	Level        string `yaml:"level" valid:"required,in(trace|debug|info|warn|error|fatal|panic)"`
	ReportCaller bool   `yaml:"reportCaller" valid:"type(bool)"`
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

type Metrics struct {
	Enable      bool   `yaml:"enable,omitempty"`
	Scheme      string `yaml:"scheme,omitempty"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty"`
	Port        int    `yaml:"port,omitempty"`
	Namespace   string `yaml:"namespace,omitempty"`
	Tls         *Tls   `yaml:"tls,omitempty"`
}

var BsfConfig *Config

// Metrics configuration helper methods
func (c *Config) AreMetricsEnabled() bool {
	if c.Configuration != nil && c.Configuration.Metrics != nil {
		return c.Configuration.Metrics.Enable
	}
	return BsfMetricsDefaultEnabled
}

func (c *Config) GetMetricsScheme() string {
	if c.Configuration != nil && c.Configuration.Metrics != nil && c.Configuration.Metrics.Scheme != "" {
		return c.Configuration.Metrics.Scheme
	}
	return BsfMetricsDefaultScheme
}

func (c *Config) GetMetricsPort() int {
	if c.Configuration != nil && c.Configuration.Metrics != nil && c.Configuration.Metrics.Port != 0 {
		return c.Configuration.Metrics.Port
	}
	return BsfMetricsDefaultPort
}

func (c *Config) GetMetricsBindingAddr() string {
	ip := "0.0.0.0"
	if c.Configuration != nil && c.Configuration.Metrics != nil && c.Configuration.Metrics.BindingIPv4 != "" {
		ip = c.Configuration.Metrics.BindingIPv4
	}
	port := c.GetMetricsPort()
	return fmt.Sprintf("%s:%d", ip, port)
}

func (c *Config) GetMetricsNamespace() string {
	if c.Configuration != nil && c.Configuration.Metrics != nil && c.Configuration.Metrics.Namespace != "" {
		return c.Configuration.Metrics.Namespace
	}
	return BsfMetricsDefaultNamespace
}

func (c *Config) GetMetricsCertKeyPath() string {
	if c.Configuration != nil && c.Configuration.Metrics != nil &&
		c.Configuration.Metrics.Tls != nil && c.Configuration.Metrics.Tls.Key != "" {
		return c.Configuration.Metrics.Tls.Key
	}
	return BsfMetricsDefaultKeyPath
}

func (c *Config) GetMetricsCertPemPath() string {
	if c.Configuration != nil && c.Configuration.Metrics != nil &&
		c.Configuration.Metrics.Tls != nil && c.Configuration.Metrics.Tls.Pem != "" {
		return c.Configuration.Metrics.Tls.Pem
	}
	return BsfMetricsDefaultPemPath
}

func (c *Config) SetLogEnable(enable bool) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Enable: enable,
			Level:  "info",
		}
	} else {
		c.Logger.Enable = enable
	}
}

func (c *Config) SetLogLevel(level string) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Level: level,
		}
	} else {
		c.Logger.Level = level
	}
}

func (c *Config) SetLogReportCaller(reportCaller bool) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Level:        "info",
			ReportCaller: reportCaller,
		}
	} else {
		c.Logger.ReportCaller = reportCaller
	}
}

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
		return fmt.Errorf("BSF config version [%s] is not supported; expected config version [%s]",
			currentVersion, BsfExpectedConfigVersion)
	}

	logger.CfgLog.Infof("BSF config version [%s]", currentVersion)

	return nil
}
