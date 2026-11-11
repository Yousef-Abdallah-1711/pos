import { mkdir, readFile } from 'node:fs/promises'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { chromium } from '@playwright/test'

const appRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..')
const brandRoot = path.join(appRoot, 'public', 'brand')
const outputRoot = path.join(brandRoot, 'png')
const exports = [
  { source: 'poos-symbol.svg', sizes: [16, 32, 64] },
  { source: 'poos-app-icon.svg', sizes: [180, 512] },
  { source: 'poos-social-profile.svg', sizes: [400] },
  { source: 'poos-logo.svg', sizes: [704] },
  { source: 'poos-logo-dark.svg', sizes: [704] },
  { source: 'poos-logo-mono.svg', sizes: [704] },
  { source: 'poos-logo-vertical.svg', sizes: [320] },
  { source: 'poos-logo-vertical-dark.svg', sizes: [320] },
  { source: 'poos-logo-vertical-mono.svg', sizes: [320] },
  { source: 'poos-wordmark.svg', sizes: [480] },
]

await mkdir(outputRoot, { recursive: true })
const browser = await chromium.launch({ headless: true })
try {
  const page = await browser.newPage({ deviceScaleFactor: 1 })
  for (const { source, sizes } of exports) {
    const svg = await readFile(path.join(brandRoot, source), 'utf8')
    const viewBox = svg.match(/viewBox="0 0 (\d+(?:\.\d+)?) (\d+(?:\.\d+)?)"/)
    if (!viewBox) throw new Error(`${source} has no measurable viewBox`)
    const sourceWidth = Number(viewBox[1])
    const sourceHeight = Number(viewBox[2])
    const base = path.basename(source, '.svg')
    for (const outputWidth of sizes) {
      const outputHeight = Math.round(outputWidth * sourceHeight / sourceWidth)
      const renderableSvg = svg.replace('<svg ', `<svg width="${outputWidth}" height="${outputHeight}" `)
      await page.setContent(renderableSvg, { waitUntil: 'load' })
      await page.locator('svg').screenshot({
        path: path.join(outputRoot, `${base}-${outputWidth}.png`),
        animations: 'disabled',
        omitBackground: true,
        scale: 'css',
      })
    }
  }
} finally {
  await browser.close()
}
