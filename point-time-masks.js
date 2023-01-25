let url = 'https://forecast-v2.metoceanapi.com/point/time'

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
