package api_test

import (
	"net/http"
	"os"
	"testing"

	"app/api"
	"app/db"
	"app/model"
	"app/repo"
	"app/worker"
	"app/ws"

	"github.com/gorilla/mux"
)

var (
	svc        *api.API
	jobId      string
	skillSetId string
	r          *mux.Router
)

func init() {
	//ws hub
	hub := ws.NewHub()
	go hub.Run()

	// Initialize connetion pool with mongo
	mongo := &db.MONGO{
		Uri:      os.Getenv("MONGO_URI"),
		Database: os.Getenv("MONGO_DATABASE"),
	}
	mongo.Dial()
	sharedRepo := &repo.Repo{
		Mongo: mongo,
	}

	//NSQ
	nSQ := workers.InitNSQ(os.Getenv("NSQ_URI"), hub, sharedRepo)
	err := nSQ.CreateHandler(nSQ.JobMatchedHandler, "job_matched", "analyze")
	if err != nil {
		panic(err)
	}

	// API share the same connection pool
	svc = &api.API{
		Repo: sharedRepo,
		Nsq:  nSQ,
	}

	// Create route with context
	r = mux.NewRouter()
	r.HandleFunc("/job/{id}", svc.GetJob)
	r.HandleFunc("/skill/{id}", svc.GetSkillSet)

	http.Handle("/", r)

}

func SetUp() {
	job := model.InitializeJob()
	job.SetTitle("Job Test")
	job.SetDescription("Test")
	jobId = job.ID

	skillset := model.InitializeSkillSet()
	skillset.SetRole("Test")
	skill := model.InitializeSkill("Test")
	skills := make([]*model.Skill, 0)
	skills = append(skills, skill)
	skillset.SetSkills(skills)
	skillSetId = skillset.ID

	db, session := svc.Repo.GetMgSession()
	defer session.Close()
	db.C("jobs").Insert(&job)
	db.C("skillsets").Insert(&skillset)
}

func TearDown(collection string) {
	db, session := svc.Repo.GetMgSession()
	db.C(collection).RemoveAll(nil)
	session.Close()
}

func TestMain(m *testing.M) {
	SetUp()
	m.Run()
	TearDown("skillsets")
	TearDown("jobs")
}
