package adapter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/free5gc/openapi/models"

	"github.com/free5gc/config/cfg-generator/global"
	nssffactory "github.com/free5gc/nssf/pkg/factory"
)

// NSSFAdapter applies global fields to nssfcfg.yaml (docs/global-mapping.md §1).
//
//   - plmn: supportedPlmnList[*], supportedNssaiInPlmnList[*].plmnId. The
//     PLMNs inside amfSetList/amfList/taList/mappingListFromPlmn describe
//     other TAs or roaming partners and are kept (§2).
//   - nrf:  nrfUri; nsiList[*].nsiInformationList[*].nrfId and
//     amfSetList[*].nrfAmfSet are rewritten only when they start with the
//     original nrfUri (other NRFs, e.g. localhost:8084, are kept).
//   - snssaiList: not applied to NSSF (decision: NSSF slice data is much
//     richer than AMF/SMF; §3).
//   - sbiScheme, logger, timeouts: see applyCommon
type NSSFAdapter struct{}

var _ Adapter[nssffactory.Config] = NSSFAdapter{}

func (NSSFAdapter) NF() NFName       { return NSSF }
func (NSSFAdapter) FileName() string { return "nssfcfg.yaml" }

func (NSSFAdapter) Apply(g *global.Global, cfg *nssffactory.Config, r *Recorder) error {
	if g == nil {
		return nil
	}
	c := cfg.Configuration
	if c == nil {
		return errMissingConfiguration
	}

	if p := g.Plmn; p != nil {
		for i := range c.SupportedPlmnList {
			base := fmt.Sprintf("supportedPlmnList[%d]", i)
			Set(r, NSSF, base+".mcc", &c.SupportedPlmnList[i].Mcc, p.Mcc)
			Set(r, NSSF, base+".mnc", &c.SupportedPlmnList[i].Mnc, p.Mnc)
		}
		for i := range c.SupportedNssaiInPlmnList {
			id := c.SupportedNssaiInPlmnList[i].PlmnId
			if id == nil {
				continue
			}
			base := fmt.Sprintf("supportedNssaiInPlmnList[%d].plmnId", i)
			Set(r, NSSF, base+".mcc", &id.Mcc, p.Mcc)
			Set(r, NSSF, base+".mnc", &id.Mnc, p.Mnc)
		}
	}

	if g.Nrf != nil && c.NrfUri != "" {
		oldURI, newURI := c.NrfUri, g.Nrf.URI()
		for i := range c.NsiList {
			for j := range c.NsiList[i].NsiInformationList {
				info := &c.NsiList[i].NsiInformationList[j]
				if v, ok := replaceURIPrefix(info.NrfId, oldURI, newURI); ok {
					Set(r, NSSF, fmt.Sprintf("nsiList[%d].nsiInformationList[%d].nrfId", i, j), &info.NrfId, v)
				}
			}
		}
		for i := range c.AmfSetList {
			set := &c.AmfSetList[i]
			if v, ok := replaceURIPrefix(set.NrfAmfSet, oldURI, newURI); ok {
				Set(r, NSSF, fmt.Sprintf("amfSetList[%d].nrfAmfSet", i), &set.NrfAmfSet, v)
			}
		}
	}

	f := commonFields{
		IdleTimeout:     &c.IdleTimeout,
		ReadIdleTimeout: &c.ReadIdleTimeout,
		PingTimeout:     &c.PingTimeout,
		NrfUri:          &c.NrfUri,
	}
	// NSSF types sbi.scheme as models.UriScheme, so it is set here rather
	// than through commonFields.SbiScheme (*string).
	if g.SbiScheme != nil && c.Sbi != nil {
		Set(r, NSSF, "sbi.scheme", &c.Sbi.Scheme, models.UriScheme(*g.SbiScheme))
	}
	if cfg.Logger == nil && g.Logger != nil {
		cfg.Logger = &nssffactory.Logger{}
	}
	if cfg.Logger != nil {
		f.LoggerEnable, f.LoggerLevel, f.LoggerReportCaller =
			&cfg.Logger.Enable, &cfg.Logger.Level, &cfg.Logger.ReportCaller
	}
	applyCommon(g, NSSF, f, r)
	return nil
}

// replaceURIPrefix returns v with the URI prefix oldURI replaced by newURI.
// It matches only whole URI prefixes: oldURI must be followed by '/' or the
// end of v, so "http://h:8000" does not match "http://h:80001/...".
func replaceURIPrefix(v, oldURI, newURI string) (string, bool) {
	if oldURI == "" || !strings.HasPrefix(v, oldURI) {
		return v, false
	}
	rest := v[len(oldURI):]
	if rest != "" && !strings.HasPrefix(rest, "/") {
		return v, false
	}
	return newURI + rest, true
}

// Validate runs NSSF's own Validate and adds the checks its model lacks:
// sbi, nrfUri and serviceNameList carry no valid tags, so NSSF accepts them
// missing (docs/nf-config-defaults.md).
func (NSSFAdapter) Validate(cfg *nssffactory.Config) error {
	var errs []error
	if _, err := cfg.Validate(); err != nil {
		errs = append(errs, err)
	}
	if c := cfg.Configuration; c != nil {
		if c.Sbi == nil {
			errs = append(errs, errors.New("configuration.sbi is required"))
		} else if !isHTTPScheme(string(c.Sbi.Scheme)) {
			errs = append(errs, fmt.Errorf("configuration.sbi.scheme %q must be http or https", c.Sbi.Scheme))
		}
		if !isAbsoluteURL(c.NrfUri) {
			errs = append(errs, fmt.Errorf("configuration.nrfUri %q must be an absolute URL", c.NrfUri))
		}
		if len(c.ServiceNameList) == 0 {
			errs = append(errs, errors.New("configuration.serviceNameList is required"))
		}
	}
	return errors.Join(errs...)
}
