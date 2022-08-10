from requests import Session
from numpy import float64
from numpy.ma import masked_array
from os import getenv
from datetime import datetime


def _convert_units(data, src_units: str, dst_units: str, units_list: dict):
    if src_units == dst_units:
        return
    conversion = units_list[src_units]['conversions'][dst_units]
    if 'multiplier' in conversion:
        data *= conversion['multiplier']
    if 'offset' in conversion:
        data += conversion['offset']


def _get_json(resp):
    if resp.status_code != 200:
        raise ValueError("{}: {}", resp.status_code, resp.text)
    return resp.json()


def _main():
    with Session() as session:
        url = "https://forecast-v2.metoceanapi.com"
        auth_headers = {"x-api-key": getenv("apikey", "MYAPIKEY")}
        units_list = _get_json(session.get(f"{url}/units/", headers=auth_headers))
        forecast = _get_json(session.post(
            f"{url}/point/time",
            headers=auth_headers,
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
        ))

    var = forecast["variables"]["air.temperature.at-2m"]
    # You will usually want to make a masked arrays with the "noData" field. 0 == Good data, any other value indicates
    # the "data" value is null, look up the field "noDataReasons" to find out why.
    temp = masked_array(var["data"], mask=var["noData"], dtype=float64)
    print(temp)
    _convert_units(data=temp, src_units=var["units"], dst_units="degreeF", units_list=units_list)
    print(temp)


if __name__ == '__main__':
    _main()
