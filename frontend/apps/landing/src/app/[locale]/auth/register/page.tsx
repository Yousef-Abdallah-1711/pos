'use client'

import { useState, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ChefHat, Eye, EyeOff, Loader2, Check, X } from 'lucide-react'
import api from '@/lib/api'
import { useAuthStore } from '@/stores/authStore'
import { toast } from 'sonner'
import { getConfig } from '@/lib/config'
import { useLanguage } from '@/contexts/LanguageContext'

type RegisterForm = {
  name: string
  email: string
  password: string
  confirmPassword: string
  restaurantName: string
  restaurantSlug: string
  phone?: string
}

export default function RegisterPage() {
  const router = useRouter()
  const setAuth = useAuthStore((state) => state.setAuth)
  const { t, dir, language } = useLanguage()
  const locale = language
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [subdomainStatus, setSubdomainStatus] = useState<'idle' | 'checking' | 'available' | 'unavailable'>('idle')
  const [subdomainMessage, setSubdomainMessage] = useState('')
  const [checkSubdomainTimeout, setCheckSubdomainTimeout] = useState<NodeJS.Timeout | null>(null)

  const registerSchema = z.object({
    name: z.string().min(2, t.auth.validation.nameMin),
    email: z.string().email(t.auth.validation.invalidEmail),
    password: z.string().min(6, t.auth.validation.passwordMin),
    confirmPassword: z.string(),
    restaurantName: z.string().min(2, t.auth.validation.restaurantNameMin),
    restaurantSlug: z.string().min(2, t.auth.validation.subdomainMin).max(63, t.auth.validation.subdomainMax),
    phone: z.string().optional(),
  }).refine((data) => data.password === data.confirmPassword, {
    message: t.auth.validation.passwordsMismatch,
    path: ['confirmPassword'],
  })

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
  })

  const restaurantSlug = watch('restaurantSlug')

  const checkSubdomainAvailability = useCallback(async (slug: string) => {
    if (!slug || slug.length < 2) {
      setSubdomainStatus('idle')
      return
    }

    setSubdomainStatus('checking')
    try {
      const response = await api.get('/auth/check-subdomain', {
        params: { slug: slug.toLowerCase() },
      })

      if (response.data.available) {
        setSubdomainStatus('available')
        setSubdomainMessage(t.auth.register.subdomainAvailable)
      } else {
        setSubdomainStatus('unavailable')
        setSubdomainMessage(response.data.message || t.auth.register.subdomainUnavailable)
      }
    } catch (error) {
      setSubdomainStatus('unavailable')
      setSubdomainMessage(t.auth.register.subdomainError)
    }
  }, [t])

  const handleSubdomainChange = useCallback((slug: string) => {
    if (checkSubdomainTimeout) {
      clearTimeout(checkSubdomainTimeout)
    }

    if (!slug || slug.length < 2) {
      setSubdomainStatus('idle')
      return
    }

    setSubdomainStatus('checking')
    const timeout = setTimeout(() => {
      checkSubdomainAvailability(slug)
    }, 500)

    setCheckSubdomainTimeout(timeout)
  }, [checkSubdomainTimeout, checkSubdomainAvailability])

  const onSubmit = async (data: RegisterForm) => {
    setIsLoading(true)
    try {
      if (subdomainStatus !== 'available') {
        toast.error(t.auth.register.subdomainRequired)
        setIsLoading(false)
        return
      }

      const response = await api.post('/auth/register', {
        name: data.name,
        email: data.email,
        password: data.password,
        restaurant_name: data.restaurantName,
        restaurant_slug: data.restaurantSlug.toLowerCase(),
        phone: data.phone,
      })

      const { user, token } = response.data
      setAuth(user, token)
      toast.success(t.auth.register.accountCreated)

      const { dashboardUrl } = getConfig()
      window.location.href = `${dashboardUrl}/auth/callback?token=${encodeURIComponent(token)}&user=${encodeURIComponent(JSON.stringify(user))}`
    } catch (error: any) {
      let message = t.auth.register.registrationFailed
      if (error.response?.data?.error) {
        message = error.response.data.error
      } else if (error.response?.status === 400) {
        message = t.auth.register.invalidInput
      } else if (error.code === 'ERR_NETWORK') {
        message = t.auth.register.networkError
      }
      toast.error(message)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div dir={dir} className="min-h-screen flex items-center justify-center bg-gradient-to-br from-orange-50 to-red-50 dark:from-gray-900 dark:to-gray-800 px-4 py-12">
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
            {t.auth.register.title}
          </p>
        </div>

        {/* Register Form */}
        <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-700 p-8">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            {/* Name */}
            <div className="space-y-2">
              <Label htmlFor="name">{t.auth.register.fullName}</Label>
              <Input id="name" type="text" placeholder={t.auth.register.fullNamePlaceholder} {...register('name')} disabled={isLoading} />
              {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
            </div>

            {/* Email */}
            <div className="space-y-2">
              <Label htmlFor="email">{t.auth.register.email}</Label>
              <Input id="email" type="email" placeholder={t.auth.register.emailPlaceholder} {...register('email')} disabled={isLoading} />
              {errors.email && <p className="text-sm text-red-500">{errors.email.message}</p>}
            </div>

            {/* Restaurant Name */}
            <div className="space-y-2">
              <Label htmlFor="restaurantName">{t.auth.register.restaurantName}</Label>
              <Input id="restaurantName" type="text" placeholder={t.auth.register.restaurantNamePlaceholder} {...register('restaurantName')} disabled={isLoading} />
              {errors.restaurantName && <p className="text-sm text-red-500">{errors.restaurantName.message}</p>}
            </div>

            {/* Website Subdomain */}
            <div className="space-y-2">
              <Label htmlFor="restaurantSlug">{t.auth.register.websiteSubdomain}</Label>
              <div className="relative">
                <Input
                  id="restaurantSlug"
                  type="text"
                  placeholder={t.auth.register.subdomainPlaceholder}
                  {...register('restaurantSlug', {
                    onChange: (e) => handleSubdomainChange(e.target.value),
                  })}
                  disabled={isLoading}
                  className={`${
                    subdomainStatus === 'available' ? 'border-green-500' : subdomainStatus === 'unavailable' ? 'border-red-500' : ''
                  }`}
                />
                <div className={`absolute ${dir === 'rtl' ? 'left-3' : 'right-3'} top-1/2 -translate-y-1/2`}>
                  {subdomainStatus === 'checking' && <Loader2 className="h-5 w-5 text-yellow-500 animate-spin" />}
                  {subdomainStatus === 'available' && <Check className="h-5 w-5 text-green-500" />}
                  {subdomainStatus === 'unavailable' && <X className="h-5 w-5 text-red-500" />}
                </div>
              </div>
              {subdomainStatus !== 'idle' && (
                <p className={`text-sm ${subdomainStatus === 'available' ? 'text-green-600' : 'text-red-600'}`}>
                  {subdomainMessage}
                </p>
              )}
              {errors.restaurantSlug && <p className="text-sm text-red-500">{errors.restaurantSlug.message}</p>}
              <p className="text-xs text-gray-500 dark:text-gray-400">
                {t.auth.register.websiteAccessible}{' '}
                <strong>{restaurantSlug ? restaurantSlug.toLowerCase() : 'your-subdomain'}.{getConfig().restaurantDomain}</strong>
              </p>
            </div>

            {/* Phone */}
            <div className="space-y-2">
              <Label htmlFor="phone">{t.auth.register.phone}</Label>
              <Input id="phone" type="tel" placeholder={t.auth.register.phonePlaceholder} {...register('phone')} disabled={isLoading} />
            </div>

            {/* Password */}
            <div className="space-y-2">
              <Label htmlFor="password">{t.auth.register.password}</Label>
              <div className="relative">
                <Input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  placeholder={t.auth.register.passwordPlaceholder}
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
              {errors.password && <p className="text-sm text-red-500">{errors.password.message}</p>}
            </div>

            {/* Confirm Password */}
            <div className="space-y-2">
              <Label htmlFor="confirmPassword">{t.auth.register.confirmPassword}</Label>
              <div className="relative">
                <Input
                  id="confirmPassword"
                  type={showConfirmPassword ? 'text' : 'password'}
                  placeholder={t.auth.register.confirmPasswordPlaceholder}
                  {...register('confirmPassword')}
                  disabled={isLoading}
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className={`absolute ${dir === 'rtl' ? 'left-3' : 'right-3'} top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200`}
                >
                  {showConfirmPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {errors.confirmPassword && <p className="text-sm text-red-500">{errors.confirmPassword.message}</p>}
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
                  {t.auth.register.creatingAccount}
                </>
              ) : (
                t.auth.register.createAccount
              )}
            </Button>
          </form>

          {/* Login Link */}
          <p className="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
            {t.auth.register.hasAccount}{' '}
            <Link href={`/${locale}/auth/login`} className="text-orange-600 font-medium hover:underline dark:text-orange-400">
              {t.auth.register.signIn}
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}
