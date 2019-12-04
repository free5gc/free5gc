//go:binary-only-package

package TimeDecode

import (
	"reflect"
	"time"

	"free5gc/lib/openapi/models"

	"github.com/mitchellh/mapstructure"
)

// Decode - Only support []map[string]interface to []models.NfProfile
func Decode(source interface{}, format string) ([]models.NfProfile, error) {}
