/*
 * NSSF Configuration Factory
 */

package factory

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/asaskevich/govalidator"

	"github.com/free5gc/nssf/internal/logger"
	"github.com/free5gc/openapi/models"
)

const (
	NssfDefaultTLSKeyLogPath   = "./log/nssfsslkey.log"
	NssfDefaultCertPemPath     = "./cert/nssf.pem"
	NssfDefaultPrivateKeyPath  = "./cert/nssf.key"
	NssfDefaultConfigPath      = "./config/nssfcfg.yaml"
	NssfSbiDefaultIPv4         = "127.0.0.31"
	NssfSbiDefaultPort         = 8000
	NssfSbiDefaultScheme       = "https"
	NssfDefaultNrfUri          = "https://127.0.0.10:8000"
	NssfNssaiavailResUriPrefix = "/nnssf-nssaiavailability/v1"
	NssfNsselectResUriPrefix   = "/nnssf-nsselection/v2"
)

type Config struct {
	Info          *Info          `yaml:"info" valid:"required"`
	Configuration *Configuration `yaml:"configuration" valid:"required"`
	Subscriptions []Subscription `yaml:"subscriptions,omitempty"`
	Logger        *Logger        `yaml:"logger" valid:"required"`
	sync.RWMutex
}

func (c *Config) Validate() (bool, error) {
	if configuration := c.Configuration; configuration != nil {
		if result, err := configuration.validate(); err != nil {
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

type Info struct {
	Version     string `yaml:"version" valid:"required,in(1.0.2)"`
	Description string `yaml:"description,omitempty" valid:"type(string)"`
}

type Configuration struct {
	NssfName                 string                  `yaml:"nssfName,omitempty"`
	Sbi                      *Sbi                    `yaml:"sbi"`
	ServiceNameList          []models.ServiceName    `yaml:"serviceNameList"`
	NrfUri                   string                  `yaml:"nrfUri"`
	NrfCertPem               string                  `yaml:"nrfCertPem,omitempty" valid:"optional"`
	SupportedPlmnList        []models.PlmnId         `yaml:"supportedPlmnList,omitempty"`
	SupportedNssaiInPlmnList []SupportedNssaiInPlmn  `yaml:"supportedNssaiInPlmnList"`
	NsiList                  []NsiConfig             `yaml:"nsiList,omitempty"`
	AmfSetList               []AmfSetConfig          `yaml:"amfSetList"`
	AmfList                  []AmfConfig             `yaml:"amfList"`
	TaList                   []TaConfig              `yaml:"taList"`
	MappingListFromPlmn      []MappingFromPlmnConfig `yaml:"mappingListFromPlmn"`
}

type Logger struct {
	Enable       bool   `yaml:"enable" valid:"type(bool)"`
	Level        string `yaml:"level" valid:"required,in(trace|debug|info|warn|error|fatal|panic)"`
	ReportCaller bool   `yaml:"reportCaller" valid:"type(bool)"`
}

func (c *Configuration) validate() (bool, error) {
	if sbi := c.Sbi; sbi != nil {
		if result, err := sbi.validate(); err != nil {
			return result, err
		}
	}

	for index, serviceName := range c.ServiceNameList {
		switch serviceName {
		case "nnssf-nsselection":
		case "nnssf-nssaiavailability":
		default:
			err := errors.New("Invalid serviceNameList[" + strconv.Itoa(index) + "]: " +
				string(serviceName) + ", should be nausf-auth.")
			return false, err
		}
	}

	for index, plmnId := range c.SupportedPlmnList {
		if result := govalidator.StringMatches(plmnId.Mcc, "^[0-9]{3}$"); !result {
			err := errors.New("Invalid plmnSupportList[" + strconv.Itoa(index) + "].Mcc: " +
				plmnId.Mcc + ", should be 3 digits interger.")
			return false, err
		}

		if result := govalidator.StringMatches(plmnId.Mnc, "^[0-9]{2,3}$"); !result {
			err := errors.New("Invalid plmnSupportList[" + strconv.Itoa(index) + "].Mnc: " +
				plmnId.Mnc + ", should be 2 or 3 digits interger.")
			return false, err
		}
	}

	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

type Sbi struct {
	Scheme models.UriScheme `yaml:"scheme"`
	// Currently only support IPv4 and thus `Ipv4Addr` field shall not be empty
	RegisterIPv4 string `yaml:"registerIPv4,omitempty" valid:"host,required"` // IP that is registered at NRF.
	// IPv6Addr string `yaml:"ipv6Addr,omitempty"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty" valid:"host,required"` // IP used to run the server in the node.
	Port        int    `yaml:"port"`
	Tls         *Tls   `yaml:"tls,omitempty" valid:"optional"`
}

func (s *Sbi) validate() (bool, error) {
	govalidator.TagMap["scheme"] = govalidator.Validator(func(str string) bool {
		return str == "https" || str == "http"
	})

	if tls := s.Tls; tls != nil {
		if result, err := tls.validate(); err != nil {
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(s)
	return result, appendInvalid(err)
}

type Tls struct {
	Pem string `yaml:"pem,omitempty" valid:"type(string),minstringlength(1),required"`
	Key string `yaml:"key,omitempty" valid:"type(string),minstringlength(1),required"`
}

func (t *Tls) validate() (bool, error) {
	result, err := govalidator.ValidateStruct(t)
	return result, err
}

func appendInvalid(err error) error {
	var errs govalidator.Errors

	if err == nil {
		return nil
	}

	es := err.(govalidator.Errors).Errors()
	for _, e := range es {
		errs = append(errs, fmt.Errorf("invalid %w", e))
	}

	return error(errs)
}

type AmfConfig struct {
	NfId                           string                                  `yaml:"nfId"`
	SupportedNssaiAvailabilityData []models.SupportedNssaiAvailabilityData `yaml:"supportedNssaiAvailabilityData"`
}

type TaConfig struct {
	Tai                  *models.Tai               `yaml:"tai"`
	AccessType           *models.AccessType        `yaml:"accessType"`
	SupportedSnssaiList  []models.ExtSnssai        `yaml:"supportedSnssaiList"`
	RestrictedSnssaiList []models.RestrictedSnssai `yaml:"restrictedSnssaiList,omitempty"`
}

type SupportedNssaiInPlmn struct {
	PlmnId              *models.PlmnId  `yaml:"plmnId"`
	SupportedSnssaiList []models.Snssai `yaml:"supportedSnssaiList"`
}

type NsiConfig struct {
	Snssai             *models.Snssai          `yaml:"snssai"`
	NsiInformationList []models.NsiInformation `yaml:"nsiInformationList"`
}

type AmfSetConfig struct {
	AmfSetId                       string                                  `yaml:"amfSetId"`
	AmfList                        []string                                `yaml:"amfList,omitempty"`
	NrfAmfSet                      string                                  `yaml:"nrfAmfSet,omitempty"`
	SupportedNssaiAvailabilityData []models.SupportedNssaiAvailabilityData `yaml:"supportedNssaiAvailabilityData"`
}

type MappingFromPlmnConfig struct {
	OperatorName    string                   `yaml:"operatorName,omitempty"`
	HomePlmnId      *models.PlmnId           `yaml:"homePlmnId"`
	MappingOfSnssai []models.MappingOfSnssai `yaml:"mappingOfSnssai"`
}

type Subscription struct {
	SubscriptionId   string                                  `yaml:"subscriptionId"`
	SubscriptionData *models.NssfEventSubscriptionCreateData `yaml:"subscriptionData"`
}

func (c *Config) GetVersion() string {
	c.RLock()
	defer c.RUnlock()

	if c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (c *Config) SetLogEnable(enable bool) {
	c.Lock()
	defer c.Unlock()

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
	c.Lock()
	defer c.Unlock()

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
	c.Lock()
	defer c.Unlock()

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

func (c *Config) GetLogEnable() bool {
	c.RLock()
	defer c.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return false
	}
	return c.Logger.Enable
}

func (c *Config) GetLogLevel() string {
	c.RLock()
	defer c.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return "info"
	}
	return c.Logger.Level
}

func (c *Config) GetLogReportCaller() bool {
	c.RLock()
	defer c.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return false
	}
	return c.Logger.ReportCaller
}
