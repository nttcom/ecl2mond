package util

import "testing"

func TestCheckAllZeroValues(t *testing.T) {
	case1 := []string{"0", "0", "0", "0"}
	result := CheckAllZeroValues(case1)

	if !result {
		t.Error("result should be true.")
	}

	case2 := []string{"0", "1", "0", "0"}
	result = CheckAllZeroValues(case2)

	if result {
		t.Error("result should be false.")
	}
}
