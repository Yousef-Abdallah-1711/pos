'use client'

import { useState } from 'react'
import Link from 'next/link'
import { isAxiosError } from 'axios'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Eye, EyeOff, Loader2 } from 'lucide-react'
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

function readLoginError(error: unknown) {
  if (typeof error !== 'object' || error === null || !('response' in error)) return null
  const response = error.response
  if (typeof response !== 'object' || response === null) return null
  const status = 'status' in response && typeof response.status === 'number' ? response.status : undefined
  const data = 'data' in response ? response.data : null
  const message = typeof data === 'object' && data !== null && 'error' in data && typeof data.error === 'string'
    ? data.error
    : undefined
  return { status, message }
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
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [showTenantSelector, setShowTenantSelector] = useState(false)
  const [selectedEmail, setSelectedEmail] = useState('')
  const [availableTenants, setAvailableTenants] = useState<Tenant[]>([])
  const [selectionToken, setSelectionToken] = useState('')

  const loginSchema = z.object({
    email: z.string().email(t.auth.validation.invalidEmail),
    password: z.string().min(6, t.auth.validation.passwordMin),
  })

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
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
    } catch (error: unknown) {
      let message = t.auth.login.loginFailed
      const responseError = readLoginError(error)
      if (responseError?.message) {
        message = responseError.message
      } else if (responseError?.status === 401) {
        message = t.auth.login.invalidCredentials
      } else if (isAxiosError(error) && error.code === 'ERR_NETWORK') {
        message = t.auth.login.networkError
      } else if (error instanceof Error && 'code' in error && error.code === 'ERR_NETWORK') {
        message = t.auth.login.networkError
      }
      toast.error(message)
    } finally {
      setIsLoading(false)
    }
  }

  const handleTenantSelected = () => {
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
              <Link className="auth-form__back" href={`/${locale}`}>
                {language === 'ar' ? 'العودة إلى POOS' : 'Back to POOS'}
              </Link>
              <h2 id="login-heading">{t.auth.login.title}</h2>
              <p className="auth-form__intro">{language === 'ar' ? 'أدخل بيانات حسابك للمتابعة.' : 'Enter your account details to continue.'}</p>

              <form method="post" onSubmit={handleSubmit(onSubmit)} className="auth-form__fields">
                <div className="auth-field">
                  <Label htmlFor="email">{t.auth.login.email}</Label>
                  <Input
                    id="email"
                    className="poos-input"
                    type="email"
                    autoComplete="username"
                    placeholder={t.auth.login.emailPlaceholder}
                    aria-invalid={Boolean(errors.email)}
                    aria-describedby={errors.email ? 'login-email-error' : undefined}
                    {...register('email')}
                    disabled={isLoading}
                  />
                  {errors.email && <p id="login-email-error" className="poos-error" role="alert">{errors.email.message}</p>}
                </div>

                <div className="auth-field">
                  <Label htmlFor="password">{t.auth.login.password}</Label>
                  <div className="auth-field__password">
                    <Input
                      id="password"
                      className="poos-input"
                      type={showPassword ? 'text' : 'password'}
                      autoComplete="current-password"
                      placeholder={t.auth.login.passwordPlaceholder}
                      aria-invalid={Boolean(errors.password)}
                      aria-describedby={errors.password ? 'login-password-error' : undefined}
                      {...register('password')}
                      disabled={isLoading}
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword((visible) => !visible)}
                      aria-label={showPassword ? (language === 'ar' ? 'إخفاء كلمة المرور' : 'Hide password') : (language === 'ar' ? 'إظهار كلمة المرور' : 'Show password')}
                      aria-pressed={showPassword}
                    >
                      {showPassword ? <EyeOff aria-hidden="true" size={20} /> : <Eye aria-hidden="true" size={20} />}
                    </button>
                  </div>
                  {errors.password && <p id="login-password-error" className="poos-error" role="alert">{errors.password.message}</p>}
                </div>

                <Button type="submit" className="poos-button auth-form__submit" disabled={isLoading} aria-live="polite">
                  {isLoading ? <><Loader2 aria-hidden="true" className="animate-spin" size={18} />{t.auth.login.signingIn}</> : t.auth.login.signIn}
                </Button>
              </form>

              <p className="auth-form__footer">
                {t.auth.login.noAccount}{' '}
                <Link href={`/${locale}/auth/register`}>{t.auth.login.createAccount}</Link>
              </p>
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
