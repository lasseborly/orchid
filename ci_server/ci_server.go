/*
The CI server of the Orchid application
*/

package ci_server

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/yaccio/orchid/core"
	"net/http"
)

/*
Runs the CI server.
This is a blocking operation
*/
func Start() {
	m := martini.Classic()

	m.Get("/logs/list", listLogs)
	m.Get("/logs/:id", getLog)
	m.Post("/runs/:id", runJob)

	m.Run()
}

/*
Endpoint for running the job with the given id
*/
func runJob(w http.ResponseWriter, r *http.Request, params martini.Params) {
	//TODO validate secret

	path := "ci"
	jobId := params["id"]

	var api core.Api
	var err error
	api, err = core.NewLocalApi(path)

	var log core.Log
	log, err = api.RunJob(jobId)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, err.Error())
		return
	}

	data, err := json.Marshal(log)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, string(data))
}

/*
Endpoint for getting the logged output of the log with the given id
*/
func getLog(w http.ResponseWriter, r *http.Request, params martini.Params) {
	//TODO validate secret
	//TODO this must be implemented with sockets (the current implementation does not work)

	path := "ci"
	logId := params["id"]

	var api core.Api
	var err error
	api, err = core.NewLocalApi(path)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	err = api.GetLogOutput(logId, w)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(200)
}

/*
Endpoint for getting all stored logs
*/
func listLogs(w http.ResponseWriter, r *http.Request) {
	//TODO validate secret

	path := "ci"
	api, err := core.NewLocalApi(path)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	logs, err := api.ListLogs()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	data, err := json.Marshal(logs)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, string(data))
}
