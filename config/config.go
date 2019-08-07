package config

import (
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	DB_TYPE     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_PORT     string
	DB_HOST     string
	DB_NAME     string
	SECRET      string
}

func GetConfig() Configuration {
	configuration := Configuration{}
	gonfig.GetConf("config/config.json", &configuration)
	return configuration
}
