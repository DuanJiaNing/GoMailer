package conf

import (
	"os"
	"strings"

	"GoMailer/log"
	"gopkg.in/yaml.v2"
)

func init() {
	path := getConfig()
	log.Info("parse config file at: ", path)

	conf = &config{}
	if f, err := os.Open(path); err != nil {
		panic(err)
	} else {
		if err = yaml.NewDecoder(f).Decode(conf); err != nil {
			panic(err)
		}
	}
}

func getConfig() string {
	defaultConfig := "app.yaml"

	if len(os.Args) == 2 {
		configFile := os.Args[1]
		if len(strings.TrimSpace(configFile)) == 0 {
			return defaultConfig
		}

		return configFile
	}

	return defaultConfig
}
