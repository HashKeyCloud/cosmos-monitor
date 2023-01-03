package okp4_db

import (
	"github.com/spf13/viper"

	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
)

type Okp4DbCli struct {
	*db.DbCli
}

func InitOkp4DbCli() (*Okp4DbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.okp4User"),
		Password: viper.GetString("postgres.okp4Password"),
		Name:     viper.GetString("postgres.okp4Name"),
		Host:     viper.GetString("postgres.okp4Host"),
		Port:     viper.GetString("postgres.okp4Port"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect okp4 database server error: ", err)
		return nil, err
	}

	return &Okp4DbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
