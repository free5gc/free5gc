package app

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

/*
type config struct {
	Path     string `yaml:"path"`
	document []byte
}*/

type context struct {
	Path     string `yaml:"path"`
	document []byte
	DbUri    string `yaml:"db_uri"`
	Logger   Logger `yaml:"logger"`
}

type Logger struct {
	AMF                AMF                `yaml:"AMF"`
	UDM                UDM                `yaml:"UDM"`
	SMF                SMF                `yaml:"SMF"`
	NAS                NAS                `yaml:"NAS"`
	FSM                FSM                `yaml:"FSM"`
	NGAP               NGAP               `yaml:"NGAP"`
	NamfComm           NamfComm           `yaml:"NamfComm"`
	NamfEventExposure  NamfEventExposure  `yaml:"NamfEventExposure"`
	NsmfPDUSession     NsmfPDUSession     `yaml:"NsmfPDUSession"`
	NudrDataRepository NudrDataRepository `yaml:"NudrDataRepository"`
	OpenApi            OpenApi            `yaml:"OpenApi"`
	Aper               Aper               `yaml:"Aper"`
	CommonConsumerTest CommonConsumerTest `yaml:"CommonConsumerTest"`
	PCF                PCF                `yaml:"PCF"`
	UDR                UDR                `yaml:"UDR"`
	NRF                NRF                `yaml:"NRF"`
	NSSF               NSSF               `yaml:"NSSF"`
	AUSF               AUSF               `yaml:"AUSF"`
	N3IWF              N3IWF              `yaml:"N3IWF"`
	WEBUI              WEBUI              `yaml:"WEBUI"`
}

type AMF struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type UDM struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type SMF struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type PCF struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type UDR struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type NRF struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type NSSF struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type AUSF struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type N3IWF struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type NAS struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type FSM struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type NGAP struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type NamfComm struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type NamfEventExposure struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type NsmfPDUSession struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type NudrDataRepository struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type OpenApi struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type Aper struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type CommonConsumerTest struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

type WEBUI struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}

var self context

// Unused code
//var context_initialized = 0

func init() {
	self = context{}
	//context_initialized = 1
}

func ContextSelf() *context {
	return &self
}

func (c *context) readFile() error {
	document, err := ioutil.ReadFile(ContextSelf().Path)
	if err != nil {
		fmt.Println("yamlFile.Get err  ", err)
		return err
	}
	ContextSelf().document = document
	return nil
}

func (c *context) parseConfig() error {
	yamlFile := ContextSelf().document
	if err := yaml.Unmarshal(yamlFile, c); err != nil {
		fmt.Printf("yaml.Umarshal error: %v", err)
	}
	return nil
}
