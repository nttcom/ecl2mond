package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/nttcom/ecl2mond/environment"
	"github.com/nttcom/ecl2mond/logging"
)

var logger = logging.NewLogger("store")

// Store is struct for string values.
type Store struct {
	Env *environment.Environment
}

// JSONData is JSON data for store data.
type JSONData struct {
	Data      map[string][]string
	CreatedAt time.Time
}

// NewStore returns Store struct.
func NewStore(env *environment.Environment) *Store {
	return &Store{
		Env: env,
	}
}

// SetData saves data.
func (store *Store) SetData(name string, data map[string][]string) error {
	serialized, err := json.Marshal(&JSONData{
		Data:      data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		logger.Error("Failed to save sample data.")
		return err
	}

	err = ioutil.WriteFile(store.getStoreFilePath(name), serialized, 0644)
	if err != nil {
		logger.Error("Failed to save sample data.")
		return err
	}

	return nil
}

// GetData returns data of name.
func (store *Store) GetData(name string) (map[string][]string, time.Time, error) {
	filePath := store.getStoreFilePath(name)

	// To confirm the presence of file.
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, time.Now(), err
	}

	serialized, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Error("Failed to load sample data.")
		return nil, time.Now(), err
	}

	var data *JSONData
	err = json.Unmarshal(serialized, &data)
	if err != nil {
		logger.Error("Failed to load sample data.")
		return nil, time.Now(), err
	}

	return data.Data, data.CreatedAt, nil
}

func (store *Store) getStoreFilePath(name string) string {
	return store.Env.RootPath + "/" + name
}
