/*
 * BSF Utility Functions
 */

package util

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
