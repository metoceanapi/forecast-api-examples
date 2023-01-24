import { argv } from 'node:process'
import { inspect } from 'node:util'
import { makeOptions, fetcher } from './index.js'
import { pointTime } from './point-time.js'
import { pointTimeWindVectors } from './point-time-wind-vectors.js'

function pretty(o) {
  return inspect(o, {showHidden: false, depth: null, colors: true})
}

function main() {
  let key = argv.pop()
  let target = pointTime
  let options = makeOptions(target.data, key)

  fetcher(target.url, options, function(data) {
    console.log('API response JSON:', pretty(data))
    if(target.cb === undefined) {
      return
    }
    let processed = target.cb(data)
    if (processed === undefined) {
      return
    }
    console.log('Processed:', pretty(processed))
  })
}

main()

// TODO show:
// vectors
// mask reasons
// base64 unpacking
// units
