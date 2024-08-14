/* This file is part of artifact-local.
* Copyright 2018- Rahul De
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
)

const DIR_NAME = "artifacts"

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ack")
}

func receive(w http.ResponseWriter, r *http.Request) {
	dir := filepath.Join(
		DIR_NAME,
		r.PathValue("group"),
		r.PathValue("name"),
		r.PathValue("runId"),
	)

	os.MkdirAll(dir, os.ModePerm)

	artifact, err := os.Create(filepath.Join(dir, r.PathValue("artifact")))
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

func delete(w http.ResponseWriter, r *http.Request) {
	artifact := filepath.Join(
		DIR_NAME,
		r.PathValue("group"),
		r.PathValue("name"),
		r.PathValue("runId"),
		r.PathValue("artifact"),
	)

	err := os.Remove(artifact)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())

		return
	}

	fmt.Fprint(w, "Ok")
}

func send(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(DIR_NAME, r.PathValue("artifactPath")))
}

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8001"
	}
	mux := http.NewServeMux()
	path := "/bob_artifact/{group}/{name}/{runId}/{artifact}"

	mux.HandleFunc("GET /ping", ping)
	mux.HandleFunc("POST "+path, receive)
	mux.HandleFunc("DELETE "+path, delete)
	mux.HandleFunc("GET /bob_artifact/{artifactPath...}", send)

	http.ListenAndServe(":"+port, mux)
}
