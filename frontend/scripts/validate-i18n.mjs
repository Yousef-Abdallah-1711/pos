import fs from 'node:fs'
import path from 'node:path'
import process from 'node:process'

const root = path.resolve(import.meta.dirname, '..')
const dictionaries = []

function collect(directory) {
  if (!fs.existsSync(directory)) return
  for (const entry of fs.readdirSync(directory, { withFileTypes: true })) {
    const full = path.join(directory, entry.name)
    if (entry.isDirectory()) collect(full)
    else if (entry.name.endsWith('.json') && !entry.name.endsWith('.bak-mojibake')) dictionaries.push(full)
  }
}

for (const app of ['dashboard', 'restaurant-website', 'landing']) {
  collect(path.join(root, 'apps', app, 'src', 'i18n', 'messages'))
}
collect(path.join(root, 'apps', 'restaurant-website', 'i18n', 'messages'))
collect(path.join(root, 'packages', 'i18n', 'locales'))

function flatten(value, prefix = '') {
  if (Array.isArray(value)) {
    const entries = [[`${prefix}.__length`, { type: 'array-length', shape: String(value.length) }]]
    for (const [index, child] of value.entries()) {
      entries.push(...flatten(child, `${prefix}[${index}]`))
    }
    return new Map(entries)
  }
  if (value && typeof value === 'object') {
    return new Map(Object.entries(value).flatMap(([key, child]) =>
      [...flatten(child, prefix ? `${prefix}.${key}` : key)]))
  }
  return new Map([[prefix, { type: typeof value, value }]])
}

const parsed = new Map()
let malformed = false
for (const file of dictionaries) {
  try {
    parsed.set(file, JSON.parse(fs.readFileSync(file, 'utf8')))
  } catch (error) {
    malformed = true
    console.error(`Malformed JSON: ${path.relative(root, file)}: ${error.message}`)
  }
}

const parity = []
for (const file of dictionaries) {
  const localeIndex = file.lastIndexOf(`${path.sep}en${path.sep}`)
  const peer = path.basename(file) === 'en.json'
    ? path.join(path.dirname(file), 'ar.json')
    : localeIndex >= 0
      ? `${file.slice(0, localeIndex)}${path.sep}ar${file.slice(localeIndex + 3)}`
      : undefined
  if (!peer) continue
  if (!parsed.has(peer)) continue
  const english = flatten(parsed.get(file))
  const arabic = flatten(parsed.get(peer))
  for (const key of new Set([...english.keys(), ...arabic.keys()])) {
    const left = english.get(key)
    const right = arabic.get(key)
    if (!left || !right || left.type !== right.type || left.shape !== right.shape) {
      parity.push(`${path.relative(root, file)} ↔ ${path.relative(root, peer)}: ${key} has a missing key or different value shape`)
      continue
    }
    if (left.type === 'string') {
      const placeholders = message => [...message.matchAll(/\{([a-zA-Z][a-zA-Z0-9_]*)\}/g)].map(match => match[1]).sort().join(',')
      if (placeholders(left.value) !== placeholders(right.value)) parity.push(`${path.relative(root, file)}: interpolation differs at ${key}`)
      const pluralCategories = message => [...message.matchAll(/(?:^|\s)(zero|one|two|few|many|other)\s*\{/g)].map(match => match[1]).sort().join(',')
      if (pluralCategories(left.value) !== pluralCategories(right.value)) parity.push(`${path.relative(root, file)}: plural forms differ at ${key}`)
    }
  }
}

console.log(`Validated ${parsed.size} JSON translation dictionaries in ${dictionaries.length} files.`)
const usageRoots = ['dashboard', 'restaurant-website', 'landing'].map(app => path.join(root, 'apps', app, 'src'))
const sourceFiles = []
function collectSource(directory) {
  if (!fs.existsSync(directory)) return
  for (const entry of fs.readdirSync(directory, { withFileTypes: true })) {
    const full = path.join(directory, entry.name)
    if (entry.isDirectory()) collectSource(full)
    else if (/\.(?:ts|tsx|js|jsx)$/.test(entry.name)) sourceFiles.push(full)
  }
}
for (const directory of usageRoots) collectSource(directory)
const usage = sourceFiles.map(file => fs.readFileSync(file, 'utf8')).join('\n')
const keyReferences = new Set([...usage.matchAll(/(?:t|translate|getTranslation)\(\s*['"]([^'"]+)['"]/g)].map(match => match[1]))
const coveredPrefixes = new Set()
for (const match of usage.matchAll(/\b(?:t|c|copy|text|messages)\.([A-Za-z_$][\w$]*(?:\.[A-Za-z_$][\w$]*)*)/g)) {
  const path = match[1]
  keyReferences.add(path)
  for (let index = path.lastIndexOf('.'); index > 0; index = path.lastIndexOf('.', index - 1)) {
    coveredPrefixes.add(path.slice(0, index))
  }
}
const unused = []
for (const [file, dictionary] of parsed) {
  const keys = flatten(dictionary)
  const isFeatureContent = file.endsWith(`${path.sep}features.json`) && file.includes(`${path.sep}landing${path.sep}`)
  const isSharedNamespace = file.includes(`${path.sep}packages${path.sep}i18n${path.sep}`)
  for (const key of keys.keys()) {
    const leaf = key.split('.').at(-1)
    const accessedThroughParent = [...coveredPrefixes].some(prefix => key.startsWith(`${prefix}.`))
    if (!isFeatureContent && !isSharedNamespace && leaf && !keyReferences.has(key) && !keyReferences.has(leaf) && !accessedThroughParent) {
      unused.push(`${path.relative(root, file)}:${key}`)
    }
  }
}
const reportPath = process.env.I18N_REPORT_PATH || path.join('/tmp', 'servio-i18n-unused-keys.report.json')
fs.writeFileSync(reportPath, `${JSON.stringify({ generatedAt: new Date().toISOString(), mode: 'heuristic source scan; review before deletion', unused }, null, 2)}\n`)
console.log(`Heuristic unused-key report: ${unused.length} candidate entries at ${reportPath} (dynamic calls may cause false positives; report only, no keys deleted).`)
if (parity.length) {
  console.error(`EN/AR key, value-shape, or plural parity errors (${parity.length}):`)
  for (const issue of parity) console.error(`- ${issue}`)
}
if (malformed || parity.length) process.exitCode = 1
