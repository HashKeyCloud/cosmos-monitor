package gopher_db

import (
	"github.com/spf13/viper"

	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
)

type GopherDbCli struct {
	*db.DbCli
}

func InitGopherDbCli() (*GopherDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.gopherUser"),
		Password: viper.GetString("postgres.gopherPassword"),
		Name:     viper.GetString("postgres.gopherName"),
		Host:     viper.GetString("postgres.gopherHost"),
		Port:     viper.GetString("postgres.gopherPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect gopher database server error: ", err)
		return nil, err
	}

	return &GopherDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
