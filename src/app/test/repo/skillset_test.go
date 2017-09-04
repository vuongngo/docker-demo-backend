package repo_test

import (
	"testing"

	"app/model"
)

var (
	prepSkill *Prep
)

func init() {
	prepSkill = InitPrep("skillsets")
}

func TestCreatSkillSet(t *testing.T) {
	skillset := model.InitializeSkillSet()
	skillset.SetRole("test")
	err := prepSkill.Repo.CreateSkillSet(skillset)
	if err != nil {
		t.Error(err)
	}
}

func TestGetSkillSet(t *testing.T) {
	skillset := model.InitializeSkillSet()
	skillset.SetRole("another_test")
	db, session := prepJob.Repo.GetMgSession()
	defer session.Close()
	db.C("skillsets").Insert(skillset)
	result, err := prepSkill.Repo.GetSkillSet(skillset.ID)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}
