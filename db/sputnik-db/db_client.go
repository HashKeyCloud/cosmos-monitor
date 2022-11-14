package sputnik_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"

	"github.com/spf13/viper"
)

type SputnikDbCli struct {
	*db.DbCli
}

func InitSputnikDbCli() (*SputnikDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.sputnikUser"),
		Password: viper.GetString("postgres.sputnikPassword"),
		Name:     viper.GetString("postgres.sputnikName"),
		Host:     viper.GetString("postgres.sputnikHost"),
		Port:     viper.GetString("postgres.sputnikPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect sputnik database server error: ", err)
		return nil, err
	}

	return &SputnikDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
