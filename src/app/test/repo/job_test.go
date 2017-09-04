package repo_test

import (
	"testing"

	"app/model"
)

var (
	prepJob *Prep
)

func init() {
	prepJob = InitPrep("jobs")
}

func TestCreatJob(t *testing.T) {
	job := model.InitializeJob()
	job.SetTitle("test")
	job.SetDescription("test des")
	err := prepJob.Repo.CreateJob(job)
	if err != nil {
		t.Error(err)
	}
}

func TestGetJob(t *testing.T) {
	job := model.InitializeJob()
	job.SetTitle("test")
	job.SetDescription("test des")
	db, session := prepJob.Repo.GetMgSession()
	defer session.Close()
	db.C("jobs").Insert(job)
	result, err := prepJob.Repo.GetJob(job.ID)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
