let url = 'https://forecast-v2.metoceanapi.com/point/time'

let data = {
  points: [{lon: 174.7842, lat: -37.7935}],
  variables: ['wind.speed.northward.at-10m', 'wind.speed.eastward.at-10m'],
  time: {
    from: new Date().toISOString(),
    interval: '3h',
    repeat: 3,
  }
}

/*
This demonstrates how to convert a north/east pair of wind components (sometimes referred to as u,v components)
into magnitude and direction.

Ocean current component pairs can be converted in the same way, but note that by convention,
winds come from a given direction and currents go towards that direction.
*/
export let pointTimeWindVectors = {
  url,
  data,
  cb: function(data) {
    let windSpeedNorth = data['variables']['wind.speed.northward.at-10m']
    let windSpeedEast = data['variables']['wind.speed.eastward.at-10m']

    let windSpeedScalar = [...windSpeedNorth['data']].map((north, index) => {
      if (windSpeedNorth['noData'][index] || windSpeedEast['noData'][index]) {
        return null
      }
      let east = windSpeedEast['data'][index]
      return Math.sqrt(Math.pow(north, 2) + Math.pow(east, 2))
    })

    let windDirectionDegrees = [...windSpeedNorth['data']].map((north, index) => {
      // here, you could also check the data.noDataReasons table to look up the corresponding reason if you wanted to
      if (windSpeedNorth['noData'][index] || windSpeedEast['noData'][index]) {
        return null
      }
      let east = windSpeedEast['data'][index]
      let dividend = Math.atan2(-north, -east) * (180 / Math.PI)
      return ((dividend % 360) + 360) % 360 // JavaScript % uses truncated division (like C), not floored division (like Python)
    })
    return {'wind.speed.at-10m': windSpeedScalar, 'wind.direction.at-10m': windDirectionDegrees}
  },
}
