package acrechain_db

import (
	"github.com/spf13/viper"

	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
)

type AcrechainDbCli struct {
	*db.DbCli
}

func InitAcrechainDbCli() (*AcrechainDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.acrechainUser"),
		Password: viper.GetString("postgres.acrechainPassword"),
		Name:     viper.GetString("postgres.acrechainName"),
		Host:     viper.GetString("postgres.acrechainHost"),
		Port:     viper.GetString("postgres.acrechainPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect acrechain database server error: ", err)
		return nil, err
	}

	return &AcrechainDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
