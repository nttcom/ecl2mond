package config

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config is content of config file.
type Config struct {
	Mode          string
	LogLevel      string
	AuthURL       string
	MonitoringURL string
	ResourceID    string
	TenantID      string
	UserName      string
	Password      string
	Interval      uint16
	AuthInterval  uint16
	Meters        []string
}

// DefaultConfig is defult of Config.
var DefaultConfig = &Config{
	LogLevel:     "INFO",
	Mode:         "run",
	Interval:     5,
	AuthInterval: 60,
}

func SetDefaultConfig(arg_conf_struct *Config) (ret_conf_struct *Config) {
	ret_conf_struct = arg_conf_struct
	ret_conf_struct.LogLevel = DefaultConfig.LogLevel
	ret_conf_struct.Mode = DefaultConfig.Mode
	ret_conf_struct.Interval = DefaultConfig.Interval
	ret_conf_struct.AuthInterval = DefaultConfig.AuthInterval
	return ret_conf_struct
}


// LoadConfig loads Config form config file path.
func LoadConfig(filename string) (*Config, error) {
	config := &Config{}
	config = SetDefaultConfig(config)
	if _, err := toml.DecodeFile(filename, config); err != nil {
		fmt.Println("Failed to load config file.")
		return config, err
	}

	config.Mode = DefaultConfig.Mode

	return config, nil
}

// MeterTypes returns types of meter
func (config *Config) MeterTypes() []string {
	var meterTypes []string
	encounted := map[string]bool{}
	for i := 0; i < len(config.Meters); i++ {
		meterType := strings.SplitN(config.Meters[i], ".", 2)[0]
		if !encounted[meterType] {
			encounted[meterType] = true
			meterTypes = append(meterTypes, meterType)
		}
	}
	return meterTypes
}

// Validate validates all items on config file.
func (config *Config) Validate() bool {
	result := true
	result = config.ValidateMonitoringURL() && result
	result = config.ValidateInterval() && result
	result = config.ValidateAuthURL() && result
	result = config.ValidateAuthInterval() && result
	result = config.ValidateResourceID() && result
	result = config.ValidateTenantID() && result
	result = config.ValidateUserName() && result
	result = config.ValidatePassword() && result
	result = config.ValidateLogLevel() && result

	return result
}

// ValidateMonitoringURL validates monitoringUrl.
func (config *Config) ValidateMonitoringURL() bool {
	if config.MonitoringURL == "" {
		fmt.Println("'monitoring Url' is blank.")
		return false
	}

	return true
}

// ValidateInterval validates interval.
func (config *Config) ValidateInterval() bool {
	if config.Interval < 1 || 3599 < config.Interval {
		fmt.Println("'interval' is not within 1 to 3599.")
		return false
	}

	return true
}

// ValidateAuthURL validates authUrl.
func (config *Config) ValidateAuthURL() bool {
	if config.AuthURL == "" {
		fmt.Println("'authUrl' is blank.")
		return false
	}

	return true
}

// ValidateAuthInterval validates authInterval.
func (config *Config) ValidateAuthInterval() bool {
	if config.AuthInterval < 5 || 3599 < config.AuthInterval {
		fmt.Println("'authInterval' is not within 5 to 3599.")
		return false
	}

	return true
}

// ValidateResourceID validates resourceId.
func (config *Config) ValidateResourceID() bool {
	if config.ResourceID == "" {
		fmt.Println("'resourceId' is blank.")
		return false
	}

	return true
}

// ValidateTenantID validates tenantId.
func (config *Config) ValidateTenantID() bool {
	if config.TenantID == "" {
		fmt.Println("'tenantId' is blank.")
		return false
	}

	return true
}

// ValidateUserName validates userName.
func (config *Config) ValidateUserName() bool {
	if config.UserName == "" {
		fmt.Println("'UserName' is blank.")
		return false
	}

	return true
}

// ValidatePassword validates password.
func (config *Config) ValidatePassword() bool {
	if config.Password == "" {
		fmt.Println("'password' is blank.")
		return false
	}

	return true
}


// ValidateLogLevel validates logLevel.
func (config *Config) ValidateLogLevel() bool {
	levels := []string{"TRACE", "DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
	for _, level := range levels {
		if level == config.LogLevel {
			return true
		}
	}

	if config.LogLevel == "" {
		fmt.Println("'logLevel' is blank.")
		return false
	}

	fmt.Println("'logLevel' is invalid.")
	return false
}
