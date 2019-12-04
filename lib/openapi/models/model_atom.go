package models

type Atom struct {
	Attr     string `json:"attr" bson:"attr"`
	Value    string `json:"value" bson:"value"` // TODO: AnyType
	Negative bool   `json:"negative,omitempty" bson:"negative,omitempty"`
}
