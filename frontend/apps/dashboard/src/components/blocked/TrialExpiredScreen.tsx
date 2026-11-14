'use client'

import { BlockedOverlay } from './BlockedOverlay'

interface TrialExpiredScreenProps {
  trialEndDate?: string | null
  supportEmail?: string
}

// Share the blocked flow; historical trials must not be assigned today's duration.
export function TrialExpiredScreen({ supportEmail }: TrialExpiredScreenProps) {
  return <BlockedOverlay blockReason="trial_expired" supportEmail={supportEmail} />
}
