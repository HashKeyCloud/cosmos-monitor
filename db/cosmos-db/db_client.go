package cosmos_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
	"github.com/spf13/viper"
)

type CosmosDbCli struct {
	*db.DbCli
}

func InitCosmosDbCli() (*CosmosDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.cosmosUser"),
		Password: viper.GetString("postgres.cosmosPassword"),
		Name:     viper.GetString("postgres.cosmosName"),
		Host:     viper.GetString("postgres.cosmosHost"),
		Port:     viper.GetString("postgres.cosmosPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect cosmos database server error: ", err)
		return nil, err
	}

	return &CosmosDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
