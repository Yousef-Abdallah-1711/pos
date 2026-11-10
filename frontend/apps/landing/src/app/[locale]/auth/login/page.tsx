'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from '@/components/ui/button'
import { Loader2 } from 'lucide-react'
import { FormErrorSummary, FormField, PasswordField, applyApiFieldErrors, createLoginSchema, parseApiError, fieldAccessibility } from '@pos-saas/form-foundation'
import api from '@/lib/api'
import { useAuthStore } from '@/stores/authStore'
import { toast } from 'sonner'
import { getConfig } from '@/lib/config'
import { TenantSelectorModal } from '@/components/TenantSelectorModal'
import { useLanguage } from '@/contexts/LanguageContext'
import { Navbar } from '@/components/Navbar'

type LoginForm = {
  email: string
  password: string
}

interface Tenant {
  tenant_id: number
  tenant_name: string
  user_id: number
  roles: string[]
  default_restaurant_id?: number
}

export default function LoginPage() {
  const setAuth = useAuthStore((state) => state.setAuth)
  const { t, dir, language } = useLanguage()
  const locale = language
  const [isLoading, setIsLoading] = useState(false)
  const [showTenantSelector, setShowTenantSelector] = useState(false)
  const [selectedEmail, setSelectedEmail] = useState('')
  const [availableTenants, setAvailableTenants] = useState<Tenant[]>([])
  const [selectionToken, setSelectionToken] = useState('')

  useEffect(() => {
    if (new URLSearchParams(window.location.search).get('created') === '1') {
      toast.success('Your account is ready. Sign in to continue.')
    }
  }, [])

  const loginSchema = createLoginSchema({ required: t.auth.validation.required, invalidEmail: t.auth.validation.invalidEmail, emailMax: t.auth.validation.emailMax, passwordMin: t.auth.validation.passwordMin, passwordMax: t.auth.validation.passwordMax })
  const errorFallback = language === 'ar'
    ? { validation: 'يرجى مراجعة الحقول الموضحة.', invalidCredentials: t.auth.login.invalidCredentials, rateLimited: 'محاولات كثيرة. انتظر قليلًا ثم حاول مجددًا.', unavailable: 'تسجيل الدخول غير متاح مؤقتًا. حاول مرة أخرى.', network: t.auth.login.networkError, unexpected: t.auth.login.loginFailed, conflict: t.auth.login.loginFailed }
    : { validation: 'Please review the highlighted fields.', invalidCredentials: t.auth.login.invalidCredentials, rateLimited: 'Too many attempts. Wait a moment and try again.', unavailable: 'Sign in is temporarily unavailable. Please try again.', network: t.auth.login.networkError, unexpected: t.auth.login.loginFailed, conflict: t.auth.login.loginFailed }

  const {
    register,
    handleSubmit,
    setError,
    setFocus,
    formState: { errors, isSubmitted },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
    mode: 'onTouched',
    reValidateMode: 'onChange',
  })

  const onSubmit = async (data: LoginForm) => {
    setIsLoading(true)
    try {
      const response = await api.post('/auth/login', data)

      if (response.data.multiple_tenants && response.data.tenants && response.data.tenants.length > 1) {
        setSelectedEmail(data.email)
        setAvailableTenants(response.data.tenants)
        setSelectionToken(response.data.selection_token ?? '')
        setShowTenantSelector(true)
        toast.info(t.auth.login.multipleTenants)
        setIsLoading(false)
        return
      }

      const { user } = response.data
      setAuth(user, '')
      toast.success(t.auth.login.loginSuccess)

      const { dashboardUrl } = getConfig()
      window.location.href = `${dashboardUrl}/${locale}/auth/callback`
    } catch (error) {
      const parsed = parseApiError(error, errorFallback)
      applyApiFieldErrors(parsed.fields, (field, message) => setError(field as keyof LoginForm, { type: 'server', message }))
      const firstField = Object.keys(parsed.fields)[0]
      if (firstField === 'email' || firstField === 'password') window.requestAnimationFrame(() => setFocus(firstField))
      else if (!firstField) toast.error(parsed.message)
    } finally {
      setIsLoading(false)
    }
  }

  const handleTenantSelected = (tenantId: number) => {
    setShowTenantSelector(false)
  }

  return (
    <>
      <a className="poos-skip-link" href="#main-content">{t.nav.skipToContent}</a>
      <Navbar />
      <main id="main-content" className="auth-page" dir={dir}>
        <div className="auth-layout">
          <section className="auth-story" aria-label={language === 'ar' ? 'عن POOS' : 'About POOS'}>
            <div>
              <h1>{language === 'ar' ? 'إدارة مطعمك في صورة واحدة.' : 'Your restaurant, in one view.'}</h1>
              <p>{t.footer.tagline}</p>
            </div>
            <div className="auth-story__rule">{language === 'ar' ? 'مساحة واحدة لفريق الخدمة والمطبخ.' : 'One platform for front- and back-of-house work.'}</div>
          </section>
          <section className="auth-form-panel" aria-labelledby="login-heading">
            <div className="auth-form">
              <Link className="auth-form__back" href={`/${locale}`}>{language === 'ar' ? 'العودة إلى POOS' : 'Back to POOS'}</Link>
              <h2 id="login-heading">{t.auth.login.title}</h2>
              <p className="auth-form__intro">{language === 'ar' ? 'أدخل بيانات حسابك للمتابعة.' : 'Enter your account details to continue.'}</p>
              <form method="post" noValidate onSubmit={handleSubmit(onSubmit, (invalid) => { const first = Object.keys(invalid)[0] as keyof LoginForm | undefined; if (first) window.requestAnimationFrame(() => setFocus(first)) })} className="auth-form__fields">
                <FormErrorSummary title={errorFallback.validation} errors={isSubmitted ? Object.fromEntries(Object.entries(errors).map(([key, value]) => [key, value?.message ?? errorFallback.validation])) : {}} onFocusField={(field) => setFocus(field as keyof LoginForm)} />
                <FormField id="email" label={t.auth.login.email} required error={errors.email?.message}>
                  <input type="email" autoComplete="email" maxLength={255} placeholder={t.auth.login.emailPlaceholder} {...register('email')} {...fieldAccessibility('email', errors.email?.message)} disabled={isLoading} className="poos-input" />
                </FormField>
                <FormField id="password" label={t.auth.login.password} required error={errors.password?.message}>
                  <div className="auth-field__password"><PasswordField label={t.auth.login.password} autoComplete="current-password" placeholder={t.auth.login.passwordPlaceholder} {...register('password')} {...fieldAccessibility('password', errors.password?.message)} disabled={isLoading} className="poos-input" visibleLabel={language === 'ar' ? 'إظهار كلمة المرور' : 'Show password'} hiddenLabel={language === 'ar' ? 'إخفاء كلمة المرور' : 'Hide password'} /></div>
                </FormField>
                <Button type="submit" className="poos-button auth-form__submit" disabled={isLoading} aria-live="polite">
                  {isLoading ? <><Loader2 aria-hidden="true" className="animate-spin" size={18} />{t.auth.login.signingIn}</> : t.auth.login.signIn}
                </Button>
              </form>
            </div>
          </section>
        </div>
      </main>

      {/* Tenant Selector Modal */}
      <TenantSelectorModal
        isOpen={showTenantSelector}
        email={selectedEmail}
        selectionToken={selectionToken}
        tenants={availableTenants}
        onTenantSelected={handleTenantSelected}
        onCancel={() => setShowTenantSelector(false)}
      />
    </>
  )
}
