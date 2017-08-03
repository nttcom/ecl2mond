package ecl2mond

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nttcom/ecl2mond/aggregator"
	"github.com/nttcom/ecl2mond/authenticate"
	"github.com/nttcom/ecl2mond/config"
	"github.com/nttcom/ecl2mond/environment"
	"github.com/nttcom/ecl2mond/logging"
	"github.com/nttcom/ecl2mond/meter"
	"github.com/nttcom/ecl2mond/monitoring"
)

var logger = logging.NewLogger("ecl2mond")

// Ecl2mond is sample struct.
type Ecl2mond struct {
	Config *config.Config
	Env    *environment.Environment
}

// NewEcl2mond creates new Ecl2mond struct.
func NewEcl2mond() *Ecl2mond {
	return &Ecl2mond{}
}

// Run runs main process.
func (mond *Ecl2mond) Run() error {
	mond.Env = environment.NewEnvironment()

	if err := mond.setFlagSet(); err != nil {
		return err
	}

	auth := mond.getAuthenticate()

	switch mond.Config.Mode {
	case "dry-run":
		err := mond.dryRun()
		if err != nil {
			return err
		}
	case "run-once":
		err := mond.runOnce(auth)
		if err != nil {
			return err
		}
	case "run":
		err := mond.run(auth)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mond *Ecl2mond) getAuthenticate() *authenticate.Authenticate {
	auth := authenticate.NewAuthenticate(
		mond.Config.AuthURL,
		mond.Config.TenantID,
		mond.Config.UserName,
		mond.Config.Password,
	)
	return auth
}

func (mond *Ecl2mond) dryRun() error {
	meters, err := mond.getMeters()
	if err != nil {
		return err
	}

	mond.displayMeters(meters)

	return nil
}

func (mond *Ecl2mond) runOnce(auth *authenticate.Authenticate) error {
	err := mond.saveToken(auth)
	if err != nil {
		return err
	}

	mond.monitor(auth, true)

	fmt.Fprintf(os.Stdin, "Finished.\n")

	return nil
}

func (mond *Ecl2mond) run(auth *authenticate.Authenticate) error {
	logger.Info("Start.")
	c := make(chan struct{})
	go mond.saveTokenLoop(auth, c)

	time.Sleep(10 * time.Second)

	go mond.monitorLoop(auth, c)

	<-c
	return nil
}

func (mond *Ecl2mond) monitor(auth *authenticate.Authenticate, display bool) error {
	token := mond.getToken(auth)
	meters, err := mond.getMeters()
	if err != nil {
		return err
	}

	err = mond.postMeters(meters, token, auth)
	if err != nil {
		return err
	}

	if display {
		mond.displayMeters(meters)
	}

	return nil
}

func (mond *Ecl2mond) monitorLoop(auth *authenticate.Authenticate, c chan struct{}) error {
	var startTime time.Time
	interval := time.Duration(mond.Config.Interval) * time.Minute
	nextInterval := time.Duration(0)
	for {
		select {
		case <-time.After(nextInterval):
			startTime = time.Now()
			mond.monitor(auth, false)
			logger.Info("Finished post meter process.")

			nextInterval = interval - time.Now().Sub(startTime)
		case <-c:
			return nil
		}
	}
}

func (mond *Ecl2mond) saveToken(auth *authenticate.Authenticate) error {
	err := auth.SaveToken(mond.Env)
	if err != nil {
		return err
	}
	return nil
}

func (mond *Ecl2mond) saveTokenLoop(auth *authenticate.Authenticate, c chan struct{}) error {
	var startTime time.Time
	interval := time.Duration(mond.Config.AuthInterval) * time.Minute
	nextInterval := time.Duration(0)
	for {
		select {
		case <-time.After(nextInterval):
			startTime = time.Now()
			err := mond.saveToken(auth)

			if err == nil {
				logger.Info("Authenticate token saved.")
			}

			nextInterval = interval - time.Now().Sub(startTime)
		case <-c:
			return nil
		}
	}
}

func (mond *Ecl2mond) getToken(auth *authenticate.Authenticate) string {
	return auth.GetToken(mond.Env)
}

func (mond *Ecl2mond) getMeters() ([]*meter.Meter, error) {
	aggregator := aggregator.NewAggregator(mond.Env, mond.Config, mond.Config.ResourceID)
	meters, err := aggregator.GetCollectedMeters()

	if err != nil {
		return meters, err
	}

	return meters, nil
}

func (mond *Ecl2mond) displayMeters(meters []*meter.Meter) {
	// Header
	header := `+------+---------------------------------------------+-----------------+----------------+----------+
| No   | Name                                        | Value           | Unit           | Type     |
+------+---------------------------------------------+-----------------+----------------+----------+
`
	format := "|%5d |%44s |%16s |%15s |%9s |\n"
	fmt.Fprintf(os.Stdin, header)

	for i, meter := range meters {
		fmt.Fprintf(os.Stdin, format, i+1, meter.Name, meter.Volume, meter.Unit, meter.Type)
	}
	line := "+------+--------------------------------------------+-----------------+----------------+----------+\n"
	fmt.Fprintf(os.Stdin, line)

}

func (mond *Ecl2mond) postMeters(meters []*meter.Meter, token string, auth *authenticate.Authenticate) error {
	mon := monitoring.NewMonitoring(
		mond.Config.MonitoringURL,
		mond.Config.TenantID,
		mond.Config.ResourceID,
	)

	// MonitoringへのPOST
	err := mon.PostMeters(meters, token, auth, mond.Env)
	if err != nil {
		return err
	}
	return nil
}

func (mond *Ecl2mond) setFlagSet() error {
	flag.Usage = defaultHelpMessage

	var configFilePath string
	flag.StringVar(&configFilePath, "c", mond.Env.ConfigFilePath, "")
	flag.StringVar(&configFilePath, "config", mond.Env.ConfigFilePath, "")

	var mode string
	flag.StringVar(&mode, "m", "run", "")
	flag.StringVar(&mode, "mode", "run", "")

	flag.Parse()

	conf, err := config.LoadConfig(configFilePath)
	if err != nil {
		// config file not found.
		return err
	}

	if !conf.Validate() {
		return errors.New("An invalid value was set in the Config file.")
	}

	mond.Config = conf

	err = logging.SetLogLevel(mond.Config.LogLevel)
	if err != nil {
		return err
	}

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "c", "config":
			mond.Env.ConfigFilePath = configFilePath
		case "m", "mode":
			if mode == "run" || mode == "run-once" || mode == "dry-run" {
				mond.Config.Mode = mode
			} else {
				defaultHelpMessage()
				os.Exit(1)
			}
		}
	})

	return nil
}

func defaultHelpMessage() {
	message := `Usage: ecl2mond [options]

options:
  -m, --mode {run(default)|run-once|dry-run}
      set mode
  -c, --config
      set config file
  -h, --help
      show help message

`

	fmt.Fprintf(os.Stderr, message)
}
