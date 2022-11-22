package akash_db

import (
	"github.com/spf13/viper"

	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
)

type AkashDbCli struct {
	*db.DbCli
}

func InitAkashDbCli() (*AkashDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.akashUser"),
		Password: viper.GetString("postgres.akashPassword"),
		Name:     viper.GetString("postgres.akashName"),
		Host:     viper.GetString("postgres.akashHost"),
		Port:     viper.GetString("postgres.akashPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect akash database server error: ", err)
		return nil, err
	}

	return &AkashDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
