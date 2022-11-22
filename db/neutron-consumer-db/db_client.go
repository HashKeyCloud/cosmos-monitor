package neutron_consumer_db

import (
	"github.com/spf13/viper"

	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
)

type NeutronconsumerDbCli struct {
	*db.DbCli
}

func InitNeutronconsumerDbCli() (*NeutronconsumerDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.neutronconsumerUser"),
		Password: viper.GetString("postgres.neutronconsumerPassword"),
		Name:     viper.GetString("postgres.neutronconsumerName"),
		Host:     viper.GetString("postgres.neutronconsumerHost"),
		Port:     viper.GetString("postgres.neutronconsumerPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect neutron consumer chain database server error: ", err)
		return nil, err
	}

	return &NeutronconsumerDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
