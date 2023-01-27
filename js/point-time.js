let url = 'https://forecast-v2.metoceanapi.com/point/time'

// This is a very simple request, with no post-processing in this example
let data = {
  points: [{lon: 174.7842, lat: -37.7935}],
  variables: ['wave.height'],
  time: {
    from: new Date().toISOString(),
    interval: '3h',
    repeat: 3,
  }
}

export let pointTime = { url, data }
