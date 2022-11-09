package main

import (
	"cosmosmonitor/config"
	"cosmosmonitor/log"
	"cosmosmonitor/monitor"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	var cfg = pflag.StringP("config", "c", "config/.conf.yaml", "config file path.")
	pflag.Parse()
	err := config.InitConfig(*cfg)
	if err != nil {
		fmt.Println("read config file error:", err)
		return
	}
	log.InitLogger(log.Logger, viper.GetString("log.path"))
	log.InitLogger(log.DBLogger, viper.GetString("log.path"))
	log.InitLogger(log.RPCLogger, viper.GetString("log.path"))
	log.InitLogger(log.MailLogger, viper.GetString("log.path"))
	log.InitJsonLogger(log.EventLogger, viper.GetString("log.eventlogpath"))
}

func main() {
	logger := log.Logger.WithField("module", "main")
	logger.Info("Successfully read config file!")
	m, err := monitor.NewMonitor()
	if err != nil {
		logger.Error("Failed to initialize monitoring client")
	}
	logger.Info("Starting Monitor")
	go m.Start()
	go m.SendEmail()
	m.WaitInterrupted()
}
