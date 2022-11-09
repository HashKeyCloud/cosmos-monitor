package juno_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
	"github.com/spf13/viper"
)

type JunoDbCli struct {
	*db.DbCli
}

func InitJunoDbCli() (*JunoDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.junoUser"),
		Password: viper.GetString("postgres.junoPassword"),
		Name:     viper.GetString("postgres.junoName"),
		Host:     viper.GetString("postgres.junoHost"),
		Port:     viper.GetString("postgres.junoPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect juno database server error: ", err)
		return nil, err
	}

	return &JunoDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
