package models

type Dnf struct {
	dnfUnits []DnfUnit `json:"dnfUints" bson:"dnfUnits"`
}
