package config

import (
	"strings"

	"github.com/spf13/viper"

	"cosmosmonitor/log"
)

func InitConfig(cfg string) error {
	if cfg == "" {
		viper.AddConfigPath("config")
		viper.SetConfigName("conf")
	} else {
		viper.SetConfigFile(cfg)
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix(viper.GetString("env_prefix"))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func GetoperatorAddrs() []string {
	operatorAddr := viper.GetString("alert.operatorAddr")
	operatorAddrs := strings.Split(operatorAddr, ",")
	if len(operatorAddrs) == 0 {
		logger.Error("failed get operatorAddr")
	}
	return operatorAddrs
}

var logger = log.RPCLogger.WithField("module", "config")