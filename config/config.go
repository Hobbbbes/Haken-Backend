package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

//Config contains important information for the running of the Programm
type Config struct {
	DBName     string `yaml:"dbname"`
	DBPassword string `yaml:"dbpassword"`
	DBUser     string `yaml:"dbuser"`
	Port       string `yaml:"port"`
}

//ReadConfig reads fills a config struct with information from the given file
func ReadConfig(path string) Config {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var c Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		panic(err)
	}
	return c
}
