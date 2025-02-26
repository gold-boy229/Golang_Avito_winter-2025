package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port             string `mapstructure:"port"`
	ConnectionString string `mapstructure:"connection_string"`
}

var AppConfig *Config

func LoadAppConfig() {
	setConfigFileParametrs()

	if err := readConfigFile(); err != nil {
		log.Fatal(err)
	}

	if err := loadAppConfig(); err != nil {
		log.Fatal(err)
	}

	log.Println("Server Configurations were loaded successfully")
}

func setConfigFileParametrs() {
	viper.AddConfigPath("./config/")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
}

func readConfigFile() error {
	return viper.ReadInConfig()
}

func loadAppConfig() error {
	return viper.Unmarshal(&AppConfig)
}
