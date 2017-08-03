package ecl2mond

import (
	"github.com/nttcom/ecl2mond/config"
	"github.com/nttcom/ecl2mond/environment"
)

func getEcl2mond() *Ecl2mond {
	ma := NewEcl2mond()
	ma.Config = getConfig()
	ma.Env = environment.NewEnvironment()
	return ma
}

func getConfig() *config.Config {
	defaultConfig := &config.Config{
		Mode:          "run",
		LogLevel:      "INFO",
		AuthURL:       "http://auth.example.com/",
		MonitoringURL: "http://monitoring.example.com/",
		ResourceID:    "aaaaaaaaaaaaaaa",
		TenantID:      "bbbbbbbbbbbbbbbbb",
		UserName:      "ccccccccccccccccc",
		Password:      "ddddddddddddddddd",
	}
	return defaultConfig
}
