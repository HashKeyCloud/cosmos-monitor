package secret_db

import (
	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"

	"github.com/spf13/viper"
)

type SecretDbCli struct {
	*db.DbCli
}

func InitSecretDbCli() (*SecretDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.secretUser"),
		Password: viper.GetString("postgres.secretPassword"),
		Name:     viper.GetString("postgres.secretName"),
		Host:     viper.GetString("postgres.secretHost"),
		Port:     viper.GetString("postgres.secretPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect secret database server error: ", err)
		return nil, err
	}

	return &SecretDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
