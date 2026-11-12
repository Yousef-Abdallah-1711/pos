'use client'

import Link from 'next/link'
import { MessageCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface WhatsAppButtonProps {
  phoneNumber?: string
  show?: boolean
  message?: string
}

export function WhatsAppButton({
  phoneNumber,
  show = true,
  message = '?????? ???? ?? ????',
}: WhatsAppButtonProps) {
  if (!show || !phoneNumber) {
    return null
  }

  const normalized = phoneNumber.replace(/\D/g, '')
  const href = `https://wa.me/${normalized}?text=${encodeURIComponent(message)}`

  return (
    <Button asChild className="h-12 gap-2 bg-green-600 px-4 text-white hover:bg-green-700">
      <Link href={href} target="_blank" rel="noreferrer" aria-label="Order via WhatsApp">
        <MessageCircle className="h-5 w-5" />
        <span className="hidden md:inline">Order via WhatsApp</span>
      </Link>
    </Button>
  )
}
