import { expect, test } from '@playwright/test'
import { fillRegistration, openRegistration, submitRegistration } from './fixtures'

const widths = [320, 375, 390, 768, 1440]

test('keeps registration responsive and labelled in English LTR and Arabic RTL', async ({ page }) => {
  for (const locale of ['en', 'ar']) {
    await openRegistration(page, locale)
    await expect(page.locator('main')).toHaveAttribute('dir', locale === 'ar' ? 'rtl' : 'ltr')
    for (const width of widths) {
      await page.setViewportSize({ width, height: 900 })
      const dimensions = await page.evaluate(() => {
        const viewport = document.documentElement.clientWidth
        const controls = Array.from(document.querySelectorAll('input, button'))
        return {
          viewport,
          page: document.documentElement.scrollWidth,
          clippedControls: controls.filter((control) => {
            const bounds = control.getBoundingClientRect()
            return bounds.left < 0 || bounds.right > viewport
          }).map((control) => control.id || control.textContent?.trim() || control.tagName),
        }
      })
      expect(dimensions.page, `${locale} page overflows at ${width}px`).toBeLessThanOrEqual(dimensions.viewport)
      expect(dimensions.clippedControls, `${locale} fields or buttons are clipped at ${width}px`).toEqual([])
      await expect(page.locator('input')).toHaveCount(7)
      const allInputsHaveLabels = await page.locator('input').evaluateAll((inputs) => inputs.every((input) => (input as HTMLInputElement).labels?.length === 1))
      expect(allInputsHaveLabels).toBe(true)
    }
  }
})

test('focuses and associates the first invalid field, supports keyboard focus, and labels password controls', async ({ page }) => {
  await openRegistration(page)
  const name = page.locator('#name')
  await name.focus()
  await page.keyboard.press('Tab')
  await expect(page.locator('#email')).toBeFocused()
  const passwordToggle = page.getByRole('button', { name: 'Show password: Password' })
  await page.locator('#password').focus()
  await page.keyboard.press('Tab')
  await expect(passwordToggle).toBeFocused()
  await page.keyboard.press('Enter')
  const hidePassword = page.getByRole('button', { name: 'Hide password: Password' })
  await expect(hidePassword).toHaveAttribute('aria-pressed', 'true')
  await expect(page.locator('#password')).toHaveAttribute('type', 'text')
  await submitRegistration(page)
  await expect(name).toBeFocused()
  await expect(name).toHaveAttribute('aria-invalid', 'true')
  await expect(name).toHaveAttribute('aria-describedby', /name-error/)
  await expect(page.getByRole('alert').first()).toBeVisible()
})

test('announces subdomain availability through a status region', async ({ page }) => {
  await openRegistration(page)
  await page.getByLabel(/website subdomain/i).fill(`a11y-${Date.now()}`)
  await expect(page.getByRole('status').filter({ hasText: /checking|available/i })).toBeVisible()
  await expect(page.getByRole('status').filter({ hasText: /available/i })).toBeVisible()
})
