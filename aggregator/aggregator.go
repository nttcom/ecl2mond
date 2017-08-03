package aggregator

import (
	"fmt"

	"github.com/nttcom/ecl2mond/collector"
	"github.com/nttcom/ecl2mond/config"
	"github.com/nttcom/ecl2mond/environment"
	"github.com/nttcom/ecl2mond/meter"
	"github.com/nttcom/ecl2mond/logging"
)

var logger = logging.NewLogger("aggregator")

// Aggregator is collectoring meters struct.
type Aggregator struct {
	Env        *environment.Environment
	Config     *config.Config
	ResourceID string
}

// NewAggregator returns Aggregator struct.
func NewAggregator(env *environment.Environment, config *config.Config, resourceID string) *Aggregator {
	return &Aggregator{
		Env:        env,
		Config:     config,
		ResourceID: resourceID,
	}
}

// GetCollectedMeters returns collected meters.
func (ag *Aggregator) GetCollectedMeters() ([]*meter.Meter, error) {
	var meters []*meter.Meter
	for _, collector := range ag.GetCollectors() {
		tempMeters, err := collector.Collect()
		if err == nil {
			meters = append(meters, tempMeters...)
		}
	}
	return ag.getRequireMeters(meters), nil
}

// GetCollectors returns collector's interfaces.
func (ag *Aggregator) GetCollectors() []collector.Collector {
	var collectors []collector.Collector
	for _, meterType := range ag.Config.MeterTypes() {
		collector, ok := ag.collectorMapping()[meterType]
		if ok {
			collectors = append(collectors, collector)
		} else {
			logger.Error(fmt.Sprintf("Invalid meter type. Name: %s", meterType))
		}
	}

	return collectors
}

func (ag *Aggregator) getRequireMeters(meters []*meter.Meter) []*meter.Meter {
	var ret []*meter.Meter

	for _, meterName := range ag.Config.Meters {
		i := ag.getMeterIndexByName(meterName, meters)
		if i != -1 {
			ret = append(ret, meters[i])
		}
	}

	return ret
}

func (ag *Aggregator) getMeterIndexByName(name string, meters []*meter.Meter) int {
	for i, meter := range meters {
		if meter.Name == name {
			return i
		}
	}

	return -1
}
