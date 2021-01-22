package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type Struck struct {
	Redis struct {
		User           string        `yaml:"user"`
		Db             int           `yaml:"db"`
		Host           string        `yaml:"host"`
		Password       string        `yaml:"password"`
		Port           string        `yaml:"port"`
		MaxIdle        int           `yaml:"max_idle"`
		ConnectTimeOut time.Duration `yaml:"connect_time_out"`
		WriteTimeOut   time.Duration `yaml:"write_time_out"`
		ReadTimeOut    time.Duration `yaml:"read_time_out"`
	} `yaml:"redis"`
}

var (
	Config = &Struck{}
)

func init() {
	readYaml()
}

// ReadYaml ..
func readYaml() {
	f, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err.Error())
	}
	err = yaml.Unmarshal(f, Config)
	if err != nil {
		log.Fatal(err.Error())
	}
}
