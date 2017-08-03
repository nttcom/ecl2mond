package util

// CheckAllZeroValues checks values are zero strings.
func CheckAllZeroValues(values []string) bool {
	for _, value := range values {
		if value != "0" {
			return false
		}
	}
	return true
}
