const { chromium } = require('playwright')

;(async () => {
  const browser = await chromium.launch({ executablePath: 'C:/Program Files/Google/Chrome/Application/chrome.exe' })
  try {
    const page = await browser.newPage()
    await page.goto('http://localhost:3004/en')
    await page.evaluate(() => localStorage.setItem('servio-marketing-theme', 'dark'))
    await page.reload()
    await page.addScriptTag({ path: require.resolve('axe-core/axe.min.js') })
    const violations = await page.evaluate(async () => (await window.axe.run(document, {
      runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag21aa', 'wcag22aa'] },
    })).violations.map(item => ({ id: item.id, nodes: item.nodes.map(node => node.target) })))
    console.log(JSON.stringify(violations))
  } finally {
    await browser.close()
  }
})().catch(error => { console.error(error); process.exitCode = 1 })
