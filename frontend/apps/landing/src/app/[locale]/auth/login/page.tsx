'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ChefHat, Eye, EyeOff, Loader2 } from 'lucide-react'
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
  const router = useRouter()
  const setAuth = useAuthStore((state) => state.setAuth)
  const { t, dir, language } = useLanguage()
  const locale = language
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [showTenantSelector, setShowTenantSelector] = useState(false)
  const [selectedEmail, setSelectedEmail] = useState('')
  const [availableTenants, setAvailableTenants] = useState<Tenant[]>([])

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
        setShowTenantSelector(true)
        toast.info(t.auth.login.multipleTenants)
        setIsLoading(false)
        return
      }

      const { user, token } = response.data
      setAuth(user, token)
      toast.success(t.auth.login.loginSuccess)

      const { dashboardUrl } = getConfig()
      window.location.href = `${dashboardUrl}/auth/callback?token=${encodeURIComponent(token)}&user=${encodeURIComponent(JSON.stringify(user))}`
    } catch (error: any) {
      let message = t.auth.login.loginFailed
      if (error.response?.data?.error) {
        message = error.response.data.error
      } else if (error.response?.status === 401) {
        message = t.auth.login.invalidCredentials
      } else if (error.code === 'ERR_NETWORK') {
        message = t.auth.login.networkError
      }
      toast.error(message)
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
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
              {/* Email */}
              <div className="space-y-2">
                <Label htmlFor="email">{t.auth.login.email}</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder={t.auth.login.emailPlaceholder}
                  {...register('email')}
                  disabled={isLoading}
                />
                {errors.email && (
                  <p className="text-sm text-red-500">{errors.email.message}</p>
                )}
              </div>

              {/* Password */}
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <Label htmlFor="password">{t.auth.login.password}</Label>
                  <Link
                    href={`/${locale}/auth/forgot-password`}
                    className="text-sm text-orange-600 hover:underline dark:text-orange-400"
                  >
                    {t.auth.login.forgotPassword}
                  </Link>
                </div>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? 'text' : 'password'}
                    placeholder={t.auth.login.passwordPlaceholder}
                    {...register('password')}
                    disabled={isLoading}
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className={`absolute ${dir === 'rtl' ? 'left-3' : 'right-3'} top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200`}
                  >
                    {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                  </button>
                </div>
                {errors.password && (
                  <p className="text-sm text-red-500">{errors.password.message}</p>
                )}
              </div>

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
        tenants={availableTenants}
        onTenantSelected={handleTenantSelected}
      />
    </>
  )
}
