//go:binary-only-package

package tlv

import (
	"reflect"
)

func isNumber(typ reflect.Type) bool {}

func isRefType(typ reflect.Type) bool {}

func hasValue(value reflect.Value) bool {}
