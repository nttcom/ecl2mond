package authenticate

import (
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nttcom/ecl2mond/environment"
)

func TestNewAuthenticate(t *testing.T) {
	auth := NewAuthenticate("http://auth.example.com", "TenantID", "userName", "password")

	if auth.URL != "http://auth.example.com" {
		t.Errorf("auth.URL should be http://auth.example.com but %v", auth.URL)
	}

	if auth.TenantID != "TenantID" {
		t.Errorf("auth.TenantID should be TenantID but %v", auth.TenantID)
	}

	if auth.UserName != "userName" {
		t.Errorf("auth.UserName should be userName but %v", auth.UserName)
	}

	if auth.Password != "password" {
		t.Errorf("auth.Password should be password but %v", auth.Password)
	}
}

func TestGetRequestData(t *testing.T) {
	auth := NewAuthenticate("http://auth.example.com", "TenantID", "userName", "password")

	b := auth.getRequestBody()

	if b != `{"auth":{"identity":{"methods":["password"],"password":{"user":{"domain":{"id":"default"},"name":"userName","password":"password"}}},"scope":{"project":{"id":"TenantID"}}}}` {
		t.Errorf("b should be hoge but %v", b)
	}
}

func TestGetUrl(t *testing.T) {
	auth := NewAuthenticate("http://auth.example.com", "TenantID", "userName", "password")
	if auth.getURI() != "http://auth.example.com/v3/auth/tokens" {
		t.Errorf("auth.getURI() should return http://auth.example.com/v3/auth/tokens but %v", auth.getURI())
	}

	auth = NewAuthenticate("http://auth.example.com/", "TenantID", "userName", "password")
	if auth.getURI() != "http://auth.example.com/v3/auth/tokens" {
		t.Errorf("auth.getURI() should return http://auth.example.com/v3/auth/tokens but %v", auth.getURI())
	}
}

// ローカル環境からKeyStoneへの疎通テスト
func NotTestGetTokenFromAuthenticateServer(t *testing.T) {
	auth := NewAuthenticate("http://10.128.49.10:5000", "f7da5b239b294ed784277cc87ff1366f", "AjzkKOvprS4fFqjXSbc3tn8pcu8irIkV", "BmmekGYhPqSWO6EY")
	token, err := auth.getTokenFromAuthenticateServer()

	if err != nil {
		t.Error("should not raise error.")
	}

	if token == "" {
		t.Errorf("should not be presnet but %v", token)
	}
}

func NotTestToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	env := environment.NewEnvironment()
	os.Remove(env.RootPath + "/token")
	defer os.Remove(env.RootPath + "/token")

	httpmock.RegisterResponder("POST", "http://auth.example.com/v3/auth/tokens",
		func(req *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(201, `{}`)
			resp.Header.Set("X-Subject-Token", "newToken")
			return resp, nil
		},
	)

	auth := NewAuthenticate("http://auth.example.com/", "tenantId", "userName", "password")
	err := auth.SaveToken(env)

	if err != nil {
		t.Error("Authenticate Server access error.")
	}

	newToken := auth.GetToken(env)

	if newToken != "newToken" {
		t.Errorf("Token does not save")
	}
}

func NotTestSaveToken(t *testing.T) {
	env := environment.NewEnvironment()
	auth := NewAuthenticate("http://10.128.49.10:5000", "f7da5b239b294ed784277cc87ff1366f", "AjzkKOvprS4fFqjXSbc3tn8pcu8irIkV", "BmmekGYhPqSWO6EY")
	err := auth.SaveToken(env)

	if err != nil {
		t.Errorf("should not raise error %v", err)
	}
}

func NotTestGetToken(t *testing.T) {
	env := environment.NewEnvironment()
	env.AuthToken = "test"

	auth := NewAuthenticate("http://10.128.49.10:5000", "f7da5b239b294ed784277cc87ff1366f", "AjzkKOvprS4fFqjXSbc3tn8pcu8irIkV", "BmmekGYhPqSWO6EY")

	token := auth.GetToken(env)

	if token != "test" {
		t.Errorf("Token should be test but %v", token)
	}
}
