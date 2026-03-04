'use client'

import { useTranslations } from 'next-intl'

interface Props {
  checklist: Array<{ key: string; label: string; done: boolean }>
  onToggle: (key: string) => void
}

export function ExitChecklist({ checklist, onToggle }: Props) {
  const t = useTranslations('hrAdvanced')
  return (
    <div className="rounded border bg-white p-4 dark:bg-gray-800">
      <h3 className="mb-3 font-semibold">{t('exitChecklist.title')}</h3>
      <div className="space-y-2">
        {checklist.map((item) => (
          <label key={item.key} className="flex items-center gap-2 text-sm">
            <input type="checkbox" checked={item.done} onChange={() => onToggle(item.key)} />
            <span>{item.label}</span>
          </label>
        ))}
      </div>
    </div>
  )
}
