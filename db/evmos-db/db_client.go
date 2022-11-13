package evmos_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
	"github.com/spf13/viper"
)

type EvmosDbCli struct {
	*db.DbCli
}

func InitEvmosDbCli() (*EvmosDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.evmosUser"),
		Password: viper.GetString("postgres.evmosPassword"),
		Name:     viper.GetString("postgres.evmosName"),
		Host:     viper.GetString("postgres.evmosHost"),
		Port:     viper.GetString("postgres.evmosPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect evmos database server error: ", err)
		return nil, err
	}

	return &EvmosDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
