package meter

import "testing"

func TestNewMeter(t *testing.T) {
	meter := NewMeter("meter.name", "resourceId", "percent", "delta")
	if meter.Name != "meter.name" {
		t.Errorf("should be meter.name but %v", meter.Name)
	}

	if meter.ResourceID != "resourceId" {
		t.Errorf("should be resourceId but %v", meter.ResourceID)
	}

	if meter.Unit != "percent" {
		t.Errorf("should be percent but %v", meter.Unit)
	}

	if meter.Type != "delta" {
		t.Errorf("should be delta but %v", meter.Type)
	}
}

func TestGetMeterType(t *testing.T) {
	meter := NewMeter("cpuusage.user", "resourceId", "percent", "delta")

	if meter.GetMeterType() != "cpuusage" {
		t.Errorf("should be cpuusage but %v", meter.GetMeterType())
	}
}

func TestSetVolume(t *testing.T) {
	meter := NewMeter("meter.name", "resourceId", "percent", "delta")
	meter.SetVolume("1.02")

	if meter.Volume != "1.02" {
		t.Errorf("should be 1.01 but %v", meter.Volume)
	}
}
