'use client'

import { useCallback, useEffect, useState } from 'react'
import { trackFunnel } from '@/lib/funnel'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { AsyncAvailabilityField, FormErrorSummary, FormField, PasswordField, SubmitButton, applyApiFieldErrors, createRegisterSchema, isSubdomainFormatValid, parseApiError, useDebouncedAvailability, fieldAccessibility } from '@pos-saas/form-foundation'
import api from '@/lib/api'
import { useAuthStore } from '@/stores/authStore'
import { getConfig } from '@/lib/config'
import { useLanguage } from '@/contexts/LanguageContext'
import { Navbar } from '@/components/Navbar'

type RegisterForm = { name: string; email: string; password: string; confirmPassword: string; restaurantName: string; restaurantSlug: string; phone: string }

export default function RegisterPage() {
  const [billing, setBilling] = useState<'monthly' | 'annual'>('monthly')
  useEffect(() => {
    const selected = new URLSearchParams(window.location.search).get('billing') === 'annual' ? 'annual' : 'monthly'
    setBilling(selected)
    trackFunnel('trial_registration_started', { locale: window.location.pathname.startsWith('/ar') ? 'ar' : 'en', source: 'registration', billing: selected })
  }, [])
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
  const fallback = language === 'ar'
    ? { validation: 'يرجى مراجعة الحقول الموضحة.', invalidCredentials: 'البريد الإلكتروني أو كلمة المرور غير صحيحة.', rateLimited: 'محاولات كثيرة. انتظر قليلًا ثم حاول مجددًا.', unavailable: 'الخدمة غير متاحة مؤقتًا. حاول مرة أخرى.', network: 'تعذر الاتصال. تحقق من الاتصال ثم حاول مجددًا.', unexpected: 'تعذر إنشاء الحساب. حاول مرة أخرى.', conflict: 'هذا البريد الإلكتروني أو رابط الموقع مستخدم بالفعل.' }
    : { validation: 'Please review the highlighted fields.', invalidCredentials: 'Email or password is incorrect.', rateLimited: 'Too many attempts. Wait a moment and try again.', unavailable: 'The service is temporarily unavailable. Please try again.', network: 'Unable to connect. Check your connection and try again.', unexpected: 'We could not create your account. Please try again.', conflict: 'That email or website address is already in use.' }
  const messagesByState = language === 'ar'
    ? { invalid: messages.invalidSubdomain, checking: 'جارٍ التحقق من التوفر…', available: text.subdomainAvailable, taken: text.subdomainUnavailable, reserved: 'رابط الموقع هذا محجوز.', error: 'تعذر التحقق من التوفر. تحقق من الاتصال ثم حاول مجددًا.' }
    : { invalid: messages.invalidSubdomain, checking: 'Checking availability…', available: text.subdomainAvailable, taken: text.subdomainUnavailable, reserved: 'This website address is reserved.', error: 'Availability could not be checked. Check your connection and try again.' }
  const focusField = (field: keyof RegisterForm) => window.requestAnimationFrame(() => setFocus(field))

  const onSubmit = async (data: RegisterForm) => {
    if (availability.state !== 'available' || availability.checkedValue !== data.restaurantSlug.trim().toLowerCase()) {
      setError('restaurantSlug', { type: 'validate', message: availability.state === 'error' ? messagesByState.error : text.subdomainRequired })
      focusField('restaurantSlug')
      return
    }
    clearErrors()
    try {
      const response = await api.post('/auth/register', { name: data.name, email: data.email, password: data.password, restaurant_name: data.restaurantName, restaurant_slug: data.restaurantSlug, phone: data.phone || undefined, billing_period: billing })
      trackFunnel('trial_registration_completed', { locale: language, source: 'registration', billing })
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
        focusField(fieldName as keyof RegisterForm)
      } else setError('root.server', { type: parsed.code, message: parsed.message })
    }
  }

  const onInvalid = (invalid: Record<string, unknown>) => {
    const firstField = Object.keys(invalid)[0] as keyof RegisterForm | undefined
    if (firstField) focusField(firstField)
  }
  const inputClass = 'poos-input'
  const common = { disabled: isSubmitting, className: inputClass }
  const fieldError = (key: keyof RegisterForm) => errors[key]?.message
  const { restaurantDomain } = getConfig()

  return <>
    <a className="poos-skip-link" href="#main-content">{t.nav.skipToContent}</a>
    <Navbar />
    <main id="main-content" dir={dir} className="auth-page">
      <div className="auth-layout">
        <section className="auth-story" aria-label={language === 'ar' ? 'عن SERVIO' : 'About SERVIO'}>
          <div>
            <h1>{language === 'ar' ? 'ابدأ مساحة عمل مطعمك.' : 'Set up your restaurant workspace.'}</h1>
            <p>{t.footer.tagline}</p>
          </div>
          <div className="auth-story__rule">{language === 'ar' ? 'منصة واحدة لفريق الخدمة والمطبخ.' : 'One platform for front- and back-of-house work.'}</div>
        </section>
        <section className="auth-form-panel" aria-labelledby="register-heading">
          <div className="auth-form">
            <Link className="auth-form__back" href={`/${language}`}>{language === 'ar' ? 'العودة إلى SERVIO' : 'Back to SERVIO'}</Link>
            <h2 id="register-heading">{text.title}</h2>
            <p className="auth-form__intro">{language === 'ar' ? 'أدخل بياناتك وبيانات مطعمك لبدء إنشاء الحساب.' : 'Enter your details and restaurant information to create an account.'}</p>
            <form method="post" noValidate onSubmit={handleSubmit(onSubmit, onInvalid)} className="auth-form__fields">
              <p className="mb-6 text-sm text-slate-600">{language === 'ar' ? `تجربة لمدة 14 يومًا. تفضيل الفوترة: ${billing === 'annual' ? 'سنوي' : 'شهري'}. لا يتم تحصيل أي مبلغ عند التسجيل.` : `14-day trial. Billing preference: ${billing}. Registration does not charge you.`}</p>
              <FormErrorSummary title={fallback.validation} errors={isSubmitted ? Object.fromEntries(Object.entries(errors).filter(([key]) => key !== 'root').map(([key, value]) => [key, value?.message ?? 'Review this field'])) : {}} onFocusField={(key) => setFocus(key as keyof RegisterForm)} />
              {errors.root?.server?.message && <p role="alert" className="poos-error">{errors.root.server.message}</p>}
              <FormField id="name" label={text.fullName} required error={fieldError('name')}><input {...register('name')} {...fieldAccessibility('name', fieldError('name'))} autoComplete="name" maxLength={255} placeholder={text.fullNamePlaceholder} {...common} /></FormField>
              <FormField id="email" label={text.email} required error={fieldError('email')}><input {...register('email')} {...fieldAccessibility('email', fieldError('email'))} type="email" autoComplete="email" maxLength={255} placeholder={text.emailPlaceholder} {...common} /></FormField>
              <FormField id="restaurantName" label={text.restaurantName} required error={fieldError('restaurantName')}><input {...register('restaurantName')} {...fieldAccessibility('restaurantName', fieldError('restaurantName'))} autoComplete="organization" maxLength={255} placeholder={text.restaurantNamePlaceholder} {...common} /></FormField>
              <FormField id="restaurantSlug" label={text.websiteSubdomain} required description={`${language === 'ar' ? 'عنوان الموقع المطلوب؛ تفعيل الموقع خطوة إعداد مستقلة:' : 'Requested address; storefront activation is a separate setup step:'} ${restaurantSlug || 'your-subdomain'}.${restaurantDomain}`} error={fieldError('restaurantSlug')}>
                <input {...register('restaurantSlug')} {...fieldAccessibility('restaurantSlug', fieldError('restaurantSlug'))} autoCapitalize="none" autoCorrect="off" maxLength={63} placeholder={text.subdomainPlaceholder} {...common} />
                <AsyncAvailabilityField state={availability.state} messages={messagesByState} />
              </FormField>
              <FormField id="phone" label={text.phone} optionalLabel={t.auth.validation.optional} error={fieldError('phone')}><input {...register('phone')} {...fieldAccessibility('phone', fieldError('phone'))} type="tel" autoComplete="tel" placeholder={text.phonePlaceholder} {...common} /></FormField>
              <FormField id="password" label={text.password} required description={language === 'ar' ? '6 أحرف على الأقل، وبحد أقصى 72 بايتًا.' : 'At least 6 characters; maximum 72 UTF-8 bytes.'} error={fieldError('password')}><div className="auth-field__password"><PasswordField {...register('password')} {...fieldAccessibility('password', fieldError('password'))} autoComplete="new-password" placeholder={text.passwordPlaceholder} label={text.password} visibleLabel={language === 'ar' ? 'إظهار كلمة المرور' : 'Show password'} hiddenLabel={language === 'ar' ? 'إخفاء كلمة المرور' : 'Hide password'} {...common} /></div></FormField>
              <FormField id="confirmPassword" label={text.confirmPassword} required error={fieldError('confirmPassword')}><div className="auth-field__password"><PasswordField {...register('confirmPassword')} {...fieldAccessibility('confirmPassword', fieldError('confirmPassword'))} autoComplete="new-password" placeholder={text.confirmPasswordPlaceholder} label={text.confirmPassword} visibleLabel={language === 'ar' ? 'إظهار كلمة المرور' : 'Show password'} hiddenLabel={language === 'ar' ? 'إخفاء كلمة المرور' : 'Hide password'} {...common} /></div></FormField>
              <SubmitButton isSubmitting={isSubmitting} loadingLabel={text.creatingAccount} className="poos-button auth-form__submit">{text.createAccount}</SubmitButton>
            </form>
          </div>
        </section>
      </div>
    </main>
  </>
}
