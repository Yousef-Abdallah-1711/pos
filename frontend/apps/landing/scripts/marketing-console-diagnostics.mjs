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
