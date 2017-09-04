package model

import (
	"net/http"
	"time"

	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
)

type Job struct {
	ID          string    `bson:"_id" json:"_id" val_id:"required"`
	Title       string    `bson:"title" json:"title" valid:"required"`
	Description string    `bson:"description" json:"description" valid:"required"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt" valid:"required"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt" valid:"required"`
}

func InitializeJob() *Job {
	return &Job{
		ID:        bson.NewObjectId().Hex(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (job *Job) SetTitle(title string) {
	job.Title = title
}

func (job *Job) SetDescription(description string) {
	job.Description = description
}

type JobReq struct {
	Title       string `json:"title" xml:"title" form:"title" valid:"required"`
	Description string `json:"description" xml:"description" form:"description" valid:"required"`
}

func (l *JobReq) FieldMap(r *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&l.Title:       "title",
		&l.Description: "description",
	}
}

type JobMatchedMes struct {
	Id      string  `json:"id" xml:"id" form:"id" valid:"required"`
	Session string  `json:"session" xml:"session" form:"session" valsession:"required"`
	Ratio   float32 `json:"ratio" xml:"ratio" form:"ratio" valratio:"required"`
}
