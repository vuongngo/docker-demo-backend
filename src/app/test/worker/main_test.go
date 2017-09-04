package workers_test

import (
	"app/db"
	"app/model"
	"app/repo"
	"app/worker"
	"app/ws"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	nsq "github.com/bitly/go-nsq"
	"github.com/stretchr/testify/assert"
)

func TestJobMatch(t *testing.T) {
	assert := assert.New(t)
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

	nSQ.CreateHandler(func(message *nsq.Message) error {
		fmt.Println("Listener runnn")
		message.Touch()
		var jobMes model.JobMatchedMes
		err := json.Unmarshal(message.Body, &jobMes)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(jobMes.Id, "123", "Job received")
		message.Finish()
		return nil
	}, "job_matched", "analyze")

	producer, err := nSQ.CreateProducer()
	if err != nil {
		t.Error(err)
	}
	mes := &model.JobMatchedMes{
		Id:      "123",
		Session: "321",
		Ratio:   12.23,
	}
	m, _ := json.Marshal(mes)
	err = producer.Publish("job_matched", []byte(m))
	if err != nil {
		t.Error(err)
	}
	producer.Stop()
	time.Sleep(time.Second * 5)
}
