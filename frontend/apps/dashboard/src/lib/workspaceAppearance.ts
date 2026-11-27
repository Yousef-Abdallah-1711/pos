export const WORKSPACE_APPEARANCE_VERSION = 1

const DEFAULT_COLORS = {
  light: { background: '#F8F2E8', surface: '#EFE4D5', foreground: '#2B201C', muted: '#5A4941', border: '#D8C8B6', primary: '#A13224' },
  dark: { background: '#211B18', surface: '#2D2520', foreground: '#F2E9DC', muted: '#D0C1B1', border: '#514237', primary: '#E17861' },
} as const

type Mode = keyof typeof DEFAULT_COLORS
type ColorKey = keyof typeof DEFAULT_COLORS.light
type Colors = Record<Mode, Record<ColorKey, string>>
type Typography = { english: 'manrope' | 'system'; arabic: 'ibm-plex-arabic' | 'system'; heading: 'fraunces' | 'manrope'; baseSize: number }
export type WorkspaceAppearance = { schemaVersion: 1; colors: Colors; typography: Typography }

function object(value: unknown, name: string): Record<string, unknown> {
  if (value === null || typeof value !== 'object' || Array.isArray(value)) throw new Error(`${name} must be an object`)
  return value as Record<string, unknown>
}

function keys(value: Record<string, unknown>, allowed: readonly string[]): void {
  if (Object.keys(value).some(key => !allowed.includes(key))) throw new Error('Unknown appearance key')
}

function luminance(color: string): number {
  const channels = [1, 3, 5].map(offset => parseInt(color.slice(offset, offset + 2), 16) / 255)
    .map(channel => channel <= 0.04045 ? channel / 12.92 : ((channel + 0.055) / 1.055) ** 2.4)
  return 0.2126 * channels[0] + 0.7152 * channels[1] + 0.0722 * channels[2]
}

function contrast(first: string, second: string): number {
  const values = [luminance(first), luminance(second)].sort((a, b) => a - b)
  return (values[1] + 0.05) / (values[0] + 0.05)
}

function resolveColors(input: unknown): Colors {
  const overrides = input === undefined ? {} : object(input, 'colors')
  keys(overrides, ['light', 'dark'])
  const colors: Colors = { light: { ...DEFAULT_COLORS.light }, dark: { ...DEFAULT_COLORS.dark } }
  for (const mode of ['light', 'dark'] as const) {
    const values = overrides[mode] === undefined ? {} : object(overrides[mode], mode)
    keys(values, Object.keys(colors[mode]))
    for (const [key, value] of Object.entries(values)) {
      // Accept opaque six-digit hexadecimal colors only; never inject CSS expressions.
      if (typeof value !== 'string' || !/^#[0-9a-fA-F]{6}$/.test(value)) throw new Error('Invalid color')
      colors[mode][key as ColorKey] = value.toUpperCase()
    }
    const c = colors[mode]
    if ([c.background, c.surface].some(surface => contrast(c.foreground, surface) < 4.5 || contrast(c.muted, surface) < 4.5)) throw new Error('Insufficient text contrast')
    if (contrast(c.primary, mode === 'light' ? '#FFFFFF' : '#211B18') < 4.5) throw new Error('Insufficient primary contrast')
  }
  return colors
}

function resolveTypography(input: unknown): Typography {
  const values = input === undefined ? {} : object(input, 'typography')
  keys(values, ['english', 'arabic', 'heading', 'baseSize'])
  const english = values.english ?? 'manrope', arabic = values.arabic ?? 'ibm-plex-arabic', heading = values.heading ?? 'fraunces', baseSize = values.baseSize ?? 16
  if ((english !== 'manrope' && english !== 'system') || (arabic !== 'ibm-plex-arabic' && arabic !== 'system') || (heading !== 'fraunces' && heading !== 'manrope')) throw new Error('Unsupported font')
  if (typeof baseSize !== 'number' || !Number.isInteger(baseSize) || baseSize < 14 || baseSize > 18) throw new Error('Invalid base size')
  return { english, arabic, heading, baseSize }
}

export function resolveWorkspaceAppearance(input: unknown): WorkspaceAppearance {
  const values = object(input, 'appearance')
  keys(values, ['colors', 'typography'])
  return { schemaVersion: WORKSPACE_APPEARANCE_VERSION, colors: resolveColors(values.colors), typography: resolveTypography(values.typography) }
}

export function workspaceAppearanceVariables(appearance: WorkspaceAppearance, mode: Mode): Record<string, string> {
  const colors = appearance.colors[mode]
  return {
    '--workspace-background': colors.background,
    '--workspace-surface': colors.surface,
    '--workspace-foreground': colors.foreground,
    '--workspace-muted': colors.muted,
    '--workspace-border': colors.border,
    '--workspace-primary': colors.primary,
    '--workspace-on-primary': mode === 'light' ? '#FFFFFF' : '#211B18',
    '--workspace-base-size': `${appearance.typography.baseSize}px`,
    '--workspace-font-en': appearance.typography.english === 'manrope' ? 'var(--font-workspace-manrope), sans-serif' : 'system-ui, sans-serif',
    '--workspace-font-ar': appearance.typography.arabic === 'ibm-plex-arabic' ? 'var(--font-workspace-arabic), sans-serif' : 'system-ui, sans-serif',
    '--workspace-font-heading': appearance.typography.heading === 'fraunces' ? 'var(--font-workspace-fraunces), serif' : 'var(--font-workspace-manrope), sans-serif',
  }
}
