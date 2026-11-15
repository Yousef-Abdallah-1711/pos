export const locales = ['ar', 'en'] as const

export type Locale = (typeof locales)[number]
export type Direction = 'rtl' | 'ltr'

export const defaultLocale: Locale = 'ar'

export function isLocale(value: string): value is Locale {
  return locales.includes(value as Locale)
}

export function directionFor(locale: Locale): Direction {
  return locale === 'ar' ? 'rtl' : 'ltr'
}

export function formatNumber(locale: Locale, value: number): string {
  return new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-EG').format(value)
}

export function formatCurrency(locale: Locale, value: number, currency = 'EGP'): string {
  return new Intl.NumberFormat(locale === 'ar' ? 'ar-EG' : 'en-EG', {
    style: 'currency', currency,
  }).format(value)
}

export function formatDate(locale: Locale, value: Date): string {
  return new Intl.DateTimeFormat(locale === 'ar' ? 'ar-EG' : 'en-EG', {
    dateStyle: 'medium',
  }).format(value)
}

type LeafPaths<T> = {
  [K in keyof T & string]: T[K] extends string
    ? K
    : T[K] extends Record<string, unknown>
      ? `${K}.${LeafPaths<T[K]>}`
      : never
}[keyof T & string]

type ValueAt<T, P extends string> = P extends `${infer Head}.${infer Tail}`
  ? Head extends keyof T ? ValueAt<T[Head], Tail> : never
  : P extends keyof T ? T[P] : never

type Placeholders<S extends string> = S extends `${string}{${infer Name}}${infer Rest}`
  ? Name | Placeholders<Rest>
  : never

type ValuesFor<T, P extends string> = ValueAt<T, P> extends string
  ? [Placeholders<ValueAt<T, P> & string>] extends [never]
    ? undefined
    : Record<Placeholders<ValueAt<T, P> & string>, string | number>
  : undefined

export function translate<
  const Dictionary extends Record<string, unknown>,
  const Key extends LeafPaths<Dictionary>,
>(
  dictionary: Dictionary,
  key: Key,
  ...args: ValuesFor<Dictionary, Key> extends undefined ? [] : [ValuesFor<Dictionary, Key>]
): string {
  const values = args[0]
  const value = key.split('.').reduce<unknown>((node, part) => {
    if (!node || typeof node !== 'object' || !Object.hasOwn(node, part)) return undefined
    return (node as Record<string, unknown>)[part]
  }, dictionary)
  if (typeof value !== 'string') throw new Error(`Missing translation: ${key}`)
  const interpolations = values as Record<string, string | number> | undefined
  return value.replace(/\{([a-zA-Z][a-zA-Z0-9_]*)\}/g, (token, name: string) => {
    const replacement = interpolations?.[name]
    if (replacement === undefined) throw new Error(`Missing interpolation value: ${key}.${name}`)
    return String(replacement)
  })
}

export function selectPlural<T extends Record<string, string>>(
  locale: Locale,
  forms: T,
  count: number,
): string {
  const category = new Intl.PluralRules(locale).select(count)
  const form = forms[category] ?? forms.other
  if (!form) throw new Error(`Missing plural form: ${category} or other`)
  return form.replace('{count}', formatNumber(locale, count))
}
