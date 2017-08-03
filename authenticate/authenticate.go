package authenticate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nttcom/ecl2mond/environment"
	"github.com/nttcom/ecl2mond/logging"
)

var logger = logging.NewLogger("authenticate")

// Authenticate is data for Auth server access.
type Authenticate struct {
	URL      string
	Path     string
	TenantID string
	UserName string
	Password string
}

// RequestJSONData is Request body json data.
type RequestJSONData struct {
	Auth AuthJSONData `json:"auth"`
}

// AuthJSONData is json data in RequestJSONData.
type AuthJSONData struct {
	Identity IdentityJSONData `json:"identity"`
	Scope    ScopeJSONData    `json:"scope"`
}

// IdentityJSONData is json data in RequestJSONData.
type IdentityJSONData struct {
	Methods  []string         `json:"methods"`
	Password PasswordJSONData `json:"password"`
}

// PasswordJSONData is json data in RequestJSONData.
type PasswordJSONData struct {
	User UserJSONData `json:"user"`
}

// UserJSONData is json data in RequestJSONData.
type UserJSONData struct {
	Domain   DomainJSONData `json:"domain"`
	Name     string         `json:"name"`
	Password string         `json:"password"`
}

// DomainJSONData is json data in RequestJSONData.
type DomainJSONData struct {
	ID string `json:"id"`
}

// ScopeJSONData is json data in RequestJSONData.
type ScopeJSONData struct {
	Project ProjectJSONData `json:"project"`
}

// ProjectJSONData is json data in RequestJSONData.
type ProjectJSONData struct {
	ID string `json:"id"`
}

// NewAuthenticate creates Authenticate.
func NewAuthenticate(url string, tenantID string, userName string, password string) *Authenticate {
	return &Authenticate{
		URL:      url,
		Path:     "/v3/auth/tokens",
		TenantID: tenantID,
		UserName: userName,
		Password: password,
	}
}

// SaveToken saves new token.
func (auth *Authenticate) SaveToken(env *environment.Environment) error {
	// Attempt = 5
	var token string
	var err error
	attempt := 5

	for attempt > 0 {
		attempt--
		token, err = auth.getTokenFromAuthenticateServer()

		if err == nil {
			break
		}

		logger.Debug("Could not access authenticate server. retry..")
		time.Sleep(1 * time.Second)
	}

	if attempt == 0 {
		return err
	}

	// Save
	env.AuthToken = token

	return nil
}

// GetToken returns token
func (auth *Authenticate) GetToken(env *environment.Environment) string {
	return env.AuthToken
}

// getTokenFromAuthenticateServer returns new token from Authenticate server.
func (auth *Authenticate) getTokenFromAuthenticateServer() (string, error) {
	// Send POST to Authenticate Server
	req, err := http.NewRequest(
		"POST",
		auth.getURI(),
		bytes.NewBuffer([]byte(auth.getRequestBody())),
	)
	if err != nil {
		logger.Error("Failed to connect Authenticate API.")
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Duration(30 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to connect Authenticate API.")
		return "", err
	}

	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 {
		logger.Error("An error response is returned from the authentication API.")

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		logger.Error(fmt.Sprintf("Response: %s", body))
		return "", errors.New("Authenticate Falied.")
	}

	// Get token from response header
	token := resp.Header.Get("X-Subject-Token")
	return token, nil
}

func (auth *Authenticate) getURI() string {
	var url string

	if auth.URL[len(auth.URL)-1] == '/' {
		url = auth.URL[0 : len(auth.URL)-1]
	} else {
		url = auth.URL
	}

	return url + auth.Path
}

// GetRequestBody returns request body.
func (auth *Authenticate) getRequestBody() string {
	data := &RequestJSONData{
		Auth: AuthJSONData{
			Identity: IdentityJSONData{
				Methods: []string{"password"},
				Password: PasswordJSONData{
					User: UserJSONData{
						Domain: DomainJSONData{
							ID: "default",
						},
						Name:     auth.UserName,
						Password: auth.Password,
					},
				},
			},
			Scope: ScopeJSONData{
				Project: ProjectJSONData{
					ID: auth.TenantID,
				},
			},
		},
	}
	b, _ := json.Marshal(data)

	return string(b)
}
