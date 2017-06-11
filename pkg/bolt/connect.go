package bolt

import (
	"path"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

func Open(dbDir string) (*bolt.DB, error) {
	filePath := path.Join(dbDir, "sprintbot.db")
	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "attempted to open boltdb sprintbot.db ")
	}
	return db, nil
}
