package model

import (
	"net/http"
	"time"

	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
)

type SkillSet struct {
	ID        string    `bson:"_id" json:"_id" val_id:"required"`
	Role      string    `bson:"role" json:"role" valid:"required"`
	Skills    []*Skill  `json:"skills"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt" valid:"required"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt" valid:"required"`
}

type Skill struct {
	ID   string `bson:"id" json:"id" valid:"required"`
	Name string `bson:"name" json:"name" valid:"required"`
}

func InitializeSkillSet() *SkillSet {
	return &SkillSet{
		ID:        bson.NewObjectId().Hex(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (skillset *SkillSet) SetRole(role string) {
	skillset.Role = role
}

func InitializeSkill(text string) *Skill {
	return &Skill{
		ID:   bson.NewObjectId().Hex(),
		Name: text,
	}
}

func (skillset *SkillSet) SetSkills(skills []*Skill) {
	skillset.Skills = skills
}

type SkillSetReq struct {
	Role   string   `json:"role" xml:"role" form:"role" valid:"required"`
	Skills []string `json:"skills" xml:"skills" form:"skills" valid:"required"`
}

func (l *SkillSetReq) FieldMap(r *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&l.Role:   "role",
		&l.Skills: "skills",
	}
}

type SkillMes struct {
	Id      string `json:"id" xml:"id" form:"id" valid:"required"`
	Session string `json:"session" xml:"session" form:"session" valsession:"required"`
}
