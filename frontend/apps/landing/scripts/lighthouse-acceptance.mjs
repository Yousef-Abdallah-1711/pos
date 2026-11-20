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
const runsPerLocale = 5
function median(values) {
 const sorted = [...values].sort((left, right) => left - right)
 if (!sorted.length || sorted.some(value => !Number.isFinite(value))) throw new Error('Lighthouse returned an invalid metric sample')
 const middle = Math.floor(sorted.length / 2)
 return sorted.length % 2 ? sorted[middle] : (sorted[middle - 1] + sorted[middle]) / 2
}
await fs.mkdir(output, { recursive: true })
const chrome = await launch({ chromePath: chromium.executablePath(), chromeFlags: ['--headless', '--no-sandbox', '--disable-dev-shm-usage', ...(process.env.SERVIO_DISPOSABLE_BROWSER === '1' ? ['--ignore-certificate-errors'] : [])] })
const results = []
try {
 for (const locale of ['en', 'ar']) {
  const samples = []
  for (let run = 1; run <= runsPerLocale; run++) {
   const result = await lighthouse(`${base}/${locale}`, { port: chrome.port, output: 'json', logLevel: 'error', onlyCategories: ['performance', 'accessibility', 'seo'] })
   if (!result) throw new Error(`No Lighthouse result for ${locale}, run ${run}`)
   await fs.writeFile(path.join(output, `${locale}-${String(run).padStart(2, '0')}.json`), result.report)
   const report = result.lhr
   samples.push({
    performance: report.categories.performance.score,
    accessibility: report.categories.accessibility.score,
    seo: report.categories.seo.score,
    lcpMs: report.audits['largest-contentful-paint'].numericValue,
    cls: report.audits['cumulative-layout-shift'].numericValue,
    warnings: report.runWarnings || [],
   })
  }
  const performanceScores = samples.map(sample => sample.performance)
  const accessibilityScores = samples.map(sample => sample.accessibility)
  const layoutShifts = samples.map(sample => sample.cls)
  results.push({
   locale,
   sampleCount: samples.length,
   performance: median(performanceScores),
   performanceMinimum: Math.min(...performanceScores),
   performanceScores,
   accessibility: Math.min(...accessibilityScores),
   seo: median(samples.map(sample => sample.seo)),
   lcpMs: median(samples.map(sample => sample.lcpMs)),
   cls: Math.max(...layoutShifts),
   samples,
  })
 }
} finally {
 await chrome.kill()
 await fs.writeFile(path.join(output, 'results.json'), JSON.stringify({ note: 'Standard Lighthouse mobile simulation. Staging intentionally remains noindex.', results }, null, 2))
}
if (results.some(result => result.performance < .90 || result.accessibility < .95 || result.cls >= .1)) {
 throw new Error('Lighthouse median performance, per-run accessibility, or worst-run layout stability threshold failed')
}
