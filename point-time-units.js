let url = 'https://forecast-v2.metoceanapi.com/point/time'

let data = {
  points: [{lon: 174.7842, lat: -37.7935}],
  variables: ['air.temperature.at-2m'],
  time: {
    from: new Date().toISOString(),
    interval: '3h',
    repeat: 3,
  }
}

function multiplierOffsetToDegreesC(srcUnits) {
    switch (srcUnits) {
    case 'degreeK':
      return {multiplier: 1, offset: -273.15}
    case 'degreeC':
      return {multiplier: 1, offset: 0}
    case 'degreeF':
      return {multiplier: 0.555556, offset: -17.777778}
    }
  return {error: 'Can\'t convert data to degreesC from ' + srcUnits}
}

/*
This demonstrates how to detect the units a variable's data are in,
and one way to convert them.
Some conversion information can also be retrieved from the
https://forecast-v2.metoceanapi.com/units/ endpoint, which returns JSON
with the same multiplier/offset format used here
*/
export let pointTimeUnits = {
  url,
  data,
  cb: function(data) {
    let temp = data['variables']['air.temperature.at-2m']
    const transform = multiplierOffsetToDegreesC(temp['units'])
    if (transform.error) {
      return transform.error
    }

    return {'air.temperature.at-2m': {
      data: temp['data'].map((datum, index) => {
        if (temp['noData'][index]) {
          return null
        }
        return transform.multiplier * datum + transform.offset
      }),
      'noData': temp['noData'],
      'units': 'degreeC',
    }}
  }
}
