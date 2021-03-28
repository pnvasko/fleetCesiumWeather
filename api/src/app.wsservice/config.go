package app_wsservice

import (
	base "app.base"
	"github.com/spf13/viper"
)

type Config struct {
	Name           string
	Debug          bool
	LogPath        string
	WorkerPoolSize int
	HttpConfig     base.HttpConfig
	Db             base.DbConfig
}

func newConfig() *Config {
	return &Config{}
}

func (conf *Config) Load(path string) error {
	var cc Config
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&cc); err != nil {
		return err
	}
	*conf = cc
	return nil
}
