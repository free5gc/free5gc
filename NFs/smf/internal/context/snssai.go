package context

import (
	"net"
	"strings"

	"github.com/free5gc/openapi/models"
)

type SNssai struct {
	Sst int32
	Sd  string
}

// Equal return true if two S-NSSAI is equal
func (s *SNssai) Equal(target *SNssai) bool {
	return s.Sst == target.Sst && strings.EqualFold(s.Sd, target.Sd)
}

func (s *SNssai) EqualModelsSnssai(target *models.Snssai) bool {
	return s.Sst == target.Sst && strings.EqualFold(s.Sd, target.Sd)
}

type SnssaiUPFInfo struct {
	SNssai  *SNssai
	DnnList []*DnnUPFInfoItem
}

// DnnUpfInfoItem presents UPF dnn information
type DnnUPFInfoItem struct {
	Dnn             string
	DnaiList        []string
	PduSessionTypes []models.PduSessionType
	UeIPPools       []*UeIPPool
	StaticIPPools   []*UeIPPool
}

// ContainsDNAI return true if the this dnn Info contains the specify DNAI
func (d *DnnUPFInfoItem) ContainsDNAI(targetDnai string) bool {
	if targetDnai == "" {
		return len(d.DnaiList) == 0
	}
	for _, dnai := range d.DnaiList {
		if dnai == targetDnai {
			return true
		}
	}
	return false
}

// ContainsIPPool returns true if the ip pool of this upf dnn info contains the `ip`
func (d *DnnUPFInfoItem) ContainsIPPool(ip net.IP) bool {
	if ip == nil {
		return true
	}
	for _, ipPool := range d.UeIPPools {
		if ipPool.ueSubNet.Contains(ip) {
			return true
		}
	}
	return false
}
