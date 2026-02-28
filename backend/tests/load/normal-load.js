/**
 * LOCAL/TEST-ONLY SCRIPT
 *
 * Defaults to localhost for local/staging load checks.
 * Override target safely with:
 *   BASE_URL=https://api.staging.example.com k6 run backend/tests/load/normal-load.js
 *
 * k6 Load Testing - Normal Load Scenario
 *
 * Normal load test with 100 concurrent users
 * Represents expected daily traffic
 *
 * Run: k6 run normal-load.js
 */

import http from 'k6/http'
import { check, group, sleep } from 'k6'
import { Rate, Trend, Counter } from 'k6/metrics'

const errorRate = new Rate('errors')
const httpDuration = new Trend('http_req_duration')
const httpRequests = new Counter('http_requests')

const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/$/, '')
const USERS = 100

export const options = {
  stages: [
    { duration: '2m', target: USERS },      // Ramp up
    { duration: '5m', target: USERS },      // Sustain
    { duration: '2m', target: USERS / 2 },  // Ramp down
    { duration: '1m', target: 0 },          // Final ramp down
  ],
  thresholds: {
    'http_req_duration': ['p(95)<800', 'p(99)<2000'],
    'errors': ['rate<0.05'],
  },
}

export default function () {
  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${__ENV.AUTH_TOKEN || ''}`,
  }

  // Simulate realistic user flow
  group('Dashboard Navigation', () => {
    // Dashboard home
    let res = http.get(`${BASE_URL}/dashboard`, { headers })
    check(res, {
      'dashboard loads': (r) => r.status === 200,
      'fast load': (r) => r.timings.duration < 2000,
    })
    httpDuration.add(res.timings.duration)
    httpRequests.add(1)
    errorRate.add(res.status !== 200)
  })

  sleep(1)

  // Product browse
  group('Product Browse', () => {
    let res = http.get(
      `${BASE_URL}/api/v1/products?page=1&limit=20&sort=trending`,
      { headers }
    )
    check(res, {
      'get products': (r) => r.status === 200,
      'response time': (r) => r.timings.duration < 1000,
    })
    httpDuration.add(res.timings.duration)
    httpRequests.add(1)
    errorRate.add(res.status !== 200)
  })

  sleep(2)

  // Search products
  group('Search Products', () => {
    const queries = ['electronics', 'laptop', 'phone', 'tablet', 'monitor']
    const query = queries[Math.floor(Math.random() * queries.length)]

    let res = http.get(
      `${BASE_URL}/api/v1/search?q=${query}&limit=30`,
      { headers }
    )
    check(res, {
      'search works': (r) => r.status === 200,
      'search fast': (r) => r.timings.duration < 1500,
    })
    httpDuration.add(res.timings.duration)
    httpRequests.add(1)
    errorRate.add(res.status !== 200)
  })

  sleep(2)

  // Get categories
  group('Categories', () => {
    let res = http.get(`${BASE_URL}/api/v1/categories`, { headers })
    check(res, {
      'get categories': (r) => r.status === 200,
      'cached': (r) => r.timings.duration < 100,
    })
    httpDuration.add(res.timings.duration)
    httpRequests.add(1)
    errorRate.add(res.status !== 200)
  })

  sleep(1)

  // View product detail
  group('Product Detail', () => {
    const productIds = [1, 2, 3, 4, 5, 10, 15, 20]
    const id = productIds[Math.floor(Math.random() * productIds.length)]

    let res = http.get(`${BASE_URL}/api/v1/products/${id}`, { headers })
    check(res, {
      'product detail': (r) => r.status === 200,
      'detail fast': (r) => r.timings.duration < 500,
    })
    httpDuration.add(res.timings.duration)
    httpRequests.add(1)
    errorRate.add(res.status !== 200)
  })

  sleep(3)

  // Analytics
  group('Analytics', () => {
    let res = http.get(
      `${BASE_URL}/api/v1/analytics?period=week`,
      { headers }
    )
    check(res, {
      'analytics loads': (r) => r.status === 200,
      'analytics time': (r) => r.timings.duration < 2000,
    })
    httpDuration.add(res.timings.duration)
    httpRequests.add(1)
    errorRate.add(res.status !== 200)
  })

  sleep(2)
}
