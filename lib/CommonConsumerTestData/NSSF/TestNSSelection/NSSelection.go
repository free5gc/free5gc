/*
 * NSSF NS Selection
 *
 * Utility for Client API
 */

package TestNSSelection

import (
	"encoding/json"

	"github.com/antihax/optional"
)

func ConvertQueryParamIntf(intf interface{}) optional.Interface {
	e, _ := json.Marshal(intf)
	return optional.NewInterface(string(e))
}
