//go:binary-only-package

package tlv

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"reflect"
	"strconv"
)

func Marshal(v interface{}) ([]byte, error) {}

func makeTLV(tag int, value []byte) []byte {}

func buildTLV(tag int, v interface{}) ([]byte, error) {}
