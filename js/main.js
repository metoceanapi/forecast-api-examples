import { argv } from 'node:process'
import { inspect } from 'node:util'
import { makeOptions, fetcher } from './index.js'
import { pointTime } from './point-time.js'
import { pointTimeWindVectors } from './point-time-wind-vectors.js'
import { pointTimeMasks } from './point-time-masks.js'
import { pointTimeBase64 } from './point-time-base64.js'

function pretty(o) {
  return inspect(o, {showHidden: false, depth: null, colors: true})
}

function main() {
  let examples = {
    pointTime,
    pointTimeWindVectors,
    pointTimeMasks,
    pointTimeBase64
  }
  let key = argv.pop()
  let exampleName = argv.pop()
  if (!key || !exampleName) {
    console.log('Arguments <example> <key> are required')
    return
  }
  let example = examples[exampleName]
  if (!example) {
    console.log('<example> must be one of:', Object.keys(examples).reduce((acc, elem) => acc + ', ' + elem))
    return
  }

  let options = makeOptions(example.data, key)
  fetcher(example.url, options, function(data) {
    console.log('API response JSON:', pretty(data))
    if (!example.cb) {
      return
    }
    let processed = example.cb(data)
    if (!processed) {
      return
    }
    console.log('Processed:', pretty(processed))
  })
}

main()

// TODO comments
// TODO selector for browser
// TODO show:
// units?
