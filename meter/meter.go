package meter

import "strings"

// Meter is monitoring metrics data.
type Meter struct {
	Name       string
	ResourceID string
	Unit       string
	Type       string
	Volume     string
}

// NewMeter creates Meter.
func NewMeter(name string, resourceID string, unit string, meterType string) *Meter {
	return &Meter{
		Name:       name,
		ResourceID: resourceID,
		Unit:       unit,
		Type:       meterType,
	}
}

// SetVolume set meter's volume.
func (meter *Meter) SetVolume(volume string) {
	meter.Volume = volume
}

// GetMeterType gets MeterType.
func (meter *Meter) GetMeterType() string {
	return strings.Split(meter.Name, ".")[0]
}
