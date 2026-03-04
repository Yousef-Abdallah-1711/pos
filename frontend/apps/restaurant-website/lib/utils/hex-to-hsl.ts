const hexToHslCache = new Map<string, string>()

function expandShortHex(hex: string): string {
  if (!/^#[0-9a-fA-F]{3}$/.test(hex)) {
    return hex
  }

  const chars = hex.slice(1).split('')
  return `#${chars.map((char) => `${char}${char}`).join('')}`
}

function normalizeHex(hex: string | null | undefined): string {
  if (typeof hex !== 'string') {
    return '#000000'
  }

  const trimmed = hex.trim()
  if (!trimmed) {
    return '#000000'
  }

  const withHash = trimmed.startsWith('#') ? trimmed : `#${trimmed}`
  return expandShortHex(withHash).toLowerCase()
}

export function hexToHSL(hex: string | null | undefined): string {
  const normalized = normalizeHex(hex)
  const match = /^#([0-9a-f]{6})$/.exec(normalized)
  if (!match) {
    return '0 0% 0%'
  }

  const r = parseInt(match[1].slice(0, 2), 16) / 255
  const g = parseInt(match[1].slice(2, 4), 16) / 255
  const b = parseInt(match[1].slice(4, 6), 16) / 255

  const max = Math.max(r, g, b)
  const min = Math.min(r, g, b)
  const delta = max - min

  let h = 0
  const l = (max + min) / 2
  let s = 0

  if (delta !== 0) {
    s = delta / (1 - Math.abs(2 * l - 1))
    switch (max) {
      case r:
        h = ((g - b) / delta) % 6
        break
      case g:
        h = (b - r) / delta + 2
        break
      default:
        h = (r - g) / delta + 4
        break
    }
  }

  const hue = Math.round((h * 60 + 360) % 360)
  const saturation = Math.round(s * 100)
  const lightness = Math.round(l * 100)
  return `${hue} ${saturation}% ${lightness}%`
}

export function memoHexToHSL(hex: string | null | undefined): string {
  const normalized = normalizeHex(hex)
  const cached = hexToHslCache.get(normalized)
  if (cached) {
    return cached
  }

  const value = hexToHSL(normalized)
  hexToHslCache.set(normalized, value)
  return value
}
