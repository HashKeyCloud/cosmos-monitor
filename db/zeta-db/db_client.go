package zeta_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"

	"github.com/spf13/viper"
)

type ZetaDbCli struct {
	*db.DbCli
}

func InitZetaDbCli() (*ZetaDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.zetaUser"),
		Password: viper.GetString("postgres.zetaPassword"),
		Name:     viper.GetString("postgres.zetaName"),
		Host:     viper.GetString("postgres.zetaHost"),
		Port:     viper.GetString("postgres.zetaPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect zeta database server error: ", err)
		return nil, err
	}

	return &ZetaDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
