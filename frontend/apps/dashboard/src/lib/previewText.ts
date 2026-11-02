/**
 * Editor previews show titles that can be stored either as a plain string or as a bilingual
 * object ({ en, ar }). This returns the English text for either shape, or the fallback.
 */
export function previewText(value: unknown, fallback = ''): string {
  if (typeof value === 'string') return value || fallback
  if (value && typeof value === 'object' && 'en' in value) {
    const en = (value as { en?: unknown }).en
    if (typeof en === 'string' && en) return en
  }
  return fallback
}
