// +build linux

package linux

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/nttcom/ecl2mond/environment"
	"github.com/nttcom/ecl2mond/logging"
	"github.com/nttcom/ecl2mond/meter"
)

var memoryLogger = logging.NewLogger("memory")

var items = []string{
	"MemTotal",
	"MemFree",
	"Buffers",
	"Cached",
	"SwapCached",
	"Active",
	"Inactive",
	"SwapTotal",
	"SwapFree",
}

// MemoryCollector is collector for memory.
type MemoryCollector struct {
	Env        *environment.Environment
	ResourceID string
}

// Collect gets meters for memory.
func (collector *MemoryCollector) Collect() ([]*meter.Meter, error) {
	lines, err := collector.getProcMemInfoLines()
	if err != nil {
		memoryLogger.Error("Failed to read /proc/meminfo file.")
		return nil, err
	}

	values, err := collector.parseProcMeminfo(lines)
	if err != nil {
		memoryLogger.Error("Failed to parse /proc/meminfo file.")
		return nil, err
	}

	meters := collector.getBasedMeters()
	for i, meter := range meters {
		meter.SetVolume(values[i])
	}

	return meters, nil
}

func (collector *MemoryCollector) getProcMemInfoLines() ([]string, error) {
	out, err := ioutil.ReadFile("/proc/meminfo")
	memoryLogger.Debug(fmt.Sprintf("/proc/meminfo:\n%s", out))
	if err != nil {
		return nil, err
	}

	return strings.Split(string(out), "\n"), nil
}

func (collector *MemoryCollector) parseProcMeminfo(lines []string) ([]string, error) {

	data := make(map[string]string)

	for _, line := range lines {
		name, value := collector.getMeminfoByLine(line)

		if collector.checkExistItem(name) {
			data[name] = value
		}
	}

	values := collector.getMemoryValues(data)
	usedTotal, err := collector.getUsedTotalValue(data)
	if err != nil {
		return nil, err
	}
	values = append(values, usedTotal)

	return values, nil
}

func (collector *MemoryCollector) getUsedTotalValue(data map[string]string) (string, error) {
	memTotal, err := strconv.ParseFloat(data["MemTotal"], 64)
	if err != nil {
		return "", err
	}

	memFree, err := strconv.ParseFloat(data["MemFree"], 64)
	if err != nil {
		return "", err
	}

	buffers, err := strconv.ParseFloat(data["Buffers"], 64)
	if err != nil {
		return "", err
	}

	cached, err := strconv.ParseFloat(data["Cached"], 64)
	if err != nil {
		return "", err
	}

	usedTotal := strconv.FormatFloat(memTotal-memFree-buffers-cached, 'f', 0, 64)
	return usedTotal, nil
}

func (collector *MemoryCollector) getMemoryValues(data map[string]string) []string {
	var values []string
	for _, item := range items {
		values = append(values, data[item])
	}

	return values
}

func (collector *MemoryCollector) getMeminfoByLine(line string) (string, string) {
	cols := strings.Fields(line)

	if len(cols) < 2 {
		return "", ""
	}

	name := cols[0][0:(len(cols[0]) - 1)]
	value := cols[1]
	return name, value
}

func (collector *MemoryCollector) checkExistItem(name string) bool {
	for _, item := range items {
		if item == name {
			return true
		}
	}
	return false
}

func (collector *MemoryCollector) getBasedMeters() []*meter.Meter {
	resourceID := collector.ResourceID
	return []*meter.Meter{
		meter.NewMeter("memory.memtotal.kilobytes", resourceID, "kilobyte", "gauge"),
		meter.NewMeter("memory.memfree.kilobytes", resourceID, "kilobyte", "gauge"),
		meter.NewMeter("memory.buffers.kilobytes", resourceID, "kilobyte", "gauge"),
		meter.NewMeter("memory.cached.kilobytes", resourceID, "kilobyte", "gauge"),
		meter.NewMeter("memory.swapcached.killobytes", resourceID, "kilobyte", "gauge"),
		meter.NewMeter("memory.active.killobytes", resourceID, "kilobyte", "gauge"),
		meter.NewMeter("memory.inactive.killobytes", resourceID, "kilobyte", "gauge"),
		meter.NewMeter("memory.swaptotal.killobytes", resourceID, "kilobyte", "gauge"),
		meter.NewMeter("memory.swapfree.killobytes", resourceID, "kilobyte", "gauge"),
		meter.NewMeter("memory.usedtotal.killobytes", resourceID, "kilobyte", "gauge"),
	}
}
