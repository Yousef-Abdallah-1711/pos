import { expect, test } from '@playwright/test'
import { routeContainerApiToCandidate } from './api-routing'

const api = process.env.WS05_API_URL || 'http://localhost:18082/api/v1'
const password = process.env.WS05_TEST_PASSWORD || 'WS05-Disposable-Only-2026!'
type ShiftRecord = { value: number; startTime: number; sources: string[] }
type LayoutShiftRecord = PerformanceEntry & { value: number; hadRecentInput: boolean; sources?: Array<{ node: Node | null }> }

test('measure cold authenticated dashboard CLS', async ({ page }) => {
  await routeContainerApiToCandidate(page)
  await page.addInitScript(() => {
    const shifts: Array<{ value: number; startTime: number; sources: string[] }> = []
    new PerformanceObserver((list) => {
      for (const entry of list.getEntries() as LayoutShiftRecord[]) {
        if (entry.hadRecentInput) continue
        shifts.push({ value: entry.value, startTime: entry.startTime, sources: (entry.sources || []).map((source) => {
          const node = source.node as HTMLElement | null
          return node ? `${node.tagName}.${String(node.className).slice(0, 100)}` : 'unknown'
        }) })
      }
    }).observe({ type: 'layout-shift', buffered: true })
    Object.defineProperty(window, '__clsShifts', { value: shifts })
  })
  const response = await page.request.post(`${api}/auth/login`, {
    data: { email: 'phase2a-owner-a@example.test', password },
    headers: { Origin: 'http://localhost:13112' },
  })
  expect(response.ok()).toBeTruthy()
  await page.goto('/en/dashboard', { waitUntil: 'domcontentloaded', timeout: 90_000 })
  await expect(page.locator('main h1')).toBeVisible({ timeout: 45_000 })
  await page.waitForTimeout(6000)
  const result = await page.evaluate(() => ({
    shifts: (window as Window & { __clsShifts?: ShiftRecord[] }).__clsShifts || [],
    sidebar: document.querySelector('aside')?.getBoundingClientRect().height,
  }))
  let windowStart = 0
  let previousStart = 0
  let windowValue = 0
  let cls = 0
  for (const shift of result.shifts) {
    if (shift.startTime - previousStart > 1000 || shift.startTime - windowStart > 5000) {
      windowStart = shift.startTime
      windowValue = 0
    }
    windowValue += shift.value
    previousStart = shift.startTime
    cls = Math.max(cls, windowValue)
  }
  console.log(`DASHBOARD_CLS_RESULT ${JSON.stringify({ ...result, cls })}`)
  expect(cls).toBeLessThan(0.1)
})
