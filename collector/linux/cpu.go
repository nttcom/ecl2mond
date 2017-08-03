// +build linux

package linux

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nttcom/ecl2mond/environment"
	"github.com/nttcom/ecl2mond/logging"
	"github.com/nttcom/ecl2mond/meter"
	"github.com/nttcom/ecl2mond/store"
)

var cpuLogger = logging.NewLogger("cpu")

// CPUCollector is collector for CPU usage.
type CPUCollector struct {
	Env        *environment.Environment
	ResourceID string
}

// Collect gets meters for CPU usage.
func (collector *CPUCollector) Collect() ([]*meter.Meter, error) {
	store := store.NewStore(collector.Env)
	data, createdAt, err := store.GetData("cpu")
	previousValues := data["cpu"]
	notExistPrevious := err != nil
	expiredPrevious := time.Now().Sub(createdAt) >= 25*time.Hour

	lines, err := collector.getProcStatLines()
	if err != nil {
		cpuLogger.Error("Failed to read /proc/stat file.")
		return nil, err
	}

	values := collector.parseProcStat(lines)

	cpuLogger.Debug(fmt.Sprintf("previousValues: %s", previousValues))
	cpuLogger.Debug(fmt.Sprintf("currentValue: %s", values))

	// 今回の収集内容を保存
	store.SetData("cpu", map[string][]string{"cpu": values})
	meters := collector.getCPUBasedMeters(len(values))

	if notExistPrevious || expiredPrevious || len(previousValues) != len(values) {
		if notExistPrevious {
			cpuLogger.Info("No previous samples.")
		} else if expiredPrevious {
			cpuLogger.Info("Previous data has expired(25h).")
		}

		return meters, nil
	}

	// 計算
	previousTotal, err := collector.getTotal(previousValues)
	if err != nil {
		cpuLogger.Error("Failed to parse /proc/stat file.")
		return nil, err
	}
	cpuLogger.Debug(fmt.Sprintf("previousTotal: %f", previousTotal))
	currentTotal, err := collector.getTotal(values)
	if err != nil {
		cpuLogger.Error("Failed to parse /proc/stat file.")
		return nil, err
	}
	cpuLogger.Debug(fmt.Sprintf("currentTotal: %f", currentTotal))

	cpuNum := collector.getCPUNum(lines)
	cpuLogger.Debug(fmt.Sprintf("cpuNum: %d", cpuNum))

	for i, meter := range meters {
		previousValue, err := strconv.ParseFloat(previousValues[i], 64)
		if err != nil {
			cpuLogger.Error("Failed to parse /proc/stat file.")
			return nil, err
		}

		currentValue, err := strconv.ParseFloat(values[i], 64)
		if err != nil {
			cpuLogger.Error("Failed to parse /proc/stat file.")
			return nil, err
		}

		volume := (currentValue - previousValue) * 100.0 * float64(cpuNum) / (currentTotal - previousTotal)
		meter.SetVolume(strconv.FormatFloat(volume, 'f', 1, 64))
	}

	return meters, nil
}

func (collector *CPUCollector) getProcStatLines() ([]string, error) {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	cpuLogger.Debug(fmt.Sprintf("/proc/stat:\n%s", data))

	return strings.Split(string(data), "\n"), nil
}

func (collector *CPUCollector) parseProcStat(lines []string) []string {
	cpuline := lines[0]
	return strings.Fields(cpuline)[1:]
}

func (collector *CPUCollector) getCPUNum(lines []string) uint {
	var cpuNum uint
	regexp := regexp.MustCompile(`^cpu\d+\s`)

	for _, line := range lines {
		if regexp.MatchString(line) {
			cpuNum++
		}
	}

	return cpuNum
}

func (collector *CPUCollector) getTotal(values []string) (float64, error) {
	var total float64

	for _, strValue := range values {
		value, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return 0, err
		}

		total += value
	}
	return total, nil
}

func (collector *CPUCollector) getCPUBasedMeters(num int) []*meter.Meter {
	index := 9
	if num < index {
		index = num
	}

	resourceID := collector.ResourceID
	base := []*meter.Meter{
		meter.NewMeter("cpu.user.percents", resourceID, "percent", "delta"),
		meter.NewMeter("cpu.nice.percents", resourceID, "percent", "delta"),
		meter.NewMeter("cpu.system.percents", resourceID, "percent", "delta"),
		meter.NewMeter("cpu.idle.percents", resourceID, "percent", "delta"),
		meter.NewMeter("cpu.iowait.percents", resourceID, "percent", "delta"),
		meter.NewMeter("cpu.irq.percents", resourceID, "percent", "delta"),
		meter.NewMeter("cpu.softirq.percents", resourceID, "percent", "delta"),
		meter.NewMeter("cpu.steal.percents", resourceID, "percent", "delta"),
		meter.NewMeter("cpu.guest.percents", resourceID, "percent", "delta"),
	}

	return base[:index]
}
