import assert from 'node:assert/strict'
import test from 'node:test'
import { evaluateLighthouseThresholds, failedLighthouseLocales } from './lighthouse-thresholds.mjs'

test('accepts scores at the unchanged floors and CLS strictly below its maximum', () => {
 const rows = evaluateLighthouseThresholds([
  { locale: 'en', performance: .90, accessibility: .95, lcpMs: 900, cls: .099 },
 ])

 assert.deepEqual(failedLighthouseLocales(rows), [])
})

test('rejects scores below the floors and CLS at its exclusive maximum', () => {
 const rows = evaluateLighthouseThresholds([
  { locale: 'en', performance: .899, accessibility: .949, lcpMs: 900, cls: .1 },
 ])

 assert.deepEqual(failedLighthouseLocales(rows), [
  { locale: 'en', performance: .899, accessibility: .949, cls: .1 },
 ])
})

test('fails closed when any gated metric is missing or non-finite', () => {
 const rows = evaluateLighthouseThresholds([
  { locale: 'en', performance: undefined, accessibility: .99, cls: 0 },
  { locale: 'ar', performance: .99, accessibility: Number.NaN, cls: 0 },
  { locale: 'fr', performance: .99, accessibility: .99, cls: Number.POSITIVE_INFINITY },
 ])

 assert.deepEqual(rows.map(row => [row.performancePass, row.accessibilityPass, row.layoutStabilityPass]), [
  [false, true, true],
  [true, false, true],
  [true, true, false],
 ])
 assert.deepEqual(failedLighthouseLocales(rows).map(row => row.locale), ['en', 'ar', 'fr'])
})
