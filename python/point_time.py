from requests import post
from os import environ
from datetime import datetime, timezone
import json


# NOTE: don't for get to set "apikey" env, or the default below.
resp = post(
    "https://forecast-v2.metoceanapi.com/point/time",
    headers={"x-api-key": environ["apikey"]},
    json={
        "points": [{
            "lon": 174.7842,
            "lat": -37.7935
        }],
        "variables": [
            "cloud.cover"
        ],
        "time": {
            "from": "{:%Y-%m-%dT00:00:00Z}".format(datetime.now(timezone.utc)),
            "interval": "3h",
            "repeat": 2
        }
    }
)

if resp.status_code != 200:
    raise ValueError("{}: {}", resp.status_code, resp.text)

print(json.dumps(resp.json(), indent=1))
