const EGYPTIAN_POUND = 'EGP'

export function formatArabicNumber(num: number): string {
  return num
    .toString()
    .replace(/\d/g, (digit) => '٠١٢٣٤٥٦٧٨٩'[Number(digit)])
}

export function formatPrice(amount: number, locale: string): string {
  const isArabic = locale.startsWith('ar')
  const normalized = Number.isFinite(amount) ? amount : 0

  const formatter = new Intl.NumberFormat(isArabic ? 'ar-EG' : 'en-EG', {
    minimumFractionDigits: normalized % 1 === 0 ? 0 : 2,
    maximumFractionDigits: 2,
  })

  if (isArabic) {
    return `${formatter.format(normalized)} ج.م`
  }

  return `${formatter.format(normalized)} ${EGYPTIAN_POUND}`
}

export function formatDate(date: string, locale: string): string {
  const parsed = new Date(date)
  if (Number.isNaN(parsed.getTime())) {
    return date
  }

  const isArabic = locale.startsWith('ar')
  return new Intl.DateTimeFormat(isArabic ? 'ar-EG' : 'en-EG', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  }).format(parsed)
}

export function formatRating(rating: number, count: number, locale: string): string {
  const safeRating = Number.isFinite(rating) ? rating : 0
  const safeCount = Number.isFinite(count) ? count : 0
  const isArabic = locale.startsWith('ar')

  if (isArabic) {
    return `${formatArabicNumber(Number(safeRating.toFixed(1)))}⭐ (${formatArabicNumber(Math.round(safeCount))} تقييم)`
  }

  return `${safeRating.toFixed(1)}⭐ (${Math.round(safeCount)} reviews)`
}
