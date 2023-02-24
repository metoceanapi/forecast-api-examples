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

const bytesPerUint32 = 4

function unpack(variable, reasonsByCode) {
  let b64data = variable['data']
  let bytes = Uint8Array.from(atob(b64data), c => c.charCodeAt(0))
  let view = new DataView(bytes.buffer)
  // If you know you are running on a little-endian platform,
  // you can simply use Uint32Array and Float32Array on bytes.buffer
  return Array.from({length: bytes.length / bytesPerUint32}).map((_, index) => {
    let offset = index * bytesPerUint32
    let integer = view.getUint32(offset, true) // true indicates little endian
    let reason = reasonsByCode.get(integer)
    return reason ? reason : view.getFloat32(offset, true)
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
