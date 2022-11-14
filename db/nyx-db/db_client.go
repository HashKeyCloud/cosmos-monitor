package nyx_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"

	"github.com/spf13/viper"
)

type NyxDbCli struct {
	*db.DbCli
}

func InitNyxDbCli() (*NyxDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.nyxUser"),
		Password: viper.GetString("postgres.nyxPassword"),
		Name:     viper.GetString("postgres.nyxName"),
		Host:     viper.GetString("postgres.nyxHost"),
		Port:     viper.GetString("postgres.nyxPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect nyx database server error: ", err)
		return nil, err
	}

	return &NyxDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
