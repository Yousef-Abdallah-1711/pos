import { test } from 'node:test'
import assert from 'node:assert/strict'
import { resolveWorkspaceAppearance, workspaceAppearanceVariables } from './workspaceAppearance.ts'

test('inherits Coral defaults without mutating them across tenant resolutions', () => {
  const first = resolveWorkspaceAppearance({ colors: { light: { primary: '#722116' } } })
  const second = resolveWorkspaceAppearance({})
  assert.equal(first.colors.light.primary, '#722116')
  assert.equal(second.colors.light.primary, '#A13224')
  assert.equal(second.colors.dark.background, '#211B18')
})

test('rejects unknown tokens, unsafe colors, unsupported fonts and typography bounds', () => {
  for (const value of [
    { colors: { light: { surprise: '#ffffff' } } },
    { colors: { light: { primary: 'url(https://example.com)' } } },
    { typography: { english: 'https://example.com/font.woff2' } },
    { typography: { baseSize: 500 } },
    { extra: true },
  ]) assert.throws(() => resolveWorkspaceAppearance(value))
})

test('rejects inaccessible resolved text combinations', () => {
  assert.throws(() => resolveWorkspaceAppearance({ colors: { light: { foreground: '#F8F2E8' } } }))
  assert.throws(() => resolveWorkspaceAppearance({ colors: { light: { background: '#3B3B3B', surface: '#3B3B3B', foreground: '#FFFFFF', muted: '#FFFFFF', primary: '#3B3B3B' } } }))
})

test('maps the selected mode to scoped semantic variables', () => {
  const style = workspaceAppearanceVariables(resolveWorkspaceAppearance({}), 'dark')
  assert.equal(style['--workspace-background'], '#211B18')
  assert.equal(style['--workspace-primary'], '#E17861')
  assert.equal(style['--workspace-base-size'], '16px')
})
