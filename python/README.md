In this repo you will also find some examples of how to use the Forecast-API (this is a work in progress still).  

Python:  
* [point_time.py](python/point_time.py) shows you the response for a basic point time response
* [point_time_masks.py](python/point_time_masks.py) shows you how to deal with masks / no-data values and the reason for why data was NULL.
* [point_time_units.py](python/point_time_units.py) basic unit converstion for the variable's data.
* [point_time_wind_vectors.py](python/point_time_wind_vectors.py) how to convert vectors like wind and ocean currents into magnitude and direction.
* [point_time_rain.py](python/point_time_rain.py) derives the chances of rain to text from the variable `precipitation.rate`.

To run the examples you will need to set the ENV `apikey` e.g.
```
apikey=MYAPIKEY python3 python/point_time_wind_vectors.py
```

