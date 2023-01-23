export function makeOptions(data, key) {
  return {
    method: 'post',
    body: JSON.stringify(data),
    headers: {
      'Content-Type': 'application/json',
      'x-api-key': key,
    }
  }
}

export async function fetcher(url, options, cb) {
  await fetch(url, options)
    .then(response => {
      console.log('API response status:', response.status);
      return response.json();
    }).then(json => {
      console.log('API response JSON:', json);
      if (cb !== undefined) {
        cb(json)
      }
    });
}
