from requests import post
from numpy import float64, pi
from numpy.ma import masked_array, mod, arctan2, sqrt, power
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
            "wind.speed.northward.at-10m",
            "wind.speed.eastward.at-10m"
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
wind_north = data["variables"]["wind.speed.northward.at-10m"]
wind_east = data["variables"]["wind.speed.eastward.at-10m"]

# You will usually want to make a masked arrays with the "noData" field. 0 == Good data, any other value indicates
# the "data" value is null, look up the field "noDataReasons" to find out why.
wind_v = masked_array(wind_north["data"], mask=wind_north["noData"], dtype=float64)
wind_u = masked_array(wind_east["data"], mask=wind_east["noData"], dtype=float64)

# To calculate speed and direction from these variables you can use the following formulas.
# It's worth noting that ocean current and wind vectors are calculated the same way, but winds are coming from that
# direction and ocean currents are going to that direction.
wind_speed = sqrt(power(wind_u, 2) + power(wind_v, 2))
wind_direction = mod(arctan2(-wind_u, -wind_v) * (180 / pi), 360)

print(wind_speed)
print(wind_direction)
