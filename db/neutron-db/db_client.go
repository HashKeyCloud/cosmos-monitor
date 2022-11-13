package neutron_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
	"github.com/spf13/viper"
)

type NeutronDbCli struct {
	*db.DbCli
}

func InitNeutronDbCli() (*NeutronDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.neutronUser"),
		Password: viper.GetString("postgres.neutronPassword"),
		Name:     viper.GetString("postgres.neutronName"),
		Host:     viper.GetString("postgres.neutronHost"),
		Port:     viper.GetString("postgres.neutronPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect neutron database server error: ", err)
		return nil, err
	}

	return &NeutronDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
