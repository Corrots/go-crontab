package common

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig(filename *string) error {
	viper.SetConfigFile(*filename)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//viper.AddConfigPath("./")
			viper.SetConfigName("conf")
			viper.SetConfigType("yaml")
		} else {
			return fmt.Errorf("read config file err: %v\n", err)
		}
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Printf("config file changed: %v\n", in.String())
	})
	return nil
}
