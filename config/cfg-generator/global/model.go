// Package global defines the user-facing global configuration block and its
// validation. Values provided here are applied to every NF config by the
// per-NF adapters; nil fields mean "not provided, keep the NF value".
package global

import (
	"net"
	"strconv"
	"time"

	"github.com/free5gc/openapi/models"
)

// Global is the `global:` block of the integrated input file.
//
// Semantics (docs/adapter-data-model.md §2):
//   - a nil scalar or nil list means "not provided": NF values are kept;
//   - an empty list (`[]`) is rejected by Validate, so a typo cannot wipe
//     NF settings.
//
// The NF paths each field maps to are listed in docs/global-mapping.md.
type Global struct {
	// Plmn is the home PLMN.
	Plmn *models.PlmnId `yaml:"plmn,omitempty"`
	// Tac is the tracking area code, a 6-digit hex string such as "000001".
	Tac *string `yaml:"tac,omitempty"`
	// Snssais are the slices. Applied by value mapping: the i-th distinct
	// S-NSSAI found in an NF config is replaced by Snssais[i].
	Snssais []models.Snssai `yaml:"snssaiList,omitempty"`
	// Dnns are the data network names, applied by value mapping like Snssais.
	Dnns []string `yaml:"dnnList,omitempty"`
	// Nrf is the NRF SBI endpoint; NFs get it as nrfUri.
	Nrf *NrfRef `yaml:"nrf,omitempty"`
	// SbiScheme is the SBI scheme of every NF: "http" or "https".
	SbiScheme *string `yaml:"sbiScheme,omitempty"`
	// MongoDB overrides the database of every NF that has one, field by field.
	MongoDB *MongoDB `yaml:"mongodb,omitempty"`
	// Logger overrides the logger block of every NF, field by field.
	Logger *Logger `yaml:"logger,omitempty"`
	// Timeouts overrides the SBI timeouts of every NF, field by field.
	Timeouts *Timeouts `yaml:"timeouts,omitempty"`
}

// NrfRef is the NRF SBI endpoint.
type NrfRef struct {
	Scheme string `yaml:"scheme"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

// URI returns the NRF URI in the form NFs expect, e.g. "http://127.0.0.10:8000".
// IPv6 hosts are bracketed.
func (n NrfRef) URI() string {
	return n.Scheme + "://" + net.JoinHostPort(n.Host, strconv.Itoa(n.Port))
}

// MongoDB fields are pointers so each can be overridden on its own.
type MongoDB struct {
	Name *string `yaml:"name,omitempty"`
	Url  *string `yaml:"url,omitempty"`
}

// Logger fields are pointers so each can be overridden on its own.
type Logger struct {
	Enable       *bool   `yaml:"enable,omitempty"`
	Level        *string `yaml:"level,omitempty"`
	ReportCaller *bool   `yaml:"reportCaller,omitempty"`
}

// Timeouts map to the NF fields idleTimeout, readIdleTimeout and pingTimeout.
// Each is applied only when provided. Omitting them in an NF config does not
// mean "use the sample value": AMF falls back to 1ms / 1s / 2s
// (docs/nf-config-defaults.md).
type Timeouts struct {
	Idle     *time.Duration `yaml:"idle,omitempty"`
	ReadIdle *time.Duration `yaml:"readIdle,omitempty"`
	Ping     *time.Duration `yaml:"ping,omitempty"`
}
