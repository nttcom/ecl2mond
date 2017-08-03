// +build linux

package linux

import (
	"testing"

	"github.com/nttcom/ecl2mond/environment"
)

func NewNetworkCollector(env *environment.Environment, resourceID string) *NetworkCollector {
	return &NetworkCollector{
		Env:        env,
		ResourceID: resourceID,
	}
}

func TestGetNetworkMeters(t *testing.T) {
	env := environment.NewEnvironment()
	env.RootPath = "./"
	collector := NewNetworkCollector(env, "resourceId")

	_, err := collector.Collect()

	if err != nil {
		t.Error("should not raise error.")
	}

}

func TestParseProcNetDev(t *testing.T) {
	lines := []string{
		"Inter-|   Receive                                                |  Transmit",
		" face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed",
		"    lo: 2318141   35273    0    0    0     0          0         0  2318141   35273    0    0    0     0       0          0",
		"  eth0: 61058551   44173    0    0    0     0          0         0   422875    5776    0    0    0     0       0          0",
		"  eth1: 284076872 1465803    0    0    0     0          0       281 51088488  168618    0    0    0     0       0          0",
		"docker0:       0       0    0    0    0     0          0         0      0       0    0    0    0     0       0          0",
	}

	env := environment.NewEnvironment()
	collector := NewNetworkCollector(env, "resourceId")
	data := collector.parseProcNetDev(lines)

	eth0 := data["eth0"]
	if len(eth0) != 16 {
		t.Errorf("eth0 length should be 16, but %v", eth0)
	}

	expected := []string{"61058551", "44173", "0", "0", "0", "0", "0", "0", "422875", "5776", "0", "0", "0", "0", "0", "0"}
	for i, value := range eth0 {
		if value != expected[i] {
			t.Errorf("value should be %v, but %v", expected[i], value)
		}
	}

	eth1 := data["eth1"]
	if len(eth1) != 16 {
		t.Errorf("eth1 length should be 16, but %v", eth1)
	}

	expected = []string{"284076872", "1465803", "0", "0", "0", "0", "0", "281", "51088488", "168618", "0", "0", "0", "0", "0", "0"}
	for i, value := range eth1 {
		if value != expected[i] {
			t.Errorf("value should be %v, but %v", expected[i], value)
		}
	}

	docker := data["docker0"]
	if len(docker) != 0 {
		t.Errorf("docker length should be 0, but %v", docker)
	}
}
