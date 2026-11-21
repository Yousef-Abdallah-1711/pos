import { createRequire } from 'node:module'
import { pathToFileURL } from 'node:url'
import fs from 'node:fs/promises'
import path from 'node:path'
import { evaluateLighthouseThresholds, failedLighthouseLocales } from './lighthouse-thresholds.mjs'
const require = createRequire(import.meta.url)
const { chromium } = require('playwright')
const { default: lighthouse } = await import(pathToFileURL(require.resolve('lighthouse')).href)
const { launch } = await import(pathToFileURL(require.resolve('chrome-launcher')).href)
const base = process.env.SERVIO_MARKETING_URL
if (!base) throw new Error('Set the intended SERVIO_MARKETING_URL')
if (process.env.SERVIO_ACCEPTANCE_ENV === 'ci' && new URL(base).origin !== 'https://pos.servio.test') throw new Error('CI lighthouse acceptance accepts only the disposable landing host')
const output = path.join(process.env.SERVIO_BROWSER_EVIDENCE || '/tmp/servio-lighthouse', 'lighthouse')
await fs.mkdir(output, { recursive: true })
const chrome = await launch({ chromePath: chromium.executablePath(), chromeFlags: ['--headless', '--no-sandbox', '--disable-dev-shm-usage', ...(process.env.SERVIO_DISPOSABLE_BROWSER === '1' ? ['--ignore-certificate-errors'] : [])] })
const results = []

function summarizeMetric(audit) {
 return {
  score: audit?.score ?? null,
  value: Number.isFinite(audit?.numericValue) ? Math.round(audit.numericValue * 100) / 100 : null,
  unit: audit?.numericUnit ?? null,
 }
}

function summarizePerformance(report) {
 const metrics = report.audits
 const opportunities = report.categories.performance.auditRefs
  .filter(({ id, weight }) => weight > 0 && metrics[id]?.details?.type === 'opportunity')
  .map(({ id, weight }) => ({
   id,
   weight,
   score: metrics[id].score ?? null,
   savingsMs: summarizeMetric({ numericValue: metrics[id].details.overallSavingsMs }).value,
   savingsBytes: summarizeMetric({ numericValue: metrics[id].details.overallSavingsBytes }).value,
  }))
  .sort((left, right) => (right.savingsMs ?? 0) - (left.savingsMs ?? 0) || right.weight - left.weight)
  .slice(0, 5)
 return {
  metrics: {
   firstContentfulPaint: summarizeMetric(metrics['first-contentful-paint']),
   largestContentfulPaint: summarizeMetric(metrics['largest-contentful-paint']),
   totalBlockingTime: summarizeMetric(metrics['total-blocking-time']),
   speedIndex: summarizeMetric(metrics['speed-index']),
   serverResponseTime: summarizeMetric(metrics['server-response-time']),
   cumulativeLayoutShift: summarizeMetric(metrics['cumulative-layout-shift']),
  },
  leadingOpportunities: opportunities,
 }
}

let failure
try {
 for (const locale of ['en', 'ar']) {
  const result = await lighthouse(`${base}/${locale}`, { port: chrome.port, output: 'json', logLevel: 'error', onlyCategories: ['performance', 'accessibility', 'seo'] })
  if (!result) throw new Error(`No Lighthouse result for ${locale}`)
  await fs.writeFile(path.join(output, `${locale}.json`), result.report)
  const report = result.lhr
  results.push({ locale, performance: report.categories.performance.score, accessibility: report.categories.accessibility.score, seo: report.categories.seo.score, lcpMs: report.audits['largest-contentful-paint'].numericValue, cls: report.audits['cumulative-layout-shift'].numericValue, ...summarizePerformance(report) })
 }
} catch (error) {
 failure = { name: error?.name || 'Error', message: error?.message || String(error), stack: error?.stack }
 throw error
} finally {
 await chrome.kill()
 await fs.writeFile(path.join(output, 'results.json'), JSON.stringify({ note: 'Standard Lighthouse mobile simulation. Staging intentionally remains noindex.', thresholds: { performance: 0.90, accessibility: 0.95, clsExclusiveMaximum: 0.1 }, results, failure }, null, 2))
}
const thresholdResults = evaluateLighthouseThresholds(results)
const failedLocales = failedLighthouseLocales(thresholdResults)
console.info(JSON.stringify({
 lighthouseAcceptance: 'completed',
 thresholds: { performance: .90, accessibility: .95, cls: .1 },
 results: thresholdResults,
 diagnostics: results.map(({ locale, metrics, leadingOpportunities }) => ({ locale, metrics, leadingOpportunities })),
 failedLocales,
}, null, 2))
if (failedLocales.length > 0) throw new Error('Lighthouse acceptance thresholds failed; inspect per-locale metric diagnostics above and results.json')
