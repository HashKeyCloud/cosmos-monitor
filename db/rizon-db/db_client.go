package rizon_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
	"github.com/spf13/viper"
)

type RizonDbCli struct {
	*db.DbCli
}

func InitRizonDbCli() (*RizonDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.rizonUser"),
		Password: viper.GetString("postgres.rizonPassword"),
		Name:     viper.GetString("postgres.rizonName"),
		Host:     viper.GetString("postgres.rizonHost"),
		Port:     viper.GetString("postgres.rizonPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect rizon database server error: ", err)
		return nil, err
	}

	return &RizonDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
