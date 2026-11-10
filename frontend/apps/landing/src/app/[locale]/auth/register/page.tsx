'use client'

import { useCallback } from 'react'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { ChefHat } from 'lucide-react'
import { AsyncAvailabilityField, FormErrorSummary, FormField, PasswordField, SubmitButton, applyApiFieldErrors, createRegisterSchema, isSubdomainFormatValid, parseApiError, useDebouncedAvailability, fieldAccessibility } from '@pos-saas/form-foundation'
import api from '@/lib/api'
import { useAuthStore } from '@/stores/authStore'
import { getConfig } from '@/lib/config'
import { useLanguage } from '@/contexts/LanguageContext'

type RegisterForm = { name: string; email: string; password: string; confirmPassword: string; restaurantName: string; restaurantSlug: string; phone: string }

export default function RegisterPage() {
  const setAuth = useAuthStore((state) => state.setAuth)
  const { t, dir, language } = useLanguage()
  const text = t.auth.register
  const messages = { required: t.auth.validation.required, nameMin: t.auth.validation.nameMin, nameMax: t.auth.validation.nameMax, invalidEmail: t.auth.validation.invalidEmail, emailMax: t.auth.validation.emailMax, passwordMin: t.auth.validation.passwordMin, passwordMax: t.auth.validation.passwordMax, passwordMismatch: t.auth.validation.passwordsMismatch, invalidSubdomain: t.auth.validation.invalidSubdomain, reservedSubdomain: t.auth.validation.reservedSubdomain, invalidPhone: t.auth.validation.invalidPhone }
  const schema = createRegisterSchema(messages)
  const { register, handleSubmit, watch, setError, clearErrors, setFocus, formState: { errors, isSubmitting, isSubmitted } } = useForm<RegisterForm>({ resolver: zodResolver(schema), mode: 'onTouched', reValidateMode: 'onChange', defaultValues: { phone: '' } })
  const restaurantSlug = watch('restaurantSlug') ?? ''
  const checkAvailability = useCallback(async (slug: string, signal: AbortSignal) => {
    const response = await api.get('/auth/check-subdomain', { params: { slug }, signal, timeout: 8000 })
    if (response.data.available) return true
    return /reserved/i.test(String(response.data.message)) ? 'reserved' as const : false
  }, [])
  const availability = useDebouncedAvailability(restaurantSlug, isSubdomainFormatValid, checkAvailability, 500)
  const fallback = { validation: 'Please review the highlighted fields.', invalidCredentials: 'Email or password is incorrect.', rateLimited: 'Too many attempts. Wait a moment and try again.', unavailable: 'The service is temporarily unavailable. Please try again.', network: 'Unable to connect. Check your connection and try again.', unexpected: 'We could not create your account. Please try again.', conflict: 'That email or website address is already in use.' }
  const messagesByState = { invalid: messages.invalidSubdomain, checking: 'Checking availability…', available: text.subdomainAvailable, taken: text.subdomainUnavailable, reserved: 'This website address is reserved.', error: 'Availability could not be checked. Check your connection and try again.' }

  const onSubmit = async (data: RegisterForm) => {
    if (availability.state !== 'available' || availability.checkedValue !== data.restaurantSlug.trim().toLowerCase()) {
      setError('restaurantSlug', { type: 'validate', message: availability.state === 'error' ? messagesByState.error : text.subdomainRequired })
      setFocus('restaurantSlug')
      return
    }
    clearErrors()
    try {
      const response = await api.post('/auth/register', { name: data.name, email: data.email, password: data.password, restaurant_name: data.restaurantName, restaurant_slug: data.restaurantSlug, phone: data.phone || undefined })
      const { user, session_established: sessionEstablished } = response.data
      if (!sessionEstablished) {
        window.location.href = `/${language}/auth/login?created=1`
        return
      }
      setAuth(user, '')
      window.location.href = `${getConfig().dashboardUrl}/${language}/auth/callback`
    } catch (error) {
      const parsed = parseApiError(error, fallback)
      applyApiFieldErrors(parsed.fields, (field, message) => setError(field as keyof RegisterForm, { type: 'server', message }))
      const firstField = Object.keys(parsed.fields)[0]
      if (firstField) {
        const fieldName = firstField === 'restaurant_slug' ? 'restaurantSlug' : firstField === 'restaurant_name' ? 'restaurantName' : firstField
        setFocus(fieldName as keyof RegisterForm)
      } else setError('root.server', { type: parsed.code, message: parsed.message })
    }
  }

  const onInvalid = (invalid: Record<string, unknown>) => {
    const firstField = Object.keys(invalid)[0] as keyof RegisterForm | undefined
    if (firstField) setFocus(firstField)
  }
  const inputClass = 'h-11 w-full rounded-md border border-gray-300 bg-white px-3 text-gray-900 shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-orange-600 disabled:opacity-60 dark:border-gray-600 dark:bg-gray-800 dark:text-white'
  const common = { disabled: isSubmitting, className: inputClass }
  const fieldError = (key: keyof RegisterForm) => errors[key]?.message
  const { restaurantDomain } = getConfig()

  return <main dir={dir} className="flex min-h-screen items-center justify-center bg-gradient-to-br from-orange-50 to-red-50 px-4 py-8 dark:from-gray-900 dark:to-gray-800">
    <div className="w-full max-w-md min-w-0">
      <header className="mb-6 text-center"><Link href={`/${language}`} className="mb-3 inline-flex items-center gap-2.5"><span className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-orange-500 to-red-600"><ChefHat className="h-6 w-6 text-white" /></span><span className="text-2xl font-extrabold text-gray-900 dark:text-white">Restaurant<span className="text-orange-500">POS</span></span></Link><p className="text-gray-600 dark:text-gray-400">{text.title}</p></header>
      <section className="rounded-2xl border border-gray-100 bg-white p-5 shadow-xl dark:border-gray-700 dark:bg-gray-800 sm:p-8">
        <form noValidate onSubmit={handleSubmit(onSubmit, onInvalid)} className="space-y-4">
          {isSubmitted && Object.keys(errors).length > 0 && <FormErrorSummary title={fallback.validation} errors={Object.fromEntries(Object.entries(errors).filter(([key]) => key !== 'root').map(([key, value]) => [key, value?.message ?? 'Review this field']))} onFocusField={(key) => setFocus(key as keyof RegisterForm)} />}
          {errors.root?.server?.message && <p role="alert" className="rounded-md bg-red-50 p-3 text-sm text-red-800">{errors.root.server.message}</p>}
          <FormField id="name" label={text.fullName} required error={fieldError('name')}><input {...register('name')} {...fieldAccessibility('name', fieldError('name'))} autoComplete="name" maxLength={255} placeholder={text.fullNamePlaceholder} {...common} /></FormField>
          <FormField id="email" label={text.email} required error={fieldError('email')}><input {...register('email')} {...fieldAccessibility('email', fieldError('email'))} type="email" autoComplete="email" maxLength={255} placeholder={text.emailPlaceholder} {...common} /></FormField>
          <FormField id="restaurantName" label={text.restaurantName} required error={fieldError('restaurantName')}><input {...register('restaurantName')} {...fieldAccessibility('restaurantName', fieldError('restaurantName'))} autoComplete="organization" maxLength={255} placeholder={text.restaurantNamePlaceholder} {...common} /></FormField>
          <FormField id="restaurantSlug" label={text.websiteSubdomain} required description={`${text.websiteAccessible} ${restaurantSlug || 'your-subdomain'}.${restaurantDomain}`} error={fieldError('restaurantSlug')}>
            <input {...register('restaurantSlug')} {...fieldAccessibility('restaurantSlug', fieldError('restaurantSlug'))} autoCapitalize="none" autoCorrect="off" maxLength={63} placeholder={text.subdomainPlaceholder} {...common} />
            <AsyncAvailabilityField state={availability.state} messages={messagesByState} />
          </FormField>
          <FormField id="phone" label={text.phone} optionalLabel={t.auth.validation.optional} error={fieldError('phone')}><input {...register('phone')} {...fieldAccessibility('phone', fieldError('phone'))} type="tel" autoComplete="tel" placeholder={text.phonePlaceholder} {...common} /></FormField>
          <FormField id="password" label={text.password} required description="At least 6 characters; maximum 72 UTF-8 bytes." error={fieldError('password')}><PasswordField {...register('password')} {...fieldAccessibility('password', fieldError('password'))} autoComplete="new-password" placeholder={text.passwordPlaceholder} label={text.password} {...common} /></FormField>
          <FormField id="confirmPassword" label={text.confirmPassword} required error={fieldError('confirmPassword')}><PasswordField {...register('confirmPassword')} {...fieldAccessibility('confirmPassword', fieldError('confirmPassword'))} autoComplete="new-password" placeholder={text.confirmPasswordPlaceholder} label={text.confirmPassword} {...common} /></FormField>
          <SubmitButton isSubmitting={isSubmitting} loadingLabel={text.creatingAccount} className="w-full rounded-md bg-gradient-to-r from-orange-500 to-red-600 px-4 text-white hover:from-orange-600 hover:to-red-700">{text.createAccount}</SubmitButton>
        </form>
        <p className="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">{text.hasAccount} <Link href={`/${language}/auth/login`} className="font-medium text-orange-600 hover:underline dark:text-orange-400">{text.signIn}</Link></p>
      </section>
    </div>
  </main>
}
