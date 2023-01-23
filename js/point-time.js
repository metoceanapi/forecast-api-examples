let url = 'https://forecast-v2.metoceanapi.com/point/time';

let data = {
  points: [{lon: 174.7842, lat: -37.7935}],
  variables: ['wave.height'],
  time: {
    from: '2023-01-16T00:00:00Z',
    interval: '3h',
    repeat: 3,
  }
};

let options = {
  method: 'post',
  body: JSON.stringify(data),
  headers: {
    'Content-Type': 'application/json',
    'x-api-key': ''
  }
};

await fetch(url, options)
    .then(response => {
      console.log('API response status:', response.status);
      return response.json();
    }).then(json => {
      console.log('API response JSON:', json);
      document.getElementById('response').innerHTML = JSON.stringify(json);
    });
