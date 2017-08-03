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

var networkLogger = logging.NewLogger("network")

// NetworkCollector is collector for network.
type NetworkCollector struct {
	Env        *environment.Environment
	ResourceID string
}

// Collect gets meters for network.
func (collector *NetworkCollector) Collect() ([]*meter.Meter, error) {
	store := store.NewStore(collector.Env)
	previousData, createdAt, err := store.GetData("network")
	notExistPrevious := err != nil
	expiredPrevious := time.Now().Sub(createdAt) >= 25*time.Hour

	lines, err := collector.getProcNetDevLines()
	if err != nil {
		networkLogger.Error("Failed to read /proc/net/dev file.")
		return nil, err
	}

	currentData := collector.parseProcNetDev(lines)
	store.SetData("network", currentData)

	networkLogger.Debug(fmt.Sprintf("previousData: %s", previousData))
	networkLogger.Debug(fmt.Sprintf("currentData: %s", currentData))

	var allMeters []*meter.Meter
	for interfaceName, currentValues := range currentData {
		meters := collector.getBasedMetersByInterfaceName(interfaceName)
		previousValues := previousData[interfaceName]

		if notExistPrevious || expiredPrevious || len(previousValues) != len(currentValues) {
			allMeters = append(allMeters, meters...)
			continue
		}

		for i, meter := range meters {
			previousValue, err := strconv.ParseFloat(previousValues[i], 64)
			if err != nil {
				networkLogger.Error("Failed to parse /proc/net/dev file.")
				return nil, err
			}

			currentValue, err := strconv.ParseFloat(currentValues[i], 64)
			if err != nil {
				networkLogger.Error("Failed to parse /proc/net/dev file.")
				return nil, err
			}

			volume := currentValue - previousValue
			meter.SetVolume(strconv.FormatFloat(volume, 'f', 1, 64))
			allMeters = append(allMeters, meter)
		}
	}

	if notExistPrevious {
		networkLogger.Info("No previous samples.")
	} else if expiredPrevious {
		networkLogger.Info("Previous data has expired(25h).")
	}

	return allMeters, nil
}

func (collector *NetworkCollector) getProcNetDevLines() ([]string, error) {
	data, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, err
	}

	networkLogger.Debug(fmt.Sprintf("/proc/net/dev:\n%s", data))

	return strings.Split(string(data), "\n"), nil
}

func (collector *NetworkCollector) parseProcNetDev(lines []string) map[string][]string {
	data := make(map[string][]string)

	for _, line := range lines {
		sp := strings.Split(line, ":")
		if len(sp) != 2 {
			continue
		}

		name := strings.Trim(sp[0], " ")
		if name == "lo" {
			continue
		}

		values := strings.Fields(sp[1])

		if util.CheckAllZeroValues(values) {
			continue
		}

		data[name] = values
	}

	return data
}

func (collector *NetworkCollector) getBasedMetersByInterfaceName(name string) []*meter.Meter {
	resourceID := collector.ResourceID
	return []*meter.Meter{
		meter.NewMeter(fmt.Sprintf("network.%s.receive.bytes", name), resourceID, "byte", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.receive.packets.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.receive.errs.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.receive.drop.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.receive.fifo.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.receive.frame.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.receive.compressed.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.receive.multicast.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.transmit.bytes", name), resourceID, "byte", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.transmit.packets.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.transmit.errs.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.transmit.drop.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.transmit.fifo.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.transmit.colls.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.transmit.carrier.count", name), resourceID, "count", "delta"),
		meter.NewMeter(fmt.Sprintf("network.%s.transmit.compressed.count", name), resourceID, "count", "delta"),
	}
}
