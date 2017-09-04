package workers

import (
	"app/model"
	"app/ws"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bitly/go-nsq"
)

//JobMatched distribute message
func (n *NSQ) JobMatchedHandler(message *nsq.Message) error {
	message.Touch()
	fmt.Printf("Got a message: %v", string(message.Body[:]))
	// Decode message
	var jobMes model.JobMatchedMes
	err := json.Unmarshal(message.Body, &jobMes)
	if err != nil {
		log.Println(err)
	}
	// Find job
	job, err := n.Repo.GetJob(jobMes.Id)
	if err != nil {
		log.Println(err)
	}
	res, err := json.Marshal(struct {
		Job   *model.Job `json: "job"`
		Ratio float32    `json: "ratio"`
	}{
		Job:   job,
		Ratio: jobMes.Ratio,
	})
	if err != nil {
		log.Println(err)
	}

	n.Hub.Emit <- &ws.Emitter{Ids: []string{jobMes.Session}, Message: res}

	message.Finish()
	return nil
}
