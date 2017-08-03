// +build linux

package linux

import (
	"testing"

	"github.com/nttcom/ecl2mond/environment"
)

func NewLoadavgCollector(env *environment.Environment, resourceID string) *LoadavgCollector {
	return &LoadavgCollector{
		Env:        env,
		ResourceID: resourceID,
	}
}

func TestGetLoadavgMeters(t *testing.T) {
	env := environment.NewEnvironment()
	collector := NewLoadavgCollector(env, "resourceId")
	meters, err := collector.Collect()

	if err != nil {
		t.Error("should not raise error.")
	}

	if len(meters) != 3 {
		t.Errorf("meters length should be 3, but %v", meters)
	}
}

func TestParseLoadavg(t *testing.T) {
	procLoadavg := "0.06 0.04 0.01 1/179 24573"

	env := environment.NewEnvironment()
	collector := NewLoadavgCollector(env, "resourceId")
	values := collector.parseProcLoadavg(procLoadavg)

	if len(values) != 3 {
		t.Errorf("Values length should be 3, but %v", values)
	}

	expectedValues := []string{"0.06", "0.04", "0.01"}

	for i := 0; i < len(values); i++ {
		if values[i] != expectedValues[i] {
			t.Errorf("values[i] should be %v, but %v", expectedValues[i], values[i])
		}
	}
}
