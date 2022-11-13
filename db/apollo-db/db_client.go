package apollo_db

import (
	"github.com/spf13/viper"

	"cosmosmonitor/db"
	"cosmosmonitor/log"
	"cosmosmonitor/types"
)

type ApolloDbCli struct {
	*db.DbCli
}

func InitApolloDbCli() (*ApolloDbCli, error) {
	dbConf := &types.DatabaseConfig{
		Username: viper.GetString("postgres.apolloUser"),
		Password: viper.GetString("postgres.apolloPassword"),
		Name:     viper.GetString("postgres.apolloName"),
		Host:     viper.GetString("postgres.apolloHost"),
		Port:     viper.GetString("postgres.apolloPort"),
	}
	dbCli, err := db.InitDB(dbConf)
	if err != nil {
		logger.Error("connect apollo database server error: ", err)
		return nil, err
	}

	return &ApolloDbCli{
		DbCli: dbCli,
	}, nil
}

var logger = log.DBLogger.WithField("module", "db")
