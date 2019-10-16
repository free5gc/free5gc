package models

type ComplexQuery struct {
	CNf *Cnf `json: "cnf,omitempty" bson:"cnf,omitempty"`
	DNf *Dnf `json: "dnf,omitempty" bson:"dnf,omitempty"`
}
