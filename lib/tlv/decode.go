//go:binary-only-package

package tlv

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func Unmarshal(b []byte, v interface{}) error {}

func decodeValue(b []byte, v interface{}) (err error) {}

func parseTLV(b []byte) (fragments, error) {}
