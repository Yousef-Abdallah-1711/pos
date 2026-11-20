import {
  integerEnvironment,
  protectedSessionCheck,
  requireDedicatedEnvironment,
  requireSessionPool,
} from './lib.js'

const TARGET_RPS = integerEnvironment('TARGET_RPS', 10, 1, 2000)
const MAX_VUS = integerEnvironment('MAX_VUS', 20, 1, 10000)
const PRE_ALLOCATED_VUS = integerEnvironment('PRE_ALLOCATED_VUS', Math.min(MAX_VUS, 20), 1, MAX_VUS)

requireDedicatedEnvironment(MAX_VUS, TARGET_RPS)
requireSessionPool(MAX_VUS)

export const options = {
  discardResponseBodies: true,
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)'],
  scenarios: {
    protectedSession: {
      executor: 'constant-arrival-rate',
      rate: TARGET_RPS,
      timeUnit: '1s',
      duration: __ENV.DURATION || '1m',
      preAllocatedVUs: PRE_ALLOCATED_VUS,
      maxVUs: MAX_VUS,
      gracefulStop: '10s',
    },
  },
}

export default function () {
  protectedSessionCheck()
}
