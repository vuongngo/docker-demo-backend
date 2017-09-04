package repo

import (
	"app/model"
)

func (repo *Repo) CreateJob(job *model.Job) error {
	db, session := repo.GetMgSession()
	defer session.Close()
	err := db.C("jobs").Insert(job)
	return err
}

func (repo *Repo) GetJob(id string) (*model.Job, error) {
	db, session := repo.GetMgSession()
	defer session.Close()
	var job model.Job
	err := db.C("jobs").FindId(id).One(&job)
	if err != nil {
		return &model.Job{}, err
	}
	return &job, nil
}
