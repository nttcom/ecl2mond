// +build linux

package linux

import (
	"testing"

	"github.com/nttcom/ecl2mond/environment"
	"github.com/nttcom/ecl2mond/meter"
)

func NewDiskCollector(env *environment.Environment, resourceID string) *DiskCollector {
	return &DiskCollector{
		Env:        env,
		ResourceID: resourceID,
	}
}

func TestGetDiskMeters(t *testing.T) {
	env := environment.NewEnvironment()
	env.RootPath = "./"
	collector := NewDiskCollector(env, "resourceId")

	_, err := collector.Collect()

	if err != nil {
		t.Errorf("should not raise error. %s", err)
	}
}

func TestParseProcDiskStat(t *testing.T) {
	lines := []string{
		"   1       0 ram0 0 0 0 0 0 0 0 0 0 0 0",
		"   1       1 ram1 0 0 0 0 0 0 0 0 0 0 0",
		"   1       2 ram2 0 0 0 0 0 0 0 0 0 0 0",
		"   1       3 ram3 0 0 0 0 0 0 0 0 0 0 0",
		"   1       4 ram4 0 0 0 0 0 0 0 0 0 0 0",
		"   1       5 ram5 0 0 0 0 0 0 0 0 0 0 0",
		"   1       6 ram6 0 0 0 0 0 0 0 0 0 0 0",
		"   1       7 ram7 0 0 0 0 0 0 0 0 0 0 0",
		"   1       8 ram8 0 0 0 0 0 0 0 0 0 0 0",
		"   1       9 ram9 0 0 0 0 0 0 0 0 0 0 0",
		"   1      10 ram10 0 0 0 0 0 0 0 0 0 0 0",
		"   1      11 ram11 0 0 0 0 0 0 0 0 0 0 0",
		"   1      12 ram12 0 0 0 0 0 0 0 0 0 0 0",
		"   1      13 ram13 0 0 0 0 0 0 0 0 0 0 0",
		"   1      14 ram14 0 0 0 0 0 0 0 0 0 0 0",
		"   1      15 ram15 0 0 0 0 0 0 0 0 0 0 0",
		"   7       0 loop0 0 0 0 0 0 0 0 0 0 0 0",
		"   7       1 loop1 0 0 0 0 0 0 0 0 0 0 0",
		"   7       2 loop2 0 0 0 0 0 0 0 0 0 0 0",
		"   7       3 loop3 0 0 0 0 0 0 0 0 0 0 0",
		"   7       4 loop4 0 0 0 0 0 0 0 0 0 0 0",
		"   7       5 loop5 0 0 0 0 0 0 0 0 0 0 0",
		"   7       6 loop6 0 0 0 0 0 0 0 0 0 0 0",
		"   7       7 loop7 0 0 0 0 0 0 0 0 0 0 0",
		"   8       0 sda 83402 64862 4328066 77552 96113 446026 4241580 410535 0 79546 487009",
		"   8       1 sda1 763 749 6216 215 14 1 84 11 0 219 225",
		"   8       2 sda2 82487 64113 4320634 77303 83488 446025 4241496 409243 0 78153 485471",
		" 253       0 dm-0 144438 0 4297994 165933 439219 0 3513056 6160484 0 77018 6326417",
		" 253       1 dm-1 2666 0 21328 4390 91055 0 728440 2198816 0 8561 2203207",
		" 253       2 dm-2 0 0 0 0 0 0 0 0 0 0 0",
	}

	env := environment.NewEnvironment()
	collector := NewDiskCollector(env, "resourceId")
	result := collector.parseProcDiskStat(lines)

	// sda
	sdaValues := result["sda"]
	exptectedValues := []string{"83402", "64862", "4328066", "77552", "96113", "446026", "4241580", "410535", "0", "79546", "487009"}
	if len(sdaValues) != 11 {
		t.Errorf("sdaValues length should be 11, but %v", sdaValues)
	}

	for i, sdaValue := range sdaValues {
		if sdaValue != exptectedValues[i] {
			t.Errorf("sdaValue should be %v, but %v", exptectedValues[i], sdaValue)
		}
	}

	// sda1
	sda1Values := result["sda1"]
	exptectedValues = []string{"763", "749", "6216", "215", "14", "1", "84", "11", "0", "219", "225"}
	if len(sda1Values) != 11 {
		t.Errorf("sda1Values length should be 11, but %v", sda1Values)
	}
	for i, sda1Value := range sda1Values {
		if sda1Value != exptectedValues[i] {
			t.Errorf("sda1Value should be %v, but %v", exptectedValues[i], sda1Value)
		}
	}

	// sda2
	sda2Values := result["sda2"]
	exptectedValues = []string{"82487", "64113", "4320634", "77303", "83488", "446025", "4241496", "409243", "0", "78153", "485471"}
	if len(sda2Values) != 11 {
		t.Errorf("sda2Values length should be 11, but %v", sda2Values)
	}
	for i, sda2Value := range sda2Values {
		if sda2Value != exptectedValues[i] {
			t.Errorf("sda2Value should be %v, but %v", exptectedValues[i], sda2Value)
		}
	}

	// dm-0
	dm0Values := result["dm-0"]
	exptectedValues = []string{"144438", "0", "4297994", "165933", "439219", "0", "3513056", "6160484", "0", "77018", "6326417"}
	if len(dm0Values) != 11 {
		t.Errorf("dm0Values length should be 11, but %v", dm0Values)
	}
	for i, dm0Value := range dm0Values {
		if dm0Value != exptectedValues[i] {
			t.Errorf("dm0Value should be %v, but %v", exptectedValues[i], dm0Value)
		}
	}

	// dm-1
	dm1Values := result["dm-1"]
	exptectedValues = []string{"2666", "0", "21328", "4390", "91055", "0", "728440", "2198816", "0", "8561", "2203207"}
	if len(dm1Values) != 11 {
		t.Errorf("dm1Values length should be 11, but %v", dm1Values)
	}
	for i, dm1Value := range dm1Values {
		if dm1Value != exptectedValues[i] {
			t.Errorf("dm1Value should be %v, but %v", exptectedValues[i], dm1Value)
		}
	}

	// dm-2
	dm2Values := result["dm-2"]
	if len(dm2Values) != 0 {
		t.Errorf("dm2Values should be empry, but %v", dm2Values)
	}
}

func TestGetBasedMetersByDevice(t *testing.T) {
	env := environment.NewEnvironment()
	collector := NewDiskCollector(env, "resourceId")

	meters := collector.getBasedMetersByDevice("sda")

	if len(meters) != 11 {
		t.Errorf("meters length should be 11, but %v", meters)
	}

	expectedMeters := []*meter.Meter{
		meter.NewMeter("disk.sda.reads.completed.count", "resourceId", "count", "delta"),
		meter.NewMeter("disk.sda.reads.merged.count", "resourceId", "count", "delta"),
		meter.NewMeter("disk.sda.reads.sectors.count", "resourceId", "count", "delta"),
		meter.NewMeter("disk.sda.reads.milliseconds", "resourceId", "millisecond", "delta"),
		meter.NewMeter("disk.sda.writes.completed.count", "resourceId", "count", "delta"),
		meter.NewMeter("disk.sda.writes.merged.count", "resourceId", "count", "delta"),
		meter.NewMeter("disk.sda.writes.sectors.count", "resourceId", "count", "delta"),
		meter.NewMeter("disk.sda.writes.milliseconds", "resourceId", "millisecond", "delta"),
		meter.NewMeter("disk.sda.currently.ios.count", "resourceId", "count", "gauge"),
		meter.NewMeter("disk.sda.ios.milliseconds", "resourceId", "millisecond", "delta"),
		meter.NewMeter("disk.sda.weighted.ios.milliseconds", "resourceId", "millisecond", "delta"),
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
