/*
 * AMF Configuration Factory
 */

package factory

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/davecgh/go-spew/spew"

	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/openapi/models"
)

const (
	AmfDefaultTLSKeyLogPath   = "./log/amfsslkey.log"
	AmfDefaultCertPemPath     = "./cert/amf.pem"
	AmfDefaultPrivateKeyPath  = "./cert/amf.key"
	AmfDefaultConfigPath      = "./config/amfcfg.yaml"
	AmfSbiDefaultIPv4         = "127.0.0.18"
	AmfSbiDefaultPort         = 8000
	AmfSbiDefaultScheme       = "https"
	AmfDefaultNrfUri          = "https://127.0.0.10:8000"
	sctpDefaultNumOstreams    = 3
	sctpDefaultMaxInstreams   = 5
	sctpDefaultMaxAttempts    = 2
	sctpDefaultMaxInitTimeout = 2
	ngapDefaultPort           = 38412
	AmfCallbackResUriPrefix   = "/namf-callback/v1"
	AmfCommResUriPrefix       = "/namf-comm/v1"
	AmfEvtsResUriPrefix       = "/namf-evts/v1"
	AmfLocResUriPrefix        = "/namf-loc/v1"
	AmfMtResUriPrefix         = "/namf-mt/v1"
	AmfOamResUriPrefix        = "/namf-oam/v1"
	AmfMbsComResUriPrefix     = "/namf-mbs-comm/v1"
	AmfMbsBCResUriPrefix      = "/namf-mbs-bc/v1"
)

type Config struct {
	Info          *Info          `yaml:"info" valid:"required"`
	Configuration *Configuration `yaml:"configuration" valid:"required"`
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
	Version     string `yaml:"version,omitempty" valid:"required,in(1.0.9)"`
	Description string `yaml:"description,omitempty" valid:"type(string)"`
}

type Configuration struct {
	AmfName                string            `yaml:"amfName,omitempty" valid:"required, type(string)"`
	NgapIpList             []string          `yaml:"ngapIpList,omitempty" valid:"required"`
	NgapPort               int               `yaml:"ngapPort,omitempty" valid:"optional,port"`
	Sbi                    *Sbi              `yaml:"sbi,omitempty" valid:"required"`
	ServiceNameList        []string          `yaml:"serviceNameList,omitempty" valid:"required"`
	ServedGumaiList        []models.Guami    `yaml:"servedGuamiList,omitempty" valid:"required"`
	SupportTAIList         []models.Tai      `yaml:"supportTaiList,omitempty" valid:"required"`
	PlmnSupportList        []PlmnSupportItem `yaml:"plmnSupportList,omitempty" valid:"required"`
	SupportDnnList         []string          `yaml:"supportDnnList,omitempty" valid:"required"`
	SupportLadnList        []Ladn            `yaml:"supportLadnList,omitempty" valid:"optional"`
	NrfUri                 string            `yaml:"nrfUri,omitempty" valid:"required, url"`
	NrfCertPem             string            `yaml:"nrfCertPem,omitempty" valid:"optional"`
	Security               *Security         `yaml:"security,omitempty" valid:"required"`
	NetworkName            NetworkName       `yaml:"networkName,omitempty" valid:"required"`
	NgapIE                 *NgapIE           `yaml:"ngapIE,omitempty" valid:"optional"`
	NasIE                  *NasIE            `yaml:"nasIE,omitempty" valid:"optional"`
	T3502Value             int               `yaml:"t3502Value,omitempty" valid:"required, type(int)"`
	T3512Value             int               `yaml:"t3512Value,omitempty" valid:"required, type(int)"`
	Non3gppDeregTimerValue int               `yaml:"non3gppDeregTimerValue,omitempty" valid:"-"`
	T3513                  TimerValue        `yaml:"t3513" valid:"required"`
	T3522                  TimerValue        `yaml:"t3522" valid:"required"`
	T3550                  TimerValue        `yaml:"t3550" valid:"required"`
	T3560                  TimerValue        `yaml:"t3560" valid:"required"`
	T3565                  TimerValue        `yaml:"t3565" valid:"required"`
	T3570                  TimerValue        `yaml:"t3570" valid:"required"`
	T3555                  TimerValue        `yaml:"t3555" valid:"required"`
	Locality               string            `yaml:"locality,omitempty" valid:"type(string),optional"`
	SCTP                   *Sctp             `yaml:"sctp,omitempty" valid:"optional"`
	DefaultUECtxReq        bool              `yaml:"defaultUECtxReq,omitempty" valid:"type(bool),optional"`
}

type Logger struct {
	Enable       bool   `yaml:"enable" valid:"type(bool)"`
	Level        string `yaml:"level" valid:"required,in(trace|debug|info|warn|error|fatal|panic)"`
	ReportCaller bool   `yaml:"reportCaller" valid:"type(bool)"`
}

func (c *Configuration) validate() (bool, error) {
	if c.NgapIpList != nil {
		var errs govalidator.Errors
		for _, v := range c.NgapIpList {
			if result := govalidator.IsHost(v); !result {
				err := fmt.Errorf("invalid NgapIpList: %s, value should be in the form of IP", v)
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return false, error(errs)
		}
	}

	if c.Sbi != nil {
		if _, err := c.Sbi.validate(); err != nil {
			return false, err
		}
	}

	if c.ServiceNameList != nil {
		var errs govalidator.Errors
		for _, v := range c.ServiceNameList {
			if v != "namf-comm" && v != "namf-evts" && v != "namf-mt" && v != "namf-loc" && v != "namf-oam" {
				err := fmt.Errorf("invalid ServiceNameList: %s,"+
					" value should be namf-comm or namf-evts or namf-mt or namf-loc or namf-oam", v)
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return false, error(errs)
		}
	}

	if c.ServedGumaiList != nil {
		var errs govalidator.Errors
		for _, v := range c.ServedGumaiList {
			if v.PlmnId == nil {
				return false, fmt.Errorf("ServedGumaiList: PlmnId is nil")
			}
			mcc := v.PlmnId.Mcc
			if result := govalidator.StringMatches(mcc, "^[0-9]{3}$"); !result {
				err := fmt.Errorf("invalid mcc: %s, should be a 3-digit number", mcc)
				errs = append(errs, err)
			}

			mnc := v.PlmnId.Mnc
			if result := govalidator.StringMatches(mnc, "^[0-9]{2,3}$"); !result {
				err := fmt.Errorf("invalid mnc: %s, should be a 2 or 3-digit number", mnc)
				errs = append(errs, err)
			}

			amfId := v.AmfId
			if result := govalidator.StringMatches(amfId, "^[A-Fa-f0-9]{6}$"); !result {
				err := fmt.Errorf("invalid amfId: %s,"+
					" should be 3 bytes hex string, range: 000000~FFFFFF", amfId)
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return false, error(errs)
		}
	}

	if c.SupportTAIList != nil {
		var errs govalidator.Errors
		for _, v := range c.SupportTAIList {
			if v.PlmnId == nil {
				return false, fmt.Errorf("SupportTAIList: PlmnId is nil")
			}
			mcc := v.PlmnId.Mcc
			if result := govalidator.StringMatches(mcc, "^[0-9]{3}$"); !result {
				err := fmt.Errorf("invalid mcc: %s, should be a 3-digit number", mcc)
				errs = append(errs, err)
			}

			mnc := v.PlmnId.Mnc
			if result := govalidator.StringMatches(mnc, "^[0-9]{2,3}$"); !result {
				err := fmt.Errorf("invalid mnc: %s, should be a 2 or 3-digit number", mnc)
				errs = append(errs, err)
			}

			tac := v.Tac
			if result := govalidator.StringMatches(tac, "^[A-Fa-f0-9]{6}$"); !result {
				err := fmt.Errorf("invalid tac: %s, should be 3 bytes hex string, range: 000000~FFFFFF", tac)
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return false, error(errs)
		}
	}

	if c.PlmnSupportList != nil {
		var errs govalidator.Errors
		for _, v := range c.PlmnSupportList {
			if _, err := v.validate(); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return false, error(errs)
		}
	}

	if c.SupportLadnList != nil {
		var errs govalidator.Errors
		for _, v := range c.SupportLadnList {
			if _, err := v.validate(); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return false, error(errs)
		}
	}

	if c.Security != nil {
		if _, err := c.Security.validate(); err != nil {
			return false, err
		}
	}

	if n3gppVal := &(c.Non3gppDeregTimerValue); n3gppVal == nil {
		err := fmt.Errorf("invalid Non3gppDeregTimerValue: value is required")
		return false, err
	}

	if _, err := c.NetworkName.validate(); err != nil {
		return false, err
	}

	if c.NgapIE != nil {
		if _, err := c.NgapIE.validate(); err != nil {
			return false, err
		}
	}

	if c.NasIE != nil {
		if _, err := c.NasIE.validate(); err != nil {
			return false, err
		}
	}

	if _, err := c.T3513.validate(); err != nil {
		return false, err
	}

	if _, err := c.T3522.validate(); err != nil {
		return false, err
	}

	if _, err := c.T3550.validate(); err != nil {
		return false, err
	}

	if _, err := c.T3560.validate(); err != nil {
		return false, err
	}

	if _, err := c.T3565.validate(); err != nil {
		return false, err
	}

	if _, err := c.T3570.validate(); err != nil {
		return false, err
	}

	if _, err := c.T3555.validate(); err != nil {
		return false, err
	}

	if c.SCTP != nil {
		if _, err := c.SCTP.validate(); err != nil {
			return false, err
		}
	}

	if _, err := govalidator.ValidateStruct(c); err != nil {
		return false, appendInvalid(err)
	}

	return true, nil
}

type Sbi struct {
	Scheme       string `yaml:"scheme" valid:"required,scheme"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty" valid:"required,host"` // IP that is registered at NRF.
	BindingIPv4  string `yaml:"bindingIPv4,omitempty" valid:"required,host"`  // IP used to run the server in the node.
	Port         int    `yaml:"port,omitempty" valid:"required,port"`
	Tls          *Tls   `yaml:"tls,omitempty" valid:"optional"`
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

	if _, err := govalidator.ValidateStruct(s); err != nil {
		return false, appendInvalid(err)
	}

	return true, nil
}

type Tls struct {
	Pem string `yaml:"pem,omitempty" valid:"type(string),minstringlength(1),required"`
	Key string `yaml:"key,omitempty" valid:"type(string),minstringlength(1),required"`
}

func (t *Tls) validate() (bool, error) {
	result, err := govalidator.ValidateStruct(t)
	return result, err
}

type Security struct {
	IntegrityOrder []string `yaml:"integrityOrder,omitempty" valid:"-"`
	CipheringOrder []string `yaml:"cipheringOrder,omitempty" valid:"-"`
}

func (s *Security) validate() (bool, error) {
	var errs govalidator.Errors

	if s.IntegrityOrder != nil {
		for _, val := range s.IntegrityOrder {
			if result := govalidator.Contains(val, "NIA"); !result {
				err := fmt.Errorf("invalid integrityOrder: %s, should be NIA-series integrity algorithms", val)
				errs = append(errs, err)
			}
		}
	}
	if s.CipheringOrder != nil {
		for _, val := range s.CipheringOrder {
			if result := govalidator.Contains(val, "NEA"); !result {
				err := fmt.Errorf("invalid cipheringOrder: %s, should be NEA-series ciphering algorithms", val)
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return false, error(errs)
	}

	return true, nil
}

type PlmnSupportItem struct {
	PlmnId     *models.PlmnId  `yaml:"plmnId" valid:"required"`
	SNssaiList []models.Snssai `yaml:"snssaiList,omitempty" valid:"required"`
}

func (p *PlmnSupportItem) validate() (bool, error) {
	var errs govalidator.Errors

	if _, err := govalidator.ValidateStruct(p); err != nil {
		return false, appendInvalid(err)
	}

	mcc := p.PlmnId.Mcc
	if result := govalidator.StringMatches(mcc, "^[0-9]{3}$"); !result {
		err := fmt.Errorf("invalid mcc: %s, should be a 3-digit number", mcc)
		errs = append(errs, err)
	}

	mnc := p.PlmnId.Mnc
	if result := govalidator.StringMatches(mnc, "^[0-9]{2,3}$"); !result {
		err := fmt.Errorf("invalid mnc: %s, should be a 2 or 3-digit number", mnc)
		errs = append(errs, err)
	}

	for _, snssai := range p.SNssaiList {
		sst := snssai.Sst
		sd := snssai.Sd
		if result := govalidator.InRangeInt(sst, 0, 255); !result {
			err := fmt.Errorf("invalid sst: %d, should be in the range of 0~255", sst)
			errs = append(errs, err)
		}
		if sd != "" {
			if result := govalidator.StringMatches(sd, "^[A-Fa-f0-9]{6}$"); !result {
				err := fmt.Errorf("invalid sd: %s, should be 3 bytes hex string, range: 000000~FFFFFF", sd)
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return false, error(errs)
	}

	return true, nil
}

type Ladn struct {
	Dnn     string       `yaml:"dnn" valid:"type(string),minstringlength(1),required"`
	TaiList []models.Tai `yaml:"taiList" valid:"required"`
}

func (l *Ladn) validate() (bool, error) {
	if _, err := govalidator.ValidateStruct(l); err != nil {
		return false, appendInvalid(err)
	}

	var errs govalidator.Errors
	for _, v := range l.TaiList {
		if v.PlmnId == nil {
			return false, fmt.Errorf("PlmnId is nil")
		}
		mcc := v.PlmnId.Mcc
		if result := govalidator.StringMatches(mcc, "^[0-9]{3}$"); !result {
			err := fmt.Errorf("invalid mcc: %s, should be a 3-digit number", mcc)
			errs = append(errs, err)
		}

		mnc := v.PlmnId.Mnc
		if result := govalidator.StringMatches(mnc, "^[0-9]{2,3}$"); !result {
			err := fmt.Errorf("invalid mnc: %s, should be a 2 or 3-digit number", mnc)
			errs = append(errs, err)
		}

		tac := v.Tac
		if result := govalidator.StringMatches(tac, "^[A-Fa-f0-9]{6}$"); !result {
			err := fmt.Errorf("invalid tac: %s, should be 3 bytes hex string, range: 000000~FFFFFF", tac)
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return false, error(errs)
	}

	return true, nil
}

type NetworkName struct {
	Full  string `yaml:"full" valid:"type(string)"`
	Short string `yaml:"short,omitempty" valid:"type(string)"`
}

func (n *NetworkName) validate() (bool, error) {
	if _, err := govalidator.ValidateStruct(n); err != nil {
		return false, appendInvalid(err)
	}

	return true, nil
}

type NgapIE struct {
	MobilityRestrictionList  *MobilityRestrictionList  `yaml:"mobilityRestrictionList,omitempty" valid:"optional"`
	MaskedIMEISV             *MaskedIMEISV             `yaml:"maskedIMEISV,omitempty" valid:"optional"`
	RedirectionVoiceFallback *RedirectionVoiceFallback `yaml:"redirectionVoiceFallback,omitempty" valid:"optional"`
}

func (n *NgapIE) validate() (bool, error) {
	if n.MobilityRestrictionList != nil {
		if _, err := n.MobilityRestrictionList.validate(); err != nil {
			return false, err
		}
	}

	if n.MaskedIMEISV != nil {
		if _, err := n.MaskedIMEISV.validate(); err != nil {
			return false, err
		}
	}

	if n.RedirectionVoiceFallback != nil {
		if _, err := n.RedirectionVoiceFallback.validate(); err != nil {
			return false, err
		}
	}

	if _, err := govalidator.ValidateStruct(n); err != nil {
		return false, appendInvalid(err)
	}

	return true, nil
}

type MobilityRestrictionList struct {
	Enable bool `yaml:"enable" valid:"type(bool)"`
}

func (m *MobilityRestrictionList) validate() (bool, error) {
	if _, err := govalidator.ValidateStruct(m); err != nil {
		return false, appendInvalid(err)
	}

	return true, nil
}

type MaskedIMEISV struct {
	Enable bool `yaml:"enable" valid:"type(bool)"`
}

func (m *MaskedIMEISV) validate() (bool, error) {
	if _, err := govalidator.ValidateStruct(m); err != nil {
		return false, appendInvalid(err)
	}

	return true, nil
}

type RedirectionVoiceFallback struct {
	Enable bool `yaml:"enable" valid:"type(bool)"`
}

func (r *RedirectionVoiceFallback) validate() (bool, error) {
	if _, err := govalidator.ValidateStruct(r); err != nil {
		return false, appendInvalid(err)
	}

	return true, nil
}

type Sctp struct {
	NumOstreams    uint `yaml:"numOstreams,omitempty" valid:"int"`
	MaxInstreams   uint `yaml:"maxInstreams,omitempty" valid:"int"`
	MaxAttempts    uint `yaml:"maxAttempts,omitempty" valid:"int"`
	MaxInitTimeout uint `yaml:"maxInitTimeout,omitempty" valid:"int"`
}

func (n *Sctp) validate() (bool, error) {
	var errs govalidator.Errors
	if _, err := govalidator.ValidateStruct(n); err != nil {
		return false, appendInvalid(err)
	}
	if n.NumOstreams > 10 || n.NumOstreams <= 0 {
		errs = append(errs, fmt.Errorf("0 < configuration.sctp.maxOsStream <=20"))
	}
	if n.MaxInstreams > 10 || n.MaxInstreams <= 0 {
		errs = append(errs, fmt.Errorf("0 < configuration.sctp.maxInputStream <=20"))
	}
	if n.MaxInitTimeout > 5 || n.MaxInitTimeout <= 0 {
		errs = append(errs, fmt.Errorf(" 0 < configuration.sctp.maxInitTimeOut<=5 "))
	}
	if n.MaxAttempts > 5 || n.MaxAttempts <= 0 {
		errs = append(errs, fmt.Errorf(" 0 < configuration.sctp.maxAttempts <=5 "))
	}
	if len(errs) > 0 {
		return false, errs
	}
	return true, nil
}

type NasIE struct {
	NetworkFeatureSupport5GS *NetworkFeatureSupport5GS `yaml:"networkFeatureSupport5GS,omitempty" valid:"optional"`
}

func (n *NasIE) validate() (bool, error) {
	if n.NetworkFeatureSupport5GS != nil {
		if _, err := n.NetworkFeatureSupport5GS.validate(); err != nil {
			return false, err
		}
	}

	if _, err := govalidator.ValidateStruct(n); err != nil {
		return false, appendInvalid(err)
	}

	return true, nil
}

type NetworkFeatureSupport5GS struct {
	Enable  bool  `yaml:"enable" valid:"type(bool)"`
	Length  uint8 `yaml:"length" valid:"type(uint8)"`
	ImsVoPS uint8 `yaml:"imsVoPS" valid:"type(uint8)"`
	Emc     uint8 `yaml:"emc" valid:"type(uint8)"`
	Emf     uint8 `yaml:"emf" valid:"type(uint8)"`
	IwkN26  uint8 `yaml:"iwkN26" valid:"type(uint8)"`
	Mpsi    uint8 `yaml:"mpsi" valid:"type(uint8)"`
	EmcN3   uint8 `yaml:"emcN3" valid:"type(uint8)"`
	Mcsi    uint8 `yaml:"mcsi" valid:"type(uint8)"`
}

func (f *NetworkFeatureSupport5GS) validate() (bool, error) {
	var errs govalidator.Errors

	if result := govalidator.InRangeInt(f.Length, 1, 3); !result {
		err := fmt.Errorf("invalid length: %d, should be in the range of 1~3", f.Length)
		errs = append(errs, err)
	}
	if result := govalidator.InRangeInt(f.ImsVoPS, 0, 1); !result {
		err := fmt.Errorf("invalid imsVoPS: %d, should be in the range of 0~1", f.ImsVoPS)
		errs = append(errs, err)
	}
	if result := govalidator.InRangeInt(f.Emc, 0, 3); !result {
		err := fmt.Errorf("invalid emc: %d, should be in the range of 0~3", f.Emc)
		errs = append(errs, err)
	}
	if result := govalidator.InRangeInt(f.Emf, 0, 3); !result {
		err := fmt.Errorf("invalid emf: %d, should be in the range of 0~3", f.Emf)
		errs = append(errs, err)
	}
	if result := govalidator.InRangeInt(f.IwkN26, 0, 1); !result {
		err := fmt.Errorf("invalid iwkN26: %d, should be in the range of 0~1", f.IwkN26)
		errs = append(errs, err)
	}
	if result := govalidator.InRangeInt(f.Mpsi, 0, 1); !result {
		err := fmt.Errorf("invalid mpsi: %d, should be in the range of 0~1", f.Mpsi)
		errs = append(errs, err)
	}
	if result := govalidator.InRangeInt(f.EmcN3, 0, 1); !result {
		err := fmt.Errorf("invalid emcN3: %d, should be in the range of 0~1", f.EmcN3)
		errs = append(errs, err)
	}
	if result := govalidator.InRangeInt(f.Mcsi, 0, 1); !result {
		err := fmt.Errorf("invalid mcsi: %d, should be in the range of 0~1", f.Mcsi)
		errs = append(errs, err)
	}
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return false, appendInvalid(err)
	}

	if len(errs) > 0 {
		return false, error(errs)
	}

	return true, nil
}

type TimerValue struct {
	Enable        bool          `yaml:"enable" valid:"type(bool)"`
	ExpireTime    time.Duration `yaml:"expireTime" valid:"type(time.Duration)"`
	MaxRetryTimes int           `yaml:"maxRetryTimes,omitempty" valid:"type(int)"`
}

func (t *TimerValue) validate() (bool, error) {
	if _, err := govalidator.ValidateStruct(t); err != nil {
		return false, appendInvalid(err)
	}

	return true, nil
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

func (c *Config) Print() {
	spew.Config.Indent = "\t"
	str := spew.Sdump(c.Configuration)
	logger.CfgLog.Infof("==================================================")
	logger.CfgLog.Infof("%s", str)
	logger.CfgLog.Infof("==================================================")
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

func (c *Config) GetSbiScheme() string {
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.Scheme != "" {
		return c.Configuration.Sbi.Scheme
	}
	return AmfSbiDefaultScheme
}

func (c *Config) GetSbiPort() int {
	if c.Configuration != nil && c.Configuration.Sbi != nil && c.Configuration.Sbi.Port != 0 {
		return c.Configuration.Sbi.Port
	}
	return AmfSbiDefaultPort
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
	return AmfSbiDefaultIPv4
}

func (c *Config) GetSbiRegisterAddr() string {
	return c.GetSbiRegisterIP() + ":" + strconv.Itoa(c.GetSbiPort())
}

func (c *Config) GetSbiUri() string {
	return c.GetSbiScheme() + "://" + c.GetSbiRegisterAddr()
}

func (c *Config) GetNrfUri() string {
	if c.Configuration != nil && c.Configuration.NrfUri != "" {
		return c.Configuration.NrfUri
	}
	return AmfDefaultNrfUri
}

func (c *Config) GetServiceNameList() []string {
	if c.Configuration != nil && len(c.Configuration.ServiceNameList) > 0 {
		return c.Configuration.ServiceNameList
	}
	return nil
}

func (c *Config) GetNgapIEMobilityRestrictionList() *MobilityRestrictionList {
	if c != nil && c.Configuration != nil && c.Configuration.NgapIE != nil {
		return c.Configuration.NgapIE.MobilityRestrictionList
	}
	return nil
}

func (c *Config) GetNgapIEMaskedIMEISV() *MaskedIMEISV {
	if c.Configuration != nil && c.Configuration.NgapIE != nil {
		return c.Configuration.NgapIE.MaskedIMEISV
	}
	return nil
}

func (c *Config) GetNgapIERedirectionVoiceFallback() *RedirectionVoiceFallback {
	if c.Configuration != nil && c.Configuration.NgapIE != nil {
		return c.Configuration.NgapIE.RedirectionVoiceFallback
	}
	return nil
}

func (c *Config) GetNasIENetworkFeatureSupport5GS() *NetworkFeatureSupport5GS {
	if c.Configuration != nil && c.Configuration.NasIE != nil {
		return c.Configuration.NasIE.NetworkFeatureSupport5GS
	}
	return nil
}

func (c *Config) GetNgapPort() int {
	if c.Configuration.NgapPort != 0 {
		return c.Configuration.NgapPort
	}
	return ngapDefaultPort
}

func (c *Config) GetSctpConfig() *Sctp {
	if c.Configuration != nil && c.Configuration.SCTP != nil {
		return c.Configuration.SCTP
	}
	return &Sctp{
		NumOstreams:    sctpDefaultNumOstreams,
		MaxInstreams:   sctpDefaultMaxInstreams,
		MaxAttempts:    sctpDefaultMaxAttempts,
		MaxInitTimeout: sctpDefaultMaxInitTimeout,
	}
}

func (c *Config) GetCertPemPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.Configuration.Sbi.Tls.Pem
}

func (c *Config) GetCertKeyPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.Configuration.Sbi.Tls.Key
}
