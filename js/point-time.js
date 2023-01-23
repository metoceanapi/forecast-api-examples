let url = 'https://forecast-v2.metoceanapi.com/point/time'

let data = {
  points: [{lon: 174.7842, lat: -37.7935}],
  variables: ['wave.height'],
  time: {
    from: '2023-01-16T00:00:00Z',
    interval: '3h',
    repeat: 3,
  }
}

export let pointTime = { url, data }
