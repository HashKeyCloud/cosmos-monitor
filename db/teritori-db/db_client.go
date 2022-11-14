package teritori_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"

	"github.com/spf13/viper"
)

type TeritoriDbCli struct {
	*db.DbCli
}

func InitTeritoriDbCli() (*TeritoriDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.teritoriUser"),
		Password: viper.GetString("postgres.teritoriPassword"),
		Name:     viper.GetString("postgres.teritoriName"),
		Host:     viper.GetString("postgres.teritoriHost"),
		Port:     viper.GetString("postgres.teritoriPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect teritori database server error: ", err)
		return nil, err
	}

	return &TeritoriDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
