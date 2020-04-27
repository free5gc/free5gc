/*
 * NRF Configuration Factory
 */

package factory

import (
	"free5gc/lib/openapi/models"
)

type Configuration struct {
	Sbi *Sbi `yaml:"sbi,omitempty"`

	MongoDBName string `yaml:"MongoDBName"`

	MongoDBUrl string `yaml:"MongoDBUrl"`

	DefaultServiceIP string `yaml:"DefaultServiceIP"`

	DefaultPlmnId models.PlmnId `yaml:"DefaultPlmnId"`

	ServiceNameList []string `yaml:"serviceNameList,omitempty"`
}
