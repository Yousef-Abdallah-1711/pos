'use client'

import { useEffect, useState } from 'react'
import { Monitor, Moon, Sun } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { useDarkMode } from '@/lib/hooks/use-dark-mode'

const modeOrder = ['light', 'dark', 'system'] as const

export function DarkModeToggle() {
  const { mode, setMode } = useDarkMode()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  const handleCycle = () => {
    const currentMode = mounted ? mode : 'system'
    const index = modeOrder.indexOf(currentMode)
    const nextMode = modeOrder[(index + 1) % modeOrder.length]
    setMode(nextMode)
  }

  const safeMode = mounted ? mode : 'system'
  const icon = safeMode === 'light' ? Sun : safeMode === 'dark' ? Moon : Monitor
  const label = safeMode === 'light' ? 'Light mode' : safeMode === 'dark' ? 'Dark mode' : 'System mode'
  const Icon = icon

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            type="button"
            variant="ghost"
            size="icon"
            onClick={handleCycle}
            className="h-12 w-12 rounded-full"
            aria-label={label}
          >
            <Icon className="h-5 w-5" />
          </Button>
        </TooltipTrigger>
        <TooltipContent>{label}</TooltipContent>
      </Tooltip>
    </TooltipProvider>
  )
}
