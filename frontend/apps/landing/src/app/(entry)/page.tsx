import { redirect } from 'next/navigation'
import { defaultLocale } from '@pos-saas/i18n'

export default function RootPage() {
  redirect(`/${defaultLocale}`)
}
