package persistence_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
	"github.com/spf13/viper"
)

type PersistenceDbCli struct {
	*db.DbCli
}

func InitPersistenceDbCli() (*PersistenceDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.persistenceUser"),
		Password: viper.GetString("postgres.persistencePassword"),
		Name:     viper.GetString("postgres.persistenceName"),
		Host:     viper.GetString("postgres.persistenceHost"),
		Port:     viper.GetString("postgres.persistencePort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect persistence database server error: ", err)
		return nil, err
	}

	return &PersistenceDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
