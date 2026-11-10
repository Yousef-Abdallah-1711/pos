const SAFE_PROTOCOLS = new Set(['data:', 'blob:'])

export function canRenderDashboardProductImage(imageUrl?: string | null) {
  const normalized = imageUrl?.trim()
  if (!normalized) {
    return false
  }

  if (normalized.startsWith('/')) {
    return true
  }

  try {
    const parsed = new URL(normalized)
    if (SAFE_PROTOCOLS.has(parsed.protocol)) {
      return true
    }

    if (typeof window === 'undefined') {
      return false
    }

    return parsed.hostname === window.location.hostname
  } catch {
    return false
  }
}
