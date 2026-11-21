export function installConsoleErrorSerialization() {
 const serializeConsoleValue = value => {
  if (value instanceof Error) {
   return {
    kind: 'error',
    name: String(value.name || 'Error'),
    message: String(value.message || ''),
    stack: typeof value.stack === 'string' ? value.stack : '',
    cause: value.cause instanceof Error ? String(value.cause.stack || value.cause.message) : value.cause == null ? null : String(value.cause),
   }
  }

  try {
   return { kind: 'value', value: typeof value === 'string' ? value : JSON.parse(JSON.stringify(value)) }
  } catch {
   return { kind: 'value', value: String(value) }
  }
 }
 const originalConsoleError = console.error.bind(console)
 console.error = (...values) => {
  const serializedValues = values.map(serializeConsoleValue)
  originalConsoleError(`SERVIO_CONSOLE_ERROR_DETAILS:${JSON.stringify(serializedValues)}`)
 }
}

export function parseSerializedConsoleError(text) {
 const prefix = 'SERVIO_CONSOLE_ERROR_DETAILS:'
 if (!text.startsWith(prefix)) return null
 try {
  const values = JSON.parse(text.slice(prefix.length))
  return Array.isArray(values) ? values : null
 } catch {
  return null
 }
}

export function parseSerializedRscFallback(values) {
 if (!Array.isArray(values) || values.length !== 2) return null
 const [message, error] = values
 if (message?.kind !== 'value' || typeof message.value !== 'string' || error?.kind !== 'error') return null
 const match = message.value.match(/^Failed to fetch RSC payload for (https?:\/\/\S+)\. Falling back to browser navigation\.$/)
 // Firefox reports an interrupted Next.js RSC stream as "Error in input stream".
 // Accept that exact browser error only when it is paired with Next.js's explicit
 // same-navigation fallback message; unrelated TypeErrors remain fatal.
 if (!match || error.name !== 'TypeError' || !/^(NetworkError when attempting to fetch resource\.|Failed to fetch|Error in input stream)$/.test(error.message || '')) return null
 return match[1]
}
