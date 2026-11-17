import { publicIteration, publicRequestRateCeiling, requireDedicatedEnvironment, integerEnvironment, SESSION_COOKIE_VALUES, requireSessionPool } from './lib.js'

const MAX_VUS = integerEnvironment('MAX_VUS', 10, 1, 10000)
requireDedicatedEnvironment(MAX_VUS, publicRequestRateCeiling(MAX_VUS))
if (SESSION_COOKIE_VALUES.length > 0) requireSessionPool(MAX_VUS)

export const options = {
  discardResponseBodies: true,
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)'],
  scenarios: {
    baseline: {
      executor: 'ramping-vus',
      startVUs: 1,
      stages: [
        { duration: __ENV.RAMP_DURATION || '30s', target: MAX_VUS },
        { duration: __ENV.HOLD_DURATION || '1m', target: MAX_VUS },
        { duration: __ENV.RAMP_DOWN_DURATION || '30s', target: 0 },
      ],
      gracefulRampDown: '10s',
    },
  },
}

export default function () {
  publicIteration()
}
