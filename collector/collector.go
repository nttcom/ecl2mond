package collector

import "github.com/nttcom/ecl2mond/meter"

// Collector is data for collection meters.
type Collector interface {
	Collect() ([]*meter.Meter, error)
}
