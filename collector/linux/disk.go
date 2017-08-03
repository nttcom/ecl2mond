// +build linux

package linux

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/nttcom/ecl2mond/environment"
	"github.com/nttcom/ecl2mond/logging"
	"github.com/nttcom/ecl2mond/meter"
	"github.com/nttcom/ecl2mond/store"
	"github.com/nttcom/ecl2mond/util"
)

var diskLogger = logging.NewLogger("disk")

// DiskCollector is collector for disk.
type DiskCollector struct {
	Env        *environment.Environment
	ResourceID string
}

// Collect gets meters for disk.
func (collector *DiskCollector) Collect() ([]*meter.Meter, error) {
	store := store.NewStore(collector.Env)
	previousData, createdAt, err := store.GetData("disk")
	notExistPrevious := err != nil
	expiredPrevious := time.Now().Sub(createdAt) >= 25*time.Hour

	lines, err := collector.getProcDiskStatLines()
	if err != nil {
		diskLogger.Error("Failed to read /proc/diskstats file.")
		return nil, err
	}

	currentData := collector.parseProcDiskStat(lines)

	store.SetData("disk", currentData)

	diskLogger.Debug(fmt.Sprintf("previousData: %s", previousData))
	diskLogger.Debug(fmt.Sprintf("currentData: %s", currentData))

	var allMeters []*meter.Meter

	for device, currentValues := range currentData {
		meters := collector.getBasedMetersByDevice(device)
		previousValues := previousData[device]

		if notExistPrevious || expiredPrevious || len(previousValues) != len(currentValues) {
			allMeters = append(allMeters, meters...)
			continue
		}

		for i, meter := range meters {
			previousValue, err := strconv.ParseFloat(previousValues[i], 64)
			if err != nil {
				diskLogger.Error("Failed to parse /proc/diskstats file.")
				return nil, err
			}

			currentValue, err := strconv.ParseFloat(currentValues[i], 64)
			if err != nil {
				diskLogger.Error("Failed to parse /proc/diskstats file.")
				return nil, err
			}

			calculatedVolume := calculatedVolume(meter, currentValue, previousValue)
			meter.SetVolume(strconv.FormatFloat(calculatedVolume, 'f', 1, 64))
			allMeters = append(allMeters, meter)
		}
	}

	if notExistPrevious {
		diskLogger.Info("No previous samples.")
	} else if expiredPrevious {
		diskLogger.Info("Previous data has expired(25h).")
	}

	return allMeters, nil
}


func calculatedVolume(meter *meter.Meter, currentValue float64, previousValue float64) (ret_volume float64) {

	if meter.Type == "gauge" {
		ret_volume = currentValue
	} else {
		ret_volume = currentValue - previousValue
	}
	return ret_volume
}

func (collector *DiskCollector) getProcDiskStatLines() ([]string, error) {
	data, err := ioutil.ReadFile("/proc/diskstats")
	if err != nil {
		return nil, err
	}

	diskLogger.Debug(fmt.Sprintf("/proc/diskstats: %s", data))

	return strings.Split(string(data), "\n"), nil
}

func (collector *DiskCollector) parseProcDiskStat(lines []string) map[string][]string {
	result := make(map[string][]string)

	for _, line := range lines {
		cols := strings.Fields(line)

		if len(cols) < 3 {
			continue
		}

		name := cols[2]
		values := cols[3:]

		if util.CheckAllZeroValues(values) {
			continue
		}

		result[name] = values
	}

	return result
}

func (collector *DiskCollector) getBasedMetersByDevice(device string) []*meter.Meter {
	resourceID := collector.ResourceID
	return []*meter.Meter{
		meter.NewMeter(fmt.Sprintf("disk.%s.reads.completed.count", device), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("disk.%s.reads.merged.count", device), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("disk.%s.reads.sectors.count", device), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("disk.%s.reads.milliseconds", device), resourceID, "millisecond", "delta"),
		meter.NewMeter(fmt.Sprintf("disk.%s.writes.completed.count", device), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("disk.%s.writes.merged.count", device), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("disk.%s.writes.sectors.count", device), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("disk.%s.writes.milliseconds", device), resourceID, "millisecond", "delta"),
		meter.NewMeter(fmt.Sprintf("disk.%s.currently.ios.count", device), resourceID, "count", "gauge"),
		meter.NewMeter(fmt.Sprintf("disk.%s.ios.milliseconds", device), resourceID, "millisecond", "delta"),
		meter.NewMeter(fmt.Sprintf("disk.%s.weighted.ios.milliseconds", device), resourceID, "millisecond", "delta"),
	}
}
