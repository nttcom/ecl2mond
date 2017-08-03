package monitoring

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/nttcom/ecl2mond/logging"
	"github.com/nttcom/ecl2mond/meter"
	"github.com/nttcom/ecl2mond/authenticate"
	"github.com/nttcom/ecl2mond/environment"
)

var logger = logging.NewLogger("monitoring")

const customResourceMaxErrorMessage = "Custom resource is over than the creation limit."
const other40XErrorMessage = "Response code : 40X"
const authErrorMessage = "Authentication Error. (Token expired)"

// Monitoring is data for Monitoring server access.
type Monitoring struct {
	URL        string
	Path       string
	TenantID   string
	ResourceID string
}

// NewMonitoring create Monitoring.
func NewMonitoring(url string, tenantID string, resourceID string) *Monitoring {
	return &Monitoring{
		URL:        url,
		Path:       "/v2/meters",
		TenantID:   tenantID,
		ResourceID: resourceID,
	}
}

// PostMeters accesses to Monitoring Server for creating custom meters.
func (mon *Monitoring) PostMeters(meters []*meter.Meter, token string, auth *authenticate.Authenticate, env *environment.Environment) error {
	logger.Debug("Post Meters start.")
	var attempt int
	for _, meter := range meters {
		logger.Debug(fmt.Sprintf("Post meter start. Name: %s", meter.Name))
		if meter.Volume == "" {
			logger.Debug(fmt.Sprintf("Meter's volume is empty. Skip. Name: %s", meter.Name))
			continue
		}
		attempt = 5
		for attempt > 0 {
			attempt--
			err := mon.postMeter(meter, token, env)

			if err == nil {
				break
			}

			if err.Error() == customResourceMaxErrorMessage {
				// カスタムリソース作成上限に抵触した場合は、
				// 後続のメーター送信も同エラーとなるため継続しない
				return err
			}

			if err.Error() == other40XErrorMessage {
				// その他400系エラーは当該メーターのリトライは行わない
				break
			}

			if err.Error() == authErrorMessage {
 				auth.SaveToken(env)
 				token = auth.GetToken(env)
 				logger.Info("Authenticate token saved. (Not interval)")
 			}

			logger.Debug("Could not post meter to monitoring server. Retry..")
			time.Sleep(1 * time.Second)
		}

		if attempt == 0 {
			attempt = 5
			logger.Error(fmt.Sprintf("Post meter failed. Name: %s", meter.Name))
			continue
		}
	}

	return nil
}

// PostMeter accesses to Monitoring Server for creating a custom meter.
func (mon *Monitoring) postMeter(meter *meter.Meter, token string, env *environment.Environment) error {
	req, err := http.NewRequest(
		"POST",
		mon.getPostMetersURI(meter),
		bytes.NewBuffer([]byte(mon.getPostMetersRequestJSON(meter))),
	)

	if err != nil {
		logger.Error("Failed to connect Monitoring API.")
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("User-Agent", "CustomMeterAgent/" + env.Version + " " + mon.ResourceID)

	client := &http.Client{Timeout: time.Duration(30 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to connect Monitoring API.")
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if resp.StatusCode == 401 {
 			logger.Info("Authentication error response is returned from the Monitoring API.")
 			logger.Info(fmt.Sprintf("Response: %s", body))
 			return errors.New(authErrorMessage)
 		}

		logger.Error(fmt.Sprintf("An error response is returned from the Monitoring API. Name: %s", meter.Name))
		logger.Error(fmt.Sprintf("Response: %s", body))

		if strings.Contains(string(body), customResourceMaxErrorMessage) {
			return errors.New(customResourceMaxErrorMessage)
		}

		if (resp.StatusCode >= 400 && resp.StatusCode <= 417) {
			return errors.New(other40XErrorMessage)
		}

		return errors.New("Post CustomMeter Error")
	}

	return nil
}

// getPostMetersURI returns POST /v2/meters uri.
func (mon *Monitoring) getPostMetersURI(meter *meter.Meter) string {
	var url string

	if mon.URL[len(mon.URL)-1] == '/' {
		url = mon.URL[0 : len(mon.URL)-1]
	} else {
		url = mon.URL
	}

	return url + mon.Path + "/" + meter.Name
}

// PostMetersRequestJSON is POST meters request json data.
type PostMetersRequestJSON struct {
	CounterName   string `json:"counter_name"`
	ResourceID    string `json:"resource_id"`
	CounterUnit   string `json:"counter_unit"`
	CounterVolume string `json:"counter_volume"`
	CounterType   string `json:"counter_type"`
}

// getPostMetersRequestJSON returns json data for request.
func (mon *Monitoring) getPostMetersRequestJSON(meter *meter.Meter) string {
	var requestJSON []*PostMetersRequestJSON
	var meterJSON *PostMetersRequestJSON

	meterJSON = &PostMetersRequestJSON{
		CounterName:   meter.Name,
		ResourceID:    mon.ResourceID,
		CounterUnit:   meter.Unit,
		CounterVolume: meter.Volume,
		CounterType:   meter.Type,
	}
	requestJSON = append(requestJSON, meterJSON)
	b, _ := json.Marshal(requestJSON)

	return string(b)
}
