import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 50 },
    { duration: '1m', target: 100 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

export default function() {
  let response = http.post('http://localhost:8081/api/v1/shorten', 
    JSON.stringify({ url: 'https://example.com/test' }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  
  check(response, {
    'status is 201': (r) => r.status === 201,
  });
  
  sleep(1);
}
