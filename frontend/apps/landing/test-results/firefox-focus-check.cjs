const { firefox } = require('playwright')

;(async () => {
  const browser = await firefox.launch()
  try {
    for (const theme of ['light', 'dark']) {
      const page = await browser.newPage({ viewport: { width: 1920, height: 1080 }, colorScheme: theme, reducedMotion: 'reduce' })
      await page.goto('http://localhost:3004/ar')
      await page.evaluate(theme => localStorage.setItem('servio-marketing-theme', theme), theme)
      await page.reload()
      for (const width of [320, 360, 375, 390, 430, 600, 768, 820, 1024, 1180, 1200, 1280, 1440, 1920]) {
        await page.setViewportSize({ width, height: width < 600 ? 844 : width === 1024 || width === 1180 ? 768 : width < 1200 ? 1024 : width === 1920 ? 1080 : 900 })
        await page.goto('http://localhost:3004/ar')
        if (width >= 1200) await page.locator('.sv-platform-trigger').click()
        else await page.locator('.sv-menu-button').click()
        await page.keyboard.press('Escape')
      }
      await page.goto('http://localhost:3004/ar')
      await page.addStyleTag({ content: '.cv-hero :focus-visible,:root.light body .sv-footer :focus-visible{outline-color:#f49680}' })
      await page.keyboard.press('Tab')
      const target = page.locator('.cv-hero .sv-button').first()
      await target.evaluate(element => element.scrollIntoView({ block: 'center', inline: 'center' }))
      await target.focus()
      await page.evaluate(() => new Promise(resolve => requestAnimationFrame(() => requestAnimationFrame(resolve))))
      const result = await target.evaluate(element => {
        const rect = element.getBoundingClientRect()
        const x = rect.x + rect.width / 2
        const y = rect.y + rect.height / 2
        const hit = document.elementFromPoint(x, y)
        return { rect: rect.toJSON(), x, y, scrollY, focused: document.activeElement === element, focusVisible: element.matches(':focus-visible'), hit: hit?.outerHTML.slice(0, 500), contains: !!hit && element.contains(hit) }
      })
      console.log(theme, JSON.stringify(result))
      await page.close()
    }
  } finally {
    await browser.close()
  }
})().catch(error => { console.error(error); process.exitCode = 1 })
