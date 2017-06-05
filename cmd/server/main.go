package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/maleck13/sprintbot/pkg/web"
	"github.com/spf13/viper"
)

var (
	logLevel    string
	port        string
	logger      *logrus.Logger
	jiraHost    string
	jiraUser    string
	jiraPass    string
	gitHubToken string
	rocketHost  string
	rocketUser  string
	rocketPass  string
	configLoc   string
)

func setupLogger() *logrus.Logger {
	switch logLevel {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	}
	return logrus.StandardLogger()
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("/etc/sprintbot")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	flag.StringVar(&logLevel, "log-level", "info", "use this to set log level: error, info, debug")
	flag.StringVar(&port, "port", "3000", "set the port to listen on. e.g 3000")
	flag.StringVar(&configLoc, "config-loc", "/etc/sprintbot", "the dir where to find the config")
	flag.Parse()
	logger = setupLogger()
	router := web.BuildRouter()
	//chat route
	{
		web.ChatRoute(router)
	}

	//http handler
	{
		port := ":3000"
		logrus.Info("starting sprintbot on  port " + port)
		httpHandler := web.BuildHTTPHandler(router)
		if err := http.ListenAndServe(port, httpHandler); err != nil {
			logger.Fatal(err)
		}
	}

}
