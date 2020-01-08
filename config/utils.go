package config

import (
	"encoding/json"
	"io/ioutil"
)

// LoadJSONFile load a json file to obj
func LoadJSONFile(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, obj)
}
