package config

import (
	"github.com/go-squads/unclog-worker/util"
	"github.com/spf13/viper"
)

var basepath = util.GetRootFolderPath()

func SetupConfig() error {
	viper.SetConfigFile(basepath + "config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
