package aggregator

import (
	"github.com/nttcom/ecl2mond/collector"
	linux_collector "github.com/nttcom/ecl2mond/collector/linux"
)

func (ag *Aggregator) collectorMapping() map[string]collector.Collector {
	return map[string]collector.Collector{
		"cpu":     &linux_collector.CPUCollector{Env: ag.Env, ResourceID: ag.ResourceID},
		"disk":    &linux_collector.DiskCollector{Env: ag.Env, ResourceID: ag.ResourceID},
		"loadavg": &linux_collector.LoadavgCollector{Env: ag.Env, ResourceID: ag.ResourceID},
		"memory":  &linux_collector.MemoryCollector{Env: ag.Env, ResourceID: ag.ResourceID},
		"network": &linux_collector.NetworkCollector{Env: ag.Env, ResourceID: ag.ResourceID},
	}
}
