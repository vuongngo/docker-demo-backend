package repo_test

import (
	"os"
	"testing"

	"app/db"
	"app/repo"

	"gopkg.in/mgo.v2"
)

type Prep struct {
	Session    *mgo.Session
	Repo       *repo.Repo
	Collection string
}

func InitPrep(collection string) *Prep {
	mongo := &db.MONGO{
		Uri:      os.Getenv("MONGO_URI"),
		Database: os.Getenv("MONGO_DATABASE"),
	}
	mongo.Dial()
	dbM := &repo.Repo{
		Mongo: mongo,
	}
	return &Prep{
		Repo:       dbM,
		Collection: collection,
	}
}

func (prep *Prep) TearDown() {
	db, session := prepJob.Repo.GetMgSession()
	defer session.Close()
	db.C(prep.Collection).RemoveAll(nil)
}

func TestMain(m *testing.M) {
	preps := []*Prep{
		prepSkill,
		prepJob,
	}
	m.Run()
	for _, pre := range preps {
		pre.TearDown()
	}
}
