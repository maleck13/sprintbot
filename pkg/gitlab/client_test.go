package gitlab_test

import (
	"testing"

	"github.com/maleck13/sprintbot/pkg/gitlab"

	"os"
)

func TestPRRevied(t *testing.T){
	//TODO
	token := os.Getenv("GITLAB_TOKEN")
	c := gitlab.NewClient(token)
}


func TestRepo(t *testing.T) {

	//TODO
}