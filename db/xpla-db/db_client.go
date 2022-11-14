package xpla_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"

	"github.com/spf13/viper"
)

type XplaDbCli struct {
	*db.DbCli
}

func InitXplaDbCli() (*XplaDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.xplaUser"),
		Password: viper.GetString("postgres.xplaPassword"),
		Name:     viper.GetString("postgres.xplaName"),
		Host:     viper.GetString("postgres.xplaHost"),
		Port:     viper.GetString("postgres.xplaPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect xpla database server error: ", err)
		return nil, err
	}

	return &XplaDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
