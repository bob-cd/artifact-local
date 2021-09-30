/* This file is part of artifact-local.
* Copyright 2018-2021 Rahul De
*
* Use of this source code is governed by an MIT-style
* license that can be found in the LICENSE file or at
* https://opensource.org/licenses/MIT.
*/

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

var DIR_NAME = "artifacts"

func Ping(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprint(w, "Ack")
}

func Receive(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	dir := filepath.Join(
		DIR_NAME,
		params.ByName("group"),
		params.ByName("name"),
		params.ByName("runId"),
	)

	os.MkdirAll(dir, os.ModePerm)

	artifact, err := os.Create(filepath.Join(dir, params.ByName("artifact")))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())

		return
	}
	defer artifact.Close()

	_, err = io.Copy(artifact, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())

		return
	}

	fmt.Fprint(w, "Ok")
}

func Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	artifact := filepath.Join(
		DIR_NAME,
		params.ByName("group"),
		params.ByName("name"),
		params.ByName("runId"),
		params.ByName("artifact"),
	)

	err := os.Remove(artifact)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())

		return
	}

	fmt.Fprint(w, "Ok")
}

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8001"
	}
	router := httprouter.New()

	router.GET("/ping", Ping)
	router.POST("/bob_artifact/:group/:name/:runId/:artifact", Receive)
	router.DELETE("/bob_artifact/:group/:name/:runId/:artifact", Delete)
	router.ServeFiles("/bob_artifact/*filepath", http.Dir(DIR_NAME))

	http.Handle("/", router)
	http.ListenAndServe(":"+port, nil)
}
