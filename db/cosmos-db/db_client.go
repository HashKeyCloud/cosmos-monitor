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
		Username: viper.GetString("postgres.user"),
		Password: viper.GetString("postgres.password"),
		Name:     viper.GetString("postgres.name"),
		Host:     viper.GetString("postgres.host"),
		Port:     viper.GetString("postgres.port"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect database server error: ", err)
		return nil, err
	}

	return &CosmosDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
