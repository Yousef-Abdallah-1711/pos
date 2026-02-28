'use client'

import { useEffect, useRef, useState } from 'react'
import { ClipboardList } from 'lucide-react'
import { KDSDisplayDensity } from '@/hooks/useKDS'
import { KDSOrder } from '@/stores/kdsStore'
import { KDSOrderCard } from './KDSOrderCard'

interface KDSGridProps {
  orders: KDSOrder[]
  locale: 'en' | 'ar'
  t: (path: string, fallback?: string) => string
  loading?: boolean
  canBump: boolean
  canRecall: boolean
  pendingOrderIds: number[]
  displayDensity?: KDSDisplayDensity
  warningMinutes?: number
  dangerMinutes?: number
  onBump: (orderId: number) => void | Promise<void>
  onRecall: (orderId: number) => void | Promise<void>
}

interface AnimatedCard {
  id: number
  order: KDSOrder
  phase: 'enter' | 'idle' | 'exit'
}

const getTimestamp = (value: string) => {
  const timestamp = new Date(value).getTime()
  return Number.isFinite(timestamp) ? timestamp : 0
}

export function KDSGrid({
  orders,
  locale,
  t,
  loading = false,
  canBump,
  canRecall,
  pendingOrderIds,
  displayDensity = 'comfortable',
  warningMinutes = 10,
  dangerMinutes = 20,
  onBump,
  onRecall,
}: KDSGridProps) {
  const [cards, setCards] = useState<AnimatedCard[]>([])
  const exitTimeoutsRef = useRef<Record<number, number>>({})
  const frameRef = useRef<number | null>(null)

  useEffect(() => {
    const sortedOrders = [...orders].sort((left, right) => {
      if (left.priority !== right.priority) {
        return right.priority - left.priority
      }
      return getTimestamp(left.received_at) - getTimestamp(right.received_at)
    })

    setCards((currentCards) => {
      const currentMap = new Map(currentCards.map((card) => [card.id, card]))
      const incomingIds = new Set(sortedOrders.map((order) => order.id))
      const nextCards: AnimatedCard[] = []

      for (const order of sortedOrders) {
        const existing = currentMap.get(order.id)
        if (!existing) {
          nextCards.push({ id: order.id, order, phase: 'enter' })
          continue
        }

        if (existing.phase === 'exit' && exitTimeoutsRef.current[order.id]) {
          window.clearTimeout(exitTimeoutsRef.current[order.id])
          delete exitTimeoutsRef.current[order.id]
        }

        nextCards.push({
          ...existing,
          order,
          phase: existing.phase === 'exit' ? 'idle' : existing.phase,
        })
      }

      for (const existing of currentCards) {
        if (!incomingIds.has(existing.id) && existing.phase !== 'exit') {
          nextCards.push({ ...existing, phase: 'exit' })
        }
      }

      return nextCards
    })
  }, [orders])

  useEffect(() => {
    const enteringIds = cards
      .filter((card) => card.phase === 'enter')
      .map((card) => card.id)

    if (enteringIds.length === 0) {
      return
    }

    frameRef.current = window.requestAnimationFrame(() => {
      setCards((currentCards) =>
        currentCards.map((card) =>
          enteringIds.includes(card.id) && card.phase === 'enter'
            ? { ...card, phase: 'idle' }
            : card
        )
      )
    })

    return () => {
      if (frameRef.current !== null) {
        window.cancelAnimationFrame(frameRef.current)
        frameRef.current = null
      }
    }
  }, [cards])

  useEffect(() => {
    for (const card of cards) {
      if (card.phase !== 'exit') {
        continue
      }
      if (exitTimeoutsRef.current[card.id]) {
        continue
      }

      exitTimeoutsRef.current[card.id] = window.setTimeout(() => {
        setCards((currentCards) => currentCards.filter((item) => item.id !== card.id))
        delete exitTimeoutsRef.current[card.id]
      }, 220)
    }
  }, [cards])

  useEffect(() => {
    return () => {
      Object.values(exitTimeoutsRef.current).forEach((timeoutId) => {
        window.clearTimeout(timeoutId)
      })

      if (frameRef.current !== null) {
        window.cancelAnimationFrame(frameRef.current)
      }
    }
  }, [])

  if (cards.length === 0 && loading) {
    return (
      <div className="rounded-xl border border-gray-200 bg-white p-6 dark:border-slate-700 dark:bg-slate-900/80">
        <div className="mb-4 h-5 w-44 animate-pulse rounded bg-gray-200 dark:bg-slate-700" />
        <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
          {Array.from({ length: 4 }).map((_, index) => (
            <div key={index} className="rounded-lg border border-gray-200 bg-gray-50 p-4 dark:border-slate-700 dark:bg-slate-800/70">
              <div className="mb-3 h-4 w-24 animate-pulse rounded bg-gray-200 dark:bg-slate-600" />
              <div className="mb-2 h-3 w-full animate-pulse rounded bg-gray-200 dark:bg-slate-700" />
              <div className="mb-2 h-3 w-3/4 animate-pulse rounded bg-gray-200 dark:bg-slate-700" />
              <div className="h-9 w-full animate-pulse rounded bg-gray-300 dark:bg-slate-600" />
            </div>
          ))}
        </div>
      </div>
    )
  }

  if (cards.length === 0) {
    return (
      <div className="rounded-xl border border-gray-200 bg-white p-10 text-center dark:border-slate-700 dark:bg-slate-900/80">
        <ClipboardList className="mx-auto mb-3 h-10 w-10 text-gray-400 dark:text-slate-500" />
        <h3 className="mb-1 text-lg font-semibold text-gray-900 dark:text-slate-100">
          {t('kds.noOrdersTitle', 'No active orders')}
        </h3>
        <p className="text-sm text-gray-600 dark:text-slate-400">
          {t('kds.noOrdersText', 'New kitchen orders will appear here as soon as they are confirmed.')}
        </p>
      </div>
    )
  }

  const densityConfig = displayDensity === 'compact'
    ? { minWidth: 280, gapClass: 'gap-3' }
    : displayDensity === 'spacious'
      ? { minWidth: 360, gapClass: 'gap-5' }
      : { minWidth: 320, gapClass: 'gap-4' }

  return (
    <section
      className={`grid ${densityConfig.gapClass}`}
      style={{ gridTemplateColumns: `repeat(auto-fill, minmax(${densityConfig.minWidth}px, 1fr))` }}
    >
      {cards.map((card) => {
        const animationClass =
          card.phase === 'enter'
            ? 'opacity-0 translate-y-3 scale-[0.99]'
            : card.phase === 'exit'
              ? 'pointer-events-none opacity-0 -translate-y-3 scale-[0.99]'
              : 'opacity-100 translate-y-0 scale-100'

        return (
          <div
            key={card.id}
            className={`transform-gpu transition-all duration-200 ease-out ${animationClass}`}
          >
            <KDSOrderCard
              order={card.order}
              locale={locale}
              t={t}
              canBump={canBump}
              canRecall={canRecall}
              isBusy={pendingOrderIds.includes(card.order.id)}
              warningMinutes={warningMinutes}
              dangerMinutes={dangerMinutes}
              onBump={onBump}
              onRecall={onRecall}
            />
          </div>
        )
      })}
    </section>
  )
}
