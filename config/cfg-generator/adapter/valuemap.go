package adapter

import (
	"fmt"
	"time"

	"github.com/free5gc/openapi/models"
)

// distinct collects unique values in first-appearance order.
type distinct[V comparable] struct {
	order []V
	seen  map[V]bool
}

func (d *distinct[V]) add(v V) {
	if d.seen == nil {
		d.seen = map[V]bool{}
	}
	if !d.seen[v] {
		d.seen[v] = true
		d.order = append(d.order, v)
	}
}

// buildValueMap implements the value-mapping rule for multi-valued
// identifiers such as S-NSSAI and DNN (docs/global-mapping.md §0):
//
//  1. the caller collects the NF's distinct values in a fixed path order;
//  2. their count must equal len(global), otherwise it is an error and the
//     generator does not guess;
//  3. the result maps found[i] -> global[i].
//
// A nil global means the field was not provided: it returns (nil, nil) and
// the caller must leave the NF untouched.
func buildValueMap[V comparable](found, global []V, what, where string) (map[V]V, error) {
	if global == nil {
		return nil, nil
	}
	if len(found) != len(global) {
		return nil, fmt.Errorf("global %s has %d entries but %s has %d distinct value(s) %v",
			what, len(global), where, len(found), found)
	}
	m := make(map[V]V, len(found))
	for i, v := range found {
		m[v] = global[i]
	}
	return m, nil
}

// mapSnssai rewrites *s through m and records per-field changes. Values not in
// m are left alone.
func mapSnssai(r *Recorder, nf NFName, path string, s *models.Snssai, m map[models.Snssai]models.Snssai) {
	if s == nil {
		return
	}
	to, ok := m[*s]
	if !ok {
		return
	}
	Set(r, nf, path+".sst", &s.Sst, to.Sst)
	Set(r, nf, path+".sd", &s.Sd, to.Sd)
}

// mapString rewrites *s through m and records the change. Values not in m are
// left alone.
func mapString(r *Recorder, nf NFName, path string, s *string, m map[string]string) {
	if to, ok := m[*s]; ok {
		Set(r, nf, path, s, to)
	}
}

// setDuration assigns a fresh pointer holding v to *dst and records the change
// if the value differs. A fresh pointer avoids aliasing global values into
// NF configs.
func setDuration(r *Recorder, nf NFName, path string, dst **time.Duration, v time.Duration) {
	old := "<unset>"
	if *dst != nil {
		if **dst == v {
			return
		}
		old = (**dst).String()
	}
	r.Changes = append(r.Changes, Change{NF: nf, Path: path, Old: old, New: v.String()})
	*dst = &v
}
