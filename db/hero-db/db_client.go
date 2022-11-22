package hero_db

import (
	"github.com/spf13/viper"

	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
)

type HeroDbCli struct {
	*db.DbCli
}

func InitHeroDbCli() (*HeroDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.heroUser"),
		Password: viper.GetString("postgres.heroPassword"),
		Name:     viper.GetString("postgres.heroName"),
		Host:     viper.GetString("postgres.heroHost"),
		Port:     viper.GetString("postgres.heroPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect hero database server error: ", err)
		return nil, err
	}

	return &HeroDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
