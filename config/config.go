package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	name string
}

func (c *Config) initWithConfig() error {

	if c.name != ""{
		viper.SetConfigFile(c.name)
	}else{
		viper.AddConfigPath("conf")
		viper.SetConfigFile("config")
	}
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}


func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Printf("Config file changed %v", in.Name)
	})
}