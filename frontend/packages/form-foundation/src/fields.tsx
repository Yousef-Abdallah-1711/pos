import { forwardRef, useId, useState, type InputHTMLAttributes, type ReactNode } from 'react'
import { Eye, EyeOff, LoaderCircle, CircleCheck, CircleX } from 'lucide-react'

export interface FormFieldProps {
  id: string
  label: string
  description?: string
  error?: string
  required?: boolean
  optionalLabel?: string
  children: ReactNode
}

export function FormField({ id, label, description, error, required, optionalLabel = 'Optional', children }: FormFieldProps) {
  const descriptionId = `${id}-description`
  const errorId = `${id}-error`
  return (
    <div className="space-y-2">
      <label htmlFor={id} className="block text-sm font-medium text-gray-800 dark:text-gray-200">
        {label}{required && <span aria-hidden="true" className="ms-1 text-red-600">*</span>}
        {!required && <span className="ms-2 text-xs font-normal text-gray-500">({optionalLabel})</span>}
      </label>
      {children}
      {description && <p id={descriptionId} className="text-xs text-gray-600 dark:text-gray-400">{description}</p>}
      {error && <p id={errorId} className="text-sm text-red-700 dark:text-red-300" role="alert">{error}</p>}
    </div>
  )
}

export function fieldAccessibility(id: string, error?: string, description?: string) {
  return {
    id,
    'aria-invalid': Boolean(error) as boolean,
    'aria-describedby': [description ? `${id}-description` : '', error ? `${id}-error` : ''].filter(Boolean).join(' ') || undefined,
  }
}

export function FormErrorSummary({ errors, title = 'Please review these fields', onFocusField }: {
  errors: Record<string, string>
  title?: string
  onFocusField?: (field: string) => void
}) {
  const entries = Object.entries(errors)
  const summaryId = useId()
  if (!entries.length) return null
  return (
    <section id={summaryId} role="alert" aria-labelledby={`${summaryId}-title`} className="rounded-md border border-red-300 bg-red-50 p-4 text-red-900 dark:border-red-800 dark:bg-red-950 dark:text-red-100">
      <h2 id={`${summaryId}-title`} className="font-semibold">{title}</h2>
      <ul className="mt-2 list-inside list-disc space-y-1">
        {entries.map(([field, message]) => (
          <li key={field}><button type="button" className="underline underline-offset-2" onClick={() => onFocusField?.(field)}>{message}</button></li>
        ))}
      </ul>
    </section>
  )
}

export const PasswordField = forwardRef<HTMLInputElement, InputHTMLAttributes<HTMLInputElement> & {
  label: string
  visibleLabel?: string
  hiddenLabel?: string
}>(function PasswordField({ label, visibleLabel = 'Show password', hiddenLabel = 'Hide password', className = '', ...inputProps }, ref) {
  const [visible, setVisible] = useState(false)
  const id = inputProps.id
  return (
    <div className="relative">
      <input {...inputProps} ref={ref} id={id} type={visible ? 'text' : 'password'} className={`h-11 w-full rounded-md border border-gray-300 bg-white px-3 pe-12 text-gray-900 shadow-sm outline-none transition focus-visible:ring-2 focus-visible:ring-orange-600 disabled:cursor-not-allowed disabled:opacity-60 dark:border-gray-600 dark:bg-gray-800 dark:text-white ${className}`} />
      <button type="button" aria-label={`${visible ? hiddenLabel : visibleLabel}: ${label}`} aria-pressed={visible} onClick={() => setVisible((current) => !current)} className="absolute inset-y-0 end-1 inline-flex min-h-11 min-w-11 items-center justify-center rounded focus-visible:outline focus-visible:outline-2 focus-visible:outline-orange-600" tabIndex={inputProps.disabled ? -1 : 0}>
        {visible ? <EyeOff aria-hidden="true" className="h-5 w-5" /> : <Eye aria-hidden="true" className="h-5 w-5" />}
      </button>
    </div>
  )
})

export function AsyncAvailabilityField({ state, messages }: {
  state: 'idle' | 'invalid' | 'checking' | 'available' | 'taken' | 'reserved' | 'error'
  messages: Partial<Record<'invalid' | 'checking' | 'available' | 'taken' | 'reserved' | 'error', string>>
}) {
  if (state === 'idle') return null
  const content = messages[state]
  if (!content) return null
  const icon = state === 'checking' ? <LoaderCircle className="h-4 w-4 animate-spin" aria-hidden="true" />
    : state === 'available' ? <CircleCheck className="h-4 w-4" aria-hidden="true" />
      : state === 'invalid' || state === 'taken' || state === 'reserved' || state === 'error' ? <CircleX className="h-4 w-4" aria-hidden="true" /> : null
  const tone = state === 'available' ? 'text-green-700 dark:text-green-300' : state === 'checking' ? 'text-amber-700 dark:text-amber-300' : 'text-red-700 dark:text-red-300'
  return <p role={state === 'error' ? 'alert' : 'status'} aria-live={state === 'error' ? 'assertive' : 'polite'} className={`flex items-center gap-2 text-sm ${tone}`}>{icon}{content}</p>
}

export function SubmitButton({ children, loadingLabel, isSubmitting, disabled, className = '' }: {
  children: ReactNode
  loadingLabel: string
  isSubmitting: boolean
  disabled?: boolean
  className?: string
}) {
  return <button type="submit" disabled={disabled || isSubmitting} aria-busy={isSubmitting} className={`inline-flex min-h-11 items-center justify-center gap-2 disabled:cursor-not-allowed disabled:opacity-60 ${className}`}>{isSubmitting && <LoaderCircle className="h-4 w-4 animate-spin" aria-hidden="true" />}{isSubmitting ? loadingLabel : children}</button>
}
