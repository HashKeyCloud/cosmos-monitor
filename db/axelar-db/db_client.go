package axelar_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"

	"github.com/spf13/viper"
)

type AxelarDbCli struct {
	*db.DbCli
}

func InitAxelarDbCli() (*AxelarDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.axelarUser"),
		Password: viper.GetString("postgres.axelarPassword"),
		Name:     viper.GetString("postgres.axelarName"),
		Host:     viper.GetString("postgres.axelarHost"),
		Port:     viper.GetString("postgres.axelarPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect axelar database server error: ", err)
		return nil, err
	}

	return &AxelarDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
