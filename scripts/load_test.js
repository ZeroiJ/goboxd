import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 10 },
    { duration: '20s', target: 50 },
    { duration: '10s', target: 100 },
    { duration: '10s', target: 0 },
  ],
};

const payload = JSON.stringify({
  language: 'py3',
  source: 'print("hello load test")',
  tests: [{ stdin: '', expected_stdout: 'hello load test' }]
});

const params = {
  headers: { 'Content-Type': 'application/json' },
};

export default function () {
  const url = __ENV.GOBXD_URL || 'http://127.0.0.1:8080/run';
  const res = http.post(url, payload, params);
  
  check(res, {
    'is status 200': (r) => r.status === 200,
    'is accepted': (r) => {
      try {
        return JSON.parse(r.body).status === 'accepted';
      } catch (e) {
        return false;
      }
    }
  });
  sleep(1);
}
