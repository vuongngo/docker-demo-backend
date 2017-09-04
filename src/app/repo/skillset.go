package repo

import (
	"app/model"
)

func (repo *Repo) CreateSkillSet(skillset *model.SkillSet) error {
	db, session := repo.GetMgSession()
	defer session.Close()
	err := db.C("skillsets").Insert(skillset)
	return err
}

func (repo *Repo) GetSkillSet(id string) (*model.SkillSet, error) {
	db, session := repo.GetMgSession()
	defer session.Close()
	var skillset model.SkillSet
	err := db.C("skillsets").FindId(id).One(&skillset)
	if err != nil {
		return &model.SkillSet{}, err
	}
	return &skillset, nil
}
