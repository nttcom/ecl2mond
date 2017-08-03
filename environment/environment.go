package environment

// Environment is environment information struct.
type Environment struct {
	RootPath        string
	ConfigFilePath string
	AuthToken       string
	Version         string
}

// NewEnvironment creates new Environment.
func NewEnvironment() *Environment {
	return &Environment{
		RootPath:       rootPath,
		ConfigFilePath: configFilePath,
		Version: "1.0.0",
	}
}

func getAgentName() string {
	return "ecl2mond"
}
