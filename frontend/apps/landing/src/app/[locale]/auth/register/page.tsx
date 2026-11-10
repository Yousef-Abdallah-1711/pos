'use client'

import { useState, useCallback, useEffect, useRef } from 'react'
import Link from 'next/link'
import { isAxiosError } from 'axios'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Eye, EyeOff, Loader2, Check, X } from 'lucide-react'
import api from '@/lib/api'
import { useAuthStore } from '@/stores/authStore'
import { toast } from 'sonner'
import { getConfig } from '@/lib/config'
import { useLanguage } from '@/contexts/LanguageContext'
import { BrandMark } from '@/components/BrandMark'

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
  const setAuth = useAuthStore((state) => state.setAuth)
  const { t, dir, language } = useLanguage()
  const locale = language
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [subdomainStatus, setSubdomainStatus] = useState<'idle' | 'checking' | 'available' | 'unavailable'>('idle')
  const [subdomainMessage, setSubdomainMessage] = useState('')
  const checkSubdomainTimeout = useRef<ReturnType<typeof setTimeout> | null>(null)
  const currentSlug = useRef('')

  useEffect(() => () => {
    if (checkSubdomainTimeout.current) clearTimeout(checkSubdomainTimeout.current)
  }, [])

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

      if (currentSlug.current !== slug.toLowerCase()) return

      if (response.data.available) {
        setSubdomainStatus('available')
        setSubdomainMessage(t.auth.register.subdomainAvailable)
      } else {
        setSubdomainStatus('unavailable')
        setSubdomainMessage(response.data.message || t.auth.register.subdomainUnavailable)
      }
    } catch (error) {
      if (currentSlug.current !== slug.toLowerCase()) return
      setSubdomainStatus('unavailable')
      setSubdomainMessage(t.auth.register.subdomainError)
    }
  }, [t])

  const handleSubdomainChange = useCallback((slug: string) => {
    if (checkSubdomainTimeout.current) clearTimeout(checkSubdomainTimeout.current)
    currentSlug.current = slug.toLowerCase()

    if (!slug || slug.length < 2) {
      setSubdomainStatus('idle')
      return
    }

    setSubdomainStatus('checking')
    checkSubdomainTimeout.current = setTimeout(() => {
      checkSubdomainAvailability(slug)
    }, 500)
  }, [checkSubdomainAvailability])

  const onSubmit = async (data: RegisterForm) => {
    setIsLoading(true)
    try {
      if (subdomainStatus !== 'available' || currentSlug.current !== data.restaurantSlug.toLowerCase()) {
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

      const { user } = response.data
      setAuth(user, '')
      toast.success(t.auth.register.accountCreated)

      const { dashboardUrl } = getConfig()
      window.location.href = `${dashboardUrl}/${locale}/auth/callback`
    } catch (error: unknown) {
      let message = t.auth.register.registrationFailed
      if (isAxiosError<{ error?: string }>(error)) {
        if (typeof error.response?.data?.error === 'string') {
          message = error.response.data.error
        } else if (error.response?.status === 400) {
          message = t.auth.register.invalidInput
        } else if (error.code === 'ERR_NETWORK') {
          message = t.auth.register.networkError
        }
      }
      toast.error(message)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <>
    <a className="poos-skip-link" href="#main-content">{t.nav.skipToContent}</a>
    <main id="main-content" dir={dir} className="auth-page">
      <div className="auth-layout">
        <section className="auth-story" aria-label={language === 'ar' ? 'عن POOS' : 'About POOS'}>
          <BrandMark href={`/${locale}`} label="POOS" />
          <div>
            <h1>{language === 'ar' ? 'ابدأ مساحة عمل مطعمك.' : 'Set up your restaurant workspace.'}</h1>
            <p>{t.footer.tagline}</p>
          </div>
          <div className="auth-story__rule">{language === 'ar' ? 'منصة واحدة لفريق الخدمة والمطبخ.' : 'One platform for front- and back-of-house work.'}</div>
        </section>

        <section className="auth-form-panel" aria-labelledby="register-heading">
          <div className="auth-form">
            <Link className="auth-form__back" href={`/${locale}`}>
              {language === 'ar' ? 'العودة إلى POOS' : 'Back to POOS'}
            </Link>
            <h2 id="register-heading">{t.auth.register.title}</h2>
            <p className="auth-form__intro">{language === 'ar' ? 'أدخل بياناتك وبيانات مطعمك لبدء إنشاء الحساب.' : 'Enter your details and restaurant information to create an account.'}</p>

          <form onSubmit={handleSubmit(onSubmit)} className="auth-form__fields">
            {/* Name */}
            <div className="auth-field">
              <Label htmlFor="name">{t.auth.register.fullName}</Label>
              <Input id="name" className="poos-input" type="text" autoComplete="name" placeholder={t.auth.register.fullNamePlaceholder} aria-invalid={Boolean(errors.name)} aria-describedby={errors.name ? 'register-name-error' : undefined} {...register('name')} disabled={isLoading} />
              {errors.name && <p id="register-name-error" className="poos-error" role="alert">{errors.name.message}</p>}
            </div>

            {/* Email */}
            <div className="auth-field">
              <Label htmlFor="email">{t.auth.register.email}</Label>
              <Input id="email" className="poos-input" type="email" autoComplete="email" placeholder={t.auth.register.emailPlaceholder} aria-invalid={Boolean(errors.email)} aria-describedby={errors.email ? 'register-email-error' : undefined} {...register('email')} disabled={isLoading} />
              {errors.email && <p id="register-email-error" className="poos-error" role="alert">{errors.email.message}</p>}
            </div>

            {/* Restaurant Name */}
            <div className="auth-field">
              <Label htmlFor="restaurantName">{t.auth.register.restaurantName}</Label>
              <Input id="restaurantName" className="poos-input" type="text" autoComplete="organization" placeholder={t.auth.register.restaurantNamePlaceholder} aria-invalid={Boolean(errors.restaurantName)} aria-describedby={errors.restaurantName ? 'register-restaurant-error' : undefined} {...register('restaurantName')} disabled={isLoading} />
              {errors.restaurantName && <p id="register-restaurant-error" className="poos-error" role="alert">{errors.restaurantName.message}</p>}
            </div>

            {/* Website Subdomain */}
            <div className="auth-field">
              <Label htmlFor="restaurantSlug">{t.auth.register.websiteSubdomain}</Label>
              <div className="auth-field__password">
                <Input
                  id="restaurantSlug"
                  className={`poos-input ${subdomainStatus === 'available' ? 'border-green-600' : subdomainStatus === 'unavailable' ? 'border-red-700' : ''}`}
                  type="text"
                  autoComplete="url"
                  placeholder={t.auth.register.subdomainPlaceholder}
                  aria-invalid={subdomainStatus === 'unavailable' || Boolean(errors.restaurantSlug)}
                  aria-describedby="register-subdomain-status register-subdomain-hint"
                  {...register('restaurantSlug', {
                    onChange: (e) => handleSubdomainChange(e.target.value),
                  })}
                  disabled={isLoading}
                />
                <div className="auth-field__status-icon" aria-hidden="true">
                  {subdomainStatus === 'checking' && <Loader2 size={19} className="animate-spin" />}
                  {subdomainStatus === 'available' && <Check size={19} className="auth-status--available" />}
                  {subdomainStatus === 'unavailable' && <X size={19} className="auth-status--unavailable" />}
                </div>
              </div>
              {subdomainStatus !== 'idle' && (
                <p id="register-subdomain-status" className={`auth-status ${subdomainStatus === 'available' ? 'auth-status--available' : subdomainStatus === 'unavailable' ? 'auth-status--unavailable' : 'poos-muted'}`} role="status" aria-live="polite">
                  {subdomainStatus === 'checking' ? (language === 'ar' ? 'جارٍ التحقق من الرابط…' : 'Checking URL availability…') : subdomainMessage}
                </p>
              )}
              {errors.restaurantSlug && <p id="register-subdomain-error" className="poos-error" role="alert">{errors.restaurantSlug.message}</p>}
              <p id="register-subdomain-hint" className="auth-form__subdomain">
                {t.auth.register.websiteAccessible}{' '}
                <strong>{restaurantSlug ? restaurantSlug.toLowerCase() : 'your-subdomain'}.{getConfig().restaurantDomain}</strong>
              </p>
            </div>

            {/* Phone */}
            <div className="auth-field">
              <Label htmlFor="phone">{t.auth.register.phone}</Label>
              <Input id="phone" className="poos-input" type="tel" autoComplete="tel" placeholder={t.auth.register.phonePlaceholder} {...register('phone')} disabled={isLoading} />
            </div>

            {/* Password */}
            <div className="auth-field">
              <Label htmlFor="password">{t.auth.register.password}</Label>
              <div className="auth-field__password">
                <Input
                  id="password"
                  className="poos-input"
                  type={showPassword ? 'text' : 'password'}
                  autoComplete="new-password"
                  placeholder={t.auth.register.passwordPlaceholder}
                  aria-invalid={Boolean(errors.password)}
                  aria-describedby={errors.password ? 'register-password-error' : undefined}
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
              {errors.password && <p id="register-password-error" className="poos-error" role="alert">{errors.password.message}</p>}
            </div>

            {/* Confirm Password */}
            <div className="auth-field">
              <Label htmlFor="confirmPassword">{t.auth.register.confirmPassword}</Label>
              <div className="auth-field__password">
                <Input
                  id="confirmPassword"
                  className="poos-input"
                  type={showConfirmPassword ? 'text' : 'password'}
                  autoComplete="new-password"
                  placeholder={t.auth.register.confirmPasswordPlaceholder}
                  aria-invalid={Boolean(errors.confirmPassword)}
                  aria-describedby={errors.confirmPassword ? 'register-confirm-password-error' : undefined}
                  {...register('confirmPassword')}
                  disabled={isLoading}
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword((visible) => !visible)}
                  aria-label={showConfirmPassword ? (language === 'ar' ? 'إخفاء تأكيد كلمة المرور' : 'Hide password confirmation') : (language === 'ar' ? 'إظهار تأكيد كلمة المرور' : 'Show password confirmation')}
                  aria-pressed={showConfirmPassword}
                >
                  {showConfirmPassword ? <EyeOff aria-hidden="true" size={20} /> : <Eye aria-hidden="true" size={20} />}
                </button>
              </div>
              {errors.confirmPassword && <p id="register-confirm-password-error" className="poos-error" role="alert">{errors.confirmPassword.message}</p>}
            </div>

            {/* Submit */}
            <Button
              type="submit"
              className="poos-button auth-form__submit"
              disabled={isLoading}
              aria-live="polite"
            >
              {isLoading ? (
                <>
                  <Loader2 aria-hidden="true" className="animate-spin" size={18} />
                  {t.auth.register.creatingAccount}
                </>
              ) : (
                t.auth.register.createAccount
              )}
            </Button>
          </form>

          <p className="auth-form__footer">
            {t.auth.register.hasAccount}{' '}
            <Link href={`/${locale}/auth/login`}>
              {t.auth.register.signIn}
            </Link>
          </p>
          </div>
        </section>
      </div>
    </main>
    </>
  )
}
