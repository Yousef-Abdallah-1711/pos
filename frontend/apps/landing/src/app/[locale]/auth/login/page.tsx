'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from '@/components/ui/button'
import { ChefHat, Loader2 } from 'lucide-react'
import { FormErrorSummary, FormField, PasswordField, applyApiFieldErrors, createLoginSchema, parseApiError, fieldAccessibility } from '@pos-saas/form-foundation'
import api from '@/lib/api'
import { useAuthStore } from '@/stores/authStore'
import { toast } from 'sonner'
import { getConfig } from '@/lib/config'
import { TenantSelectorModal } from '@/components/TenantSelectorModal'
import { useLanguage } from '@/contexts/LanguageContext'

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
      const parsed = parseApiError(error, { validation: 'Please review the highlighted fields.', invalidCredentials: t.auth.login.invalidCredentials, rateLimited: 'Too many attempts. Wait a moment and try again.', unavailable: 'Sign in is temporarily unavailable. Please try again.', network: t.auth.login.networkError, unexpected: t.auth.login.loginFailed, conflict: t.auth.login.loginFailed })
      applyApiFieldErrors(parsed.fields, (field, message) => setError(field as keyof LoginForm, { type: 'server', message }))
      const firstField = Object.keys(parsed.fields)[0]
      if (firstField === 'email' || firstField === 'password') setFocus(firstField)
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
      <div dir={dir} className="min-h-screen flex items-center justify-center bg-gradient-to-br from-orange-50 to-red-50 dark:from-gray-900 dark:to-gray-800 px-4">
        <div className="max-w-md w-full">
          {/* Logo */}
          <div className="text-center mb-8">
            <Link href={`/${locale}`} className="inline-flex items-center gap-2.5 mb-3">
              <div className="w-10 h-10 bg-gradient-to-br from-orange-500 to-red-600 rounded-xl flex items-center justify-center">
                <ChefHat className="h-6 w-6 text-white" />
              </div>
              <span className="text-2xl font-extrabold text-gray-900 dark:text-white">
                Restaurant<span className="text-orange-500">POS</span>
              </span>
            </Link>
            <p className="text-gray-600 dark:text-gray-400">
              {t.auth.login.title}
            </p>
          </div>

          {/* Login Form */}
          <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-700 p-8">
            <form noValidate onSubmit={handleSubmit(onSubmit, (invalid) => { const first = Object.keys(invalid)[0] as keyof LoginForm | undefined; if (first) setFocus(first) })} className="space-y-5">
              {isSubmitted && Object.keys(errors).length > 0 && <FormErrorSummary title="Please review these fields" errors={Object.fromEntries(Object.entries(errors).map(([key, value]) => [key, value?.message ?? 'Review this field']))} onFocusField={(field) => setFocus(field as keyof LoginForm)} />}
              {/* Email */}
              <FormField id="email" label={t.auth.login.email} required error={errors.email?.message}><input type="email" autoComplete="email" maxLength={255} placeholder={t.auth.login.emailPlaceholder} {...register('email')} {...fieldAccessibility('email', errors.email?.message)} disabled={isLoading} className="h-11 w-full rounded-md border border-gray-300 bg-white px-3 text-gray-900 outline-none focus-visible:ring-2 focus-visible:ring-orange-600 disabled:opacity-60 dark:border-gray-600 dark:bg-gray-800 dark:text-white" /></FormField>

              {/* Password */}
              <FormField id="password" label={t.auth.login.password} required error={errors.password?.message}>
                <div className="flex items-center justify-end">
                  <Link
                    href={`/${locale}/auth/forgot-password`}
                    className="text-sm text-orange-600 hover:underline dark:text-orange-400"
                  >
                    {t.auth.login.forgotPassword}
                  </Link>
                </div>
                <PasswordField label={t.auth.login.password} autoComplete="current-password" placeholder={t.auth.login.passwordPlaceholder} {...register('password')} {...fieldAccessibility('password', errors.password?.message)} disabled={isLoading} />
              </FormField>

              {/* Submit */}
              <Button
                type="submit"
                className="w-full bg-gradient-to-r from-orange-500 to-red-600 hover:from-orange-600 hover:to-red-700 text-white border-0 h-11"
                disabled={isLoading}
              >
                {isLoading ? (
                  <>
                    <Loader2 className={`${dir === 'rtl' ? 'ml-2' : 'mr-2'} h-4 w-4 animate-spin`} />
                    {t.auth.login.signingIn}
                  </>
                ) : (
                  t.auth.login.signIn
                )}
              </Button>
            </form>

            {/* Register Link */}
            <p className="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
              {t.auth.login.noAccount}{' '}
              <Link
                href={`/${locale}/auth/register`}
                className="text-orange-600 font-medium hover:underline dark:text-orange-400"
              >
                {t.auth.login.createAccount}
              </Link>
            </p>
          </div>
        </div>
      </div>

      {/* Tenant Selector Modal */}
      <TenantSelectorModal
        isOpen={showTenantSelector}
        email={selectedEmail}
        selectionToken={selectionToken}
        tenants={availableTenants}
        onTenantSelected={handleTenantSelected}
      />
    </>
  )
}
