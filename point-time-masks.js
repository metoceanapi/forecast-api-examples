let url = 'https://forecast-v2.metoceanapi.com/point/time'

/*
The first point, (174.7842, -37.7935) is in the ocean, and the second point, , is on land (175.2158, -37.7734).
In this example, we request the wave.height variable, which is not available on land.
You will see that the API returns a null value in the wave.height data array for the dry point at each timestep.
It also returns non-zero codes in the corresponding positions in the wave.height noData array.

This example transforms the data array by replacing the nulls with the reason (which should be MASK_LAND) they are null,
according to the API response.
*/
let data = {
  points: [
    {lon: 174.7842, lat: -37.7935},
    {lon: 175.2158, lat: -37.7734},
  ],
  variables: ['wave.height'],
  time: {
    from: new Date().toISOString(),
    interval: '3h',
    repeat: 3,
  }
}

export let pointTimeMasks = {
  url,
  data,
  cb: function(data) {
    let reasonsByCode = new Map(Object.entries(data['noDataReasons']).map(pair => pair.reverse()))
    let waveHeight = data['variables']['wave.height']
    return waveHeight['data'].map((datum, index) => {
      let mask = waveHeight['noData'][index]
      return mask ? reasonsByCode.get(mask) : datum
    })
  },
}
