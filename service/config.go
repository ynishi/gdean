package service

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Host          string `yaml:host`
	Port          int    `yaml:port`
	MysqlHost     string `yaml:mysql_host`
	MysqlPort     uint   `yaml:mysql_port`
	MysqlUser     string `yaml:mysql_user`
	MysqlPassword string `yaml:mysql_password`
	MysqlDbName   string `yaml:mysql_dbname`
}

var cfg Config

func InitConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.AddConfigPath("$HOME/.gdean/config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Fatal error in read config:%w\n", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("Fatal error in parse config:%w\n", err)
	}
	return &cfg, nil
}
