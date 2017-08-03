package config

import (
	"io/ioutil"
	"os"
	"testing"
)

var testConfig = `
logLevel = "ERROR"
authUrl = "http://auth.example.com"
monitoringUrl = "http://monitoring.example.com"
resourceId = "resourceId"
tenantId = "tenantId"
userName = "userName"
password = "password"
interval = 5
authInterval = 60

meters = [
  "meter.1",
	"meter.2",
	"meter.3",
	"meter2.1",
	"meter2.2",
	"meter3.1",
	"meter3.2",
]
`

var testConfigWithoutLoglevel = `
authUrl = "http://auth.example.com"
monitoringUrl = "http://monitoring.example.com"
resourceId = "resourceId"
tenantId = "tenantId"
userName = "userName"
password = "password"
`

func getTestConfig(str string, t *testing.T) *Config {
	tempFile, err := createTempConfig(str)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	defer os.Remove(tempFile.Name())

	config, err := LoadConfig(tempFile.Name())
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	return config
}

func TestLoadConfig(t *testing.T) {
	config := getTestConfig(testConfig, t)

	if config.LogLevel != "ERROR" {
		t.Errorf("should be ERROR but %s", config.LogLevel)
	}

	if config.AuthURL != "http://auth.example.com" {
		t.Errorf("should be http://auth.example.com but %s", config.AuthURL)
	}

	if config.MonitoringURL != "http://monitoring.example.com" {
		t.Errorf("should be http://monitoring.example.com but %s", config.MonitoringURL)
	}

	if config.TenantID != "tenantId" {
		t.Errorf("should be tenantId but %s", config.TenantID)
	}

	if config.UserName != "userName" {
		t.Errorf("should be userName but %s", config.UserName)
	}

	if config.Password != "password" {
		t.Errorf("should be password but %s", config.Password)
	}

	if config.ResourceID != "resourceId" {
		t.Errorf("should be resourceId but %s", config.ResourceID)
	}

	if config.Interval != 5 {
		t.Errorf("should be 5 but %v", config.Interval)
	}

	if config.AuthInterval != 60 {
		t.Errorf("should be 60 but %v", config.AuthInterval)
	}

	if config.Meters[0] != "meter.1" {
		t.Errorf("should be meter.1 but %v", config.Meters[0])
	}

	if config.Meters[1] != "meter.2" {
		t.Errorf("should be meter.2 but %v", config.Meters[1])
	}

	if config.Meters[2] != "meter.3" {
		t.Errorf("should be meter.3 but %v", config.Meters[2])
	}
}

func TestLoadConfigWithoutLogLevel(t *testing.T) {
	config := getTestConfig(testConfigWithoutLoglevel, t)

	if config.LogLevel != "INFO" {
		t.Errorf("should be INFO but %s", config.LogLevel)
	}

	if config.Interval != 5 {
		t.Errorf("should be 5 but %v", config.Interval)
	}

	if config.AuthInterval != 60 {
		t.Errorf("should be 60 but %v", config.AuthInterval)
	}

	if len(config.Meters) != 0 {
		t.Errorf("should be 0 but %v", config.Meters)
	}
}

func createTempConfig(content string) (*os.File, error) {
	tempFile, err := ioutil.TempFile("", "ecl-monitoirng-conf")
	if err != nil {
		return nil, err
	}

	if _, err := tempFile.WriteString(content); err != nil {
		os.Remove(tempFile.Name())
		return nil, err
	}
	tempFile.Sync()
	tempFile.Close()
	return tempFile, nil
}

func TestMeterTypes(t *testing.T) {
	config := getTestConfig(testConfig, t)

	meterTypes := config.MeterTypes()

	if len(meterTypes) != 3 {
		t.Errorf("meterTypes length should be 3, but %v", meterTypes)
	}

	if meterTypes[0] != "meter" {
		t.Errorf("meterTypes[0] should be meter, but %v", meterTypes[0])
	}

	if meterTypes[1] != "meter2" {
		t.Errorf("meterTypes[1] should be meter2, but %v", meterTypes[1])
	}

	if meterTypes[2] != "meter3" {
		t.Errorf("meterTypes[2] should be meter3, but %v", meterTypes[2])
	}
}

func TestValidate(t *testing.T) {
	var testCases = []struct {
		configFile string
		expected   bool
	}{
		{testConfig, true},
		{"", false},
	}

	for _, testCase := range testCases {
		config := getTestConfig(testCase.configFile, t)
		if config.Validate() != testCase.expected {
			t.Errorf("Validate result should be %v, config: %v", testCase.expected, testCase.configFile)
		}
	}
}

func TestValidateMonitoringURL(t *testing.T) {
	var testCases = []struct {
		configFile string
		expected   bool
	}{
		{"", false},
		{`monitoringUrl = "http://monitoring.example.com/"`, true},
		{testConfig, true},
	}

	for _, testCase := range testCases {
		config := getTestConfig(testCase.configFile, t)
		if config.ValidateMonitoringURL() != testCase.expected {
			t.Errorf("ValidateMonitoringURL result should be %v, case: %v, value: %v", testCase.expected, testCase.configFile, config.MonitoringURL)
		}
	}
}

func TestValidateInterval(t *testing.T) {
	var testCases = []struct {
		configFile string
		expected   bool
	}{
		{`monitoringUrl = "http://monitoring.example.com/"`, true}, // default value 5.
		{`interval = 0`, false},
		{`interval = 1`, true},
		{`interval = 5`, true},
		{`interval = 3599`, true},
		{`interval = 3600`, false},
		{testConfig, true},
	}

	for _, testCase := range testCases {
		config := getTestConfig(testCase.configFile, t)
		if config.ValidateInterval() != testCase.expected {
			t.Errorf("ValidateInterval result should be %v, case: %v, value: %v", testCase.expected, testCase.configFile, config.Interval)
		}
	}
}

func TestValidateAuthURL(t *testing.T) {
	var testCases = []struct {
		configFile string
		expected   bool
	}{
		{`monitoringUrl = "http://monitoring.example.com/"`, false},
		{"", false},
		{`authUrl = "http://auth.example.com/"`, true},
		{testConfig, true},
	}

	for _, testCase := range testCases {
		config := getTestConfig(testCase.configFile, t)
		if config.ValidateAuthURL() != testCase.expected {
			t.Errorf("ValidateAuthURL result should be %v, case: %v, value: %v", testCase.expected, testCase.configFile, config.AuthURL)
		}
	}
}

func TestValidateAuthInterval(t *testing.T) {
	var testCases = []struct {
		configFile string
		expected   bool
	}{
		{`monitoringUrl = "http://monitoring.example.com/"`, true}, // default value 60.
		{`authInterval = 0`, false},
		{`authInterval = 1`, false},
		{`authInterval = 5`, true},
		{`authInterval = 3599`, true},
		{`authInterval = 3600`, false},
		{testConfig, true},
	}

	for _, testCase := range testCases {
		config := getTestConfig(testCase.configFile, t)
		if config.ValidateAuthInterval() != testCase.expected {
			t.Errorf("ValidateAuthInterval result should be %v, case: %v, value: %v", testCase.expected, testCase.configFile, config.ResourceID)
		}
	}
}

func TestValidateValidateResourceID(t *testing.T) {
	var testCases = []struct {
		configFile string
		expected   bool
	}{
		{`monitoringUrl = "http://monitoring.example.com/"`, false},
		{"", false},
		{`resourceId = "xxxxxxx"`, true},
		{testConfig, true},
	}

	for _, testCase := range testCases {
		config := getTestConfig(testCase.configFile, t)
		if config.ValidateResourceID() != testCase.expected {
			t.Errorf("ValidateResourceID result should be %v, case: %v, value: %v", testCase.expected, testCase.configFile, config.ResourceID)
		}
	}
}

func TestValidateTenantID(t *testing.T) {
	var testCases = []struct {
		configFile string
		expected   bool
	}{
		{`monitoringUrl = "http://monitoring.example.com/"`, false},
		{"", false},
		{`tenantId = "xxxxxxx"`, true},
		{testConfig, true},
	}

	for _, testCase := range testCases {
		config := getTestConfig(testCase.configFile, t)
		if config.ValidateTenantID() != testCase.expected {
			t.Errorf("ValidateTenantID result should be %v, case: %v, value: %v", testCase.expected, testCase.configFile, config.TenantID)
		}
	}
}

func TestValidateUserName(t *testing.T) {
	var testCases = []struct {
		configFile string
		expected   bool
	}{
		{`monitoringUrl = "http://monitoring.example.com/"`, false},
		{"", false},
		{`userName = "xxxxxxx"`, true},
		{testConfig, true},
	}

	for _, testCase := range testCases {
		config := getTestConfig(testCase.configFile, t)
		if config.ValidateUserName() != testCase.expected {
			t.Errorf("ValidateUserName result should be %v, case: %v, value: %v", testCase.expected, testCase.configFile, config.UserName)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	var testCases = []struct {
		configFile string
		expected   bool
	}{
		{`monitoringUrl = "http://monitoring.example.com/"`, false},
		{"", false},
		{`password = "xxxxxxx"`, true},
		{testConfig, true},
	}

	for _, testCase := range testCases {
		config := getTestConfig(testCase.configFile, t)
		if config.ValidatePassword() != testCase.expected {
			t.Errorf("ValidatePassword result should be %v, case: %v, value: %v", testCase.expected, testCase.configFile, config.Password)
		}
	}
}

func TestValidateLogLevel(t *testing.T) {
	var testCases = []struct {
		configFile string
		expected   bool
	}{
		{`monitoringUrl = "http://monitoring.example.com/"`, true}, // default value is INFO.
		{"", true},              // default value is INFO
		{`logLevel = ""`, false},
		{`logLevel = "WRONG"`, false},
		{`logLevel = "TRACE"`, true},
		{testConfig, true},
	}

	for _, testCase := range testCases {
		config := getTestConfig(testCase.configFile, t)
		if config.ValidateLogLevel() != testCase.expected {
			t.Errorf("ValidateLogLevel result should be %v, case: %v, value: %v", testCase.expected, testCase.configFile, config.LogLevel)
		}
	}
}
