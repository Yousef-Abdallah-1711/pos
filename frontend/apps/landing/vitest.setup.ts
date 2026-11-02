import '@testing-library/jest-dom/vitest'
import { vi } from 'vitest'

vi.mock('next/navigation', () => ({
  useRouter: vi.fn(() => ({ push: vi.fn(), replace: vi.fn(), refresh: vi.fn() })),
  usePathname: vi.fn(() => '/en/auth/login'),
  redirect: vi.fn(),
}))

vi.mock('@/stores/authStore', () => ({ useAuthStore: vi.fn() }))
vi.mock('@/lib/api', () => ({ default: { post: vi.fn() } }))
vi.mock('sonner', () => ({
  toast: { success: vi.fn(), error: vi.fn(), info: vi.fn() },
}))
vi.mock('@/contexts/LanguageContext', async () => {
  const { translations } = await import('@/lib/translations')
  return {
    useLanguage: vi.fn(() => ({ language: 'en', dir: 'ltr', t: translations.en })),
  }
})
