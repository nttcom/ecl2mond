package aggregator

import (
	"testing"

	"github.com/nttcom/ecl2mond/config"
	"github.com/nttcom/ecl2mond/environment"
)

func TestGetCollectedMeters(t *testing.T) {
	env := environment.NewEnvironment()
	config := config.DefaultConfig
	config.Meters = []string{
		"cpuusage.user",
		"cpuusage.idle",
		"memory.used",
		"memory",
		"disk.reads",
		"disk.writes",
	}

	aggregator := NewAggregator(env, config, "resourceID")
	meters, err := aggregator.GetCollectedMeters()

	if err != nil {
		t.Errorf("should not raise error.")
	}

	if len(meters) != 0 {
		t.Errorf("should be empty, but %v", meters)
	}
}
