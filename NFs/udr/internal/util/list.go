/*
 * UDR Utility - List
 *
 * Copy from NSSF and NEF utility
 */

package util

import (
	"reflect"
)

// Contain checks whether a slice contains an element
func Contain(target interface{}, slice interface{}) bool {
	arr := reflect.ValueOf(slice)
	if arr.Kind() == reflect.Slice {
		for i := 0; i < arr.Len(); i++ {
			if reflect.DeepEqual(arr.Index(i).Interface(), target) {
				return true
			}
		}
	}
	return false
}
