'use strict'

const assert = require('node:assert/strict')
const fs = require('node:fs')
const path = require('node:path')
const { pathToFileURL } = require('node:url')
const { test } = require('node:test')
const { createRequire } = require('node:module')
const { hooks } = require('../.pnpmfile.cjs')

const root = path.resolve(__dirname, '..')

function collectSourceFiles(directory) {
  return fs.readdirSync(directory, { withFileTypes: true }).flatMap((entry) => {
    const entryPath = path.join(directory, entry.name)
    if (entry.isDirectory()) return collectSourceFiles(entryPath)
    return /\.(?:[cm]?[jt]sx?)$/.test(entry.name) ? [entryPath] : []
  })
}

test('removes only the extractor watcher from the pinned next-intl release', () => {
  const packageManifest = {
    name: 'next-intl',
    version: '4.15.1',
    dependencies: {
      '@parcel/watcher': '^2.4.1',
      'next-intl-swc-plugin-extractor': '4.15.1',
    },
  }

  const result = hooks.readPackage(packageManifest)

  assert.equal(result.dependencies['@parcel/watcher'], undefined)
  assert.equal(result.dependencies['next-intl-swc-plugin-extractor'], '4.15.1')
})

test('leaves other packages and future next-intl versions unchanged', () => {
  const manifests = [
    { name: 'some-other-package', version: '1.0.0', dependencies: { '@parcel/watcher': '^2.4.1' } },
    { name: 'next-intl', version: '4.15.2', dependencies: { '@parcel/watcher': '^2.4.1' } },
  ]

  for (const manifest of manifests) {
    assert.deepEqual(hooks.readPackage(manifest), manifest)
  }
})

test('dashboard and restaurant source do not use extraction APIs', () => {
  const sourceRoots = [
    path.join(root, 'apps/dashboard/src'),
    path.join(root, 'apps/restaurant-website'),
    path.join(root, 'packages'),
  ]
  const extractionApi = /\buseExtracted\b|\bgetExtracted\b|unstable_extractMessages|next-intl\/extractor/
  const matches = sourceRoots.flatMap((sourceRoot) =>
    collectSourceFiles(sourceRoot).flatMap((file) => {
      const contents = fs.readFileSync(file, 'utf8')
      return extractionApi.test(contents) ? [path.relative(root, file)] : []
    }),
  )

  assert.deepEqual(matches, [])
})

test('next-intl plugin loads without @parcel/watcher in dashboard resolution', async () => {
  const dashboardRequire = createRequire(path.join(root, 'apps/dashboard/package.json'))

  assert.throws(() => dashboardRequire.resolve('@parcel/watcher'), { code: 'MODULE_NOT_FOUND' })
  await import(pathToFileURL(dashboardRequire.resolve('next-intl/plugin')).href)
})
