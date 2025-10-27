import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '1m', target: 200 },
    { duration: '2m', target: 500 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'],
  },
};

export default function() {
  let shortCode = 'test' + __VU;
  
  let response = http.get(`http://localhost:8081/${shortCode}`);
  
  check(response, {
    'status is 302 or 200': (r) => [200, 302].includes(r.status),
  });
  
  sleep(0.5);
}
