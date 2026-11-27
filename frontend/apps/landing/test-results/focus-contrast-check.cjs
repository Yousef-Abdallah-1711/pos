const { chromium } = require('playwright')

;(async () => {
  const browser = await chromium.launch({ executablePath: 'C:/Program Files/Google/Chrome/Application/chrome.exe' })
  try {
    for (const theme of ['light', 'dark']) {
      const page = await browser.newPage()
      await page.goto('http://localhost:3004/en')
      await page.evaluate(theme => { localStorage.setItem('servio-marketing-theme', theme) }, theme)
      await page.reload()
      await page.addStyleTag({ content: '.cv-hero :focus-visible,:root.light body .sv-footer :focus-visible{outline-color:#f49680}' })
      for (const [name, selector] of [['hero', '.cv-hero .sv-button'], ['kitchen', '#kitchen .cv-text-link'], ['final', '.cv-close .sv-button'], ['footer', '.sv-footer nav a']]) {
        const target = page.locator(selector).first()
        await target.scrollIntoViewIfNeeded()
        await target.focus()
        const result = await target.evaluate(element => {
          const luminance = color => color.match(/[\d.]+/g).slice(0, 3).map(Number).map(value => value / 255).map(value => value <= .04045 ? value / 12.92 : ((value + .055) / 1.055) ** 2.4).reduce((sum, value, index) => sum + value * [.2126, .7152, .0722][index], 0)
          let parent = element.parentElement
          while (parent && ['rgba(0, 0, 0, 0)', 'transparent'].includes(getComputedStyle(parent).backgroundColor)) parent = parent.parentElement
          const style = getComputedStyle(element)
          const levels = [luminance(style.outlineColor), luminance(parent ? getComputedStyle(parent).backgroundColor : 'rgb(255, 255, 255)')].sort((a, b) => a - b)
          const rect = element.getBoundingClientRect()
          const visible = document.elementFromPoint(rect.x + rect.width / 2, rect.y + rect.height / 2)
          return { contrast: (levels[1] + .05) / (levels[0] + .05), width: Number.parseFloat(style.outlineWidth), focusVisible: element.matches(':focus-visible'), unobscured: !!visible && element.contains(visible), outline:style.outlineColor, parentBackground:parent && getComputedStyle(parent).backgroundColor, parentClass:parent && parent.className }
        })
        console.log(theme, name, JSON.stringify(result))
      }
      await page.close()
    }
  } finally {
    await browser.close()
  }
})().catch(error => { console.error(error); process.exitCode = 1 })
