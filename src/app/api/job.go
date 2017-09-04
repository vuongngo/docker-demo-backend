package api

import (
	"encoding/json"
	"net/http"

	"app/model"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
)

/**
* GET /api/v1/job/:id
 */
func (api *API) GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["id"]
	job, err := api.Repo.GetJob(jobId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(job)
}

/**
* POST /api/v1/job
 */
func (api *API) CreateJob(w http.ResponseWriter, r *http.Request) {
	jobReq := new(model.JobReq)
	errs := binding.Bind(r, jobReq)
	if errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errs.Error()))
		return
	}

	if res, err := govalidator.ValidateStruct(jobReq); res == false {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	job := model.InitializeJob()
	job.SetTitle(jobReq.Title)
	job.SetDescription(jobReq.Description)
	err := api.Repo.CreateJob(job)
	if err != nil {
		w.WriteHeader(500)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}
