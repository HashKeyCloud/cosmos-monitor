package injective_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"

	"github.com/spf13/viper"
)

type InjectiveDbCli struct {
	*db.DbCli
}

func InitInjectiveDbCli() (*InjectiveDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.injectiveUser"),
		Password: viper.GetString("postgres.injectivePassword"),
		Name:     viper.GetString("postgres.injectiveName"),
		Host:     viper.GetString("postgres.injectiveHost"),
		Port:     viper.GetString("postgres.injectivePort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect injective database server error: ", err)
		return nil, err
	}

	return &InjectiveDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
