package conf

import (
	"fmt"
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
	defaultConfig := "app.dev.yaml"

	if len(os.Args) == 2 {
		env := os.Args[1]
		if len(strings.TrimSpace(env)) == 0 {
			return defaultConfig
		}

		return fmt.Sprintf("app.%s.yaml", env)
	}

	return defaultConfig
}
