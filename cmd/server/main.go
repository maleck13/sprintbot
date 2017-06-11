package main

import (
	"flag"
	"fmt"
	"net/http"

	"os"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/maleck13/sprintbot/pkg/bolt"
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
	rocketToken string
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
	viper.SetEnvPrefix("SB")
	viper.AutomaticEnv()
	viper.BindEnv("jira_board")
	viper.BindEnv("jira_sprint")
	viper.BindEnv("jira_user")
	viper.BindEnv("jira_pass")
	viper.BindEnv("jira_host")
	viper.BindEnv("github_token")
	viper.BindEnv("rocket_token")

	flag.StringVar(&logLevel, "log-level", "info", "use this to set log level: error, info, debug")
	flag.StringVar(&port, "port", "3000", "set the port to listen on. e.g 3000")
	flag.StringVar(&configLoc, "config-loc", "/etc/sprintbot", "the dir where to find the config")
	flag.StringVar(&gitHubToken, "github-token", "", "sets the github token")
	flag.StringVar(&rocketToken, "rocket-token", "", "sets the rocket chat auth token")
	flag.Parse()
	logger = setupLogger()
	router := web.BuildRouter()
	target := &sprintbot.Target{
		Host:     viper.GetString("jira_host"),
		UserName: viper.GetString("jira_user"),
		Password: viper.GetString("jira_pass"),
	}
	gitClient := github.NewClient(viper.GetString("github_token"))
	issueClient := jira.NewClient(target)
	//may want to refactor this
	_, err := issueClient.Login()
	if err != nil {
		logger.Fatalf("failed login to Jira %s ", err)
	}
	// db
	db, err := bolt.Open("./")
	if err != nil {
		logger.Fatalf("failed to open bolt db %s", err)
	}
	issueRepo := bolt.NewIssueRepo(db, logger)
	// sprintService
	sp := &sprintbot.Sprint{Name: viper.GetString("jira_sprint"), Board: viper.GetString("jira_board")}
	sprintService := sprint.NewService(issueClient, gitClient, issueRepo, sp, logger)
	sprintService.IgnoredRepos = []string{"RHMAPDocsNG", "fhcap", "fh-openshift-templates", "fh-core-openshift-templates"}
	//chat route
	{
		fmt.Println("sprint set to ", viper.GetString("jira_sprint"), os.Getenv("SB_JIRA_SPRINT"))
		chatUseCase := usecase.NewChat(sprintService)
		web.ChatRoute(router, chatUseCase, logger, viper.GetString("rocket_token"))
	}
	var shutDownChan = make(chan struct{})
	//start sync
	{

		go sprintService.Sync(20*time.Second, shutDownChan)
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
