package viperdump

import (
	"fmt"
	log "github.com/docktermj/go-logger/logger"
	viper "github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// Dump viper configuration contents to a string of yaml.
// See https://github.com/spf13/viper#marshalling-to-string
func yamlStringSettings() string {
	viperConfig := viper.AllSettings()
	viperByteArray, err := yaml.Marshal(viperConfig)
	if err != nil {
		log.Fatalf("Unable to marshal viper config to YAML: %v", err)
	}
	return string(viperByteArray)
}

// The Command sofware design pattern's Execute() method.
func Execute() {
	fmt.Println(yamlStringSettings())
}
