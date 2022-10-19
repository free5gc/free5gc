/*
 * NRF Configuration Factory
 */

package factory

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/asaskevich/govalidator"

	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/openapi/models"
	logger_util "github.com/free5gc/util/logger"
)

const (
	NrfDefaultCertPemPath        = "./config/cert/nrf.pem"
	NrfDefaultPrivateKeyPath     = "./config/cert/nrf.key"
	NrfDefaultRootCertPemPath    = "./config/cert/root.pem"
	NrfDefaultRootPrivateKeyPath = "./config/cert/root.key"
	NrfExpectedConfigVersion     = "1.0.2"
	NrfSbiDefaultIPv4            = "127.0.0.10"
	NrfSbiDefaultPort            = 8000
	NrfSbiDefaultScheme          = "https"
	NrfNfmResUriPrefix           = "/nnrf-nfm/v1"
	NrfDiscResUriPrefix          = "/nnrf-disc/v1"
)

type Config struct {
	Info          *Info               `yaml:"info" valid:"required"`
	Configuration *Configuration      `yaml:"configuration" valid:"required"`
	Logger        *logger_util.Logger `yaml:"logger" valid:"required"`
}

func (c *Config) Validate() (bool, error) {
	if info := c.Info; info != nil {
		if result, err := info.validate(); err != nil {
			return result, err
		}
	}

	if configuration := c.Configuration; configuration != nil {
		if result, err := configuration.validate(); err != nil {
			return result, err
		}
	}

	if logger := c.Logger; logger != nil {
		if result, err := logger.Validate(); err != nil {
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

type Info struct {
	Version     string `yaml:"version,omitempty" valid:"type(string)"`
	Description string `yaml:"description,omitempty" valid:"type(string)"`
}

func (i *Info) validate() (bool, error) {
	result, err := govalidator.ValidateStruct(i)
	return result, appendInvalid(err)
}

type Configuration struct {
	Sbi             *Sbi          `yaml:"sbi,omitempty" valid:"required"`
	MongoDBName     string        `yaml:"MongoDBName" valid:"type(string),required"`
	MongoDBUrl      string        `yaml:"MongoDBUrl" valid:"required"`
	DefaultPlmnId   models.PlmnId `yaml:"DefaultPlmnId" valid:"required"`
	ServiceNameList []string      `yaml:"serviceNameList,omitempty" valid:"required"`
}

func (c *Configuration) validate() (bool, error) {
	if sbi := c.Sbi; sbi != nil {
		if result, err := sbi.validate(); err != nil {
			return result, err
		}
	}

	defaultPlmnId := c.DefaultPlmnId
	if result := govalidator.StringMatches(defaultPlmnId.Mcc, "^[0-9]{3}$"); !result {
		err := errors.New("Invalid mcc: " + defaultPlmnId.Mcc + ", should be 3 digits string, digit: 0~9")
		return false, err
	}
	if result := govalidator.StringMatches(defaultPlmnId.Mnc, "^[0-9]{2,3}$"); !result {
		err := errors.New("Invalid mnc: " + defaultPlmnId.Mnc + ", should be 2 or 3 digits string, digit: 0~9")
		return false, err
	}

	for index, serviceName := range c.ServiceNameList {
		switch {
		case serviceName == "nnrf-nfm":
		case serviceName == "nnrf-disc":
		default:
			err := errors.New("Invalid serviceNameList[" + strconv.Itoa(index) + "]: " +
				serviceName + ", should be nnrf-nfm, nnrf-disc.")
			return false, err
		}
	}

	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

type Sbi struct {
	Scheme       string `yaml:"scheme" valid:"scheme,required"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty" valid:"host,optional"`
	// IP that is serviced or registered at another NRF.
	// IPv6Addr  string `yaml:"ipv6Addr,omitempty"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty" valid:"host,required"` // IP used to run the server in the node.
	Port        int    `yaml:"port,omitempty" valid:"port,optional"`
	Cert        *Cert  `yaml:"cert,omitempty" valid:"optional"`
	RootCert    *Cert  `yaml:"rootcert,omitempty" valid:"optional"`
	OAuth       bool   `yaml:"oauth,omitempty" valid:"optional"`
}

func (s *Sbi) validate() (bool, error) {
	govalidator.TagMap["scheme"] = govalidator.Validator(func(str string) bool {
		return str == "https" || str == "http"
	})

	result, err := govalidator.ValidateStruct(s)
	return result, appendInvalid(err)
}

type Cert struct {
	Pem string `yaml:"pem,omitempty" valid:"type(string),minstringlength(1),required"`
	Key string `yaml:"key,omitempty" valid:"type(string),minstringlength(1),required"`
}

func appendInvalid(err error) error {
	var errs govalidator.Errors

	if err == nil {
		return nil
	}

	es := err.(govalidator.Errors).Errors()
	for _, e := range es {
		errs = append(errs, fmt.Errorf("Invalid %w", e))
	}

	return error(errs)
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (c *Config) GetSbiScheme() string {
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.Scheme != "" {
		return c.Configuration.Sbi.Scheme
	}
	return NrfSbiDefaultScheme
}

func (c *Config) GetSbiPort() int {
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.Port != 0 {
		return c.Configuration.Sbi.Port
	}
	return NrfSbiDefaultPort
}

func (c *Config) GetSbiBindingIP() string {
	bindIP := "0.0.0.0"
	if c.Configuration == nil || c.Configuration.Sbi == nil {
		return bindIP
	}
	if c.Configuration.Sbi.BindingIPv4 != "" {
		if bindIP = os.Getenv(c.Configuration.Sbi.BindingIPv4); bindIP != "" {
			logger.CfgLog.Infof("Parsing ServerIPv4 [%s] from ENV Variable", bindIP)
		} else {
			bindIP = c.Configuration.Sbi.BindingIPv4
		}
	}
	return bindIP
}

func (c *Config) GetSbiBindingAddr() string {
	return c.GetSbiBindingIP() + ":" + strconv.Itoa(c.GetSbiPort())
}

func (c *Config) GetSbiRegisterIP() string {
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.RegisterIPv4 != "" {
		return c.Configuration.Sbi.RegisterIPv4
	}
	return NrfSbiDefaultIPv4
}

func (c *Config) GetSbiRegisterAddr() string {
	return c.GetSbiRegisterIP() + ":" + strconv.Itoa(c.GetSbiPort())
}

func (c *Config) GetSbiUri() string {
	return c.GetSbiScheme() + "://" + c.GetSbiRegisterAddr()
}

func (c *Config) GetOAuth() bool {
	return c.Configuration.Sbi.OAuth
}

func (c *Config) GetNrfCertPemPath() string {
	if c.Configuration.Sbi.Cert != nil {
		return c.Configuration.Sbi.Cert.Pem
	}
	return NrfDefaultCertPemPath
}

func (c *Config) GetCertBasePath() string {
	dir, _ := filepath.Split(c.GetNrfCertPemPath())
	return dir
}

func (c *Config) GetNrfPrivKeyPath() string {
	if c.Configuration.Sbi.Cert != nil {
		return c.Configuration.Sbi.Cert.Key
	}
	return NrfDefaultPrivateKeyPath
}

func (c *Config) GetRootCertPemPath() string {
	if c.Configuration.Sbi.RootCert != nil {
		return c.Configuration.Sbi.RootCert.Pem
	}
	return NrfDefaultRootCertPemPath
}

func (c *Config) GetRootPrivKeyPath() string {
	if c.Configuration.Sbi.RootCert != nil {
		return c.Configuration.Sbi.RootCert.Key
	}
	return NrfDefaultRootPrivateKeyPath
}
