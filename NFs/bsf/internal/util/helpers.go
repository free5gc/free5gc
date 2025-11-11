/*
 * BSF Utility Functions
 */

package util

import "github.com/free5gc/openapi/models"

// Helper functions for type conversion
func StringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func PtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func PtrToBindingLevel(b *models.BindingLevel) models.BindingLevel {
	if b == nil {
		return models.BindingLevel("")
	}
	return *b
}
