package helpers

import (
	"encoding/json"
	"io/ioutil"
)

// WriteJson write a json file or die trying
func WriteJson(filePath string, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	WriteData(filePath, data)
}

func WriteData(filePath string, data []byte) {
	err := ioutil.WriteFile(filePath, data, 0600)
	if err != nil {
		panic(err)
	}
}

// LoadJson load a json file or die trying
func LoadJson(filePath string, v interface{}) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}