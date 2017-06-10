package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"io/ioutil"

	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/pkg/errors"
)

func NewClient(t *sprintbot.Target) *Client {
	return &Client{
		target: t,
	}
}

type Client struct {
	target      *sprintbot.Target
	failedLogin int
}

func (c *Client) headers(req *http.Request) {
	req.Header.Add("content-type", "application/json")
	req.AddCookie(&http.Cookie{Name: "JSESSIONID", Value: c.target.Session})
}

func (c *Client) configure() http.Client {
	client := http.Client{}
	client.Timeout = time.Second * 15
	return client
}

func (c *Client) IssueHost() string {
	return c.target.Host
}

func (c *Client) FindUnresolvedOnBoard(boardName, sprint string) (*sprintbot.IssueList, error) {
	fmt.Println(c.target)
	bl, err := c.Boards(boardName)
	if err != nil {
		return nil, err
	}
	var boardID int
	for _, b := range bl.Values {
		if b.Name == boardName {
			boardID = b.ID
			break
		}
	}
	if boardID == 0 {
		return nil, errors.New("no board found for boardName " + boardName)
	}
	boardURL := fmt.Sprintf("%s/%s/%v/issue?jql=Sprint=\"%s\"", c.target.Host, "rest/agile/1.0/board",
		boardID, url.QueryEscape(sprint))
	req, err := http.NewRequest("GET", boardURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create board list request")
	}
	c.headers(req)
	client := c.configure()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do board list request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		_, err := c.Login()
		if err != nil {
			return nil, err
		}
		return c.FindUnresolvedOnBoard(boardName, sprint)
	}
	if resp.StatusCode > 300 {
		return nil, errors.New("Unexpected Jira statusCode: " + resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, " failed to read response body from Jira")
	}
	jIssues := &sprintbot.JiraIssueList{}
	if err := json.Unmarshal(data, jIssues); err != nil {
		return nil, errors.Wrap(err, "failed to decode the board issue list ")
	}
	var issues = make([]sprintbot.IssueState, len(jIssues.Issues))
	for i, j := range jIssues.Issues {
		issues[i] = j
	}
	issueList := sprintbot.NewIssueList(issues)
	return issueList, nil
}

func (c *Client) Boards(filter string) (*sprintbot.BoardList, error) {
	var list = &sprintbot.BoardList{}
	URL := fmt.Sprintf("%s/%s", c.target.Host, "rest/agile/1.0/board")
	if "" != filter {
		URL += "?name=" + url.QueryEscape(filter)
	}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create board list request")
	}
	c.headers(req)
	client := c.configure()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do board list request")
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, errors.New("Not authenticated with Jira statusCode: " + resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(list); err != nil {
		return nil, errors.Wrap(err, "failed to decode board list response at client Boards ")
	}
	return list, nil

}

func (c *Client) Login() (*sprintbot.Auth, error) {
	if c.failedLogin > 3 {
		return nil, errors.New("failed to login to Jira with credentials ")
	}
	login := `{"username":"` + c.target.UserName + `","password":"` + c.target.Password + `"}`
	authURL := fmt.Sprintf("%s/%s", c.target.Host, "rest/auth/1/session")
	req, err := http.NewRequest("POST", authURL, strings.NewReader(login))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create authentication request")
	}
	req.Header.Add("content-type", "application/json")
	authRes := &sprintbot.Auth{}
	client := c.configure()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do authentication request")
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		c.failedLogin++
		return nil, errors.New("failed to authenticate with Jira statusCode: " + resp.Status)
	}
	if resp.StatusCode != 200 {
		c.failedLogin++
		return nil, errors.New("failed to authenticate : " + resp.Status)
	}
	c.failedLogin = 0
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(authRes); err != nil {
		return nil, errors.Wrap(err, "failed to decode auth response from jira ")
	}
	c.target.Session = authRes.Session.Value
	return authRes, nil
}
