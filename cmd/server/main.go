package main

import (
	"flag"
	"fmt"
	"net/http"

	"os"

	"github.com/Sirupsen/logrus"
	"github.com/maleck13/sprintbot/pkg/github"
	"github.com/maleck13/sprintbot/pkg/jira"
	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/maleck13/sprintbot/pkg/sprintbot/sprint"
	"github.com/maleck13/sprintbot/pkg/sprintbot/usecase"
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
	jiraBoard   string
	jiraSprint  string
	gitHubToken string
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
	viper.SetEnvPrefix("SB")
	viper.AutomaticEnv()
	viper.BindEnv("jira_board")
	viper.BindEnv("jira_sprint")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	flag.StringVar(&logLevel, "log-level", "info", "use this to set log level: error, info, debug")
	flag.StringVar(&port, "port", "3000", "set the port to listen on. e.g 3000")
	flag.StringVar(&configLoc, "config-loc", "/etc/sprintbot", "the dir where to find the config")
	flag.StringVar(&jiraHost, "jira-host", "", "sets the jira host to use")
	flag.StringVar(&jiraUser, "jira-user", "", "sets the jira user")
	flag.StringVar(&jiraPass, "jira-pass", "", "sets the jira password")

	flag.Parse()
	logger = setupLogger()
	router := web.BuildRouter()
	target := &sprintbot.Target{
		Host:     jiraHost,
		UserName: jiraUser,
		Password: jiraPass,
	}
	gitClient := &github.Client{}
	issueClient := jira.NewClient(target)
	_, err = issueClient.Login()
	if err != nil {
		logger.Fatalf("failed login to Jira %s ", err)
	}
	//chat route
	{
		fmt.Println("sprint set to ", viper.GetString("jira_sprint"), os.Getenv("SB_JIRA_SPRINT"))
		sprintService := sprint.NewService(issueClient, gitClient,
			viper.GetString("jira_board"), viper.GetString("jira_sprint"))
		chatUseCase := usecase.NewChat(sprintService)
		web.ChatRoute(router, chatUseCase, logger)
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
