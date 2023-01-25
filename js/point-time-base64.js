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
  },
  outputFormat: 'base64',
}

function unpack(variable, reasonsByCode) {
  let b64data = variable['data']
  let bytes = Uint8Array.from(atob(b64data), c => c.charCodeAt(0))
  // TODO endianness -- typed arrays have platform endianness
  let integers = new Uint32Array(bytes.buffer) // TODO check signedness
  let floats = new Float32Array(bytes.buffer)

  return Array.from(integers).map((n, index) => {
    let reason = reasonsByCode.get(n)
    return reason ? reason : floats[index]
  })
}

export let pointTimeBase64 = {
  url,
  data,
  cb: function(data) {
    let reasonsByBytes = new Map(Object.entries(data['noDataMask']).map(pair => pair.reverse()))
    let waveHeight = data['variables']['wave.height']
    return unpack(waveHeight, reasonsByBytes)
  },
}
