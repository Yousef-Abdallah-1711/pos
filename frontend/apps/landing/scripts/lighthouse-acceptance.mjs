import { createRequire } from 'node:module'
import { pathToFileURL } from 'node:url'
import fs from 'node:fs/promises'
import path from 'node:path'
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
try {
 for (const locale of ['en', 'ar']) {
  const result = await lighthouse(`${base}/${locale}`, { port: chrome.port, output: 'json', logLevel: 'error', onlyCategories: ['performance', 'accessibility', 'seo'] })
  if (!result) throw new Error(`No Lighthouse result for ${locale}`)
  await fs.writeFile(path.join(output, `${locale}.json`), result.report)
  const report = result.lhr
  results.push({ locale, performance: report.categories.performance.score, accessibility: report.categories.accessibility.score, seo: report.categories.seo.score, lcpMs: report.audits['largest-contentful-paint'].numericValue, cls: report.audits['cumulative-layout-shift'].numericValue })
 }
} finally {
 await chrome.kill()
 await fs.writeFile(path.join(output, 'results.json'), JSON.stringify({ note: 'Standard Lighthouse mobile simulation. Staging intentionally remains noindex.', results }, null, 2))
}
const thresholdResults = results.map(result => ({
 locale: result.locale,
 performance: result.performance,
 performancePass: result.performance >= .90,
 accessibility: result.accessibility,
 accessibilityPass: result.accessibility >= .95,
 lcpMs: result.lcpMs,
 cls: result.cls,
 layoutStabilityPass: result.cls < .1,
}))
console.log(JSON.stringify({ thresholds: { performance: .90, accessibility: .95, cls: .1 }, results: thresholdResults }, null, 2))
if (thresholdResults.some(result => !result.performancePass || !result.accessibilityPass || !result.layoutStabilityPass)) {
 throw new Error('Lighthouse performance, accessibility, or layout stability threshold failed; inspect the per-locale metric results above')
}
