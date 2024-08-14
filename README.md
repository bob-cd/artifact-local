### Reference Bob Artifact Store

This is a simple artifact store enabling Bob to store build artifacts.

#### Requirements
- [Go](https://golang.org/dl/) 1.22+

#### Running
- `go build main.go` to compile the code and obtain a binary `main`.
- `./main` will start on port `8001` by default, set the env var `PORT` to change.

#### API

Here `{path}` represents `{pipeline-group}/{pipeline-name}/{run-id}/{artifact-name}`.

- `GET /bob_artifact/{path}`: Returns the artifact at the path if exists.
- `POST /bob_artifact/{path}`: Takes the file contents as the post body. Writes the file at the path.
- `DELETE /bob_artifact/{path}`: Deletes the file at the path.
- `GET /ping`: Responds with an `Ack`.
