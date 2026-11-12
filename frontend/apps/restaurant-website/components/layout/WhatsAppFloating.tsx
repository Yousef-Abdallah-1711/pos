'use client'

import Link from 'next/link'
import { MessageCircle } from 'lucide-react'

interface WhatsAppFloatingProps {
  phoneNumber?: string
  show?: boolean
}

export function WhatsAppFloating({ phoneNumber, show = true }: WhatsAppFloatingProps) {
  if (!show || !phoneNumber) return null

  const normalized = phoneNumber.replace(/\D/g, '')
  const href = `https://wa.me/${normalized}`

  return (
    <div className="fixed bottom-[calc(4.75rem+env(safe-area-inset-bottom))] right-4 z-50 md:bottom-6">
      <Link
        href={href}
        target="_blank"
        rel="noreferrer"
        className="group relative flex h-12 w-12 items-center justify-center rounded-full bg-green-600 text-white shadow-lg transition hover:scale-105"
        aria-label="WhatsApp"
      >
        <span className="absolute inset-0 animate-ping rounded-full bg-green-500/40" />
        <MessageCircle className="relative z-10 h-6 w-6" />
      </Link>
    </div>
  )
}
