import type { ThemeConfig, ThemeLocale, ThemeSectionConfig } from '@/types/theme'
import { DEFAULT_THEME } from '@/lib/utils/default-theme'

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function toRecord(value: unknown): Record<string, unknown> {
  return isRecord(value) ? value : {}
}

function toStringValue(value: unknown, fallback: string): string {
  if (typeof value === 'string' && value.trim().length > 0) {
    return value
  }
  return fallback
}

function toBooleanValue(value: unknown, fallback: boolean): boolean {
  if (typeof value === 'boolean') {
    return value
  }
  return fallback
}

function readLocalized(value: unknown, locale: ThemeLocale, fallback = ''): string {
  if (typeof value === 'string') {
    return value
  }
  if (!isRecord(value)) {
    return fallback
  }

  const localValue = value[locale]
  if (typeof localValue === 'string' && localValue.trim().length > 0) {
    return localValue
  }

  const englishValue = value.en
  if (typeof englishValue === 'string' && englishValue.trim().length > 0) {
    return englishValue
  }

  return fallback
}

function normalizeComponents(value: unknown, locale: ThemeLocale): ThemeSectionConfig[] {
  if (!Array.isArray(value)) {
    return DEFAULT_THEME.components
  }

  return value
    .map((item, index) => {
      const section = toRecord(item)
      return {
        id: toStringValue(section.id, `section-${index + 1}`),
        type: toStringValue(section.type, 'custom') as ThemeSectionConfig['type'],
        enabled: toBooleanValue(section.enabled, true),
        order:
          typeof section.order === 'number'
            ? section.order
            : typeof section.displayOrder === 'number'
              ? section.displayOrder
              : typeof section.display_order === 'number'
                ? section.display_order
                : index + 1,
        title: readLocalized(section.title, locale, ''),
        settings: toRecord(section.settings ?? section.config),
        componentOverrides: toRecord(
          section.componentOverrides ??
          section.component_overrides ??
          section.variants
        ),
      }
    })
    .sort((a, b) => a.order - b.order)
}

export function normalizeThemeConfig(apiData: unknown, locale: ThemeLocale): ThemeConfig {
  const raw = toRecord(apiData)
  const data = toRecord(raw.data ?? raw)

  const colors = toRecord(data.colors)
  const typography = toRecord(data.typography)
  const identity = toRecord(data.identity)
  const header = toRecord(data.header ?? data.header_config)
  const footer = toRecord(data.footer ?? data.footer_config)

  return {
    ...DEFAULT_THEME,
    id: (data.id ?? DEFAULT_THEME.id) as string | number | undefined,
    name: toStringValue(data.name, DEFAULT_THEME.name),
    slug: toStringValue(data.slug, DEFAULT_THEME.slug),
    templateId: toStringValue(data.template_id ?? data.templateId, DEFAULT_THEME.templateId ?? 'template-01'),
    headerVariant: toStringValue(data.header_variant ?? data.headerVariant, DEFAULT_THEME.headerVariant ?? 'solid') as ThemeConfig['headerVariant'],
    footerVariant: toStringValue(data.footer_variant ?? data.footerVariant, DEFAULT_THEME.footerVariant ?? 'expanded') as ThemeConfig['footerVariant'],
    colors: {
      ...DEFAULT_THEME.colors,
      ...Object.fromEntries(Object.entries(colors).map(([key, value]) => [key, String(value)])),
    },
    typography: {
      ...DEFAULT_THEME.typography,
      ...Object.fromEntries(Object.entries(typography).map(([key, value]) => [key, value as string | number])),
    },
    identity: {
      ...DEFAULT_THEME.identity,
      ...identity,
    },
    header: {
      ...DEFAULT_THEME.header,
      ...header,
      logoText: readLocalized(header.logoText ?? header.logo_text, locale, readLocalized(identity.siteTitle ?? identity.site_title, locale, DEFAULT_THEME.name)),
    },
    footer: {
      ...DEFAULT_THEME.footer,
      ...footer,
      companyName: readLocalized(footer.companyName ?? footer.company_name, locale, readLocalized(identity.siteTitle ?? identity.site_title, locale, DEFAULT_THEME.name)),
    },
    components: normalizeComponents(data.components, locale),
  }
}
