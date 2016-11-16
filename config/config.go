package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"errors"
)

type Config struct {
	InterfaceName string
	Port uint16
	LogPath string
	LogPrefix string
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
