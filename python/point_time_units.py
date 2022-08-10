from requests import post
from numpy import float64
from numpy.ma import masked_array
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
            "air.temperature.at-2m"
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
var = data["variables"]["air.temperature.at-2m"]

# You will usually want to make a masked arrays with the "noData" field. 0 == Good data, any other value indicates
# the "data" value is null, look up the field "noDataReasons" to find out why.
temperature = masked_array(var["data"], mask=var["noData"], dtype=float64)
units = var["units"]

# Convert Kelvin to Degree C.
# You have a couple of options for converting units. https://forecast-v2.metoceanapi.com/units/
# You can hand code it seen below.
# Or you can look up the units from the Forecast-API and look up the "conversions" field, Note: always apply multiplier
# before offset.
# degreeK:
#     ...
#     conversions:
#       degreeC:
#         offset: -273.15
#       degreeF:
#         offset: -459.67
#         multiplier: 1.8

print("{}: {}".format(units, temperature))

if units == "degreeK":
    temperature += -273.15

elif units == "degreeF":
    temperature *= 0.555556
    temperature += -17.777778

print("degreeC: {}".format(temperature))
