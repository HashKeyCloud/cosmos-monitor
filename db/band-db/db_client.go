package band_db

import (
	"github.com/spf13/viper"

	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
)

type BandDbCli struct {
	*db.DbCli
}

func InitBandDbCli() (*BandDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.bandUser"),
		Password: viper.GetString("postgres.bandPassword"),
		Name:     viper.GetString("postgres.bandName"),
		Host:     viper.GetString("postgres.bandHost"),
		Port:     viper.GetString("postgres.bandPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect band database server error: ", err)
		return nil, err
	}

	return &BandDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
