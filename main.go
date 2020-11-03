/* This file is part of artifact-local.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
