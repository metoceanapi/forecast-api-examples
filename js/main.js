import { argv } from 'node:process'
import { inspect } from 'node:util'
import { makeOptions, fetcher } from './index.js'
import { pointTime } from './point-time.js'

function main() {
  let key = argv.pop()
  let target = pointTime
  let options = makeOptions(target.data, key)

  fetcher(target.url, options, function(data) {
    console.log('API response JSON:', inspect(data, {showHidden: false, depth: null, colors: true}))
    if(target.cb !== undefined) {
      target.cb(data)
    }
  })
}

main()

// TODO show:
// vectors
// mask reasons
// base64 unpacking
// units
