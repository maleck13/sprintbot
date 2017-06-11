package bolt

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/maleck13/sprintbot/pkg/sprintbot"
	"github.com/pkg/errors"
)

type IssueRepo struct {
	logger *logrus.Logger
	DB     *bolt.DB
}

func NewIssueRepo(db *bolt.DB, logger *logrus.Logger) *IssueRepo {
	return &IssueRepo{
		logger: logger,
		DB:     db,
	}
}

const (
	issueBucket = "issues"
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

// FindNext will return the issues that need looking at next
func (is *IssueRepo) FindNext() (*sprintbot.NextIssues, error) {
	is.logger.Debug("finding next issues")
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
