import { argv } from 'node:process'
import { makeOptions, fetcher } from './index.js'
import { pointTime } from './point-time.js'

function main() {
  let key = argv.pop()
  let options = makeOptions(pointTime.data, key)

  fetcher(pointTime.url, options)
}

main();
