/**
 * react-hook-form types a field error as `FieldError | Merge<FieldError, ...> | string | undefined`,
 * which React cannot render. This returns the human-readable message, or undefined.
 */
export function errorText(error: unknown): string | undefined {
  if (typeof error === 'string') return error
  if (error && typeof error === 'object' && 'message' in error) {
    const message = (error as { message?: unknown }).message
    return typeof message === 'string' ? message : undefined
  }
  return undefined
}
