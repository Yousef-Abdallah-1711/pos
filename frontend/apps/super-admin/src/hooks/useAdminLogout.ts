import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { superAdminAPI } from '@/lib/api'
import { useAuthStore } from '@/stores/authStore'

export function useAdminLogout() {
  const router = useRouter()
  const [pending, setPending] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function logout() {
    if (pending) return
    setPending(true)
    setError(null)
    try {
      // Clear client state only after the server revokes the HttpOnly session.
      await superAdminAPI.logout()
      useAuthStore.getState().logout()
      localStorage.removeItem('super-admin-token')
      localStorage.removeItem('super-admin-auth-store')
      router.replace('/login')
    } catch {
      setError('Unable to sign out. Please try again.')
    } finally {
      setPending(false)
    }
  }

  return { logout, pending, error }
}
