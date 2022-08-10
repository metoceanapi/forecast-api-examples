from requests import post
from numpy import float64
from numpy.ma import masked_array, is_masked
from os import getenv
from datetime import datetime


# NOTE: don't for get to set "apikey" env, or the default below.
resp = post(
    "https://forecast-v2.metoceanapi.com/point/time",
    headers={"x-api-key": getenv("apikey", "MYAPIKEY")},
    json={
        "points": [{
          "lon": 174.7842,
          "lat": -37.7935
        }],
        "variables": [
            "precipitation.rate"
        ],
        "time": {
            "from": "{:%Y-%m-%dT00:00:00Z}".format(datetime.utcnow()),
            "interval": "3h",
            "repeat": 2
        }
    }
)

if resp.status_code != 200:
    raise ValueError("{}: {}", resp.status_code, resp.text)

data = resp.json()
var = data["variables"]["precipitation.rate"]

# You will usually want to make a masked arrays with the "noData" field. 0 == Good data, any other value indicates
# the "data" value is null, look up the field "noDataReason" to find out why.
precipitation_rate = masked_array(var["data"], mask=var["noData"], dtype=float64)

# https://en.wikipedia.org/wiki/Rain
# Light rain — when the precipitation rate is < 2.5 mm (0.098 in) per hour
# Moderate rain — when the precipitation rate is between 2.5 mm (0.098 in) - 7.6 mm (0.30 in) or 10 mm (0.39 in) per hour[106][107]
# Heavy rain — when the precipitation rate is > 7.6 mm (0.30 in) per hour,[106] or between 10 mm (0.39 in) and 50 mm (2.0 in) per hour[107]
# Violent rain — when the precipitation rate is > 50 mm (2.0 in) per hour[107]
chance_of_rain = []
for p_mm in precipitation_rate[:]:
    if is_masked(p_mm):
        chance_of_rain.append("")
    elif p_mm < 1:
        chance_of_rain.append("none")
    elif p_mm < 2.5:
        chance_of_rain.append("light")
    elif p_mm < 7.6:
        chance_of_rain.append("moderate")
    elif p_mm < 50:
        chance_of_rain.append("heavy")
    else:
        chance_of_rain.append("extreme")

print(chance_of_rain)
