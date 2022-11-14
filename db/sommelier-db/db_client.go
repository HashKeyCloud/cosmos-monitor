package sommelier_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
	"github.com/spf13/viper"
)

type SommelierDbCli struct {
	*db.DbCli
}

func InitSommelierDbCli() (*SommelierDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.sommelierUser"),
		Password: viper.GetString("postgres.sommelierPassword"),
		Name:     viper.GetString("postgres.sommelierName"),
		Host:     viper.GetString("postgres.sommelierHost"),
		Port:     viper.GetString("postgres.sommelierPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect sommelier database server error: ", err)
		return nil, err
	}

	return &SommelierDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
