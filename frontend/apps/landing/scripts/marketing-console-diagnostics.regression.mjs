import assert from 'node:assert/strict'
import { createRequire } from 'node:module'
import { installConsoleErrorSerialization, parseSerializedConsoleError } from './marketing-console-diagnostics.mjs'

const require = createRequire(import.meta.url)
const { firefox } = require('playwright')

async function main() {
 const browser = await firefox.launch({ headless: true })
 try {
  const page = await browser.newPage()
  const consoleErrors = []
  page.on('console', message => {
   if (message.type() === 'error') consoleErrors.push(message.text())
  })
  await page.addInitScript(installConsoleErrorSerialization)

  const testOrigin = 'https://console-diagnostics.invalid'
  await page.route(`${testOrigin}/**`, route => {
   const isInitial = new URL(route.request().url()).pathname === '/initial'
   return route.fulfill({
    status: 200,
    contentType: 'text/html',
    body: isInitial ? `<script>console.error(new Error("serialization survives navigation")); location.replace("${testOrigin}/after")</script>` : '<p>after navigation</p>',
   })
  })
  await page.goto(`${testOrigin}/initial`)
  await page.waitForURL(`${testOrigin}/after`, { timeout: 5000 })

  const serializedError = consoleErrors.map(parseSerializedConsoleError).find(Boolean)
  assert.ok(serializedError, 'Firefox console error should arrive with serialized details')
  assert.equal(serializedError[0]?.kind, 'error')
  assert.equal(serializedError[0]?.name, 'Error')
  assert.equal(serializedError[0]?.message, 'serialization survives navigation')
 } finally {
  await browser.close()
 }
}

main().catch(error => {
 console.error(error)
 process.exitCode = 1
})
