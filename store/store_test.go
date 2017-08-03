package store

import (
	"testing"
	"time"

	"github.com/nttcom/ecl2mond/environment"
)

func TestSetData(t *testing.T) {
	env := environment.NewEnvironment()
	env.RootPath = "./"
	store := NewStore(env)
	name := "cpu"
	expectedValues := map[string][]string{"cpu": []string{"1", "2", "3"}}

	err := store.SetData(name, expectedValues)
	if err != nil {
		t.Error("should not raise error.")
	}

	result, createdAt, err := store.GetData(name)
	values := result["cpu"]
	if err != nil {
		t.Error("should not raise error.")
	}

	if len(values) != 3 {
		t.Errorf("values length should be 3, but %v", len(values))
	}

	if time.Now().Sub(createdAt) > 1*time.Minute {
		t.Errorf("createdAt should be near by time.Now(), but %v", createdAt)
	}

	for i, value := range values {
		if value != expectedValues["cpu"][i] {
			t.Errorf("value should be %v, but %v", expectedValues["cpu"][i], value)
		}
	}
}
