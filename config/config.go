package config

import (
	"fmt"
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

func GetOperatorAddrs(project string) []string {
	operatorAddrConf := fmt.Sprintf("alert.%sOperatorAddr", project)
	operatorAddr := viper.GetString(operatorAddrConf)
	operatorAddrs := strings.Split(operatorAddr, ",")
	if len(operatorAddrs) == 0 {
		logger.Error("failed get cosmos operatorAddr")
	}
	return operatorAddrs
}

var logger = log.RPCLogger.WithField("module", "config")
