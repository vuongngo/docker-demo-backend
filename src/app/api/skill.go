package api

import (
	"app/model"
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
)

/**
* GET /api/v1/skill/:id
 */
func (api *API) GetSkillSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	skillSetId := vars["id"]
	skillSet, err := api.Repo.GetSkillSet(skillSetId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(skillSet)
}

/**
* POST /api/v1/skill
 */
func (api *API) CreateSkillSet(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	token := cookie.Value

	skillSetReq := new(model.SkillSetReq)
	errs := binding.Bind(r, skillSetReq)
	if errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errs.Error()))
		return
	}

	if res, err := govalidator.ValidateStruct(skillSetReq); res == false {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	skillSet := model.InitializeSkillSet()
	skillSet.SetRole(skillSetReq.Role)
	skills := make([]*model.Skill, 0)
	for _, s := range skillSetReq.Skills {
		newSkill := model.InitializeSkill(s)
		skills = append(skills, newSkill)
	}
	skillSet.SetSkills(skills)

	// Persisting skillset
	err = api.Repo.CreateSkillSet(skillSet)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Publish message to nsq
	// Token is assigned from node frontend
	mes := &model.SkillMes{
		Id:      skillSet.ID,
		Session: token,
	}
	mesStr, _ := json.Marshal(mes)
	go api.PublishNSQMes("skill_created", []byte(mesStr))

	// Response
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(skillSet)
}
