const DEBUG_LOGS = process.env.NEXT_PUBLIC_DEBUG_LOGS === 'true'

export const debugLog = (...args: unknown[]) => {
  if (DEBUG_LOGS) {
    console.log(...args)
  }
}

export const debugWarn = (...args: unknown[]) => {
  if (DEBUG_LOGS) {
    console.warn(...args)
  }
}

export const debugError = (...args: unknown[]) => {
  if (DEBUG_LOGS) {
    console.error(...args)
  }
}
