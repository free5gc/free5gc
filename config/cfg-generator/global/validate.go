package global

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/free5gc/openapi/models"
)

var (
	mccPattern   = regexp.MustCompile(`^[0-9]{3}$`)
	mncPattern   = regexp.MustCompile(`^[0-9]{2,3}$`)
	hex6Pattern  = regexp.MustCompile(`^[0-9A-Fa-f]{6}$`)
	schemes      = map[string]bool{"http": true, "https": true}
	loggerLevels = map[string]bool{
		"trace": true, "debug": true, "info": true, "warn": true,
		"error": true, "fatal": true, "panic": true,
	}
)

// FieldError describes one invalid global field.
type FieldError struct {
	Path   string
	Reason string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("global.%s: %s", e.Path, e.Reason)
}

// Validate checks every provided field and returns all problems found.
// A nil Global is valid: nothing will be applied.
func Validate(g *Global) []error {
	if g == nil {
		return nil
	}

	v := &validator{}
	v.plmn(g.Plmn)
	v.tac(g.Tac)
	v.snssais(g.Snssais)
	v.dnns(g.Dnns)
	v.sbiScheme(g.SbiScheme)
	v.nrf(g.Nrf)
	v.schemeConsistency(g.SbiScheme, g.Nrf)
	v.mongodb(g.MongoDB)
	v.logger(g.Logger)
	v.timeouts(g.Timeouts)
	return v.errs
}

type validator struct{ errs []error }

func (v *validator) add(path, format string, args ...any) {
	v.errs = append(v.errs, &FieldError{Path: path, Reason: fmt.Sprintf(format, args...)})
}

func (v *validator) plmn(p *models.PlmnId) {
	if p == nil {
		return
	}
	if !mccPattern.MatchString(p.Mcc) {
		v.add("plmn.mcc", "%q must be exactly 3 digits", p.Mcc)
	}
	if !mncPattern.MatchString(p.Mnc) {
		v.add("plmn.mnc", "%q must be 2 or 3 digits", p.Mnc)
	}
}

func (v *validator) tac(tac *string) {
	if tac != nil && !hex6Pattern.MatchString(*tac) {
		v.add("tac", "%q must be a 6-digit hex string, e.g. \"000001\"", *tac)
	}
}

func (v *validator) snssais(list []models.Snssai) {
	if list == nil {
		return
	}
	if len(list) == 0 {
		v.add("snssaiList", "must not be empty; omit the field to keep NF values")
		return
	}
	seen := map[models.Snssai]int{}
	for i, s := range list {
		path := fmt.Sprintf("snssaiList[%d]", i)
		if s.Sst < 0 || s.Sst > 255 {
			v.add(path+".sst", "%d must be in 0..255", s.Sst)
		}
		if s.Sd != "" && !hex6Pattern.MatchString(s.Sd) {
			v.add(path+".sd", "%q must be empty or a 6-digit hex string", s.Sd)
		}
		if j, dup := seen[s]; dup {
			v.add(path, "duplicates snssaiList[%d]", j)
		} else {
			seen[s] = i
		}
	}
}

func (v *validator) dnns(list []string) {
	if list == nil {
		return
	}
	if len(list) == 0 {
		v.add("dnnList", "must not be empty; omit the field to keep NF values")
		return
	}
	seen := map[string]int{}
	for i, d := range list {
		path := fmt.Sprintf("dnnList[%d]", i)
		if d == "" {
			v.add(path, "must not be empty")
			continue
		}
		if j, dup := seen[d]; dup {
			v.add(path, "%q duplicates dnnList[%d]", d, j)
		} else {
			seen[d] = i
		}
	}
}

func (v *validator) sbiScheme(s *string) {
	if s != nil && !schemes[*s] {
		v.add("sbiScheme", "%q must be http or https", *s)
	}
}

func (v *validator) nrf(n *NrfRef) {
	if n == nil {
		return
	}
	if !schemes[n.Scheme] {
		v.add("nrf.scheme", "%q must be http or https", n.Scheme)
	}
	if n.Host == "" {
		v.add("nrf.host", "must not be empty")
	}
	if n.Port < 1 || n.Port > 65535 {
		v.add("nrf.port", "%d must be in 1..65535", n.Port)
	}
}

// schemeConsistency: when both are given, NFs would otherwise talk to the
// NRF with a scheme different from the one they serve and expect.
func (v *validator) schemeConsistency(sbi *string, n *NrfRef) {
	if sbi == nil || n == nil || !schemes[*sbi] || !schemes[n.Scheme] {
		return
	}
	if *sbi != n.Scheme {
		v.add("nrf.scheme", "%q must match sbiScheme %q", n.Scheme, *sbi)
	}
}

func (v *validator) mongodb(m *MongoDB) {
	if m == nil {
		return
	}
	if m.Name != nil && *m.Name == "" {
		v.add("mongodb.name", "must not be empty")
	}
	if m.Url != nil && !strings.HasPrefix(*m.Url, "mongodb://") && !strings.HasPrefix(*m.Url, "mongodb+srv://") {
		v.add("mongodb.url", "%q must start with mongodb:// or mongodb+srv://", *m.Url)
	}
}

func (v *validator) logger(l *Logger) {
	if l != nil && l.Level != nil && !loggerLevels[*l.Level] {
		v.add("logger.level", "%q must be one of trace, debug, info, warn, error, fatal, panic", *l.Level)
	}
}

func (v *validator) timeouts(t *Timeouts) {
	if t == nil {
		return
	}
	for _, f := range []struct {
		path string
		d    *time.Duration
	}{{"timeouts.idle", t.Idle}, {"timeouts.readIdle", t.ReadIdle}, {"timeouts.ping", t.Ping}} {
		if f.d != nil && *f.d < 0 {
			v.add(f.path, "%s must not be negative", *f.d)
		}
	}
}
