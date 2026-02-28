/**
 * LOCAL/TEST-ONLY SCRIPT
 *
 * Defaults to localhost for local performance checks.
 * Override target safely with:
 *   BASE_URL=https://api.staging.example.com k6 run backend/tests/load/baseline.js
 *
 * k6 Load Testing - Baseline Scenario
 *
 * Baseline test with 10 concurrent users
 * Used to establish performance baseline
 *
 * Run: k6 run baseline.js
 */

import http from 'k6/http'
import { check, group, sleep } from 'k6'
import { Rate, Trend, Counter, Gauge } from 'k6/metrics'

// Custom metrics
const errorRate = new Rate('errors')
const httpDuration = new Trend('http_req_duration')
const httpRequests = new Counter('http_requests')
const concurrentUsers = new Gauge('concurrent_users')

// Configuration
const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/$/, '')
const DURATION = __ENV.DURATION || '5m'
const USERS = 10

// Test options
export const options = {
  stages: [
    { duration: '1m', target: USERS },     // Ramp up to 10 users
    { duration: '3m', target: USERS },     // Stay at 10 users
    { duration: '1m', target: 0 },         // Ramp down
  ],
  thresholds: {
    'http_req_duration': ['p(95)<500', 'p(99)<1000'], // 95th percentile < 500ms
    'errors': ['rate<0.1'],                            // Error rate < 10%
  },
}

export default function () {
  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${__ENV.AUTH_TOKEN || ''}`,
  }

  concurrentUsers.set(__ENV.CONCURRENT_USERS || 10)

  // Test 1: Homepage
  group('Homepage', () => {
    let response = http.get(`${BASE_URL}/`, { headers })
    check(response, {
      'status is 200': (r) => r.status === 200,
      'page loads': (r) => r.timings.duration < 2000,
    })
    httpDuration.add(response.timings.duration)
    httpRequests.add(1)
    errorRate.add(response.status !== 200)
  })

  sleep(1)

  // Test 2: Get Products (Static Assets Cache Hit)
  group('Products List', () => {
    let response = http.get(`${BASE_URL}/api/v1/products`, { headers })
    check(response, {
      'status is 200': (r) => r.status === 200,
      'response time < 500ms': (r) => r.timings.duration < 500,
      'has products': (r) => r.json().data?.length > 0,
    })
    httpDuration.add(response.timings.duration)
    httpRequests.add(1)
    errorRate.add(response.status !== 200)
  })

  sleep(2)

  // Test 3: Get Single Product
  group('Product Detail', () => {
    let response = http.get(`${BASE_URL}/api/v1/products/1`, { headers })
    check(response, {
      'status is 200': (r) => r.status === 200,
      'response time < 300ms': (r) => r.timings.duration < 300,
    })
    httpDuration.add(response.timings.duration)
    httpRequests.add(1)
    errorRate.add(response.status !== 200)
  })

  sleep(1)

  // Test 4: Search (Complex Query)
  group('Search', () => {
    let response = http.get(
      `${BASE_URL}/api/v1/search?q=product&limit=20`,
      { headers }
    )
    check(response, {
      'status is 200': (r) => r.status === 200,
      'response time < 1000ms': (r) => r.timings.duration < 1000,
    })
    httpDuration.add(response.timings.duration)
    httpRequests.add(1)
    errorRate.add(response.status !== 200)
  })

  sleep(2)
}

export function handleSummary(data) {
  return {
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
    './reports/baseline.json': JSON.stringify(data),
  }
}

/**
 * Text summary formatter
 */
function textSummary(data, options) {
  const indent = options.indent || ''
  let summary = '\n=== Load Test Results ===\n'

  // Metrics summary
  if (data.metrics) {
    summary += '\nMetrics:\n'

    if (data.metrics.http_req_duration) {
      const trend = data.metrics.http_req_duration.values
      summary += `${indent}HTTP Duration: avg=${trend.avg?.toFixed(2)}ms, p95=${trend['p(95)']?.toFixed(2)}ms, p99=${trend['p(99)']?.toFixed(2)}ms\n`
    }

    if (data.metrics.http_requests) {
      summary += `${indent}Total Requests: ${data.metrics.http_requests.value}\n`
    }

    if (data.metrics.errors) {
      const errorRate = (data.metrics.errors.value / (data.metrics.http_requests?.value || 1)) * 100
      summary += `${indent}Error Rate: ${errorRate.toFixed(2)}%\n`
    }
  }

  return summary
}
