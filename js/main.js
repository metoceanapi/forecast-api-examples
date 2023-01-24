import { argv } from 'node:process'
import { makeOptions, fetcher } from './index.js'
import { pointTime } from './point-time.js'

function main() {
  let key = argv.pop()
  let target = pointTime
  let options = makeOptions(target.data, key)

  fetcher(target.url, options, target.cb)
}

main()

// TODO show:
// vectors
// mask reasons
// base64 unpacking
// units
