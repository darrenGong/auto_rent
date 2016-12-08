package fetchHouse

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"errors"
)

type WebUrl struct {
	Url     string
	AreaUrl string
}

type Config struct {
	PlatUrl       map[string]WebUrl
	ServiceDir    string
	FetchDuration uint32
}

func ParseConfig(path string, config *Config) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read file: %s, err: %v\n", path, err)
		return errors.New("Failed to read file")
	}
	if err := json.Unmarshal(bytes, config); err != nil {
		log.Printf("Failed to Unmarshal file: %s, err:%v\n", path, err)
		return errors.New("Failed to Unmarshal file")
	}

	return nil
}