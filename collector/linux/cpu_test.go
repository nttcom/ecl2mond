// +build linux

package linux

import (
	"testing"

	"github.com/nttcom/ecl2mond/environment"
	"github.com/nttcom/ecl2mond/meter"
)

func NewCPUCollector(env *environment.Environment, resourceID string) *CPUCollector {
	return &CPUCollector{
		Env:        env,
		ResourceID: resourceID,
	}
}

func TestCPUMeters(t *testing.T) {
	env := environment.NewEnvironment()
	env.RootPath = "./"
	collector := NewCPUCollector(env, "resourceId")

	_, err := collector.Collect()

	if err != nil {
		t.Error("should not raise error.")
	}
}

func TestParseProcStat(t *testing.T) {
	lines := []string{
		"cpu  898 0 1180 329993 3767 128 135 0 0",
		"cpu0 462 0 662 163256 1697 128 82 0 0",
		"cpu1 435 0 517 166736 2069 0 52 0 0",
		"intr 184366 132 7 0 0 0 0 0 0 0 0 0 0 112 0 22386 0 9484 0 0 1910 735 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0",
		"ctxt 282734",
		"btime 1486701364",
		"processes 4933",
		"procs_running 1",
		"procs_blocked 0",
		"softirq 265200 0 83854 680 26544 22422 0 2 33397 1701 96600",
	}

	env := environment.NewEnvironment()
	collector := NewCPUCollector(env, "resourceId")
	values := collector.parseProcStat(lines)

	if len(values) != 9 {
		t.Errorf("values length should be 9, but %v", len(values))
	}

	expectedValues := []string{"898", "0", "1180", "329993", "3767", "128", "135", "0", "0"}
	for i, value := range values {
		if value != expectedValues[i] {
			t.Errorf("value should be %v, but %v", value, expectedValues[i])
		}
	}
}

func TestGetCpuNum(t *testing.T) {
	// Case 1
	lines := []string{
		"cpu  898 0 1180 329993 3767 128 135 0 0",
		"cpu0 462 0 662 163256 1697 128 82 0 0",
		"cpu1 435 0 517 166736 2069 0 52 0 0",
		"intr 184366 132 7 0 0 0 0 0 0 0 0 0 0 112 0 22386 0 9484 0 0 1910 735 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0",
		"ctxt 282734",
		"btime 1486701364",
		"processes 4933",
		"procs_running 1",
		"procs_blocked 0",
		"softirq 265200 0 83854 680 26544 22422 0 2 33397 1701 96600",
	}
	env := environment.NewEnvironment()
	collector := NewCPUCollector(env, "resourceId")
	cpuNum := collector.getCPUNum(lines)

	if cpuNum != 2 {
		t.Errorf("cpuNumshould be 2, but %v", cpuNum)
	}

	// Case 2
	lines = []string{
		"cpu  898 0 1180 329993 3767 128 135 0 0",
		"cpu0 462 0 662 163256 1697 128 82 0 0",
		"intr 184366 132 7 0 0 0 0 0 0 0 0 0 0 112 0 22386 0 9484 0 0 1910 735 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0",
		"ctxt 282734",
		"btime 1486701364",
		"processes 4933",
		"procs_running 1",
		"procs_blocked 0",
		"softirq 265200 0 83854 680 26544 22422 0 2 33397 1701 96600",
	}

	cpuNum = collector.getCPUNum(lines)
	if cpuNum != 1 {
		t.Errorf("cpuNum should be 1, but %v", cpuNum)
	}
}

func TestGetTotal(t *testing.T) {
	values := []string{"898", "0", "1180", "329993", "3767", "128", "135", "0", "0"}
	env := environment.NewEnvironment()
	collector := NewCPUCollector(env, "resourceId")
	total, err := collector.getTotal(values)

	if err != nil {
		t.Error("should not raise error")
	}

	if total != 336101 {
		t.Errorf("total should be 336101, but %v", total)
	}
}

func TestGetBasedMeters(t *testing.T) {
	expectedMeters := []*meter.Meter{
		meter.NewMeter("cpu.user.percents", "resourceId", "percent", "delta"),
		meter.NewMeter("cpu.nice.percents", "resourceId", "percent", "delta"),
		meter.NewMeter("cpu.system.percents", "resourceId", "percent", "delta"),
		meter.NewMeter("cpu.idle.percents", "resourceId", "percent", "delta"),
		meter.NewMeter("cpu.iowait.percents", "resourceId", "percent", "delta"),
		meter.NewMeter("cpu.irq.percents", "resourceId", "percent", "delta"),
		meter.NewMeter("cpu.softirq.percents", "resourceId", "percent", "delta"),
	}

	env := environment.NewEnvironment()
	collector := NewCPUCollector(env, "resourceId")
	meters := collector.getCPUBasedMeters(7)

	if len(meters) != 7 {
		t.Errorf("meters length should be 7, but %v", len(meters))
	}

	for i, meter := range meters {
		expectedMeter := expectedMeters[i]

		if meter.Name != expectedMeter.Name {
			t.Errorf("meter.Name should be %v, but %v", expectedMeter.Name, meter.Name)
		}

		if meter.ResourceID != expectedMeter.ResourceID {
			t.Errorf("meter.ResourceID should be %v, but %v", expectedMeter.ResourceID, meter.ResourceID)
		}

		if meter.Unit != expectedMeter.Unit {
			t.Errorf("meter.Unit should be %v, but %v", expectedMeter.Unit, meter.Unit)
		}

		if meter.Type != expectedMeter.Type {
			t.Errorf("meter.Type should be %v, but %v", expectedMeter.Type, meter.Type)
		}
	}
}

func TestGetBasedMetersOverNum(t *testing.T) {
	env := environment.NewEnvironment()
	collector := NewCPUCollector(env, "resourceId")
	meters := collector.getCPUBasedMeters(17)

	if len(meters) != 9 {
		t.Error("meters length should be 9")
	}
}
