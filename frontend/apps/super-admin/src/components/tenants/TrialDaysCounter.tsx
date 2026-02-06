'use client'

interface TrialDaysCounterProps {
  trialEndDate: string | null | undefined
  trialDays: number
  showProgress?: boolean
  className?: string
}

export function TrialDaysCounter({
  trialEndDate,
  trialDays,
  showProgress = false,
  className = '',
}: TrialDaysCounterProps) {
  if (!trialEndDate) {
    return <span className={`text-red-600 ${className}`}>No trial date</span>
  }

  const endDate = new Date(trialEndDate)
  const now = new Date()
  const daysRemaining = Math.ceil(
    (endDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)
  )

  const getColorClass = (): string => {
    if (daysRemaining > 7) return 'text-green-600'
    if (daysRemaining >= 3 && daysRemaining <= 7) return 'text-yellow-600'
    if (daysRemaining >= 1 && daysRemaining < 3) return 'text-orange-600'
    return 'text-red-600'
  }

  const getProgressColor = (): string => {
    if (daysRemaining > 7) return 'bg-green-200'
    if (daysRemaining >= 3 && daysRemaining <= 7) return 'bg-yellow-200'
    if (daysRemaining >= 1 && daysRemaining < 3) return 'bg-orange-200'
    return 'bg-red-200'
  }

  const progressPercent = Math.max(0, (daysRemaining / trialDays) * 100)

  let displayText: string

  if (daysRemaining > 7) {
    displayText = `${daysRemaining} days remaining`
  } else if (daysRemaining >= 1 && daysRemaining <= 7) {
    displayText = `${daysRemaining} day${daysRemaining !== 1 ? 's' : ''} remaining`
  } else if (daysRemaining === 0) {
    displayText = 'Expires today'
  } else if (daysRemaining < 0) {
    displayText = `Expired ${Math.abs(daysRemaining)} day${Math.abs(daysRemaining) !== 1 ? 's' : ''} ago`
  } else {
    displayText = 'Expires soon'
  }

  return (
    <div className={className}>
      <span className={`font-semibold ${getColorClass()}`}>{displayText}</span>

      {showProgress && (
        <div className="mt-2 w-full bg-gray-200 rounded-full h-2">
          <div
            className={`h-2 rounded-full transition-all duration-300 ${getProgressColor()}`}
            style={{ width: `${Math.max(0, progressPercent)}%` }}
          />
        </div>
      )}
    </div>
  )
}
