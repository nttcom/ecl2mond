// +build linux

package linux

import (
	"testing"

	"github.com/nttcom/ecl2mond/environment"
)

func NewMemoryCollector(env *environment.Environment, resourceID string) *MemoryCollector {
	return &MemoryCollector{
		Env:        env,
		ResourceID: resourceID,
	}
}

func TestGetMemoryMeters(t *testing.T) {
	env := environment.NewEnvironment()
	env.RootPath = "./"
	collector := NewMemoryCollector(env, "resourceId")

	meters, err := collector.Collect()
	if err != nil {
		t.Error("should not raise error.")
	}

	if len(meters) != 10 {
		t.Errorf("meters length should be 10, but %v", len(meters))
	}
}

func TestParseProcMeminfo(t *testing.T) {
	lines := []string{
		"MemTotal:         469176 kB",
		"MemFree:          109552 kB",
		"Buffers:           38760 kB",
		"Cached:           125772 kB",
		"SwapCached:         2884 kB",
		"Active:           116688 kB",
		"Inactive:         176320 kB",
		"Active(anon):      14928 kB",
		"Inactive(anon):   114148 kB",
		"Active(file):     101760 kB",
		"Inactive(file):    62172 kB",
		"Unevictable:           0 kB",
		"Mlocked:               0 kB",
		"SwapTotal:        950268 kB",
		"SwapFree:         581000 kB",
		"Dirty:                40 kB",
		"Writeback:             0 kB",
		"AnonPages:        126392 kB",
		"Mapped:            13536 kB",
		"Shmem:               596 kB",
		"Slab:              43188 kB",
		"SReclaimable:      17504 kB",
		"SUnreclaim:        25684 kB",
		"KernelStack:        2736 kB",
		"PageTables:        11888 kB",
		"NFS_Unstable:          0 kB",
		"Bounce:                0 kB",
		"WritebackTmp:          0 kB",
		"CommitLimit:     1184856 kB",
		"Committed_AS:    1086936 kB",
		"VmallocTotal:   34359738367 kB",
		"VmallocUsed:        4180 kB",
		"VmallocChunk:   34359728436 kB",
		"HardwareCorrupted:     0 kB",
		"AnonHugePages:         0 kB",
		"HugePages_Total:       0",
		"HugePages_Free:        0",
		"HugePages_Rsvd:        0",
		"HugePages_Surp:        0",
		"Hugepagesize:       2048 kB",
		"DirectMap4k:        8128 kB",
		"DirectMap2M:      483328 kB",
	}

	env := environment.NewEnvironment()
	collector := NewMemoryCollector(env, "resourceId")
	values, err := collector.parseProcMeminfo(lines)

	if err != nil {
		t.Error("should not be raise")
	}

	if len(values) != 10 {
		t.Errorf("values length should be 11, but %v", values)
	}

	expected := []string{"469176", "109552", "38760", "125772", "2884", "116688", "176320", "950268", "581000", "195092"}

	for i, value := range values {
		if value != expected[i] {
			t.Errorf("value should be %v, but %v", expected[i], value)
		}
	}
}

func TestGetMeminfo(t *testing.T) {
	line := "MemTotal: 469176 kB"
	env := environment.NewEnvironment()
	collector := NewMemoryCollector(env, "resourceId")

	name, value := collector.getMeminfoByLine(line)

	if name != "MemTotal" {
		t.Errorf("name should be MemTotal, but %v", name)
	}

	if value != "469176" {
		t.Errorf("value should be 469176, but %v", value)
	}
}
