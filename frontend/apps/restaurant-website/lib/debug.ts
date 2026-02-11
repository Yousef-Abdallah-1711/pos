const DEBUG = process.env.NEXT_PUBLIC_DEBUG_LOGS === 'true' || process.env.DEBUG_LOGS === 'true'

export function debugLog(...args: unknown[]) {
  if (DEBUG) console.log(...args)
}

export function debugWarn(...args: unknown[]) {
  if (DEBUG) console.warn(...args)
}

export function debugError(...args: unknown[]) {
  if (DEBUG) console.error(...args)
}
