package sprintbot

import (
	"fmt"
	"strings"
)

type JiraIssue struct {
	Expand string `json:"expand"`
	Fields struct {
		Aggregateprogress struct {
			Progress int `json:"progress"`
			Total    int `json:"total"`
		} `json:"aggregateprogress"`
		Aggregatetimeestimate         interface{} `json:"aggregatetimeestimate"`
		Aggregatetimeoriginalestimate interface{} `json:"aggregatetimeoriginalestimate"`
		Aggregatetimespent            interface{} `json:"aggregatetimespent"`
		Assignee                      struct {
			Active       bool   `json:"active"`
			DisplayName  string `json:"displayName"`
			EmailAddress string `json:"emailAddress"`
			Key          string `json:"key"`
			Name         string `json:"name"`
			Self         string `json:"self"`
			TimeZone     string `json:"timeZone"`
		} `json:"assignee"`
		Comment struct {
			Comments []struct {
				Author struct {
					Active     bool `json:"active"`
					AvatarUrls struct {
						One6x16   string `json:"16x16"`
						Two4x24   string `json:"24x24"`
						Three2x32 string `json:"32x32"`
						Four8x48  string `json:"48x48"`
					} `json:"avatarUrls"`
					DisplayName  string `json:"displayName"`
					EmailAddress string `json:"emailAddress"`
					Key          string `json:"key"`
					Name         string `json:"name"`
					Self         string `json:"self"`
					TimeZone     string `json:"timeZone"`
				} `json:"author"`
				Body         string `json:"body"`
				Created      string `json:"created"`
				ID           string `json:"id"`
				Self         string `json:"self"`
				UpdateAuthor struct {
					Active     bool `json:"active"`
					AvatarUrls struct {
						One6x16   string `json:"16x16"`
						Two4x24   string `json:"24x24"`
						Three2x32 string `json:"32x32"`
						Four8x48  string `json:"48x48"`
					} `json:"avatarUrls"`
					DisplayName  string `json:"displayName"`
					EmailAddress string `json:"emailAddress"`
					Key          string `json:"key"`
					Name         string `json:"name"`
					Self         string `json:"self"`
					TimeZone     string `json:"timeZone"`
				} `json:"updateAuthor"`
				Updated string `json:"updated"`
			} `json:"comments"`
			MaxResults int `json:"maxResults"`
			StartAt    int `json:"startAt"`
			Total      int `json:"total"`
		} `json:"comment"`
		Created string `json:"created"`
		Creator struct {
			Active     bool `json:"active"`
			AvatarUrls struct {
				One6x16   string `json:"16x16"`
				Two4x24   string `json:"24x24"`
				Three2x32 string `json:"32x32"`
				Four8x48  string `json:"48x48"`
			} `json:"avatarUrls"`
			DisplayName  string `json:"displayName"`
			EmailAddress string `json:"emailAddress"`
			Key          string `json:"key"`
			Name         string `json:"name"`
			Self         string `json:"self"`
			TimeZone     string `json:"timeZone"`
		} `json:"creator"`
		Customfield12310080 interface{} `json:"customfield_12310080"`
		Customfield12310090 interface{} `json:"customfield_12310090"`
		Customfield12310091 interface{} `json:"customfield_12310091"`
		Customfield12310120 interface{} `json:"customfield_12310120"`
		Customfield12310160 interface{} `json:"customfield_12310160"`
		Customfield12310183 interface{} `json:"customfield_12310183"`
		Customfield12310211 interface{} `json:"customfield_12310211"`
		Customfield12310213 interface{} `json:"customfield_12310213"`
		Customfield12310214 interface{} `json:"customfield_12310214"`
		Customfield12310220 []string    `json:"customfield_12310220"`
		Customfield12310241 interface{} `json:"customfield_12310241"`
		Customfield12310243 float64     `json:"customfield_12310243"`
		Customfield12310640 string      `json:"customfield_12310640"`
		Customfield12310641 string      `json:"customfield_12310641"`
		Customfield12310840 string      `json:"customfield_12310840"`
		Customfield12310940 []string    `json:"customfield_12310940"`
		Customfield12311140 interface{} `json:"customfield_12311140"`
		Customfield12311240 interface{} `json:"customfield_12311240"`
		Customfield12311640 interface{} `json:"customfield_12311640"`
		Customfield12311641 interface{} `json:"customfield_12311641"`
		Customfield12311940 string      `json:"customfield_12311940"`
		Customfield12312440 interface{} `json:"customfield_12312440"`
		Customfield12312441 interface{} `json:"customfield_12312441"`
		Customfield12312442 interface{} `json:"customfield_12312442"`
		Customfield12312640 interface{} `json:"customfield_12312640"`
		Customfield12313140 interface{} `json:"customfield_12313140"`
		Customfield12313240 interface{} `json:"customfield_12313240"`
		Customfield12313340 interface{} `json:"customfield_12313340"`
		Customfield12313440 string      `json:"customfield_12313440"`
		Customfield12313441 string      `json:"customfield_12313441"`
		Customfield12313640 string      `json:"customfield_12313640"`
		Customfield12313641 interface{} `json:"customfield_12313641"`
		Description         string      `json:"description"`
		Duedate             interface{} `json:"duedate"`
		Environment         interface{} `json:"environment"`
		FixVersions         []struct {
			Archived    bool   `json:"archived"`
			Description string `json:"description"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			ReleaseDate string `json:"releaseDate"`
			Released    bool   `json:"released"`
			Self        string `json:"self"`
		} `json:"fixVersions"`
		Flagged    bool `json:"flagged"`
		Issuelinks []struct {
			ID          string `json:"id"`
			InwardIssue struct {
				Fields struct {
					Issuetype struct {
						AvatarID    int    `json:"avatarId"`
						Description string `json:"description"`
						IconURL     string `json:"iconUrl"`
						ID          string `json:"id"`
						Name        string `json:"name"`
						Self        string `json:"self"`
						Subtask     bool   `json:"subtask"`
					} `json:"issuetype"`
					Priority struct {
						IconURL string `json:"iconUrl"`
						ID      string `json:"id"`
						Name    string `json:"name"`
						Self    string `json:"self"`
					} `json:"priority"`
					Status struct {
						Description    string `json:"description"`
						IconURL        string `json:"iconUrl"`
						ID             string `json:"id"`
						Name           string `json:"name"`
						Self           string `json:"self"`
						StatusCategory struct {
							ColorName string `json:"colorName"`
							ID        int    `json:"id"`
							Key       string `json:"key"`
							Name      string `json:"name"`
							Self      string `json:"self"`
						} `json:"statusCategory"`
					} `json:"status"`
					Summary string `json:"summary"`
				} `json:"fields"`
				ID   string `json:"id"`
				Key  string `json:"key"`
				Self string `json:"self"`
			} `json:"inwardIssue"`
			Self string `json:"self"`
			Type struct {
				ID      string `json:"id"`
				Inward  string `json:"inward"`
				Name    string `json:"name"`
				Outward string `json:"outward"`
				Self    string `json:"self"`
			} `json:"type"`
		} `json:"issuelinks"`
		Issuetype struct {
			AvatarID    int    `json:"avatarId"`
			Description string `json:"description"`
			IconURL     string `json:"iconUrl"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			Self        string `json:"self"`
			Subtask     bool   `json:"subtask"`
		} `json:"issuetype"`
		Labels     []string `json:"labels"`
		LastViewed string   `json:"lastViewed"`
		Priority   struct {
			IconURL string `json:"iconUrl"`
			ID      string `json:"id"`
			Name    string `json:"name"`
			Self    string `json:"self"`
		} `json:"priority"`
		Progress struct {
			Progress int `json:"progress"`
			Total    int `json:"total"`
		} `json:"progress"`
		Project struct {
			AvatarUrls struct {
				One6x16   string `json:"16x16"`
				Two4x24   string `json:"24x24"`
				Three2x32 string `json:"32x32"`
				Four8x48  string `json:"48x48"`
			} `json:"avatarUrls"`
			ID   string `json:"id"`
			Key  string `json:"key"`
			Name string `json:"name"`
			Self string `json:"self"`
		} `json:"project"`
		Reporter struct {
			Active     bool `json:"active"`
			AvatarUrls struct {
				One6x16   string `json:"16x16"`
				Two4x24   string `json:"24x24"`
				Three2x32 string `json:"32x32"`
				Four8x48  string `json:"48x48"`
			} `json:"avatarUrls"`
			DisplayName  string `json:"displayName"`
			EmailAddress string `json:"emailAddress"`
			Key          string `json:"key"`
			Name         string `json:"name"`
			Self         string `json:"self"`
			TimeZone     string `json:"timeZone"`
		} `json:"reporter"`
		Resolution     interface{} `json:"resolution"`
		Resolutiondate interface{} `json:"resolutiondate"`
		Sprint         struct {
			EndDate       string `json:"endDate"`
			ID            int    `json:"id"`
			Name          string `json:"name"`
			OriginBoardID int    `json:"originBoardId"`
			Self          string `json:"self"`
			StartDate     string `json:"startDate"`
			State         string `json:"state"`
		} `json:"sprint"`
		Status struct {
			Description    string `json:"description"`
			IconURL        string `json:"iconUrl"`
			ID             string `json:"id"`
			Name           string `json:"name"`
			Self           string `json:"self"`
			StatusCategory struct {
				ColorName string `json:"colorName"`
				ID        int    `json:"id"`
				Key       string `json:"key"`
				Name      string `json:"name"`
				Self      string `json:"self"`
			} `json:"statusCategory"`
		} `json:"status"`
		Subtasks             []interface{} `json:"subtasks"`
		Summary              string        `json:"summary"`
		Timeestimate         interface{}   `json:"timeestimate"`
		Timeoriginalestimate interface{}   `json:"timeoriginalestimate"`
		Timespent            interface{}   `json:"timespent"`
		Timetracking         struct{}      `json:"timetracking"`
		Updated              string        `json:"updated"`
		Versions             []struct {
			Archived    bool   `json:"archived"`
			Description string `json:"description"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			ReleaseDate string `json:"releaseDate"`
			Released    bool   `json:"released"`
			Self        string `json:"self"`
		} `json:"versions"`
		Workratio int `json:"workratio"`
	} `json:"fields"`
	Meta struct {
		PrOpen bool
	} `json:"-"`
	Id   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

func (ji *JiraIssue) ID() string {
	return ji.Id
}

func (ji *JiraIssue) PRS() []string {
	return ji.Fields.Customfield12310220
}

func (ji *JiraIssue) RemovePR(pr string) {
	na := []string{}
	for _, p := range ji.PRS() {
		if p != pr {
			na = append(na, p)
		}
	}
	ji.Fields.Customfield12310220 = na
}

func (ji *JiraIssue) Link(host string) string {
	return host + "/browse/" + ji.Key
}

func (ji *JiraIssue) Description() string {
	return ji.Fields.Summary
}

func (ji *JiraIssue) State() string {
	return ji.Fields.Status.Name
}

type JiraIssueList struct {
	Issues []*JiraIssue `json:"issues"`
}

type RocketChatCmd struct {
	Bot         bool   `json:"bot"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	IsEdited    bool   `json:"isEdited"`
	MessageID   string `json:"message_id"`
	Text        string `json:"text"`
	Timestamp   string `json:"timestamp"`
	Token       string `json:"token"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
}

func (rcmd *RocketChatCmd) Action() string {
	return strings.Replace(rcmd.Text, "sprintbot ", "", 1)
}

func (rcmd *RocketChatCmd) User() string {
	return rcmd.UserName
}

func (rcmd *RocketChatCmd) AuthToken() string {
	return rcmd.Token
}

type Sprint struct {
	Name  string
	Board string
}

type ChatResponse struct{}

type NextIssues struct {
	Message string
	Issues  []*Issue
}

type Issue struct {
	Link        string
	PRs         []string
	Description string
}

type IssueList struct {
	issues []IssueState
}

func NewIssueList(issues []IssueState) *IssueList {
	return &IssueList{
		issues: issues,
	}
}

func (il *IssueList) Issues() []IssueState {
	return il.issues
}

func (il *IssueList) FindInState(state string) *IssueList {
	var issues = []IssueState{}
	for _, i := range il.issues {
		if i.State() == state {
			issues = append(issues, i)
		}
	}
	return &IssueList{issues: issues}
}

const (
	IssueStateReadyForQA = "Ready for QA"
	IssueStateOpen       = "Open"
	IssueStateClosed     = "Closed"
	IssueStatePRSent     = "Pull Request Sent"

	CommentTypeMoveToReadyForQE = "moveToQE"
)

type ErrUnkownCMD struct {
	Message string
}

func (e *ErrUnkownCMD) Error() string {
	return fmt.Sprintf("Error unknown cmd: %s", e.Message)
}

type ErrInvalid struct {
	Message string
}

func (e *ErrInvalid) Error() string {
	return fmt.Sprintf("Error Invalid : %s", e.Message)
}

type Target struct {
	Host     string
	UserName string
	Password string
	Session  string
}

type Auth struct {
	LoginInfo struct {
		LoginCount        int    `json:"loginCount"`
		PreviousLoginTime string `json:"previousLoginTime"`
	} `json:"loginInfo"`
	Session struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"session"`
}

type BoardList struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	Total      int      `json:"total"`
	IsLast     bool     `json:"isLast"`
	Values     []*Board `json:"values"`
}

type Board struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Self string `json:"self"`
	Type string `json:"type"`
}
