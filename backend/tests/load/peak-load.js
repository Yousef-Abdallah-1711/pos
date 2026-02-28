/**
 * LOCAL/TEST-ONLY SCRIPT
 *
 * Defaults to localhost for local/staging stress checks.
 * Override target safely with:
 *   BASE_URL=https://api.staging.example.com k6 run backend/tests/load/peak-load.js
 *
 * k6 Load Testing - Peak Load Scenario
 *
 * Peak/stress test with 1000-10000 concurrent users
 * Simulates peak traffic conditions
 *
 * Run: k6 run peak-load.js --vus 1000 --duration 10m
 */

import http from 'k6/http'
import { check, group, sleep } from 'k6'
import { Rate, Trend, Counter, Histogram } from 'k6/metrics'

const errorRate = new Rate('errors')
const httpDuration = new Trend('http_req_duration')
const httpRequests = new Counter('http_requests')
const responseTime = new Histogram('response_times')

const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/$/, '')
const VUS = parseInt(__ENV.VUS || '1000')

export const options = {
  stages: [
    { duration: '2m', target: VUS / 2 },    // Ramp to 50%
    { duration: '3m', target: VUS },        // Ramp to 100%
    { duration: '5m', target: VUS },        // Sustain peak
    { duration: '2m', target: VUS / 2 },    // Ramp down
    { duration: '1m', target: 0 },          // Final ramp
  ],
  thresholds: {
    'http_req_duration': ['p(95)<2000', 'p(99)<5000'],
    'errors': ['rate<0.1'],  // Allow 10% errors under stress
  },
}

export default function () {
  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${__ENV.AUTH_TOKEN || ''}`,
  }

  // Lightweight requests (high volume)
  group('Lightweight Requests', () => {
    // Static asset check
    let res = http.get(`${BASE_URL}/static/app.js`, { headers })
    check(res, {
      'static served': (r) => r.status === 200 || r.status === 304,
    })
    httpDuration.add(res.timings.duration)
    responseTime.add(res.timings.duration)
    httpRequests.add(1)
    errorRate.add(res.status >= 500)
  })

  sleep(0.5)

  // API list endpoints
  group('List API', () => {
    let res = http.get(`${BASE_URL}/api/v1/products?limit=10`, { headers })
    check(res, {
      'list endpoint': (r) => r.status === 200,
      'fast list': (r) => r.timings.duration < 1000,
    })
    httpDuration.add(res.timings.duration)
    responseTime.add(res.timings.duration)
    httpRequests.add(1)
    errorRate.add(res.status !== 200)
  })

  sleep(0.5)

  // Heavy requests (lower frequency)
  if (Math.random() < 0.3) {
    group('Heavy Requests', () => {
      let res = http.get(
        `${BASE_URL}/api/v1/analytics?period=month&include=details`,
        { headers }
      )
      check(res, {
        'heavy request': (r) => r.status === 200,
        'acceptable time': (r) => r.timings.duration < 3000,
      })
      httpDuration.add(res.timings.duration)
      responseTime.add(res.timings.duration)
      httpRequests.add(1)
      errorRate.add(res.status !== 200)
    })
  }

  sleep(1)

  // Search requests
  if (Math.random() < 0.2) {
    group('Search', () => {
      const queries = ['product', 'electronics', 'laptop', 'phone']
      const query = queries[Math.floor(Math.random() * queries.length)]

      let res = http.get(
        `${BASE_URL}/api/v1/search?q=${query}&limit=50`,
        { headers }
      )
      check(res, {
        'search': (r) => r.status === 200,
        'search speed': (r) => r.timings.duration < 2000,
      })
      httpDuration.add(res.timings.duration)
      responseTime.add(res.timings.duration)
      httpRequests.add(1)
      errorRate.add(res.status !== 200)
    })
  }

  sleep(Math.random() * 2)
}
