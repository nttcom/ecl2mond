// +build linux

package linux

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nttcom/ecl2mond/environment"
	"github.com/nttcom/ecl2mond/logging"
	"github.com/nttcom/ecl2mond/meter"
)

var loadavgLogger = logging.NewLogger("loadavg")

// LoadavgCollector is collector for load average.
type LoadavgCollector struct {
	Env        *environment.Environment
	ResourceID string
}

// Collect gets meters for load average.
func (collector *LoadavgCollector) Collect() ([]*meter.Meter, error) {
	data, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		loadavgLogger.Error("Failed to read /proc/loadavg file.")
		return nil, err
	}

	loadavgLogger.Debug(fmt.Sprintf("/proc/loadavg:\n%s", data))

	procLoadavg := string(data)
	values := collector.parseProcLoadavg(procLoadavg)
	meters := collector.getLoadavgBaseMeters()

	for i, meter := range meters {
		meter.SetVolume(values[i])
	}

	return meters, nil
}

func (collector *LoadavgCollector) parseProcLoadavg(procLoadavg string) []string {
	cols := strings.Split(procLoadavg, " ")
	values := []string{cols[0], cols[1], cols[2]}
	return values
}

func (collector *LoadavgCollector) getLoadavgBaseMeters() []*meter.Meter {
	return []*meter.Meter{
		meter.NewMeter("loadavg.1.count", collector.ResourceID, "count", "gauge"),
		meter.NewMeter("loadavg.5.count", collector.ResourceID, "count", "gauge"),
		meter.NewMeter("loadavg.15.count", collector.ResourceID, "count", "gauge"),
	}
}
