package monitoring

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nttcom/ecl2mond/meter"
	"github.com/nttcom/ecl2mond/environment"
)

func TestNewMonitoring(t *testing.T) {
	mon := NewMonitoring("http://mon.example.com", "TenantID", "ResourceID")

	if mon.URL != "http://mon.example.com" {
		t.Errorf("should be http://mon.example.com but %v", mon.URL)
	}

	if mon.TenantID != "TenantID" {
		t.Errorf("should be TenantID but %v", mon.TenantID)
	}

	if mon.ResourceID != "ResourceID" {
		t.Errorf("should be ResourceID but %v", mon.ResourceID)
	}
}

func TestGetPostMetersRequestJSON(t *testing.T) {
	meter1 := meter.NewMeter("meter.name.1", "resourceId", "percent1", "delta1")
	meter1.SetVolume("1.0")
	mon := NewMonitoring("http://mon.example.com", "TenantID", "resourceId")
	requestJSON := mon.getPostMetersRequestJSON(meter1)

	if requestJSON != `[{"counter_name":"meter.name.1","resource_id":"resourceId","counter_unit":"percent1","counter_volume":"1.0","counter_type":"delta1"}]` {
		t.Errorf("should be JSON but %v", requestJSON)
	}
}

func TestGetPostMetersURI(t *testing.T) {
	mon := NewMonitoring("http://mon.example.com", "TenantID", "ResourceID")
	meter1 := meter.NewMeter("meter.name.1", "resourceId", "percent1", "delta1")
	meter1.SetVolume("1.0")
	uri := mon.getPostMetersURI(meter1)

	if uri != "http://mon.example.com/v2/meters/meter.name.1" {
		t.Errorf("should be http://mon.example.com/v2/meters/meter.name.1 but %v", uri)
	}

	mon = NewMonitoring("http://mon.example.com/", "TenantID", "ResourceID")
	uri = mon.getPostMetersURI(meter1)

	if uri != "http://mon.example.com/v2/meters/meter.name.1" {
		t.Errorf("should be http://mon.example.com/v2/meters/meter.name.1 but %v", uri)
	}
}

func TestPostMeter(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://mon.example.com/v2/meters/meter.name.1",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, `{}`)
			return resp, nil
		},
	)

	mon := NewMonitoring("http://mon.example.com/", "tenantId", "resourceId")
	meter1 := meter.NewMeter("meter.name.1", "resourceId", "percent", "delta")
	meter1.SetVolume("1.0")
	token := "token"
	env := environment.NewEnvironment()

	err := mon.postMeter(meter1, token, env)
	if err != nil {
		t.Errorf("Monitoring Post Meters error. %v", err)
	}
}

func NotTestPostMeter(t *testing.T) {
	token := "32c1ab15c54c42e5a30d775bf535f982"
	mon := NewMonitoring("http://10.128.51.65", "f7da5b239b294ed784277cc87ff1366f", "nova_baf6d39fef1e4678a4167175f9ebbec2")
	meter1 := meter.NewMeter("meter.name.1", "resourceId", "percent", "delta")
	meter1.SetVolume("7.0")
	env := environment.NewEnvironment()

	err := mon.postMeter(meter1, token, env)
	if err != nil {
		t.Error("Error.")
	}
}

func NotTestPostMeters(t *testing.T) {
	token := "32c1ab15c54c42e5a30d775bf535f982"
	mon := NewMonitoring("http://10.128.51.65", "f7da5b239b294ed784277cc87ff1366f", "nova_baf6d39fef1e4678a4167175f9ebbec2")
	meter1 := meter.NewMeter("meter.name.1", "resourceId", "percent", "delta")
	meter1.SetVolume("1.0")
	meter2 := meter.NewMeter("meter.name.2", "resourceId", "percent", "delta")
	meter2.SetVolume("2.0")
	meter3 := meter.NewMeter("meter.name.3", "resourceId", "percent", "delta")
	meter3.SetVolume("3.0")
	meter4 := meter.NewMeter("meter.name.4", "resourceId", "percent", "delta")
	meter4.SetVolume("4.0")
	meter5 := meter.NewMeter("meter.name.5", "resourceId", "percent", "delta")
	meter5.SetVolume("5.0")
	meters := []*meter.Meter{meter1, meter2, meter3, meter4, meter5}

	err := mon.PostMeters(meters, token, nil, nil)
	if err != nil {
		t.Error("Error.")
	}
}
