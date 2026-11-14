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
if (results.some(result => result.performance < .90 || result.accessibility < .95 || result.cls >= .1)) throw new Error('Lighthouse performance, accessibility, or layout stability threshold failed')
