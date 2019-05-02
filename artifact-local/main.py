# This file is part of artifact-local.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

import os
import pathlib
import shutil

import aiofiles
from sanic import Sanic, response

app = Sanic("artifact local")
PORT = 8001
DIR_NAME = "artifacts"


@app.route("/bob_artifact/<key:[^/].*?>", methods=["GET", "POST", "DELETE"])
async def receive(request, key):
    if request.method == "POST":
        data = request.files.get("data")

        if data is None:
            return response.text("Invalid request", status=400)

        path = os.path.join(DIR_NAME, key)

        pathlib.Path(path).parent.mkdir(parents=True, exist_ok=True)

        async with aiofiles.open(path, "wb") as artifact:
            await artifact.write(data.body)

        return response.text("Ok")
    elif request.method == "GET":
        path = os.path.join(DIR_NAME, key)

        return (
            await response.file(path)
            if os.path.exists(path)
            else response.text("No such artifact", status=404)
        )
    else:
        path = os.path.join(DIR_NAME, key)

        if not os.path.exists(path):
            return response.text("No such artifact", status=404)
        elif os.path.isfile(path):
            os.unlink(path)
        else:
            shutil.rmtree(path)

        return response.text("Ok")


@app.route("/ping")
async def handle_ping(_):
    return response.text("Ack")


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=PORT, workers=os.cpu_count() + 1)
