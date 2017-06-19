package bolt

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/pkg/errors"
)

type IssueRepo struct {
	logger  *logrus.Logger
	DB      *bolt.DB
	decoder sprintbot.IssueDecoder
}

func NewIssueRepo(db *bolt.DB, logger *logrus.Logger, decoder sprintbot.IssueDecoder) *IssueRepo {
	return &IssueRepo{
		logger:  logger,
		DB:      db,
		decoder: decoder,
	}
}

const (
	issueBucket   = "issues"
	commentBucket = "issueComments"
)

// SaveNext will save the issues that should be done next
func (is *IssueRepo) SaveNext(next *sprintbot.NextIssues) error {
	is.logger.Debug("saving next issues ")
	return is.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(issueBucket))
		if err != nil {
			return errors.Wrap(err, "when saving the next issues failed to create a bucket ")
		}
		data, err := json.Marshal(next)
		if err != nil {
			return errors.Wrap(err, "when saving the next issues, failed to marshal the target to json data ")
		}
		if err := b.Put([]byte("nextIssues"), data); err != nil {
			return errors.Wrap(err, "when saving the next issues, failed to put the target in the bolt db bucket")
		}
		return nil
	})

}

func (is *IssueRepo) SaveCommented(id string, commentID string) error {
	return is.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(commentBucket))
		if err != nil {
			return errors.Wrap(err, "when saving commentted issue failed to create a bucket ")
		}
		if err := b.Put([]byte(id), []byte(commentID)); err != nil {
			return errors.Wrap(err, " failed to save commented issue ")
		}
		return nil
	})
}

func (is *IssueRepo) FindCommentOnIssue(id string, commentID string) (string, error) {
	var ret []byte
	err := is.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(commentBucket))
		if nil == b {
			return nil
		}
		ret = b.Get([]byte(id))
		return nil
	})
	if err != nil {
		return "", errors.Wrap(err, "bolt issue repo failed to find comment on issue ")
	}
	return string(ret), err
}

// FindNext will return the issues that need looking at next
func (is *IssueRepo) FindNext() (*sprintbot.NextIssues, error) {
	var data []byte
	var ret = &sprintbot.NextIssues{}
	err := is.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(issueBucket))
		if nil == b {
			return errors.New("no such bucket found " + issueBucket)
		}
		data = b.Get([]byte("nextIssues"))
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "bolt issue repo failed to find next issues ")
	}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal the next issues ")
	}
	return ret, nil
}

func (is *IssueRepo) SaveIssuesInState(state string, issues []sprintbot.IssueState) error {
	return is.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(issueBucket))
		if err != nil {
			return errors.Wrap(err, "when saving issues in state failed to create a bucket ")
		}
		data, err := json.Marshal(issues)
		if err != nil {
			return errors.Wrap(err, "when saving  issues in state, failed to marshal the target to json data ")
		}
		if err := b.Put([]byte(state), data); err != nil {
			return errors.Wrap(err, "when saving issues in state, failed to put the target in the bolt db bucket")
		}
		return nil
	})
}

func (is *IssueRepo) FindIssuesInState(state string) ([]sprintbot.IssueState, error) {
	var data []byte
	err := is.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(issueBucket))
		if nil == b {
			return errors.New("no such bucket found " + issueBucket)
		}
		data = b.Get([]byte(state))
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "bolt issue repo failed to FindIssuesInState  ")
	}
	issues, err := is.decoder.Decode(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal the issues in state ")
	}

	return issues, err
}

func (is *IssueRepo) FindIssuesNotInState(state string) ([]sprintbot.IssueState, error) {
	var data []byte
	var issuesNotInState []sprintbot.IssueState
	err := is.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(issueBucket))
		if nil == b {
			return errors.New("no such bucket found " + issueBucket)
		}
		b.ForEach(func(k, v []byte) error {
			key := string(k)
			if key != state {
				data = b.Get([]byte(state))
				issues, err := is.decoder.Decode(data)
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal the issues in state ")
				}
				issuesNotInState = append(issuesNotInState, issues...)
			}
			return nil
		})

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "bolt issue repo failed to FindIssuesInState  ")
	}
	return issuesNotInState, err
}
