package environment

import "testing"

func TestEnvironment(t *testing.T) {
	environment := NewEnvironment()

	if environment.RootPath == "" {
		t.Error("RootPath should be present.")
	}

	if environment.ConfigFilePath == "" {
		t.Error("ConfigFilePath should be present.")
	}

	if environment.Version == "" {
		t.Error("Version should be present.")
	}
}
