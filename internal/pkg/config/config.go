package config

import (
	"dps/logger"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const (
	// Mysql is mysql driver
	Mysql = "mysql"
	// Postgres is postgres driver
	Postgres = "postgres"
)

// Store defines storage relevant info
type Store struct {
	Driver          string `yaml:"Driver" default:"mysql"`
	Host            string `yaml:"Host"`
	Port            int64  `yaml:"Port"`
	User            string `yaml:"User"`
	Password        string `yaml:"Password"`
	Db              string `yaml:"Db" default:"dps"`
	MaxOpenConns    int64  `yaml:"MaxOpenConns" default:"500"`
	MaxIdleConns    int64  `yaml:"MaxIdleConns" default:"500"`
	Schema          string `yaml:"Schema"`
	ConnMaxLifeTime int64  `yaml:"ConnMaxLifeTime" default:"5"`
}

// IsDB checks config driver is mysql or postgres
func (s *Store) IsDB() bool {
	return s.Driver == Mysql || s.Driver == Postgres
}

type App struct {
	Host string `yaml:"host" default:"127.0.0.1"`
	Port int64  `yaml:"port" default:"8080"`
}

// Config config
var Config = Type{}

// Type is the type for config of dps
type Type struct {
	App   App   `yaml:"App"`
	Store Store `yaml:"Store"`
}

// MustLoadConfig load config from env and file
func MustLoadConfig(confFile string) {
	loadFromEnv("", &Config)
	if confFile != "" {
		cont, err := ioutil.ReadFile(confFile)
		logger.FatalIfError(err)
		err = yaml.Unmarshal(cont, &Config)
		logger.FatalIfError(err)
	}
	scont, err := json.MarshalIndent(&Config, "", "  ")
	logger.FatalIfError(err)
	logger.Info(fmt.Sprintf("config file: %s loaded config is: \n%s", confFile, scont))
	err = checkConfig(&Config)
	logger.FatalfIf(err != nil, `config error: '%v'.`, err)
}
