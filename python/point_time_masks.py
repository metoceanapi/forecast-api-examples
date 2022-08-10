from requests import post
from numpy import float64
from numpy.ma import masked_array
from os import getenv
from datetime import datetime


# Below is an example of selecting two geographical points, one in the ocean -37.7935, 174.7842 and the other on
# # land -37.7734, 175.2158. The variable selected wave.height is an ocean variable. In the repsonse output a few
# things can be seen:
# * Every odd data value is null e.g. [1.099627, null, ...].
# * Every even noData value is 0 which mean the data is valid e.g. [0, 1, ...], but every odd has a noData value
# * The order of point.data in the response determind the order of the output of both data and noData.

# NOTE: don't for get to set "apikey" env, or the default below.
resp = post(
    "https://forecast-v2.metoceanapi.com/point/time",
    headers={"x-api-key": getenv("apikey", "MYAPIKEY")},
    json={
        "points": [{
          "lon": 174.7842,
          "lat": -37.7935
        }, {
          "lon": 175.2158,
          "lat": -37.7734
        }],
        "variables": [
            "wave.height"
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
wave_height = data["variables"]["wave.height"]

# You will usually want to make a masked arrays with the "noData" field. 0 == Good data, any other value indicates
# the "data" value is null, look up the field "noDataReason" to find out why.
wave_height_masked = masked_array(wave_height["data"], mask=wave_height["noData"], dtype=float64)
print(wave_height_masked)

# If you want to find the reason for noData you can checkout noDataReasons.
# Note: need to flip the key value.
no_data_reason = {v: k for k, v in data["noDataReasons"].items()}
print([no_data_reason[v] for v in wave_height["noData"]])

# the output looks like ['GOOD', 'MASK_LAND', 'GOOD', 'MASK_LAND', 'GOOD', 'MASK_LAND']
# MASK_LAND, or is near to land and we don't have a model with data that close to land.
# "GOOD": 0 will always equal zero, but no other reason's number is guaranteed.
