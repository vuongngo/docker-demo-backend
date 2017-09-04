package main

import (
	"net/http"
	"os"

	"app/api"
	"app/db"
	"app/repo"
	"app/worker"
	"app/ws"

	"github.com/gorilla/mux"
)

var (
	svc       *api.API
	websocket *ws.WebSocket
)

func init() {
	//ws hub
	hub := ws.NewHub()
	go hub.Run()

	websocket = &ws.WebSocket{
		Hub: hub,
	}

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

}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"WS",
		"GET",
		"/ws",
		websocket.ServeWs,
	},
	Route{
		"Index",
		"GET",
		"/",
		svc.Index,
	},
	Route{
		"CreateSkill",
		"POST",
		"/skill",
		svc.CreateSkillSet,
	},
	Route{
		"GetSkillSet",
		"GET",
		"/skill/{id}",
		svc.GetSkillSet,
	},
	Route{
		"job",
		"POST",
		"/job",
		svc.CreateJob,
	},
	Route{
		"job",
		"GET",
		"/job/{id}",
		svc.GetJob,
	},
}
