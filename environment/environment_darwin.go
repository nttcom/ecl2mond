package environment

import (
	"os"
	"path/filepath"
)

var rootPath = filepath.Join(os.Getenv("HOME"), "Library", getAgentName())
var configFilePath = filepath.Join(rootPath, getAgentName()+".conf")
