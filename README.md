### Reference Bob Artifact Store

This is a simple artifact store enabling Bob to store build artifacts.

#### Requirements
- Python 3.5+
- [Poetry](https://poetry.eustace.io)

#### Running
- `poetry install` to install dependencies.
- `python3 artifact-local/server.py` will start the plugin on port 8001.

#### API
- `GET /bob_artifact/<path>`: Returns the artifact at the path if exists.
- `POST /bob_artifact/<path>`: Takes the file contents as `data`. Writes the file at the path.
- `GET /ping`: Responds with an `Ack`.
