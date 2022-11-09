package provider_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
	"github.com/spf13/viper"
)

type ProviderDbCli struct {
	*db.DbCli
}

func InitProviderDbCli() (*ProviderDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.providerUser"),
		Password: viper.GetString("postgres.providerPassword"),
		Name:     viper.GetString("postgres.providerName"),
		Host:     viper.GetString("postgres.providerHost"),
		Port:     viper.GetString("postgres.providerPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect provider database server error: ", err)
		return nil, err
	}

	return &ProviderDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
