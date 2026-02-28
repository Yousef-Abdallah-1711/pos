'use client'

import { useMemo, useRef, useState, type ChangeEvent } from 'react'
import {
  Check,
  Copy,
  Download,
  ExternalLink,
  FileText,
  Globe,
  Loader2,
  QrCode,
  Wand2,
} from 'lucide-react'
import { jsPDF } from 'jspdf'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { debugError } from '@/lib/debug'
import { sanitizeQrFileName } from '@/lib/table-qr'
import {
  buildWebsiteQrTargetUrl,
  type WebsiteQrTarget,
} from '@/lib/website-qr'
import {
  QRPreview,
  type QrFrameShape,
  type QrStylePreset,
} from '../../tables/components/QRPreview'

const QR_SIZE_OPTIONS = [
  { label: 'Small', value: 128 },
  { label: 'Medium', value: 256 },
  { label: 'Large', value: 512 },
] as const

const FRAME_SHAPES: Array<{ value: QrFrameShape; label: string }> = [
  { value: 'square', label: 'Square' },
  { value: 'circle', label: 'Circle' },
  { value: 'diamond', label: 'Diamond' },
]

const STYLE_PRESETS: Array<{ value: QrStylePreset; label: string }> = [
  { value: 'standard', label: 'Standard' },
  { value: 'rounded', label: 'Rounded' },
  { value: 'dots', label: 'Circular' },
]

const ERROR_LEVELS = [
  { value: 'L', label: 'Low (~7%)' },
  { value: 'M', label: 'Medium (~15%)' },
  { value: 'Q', label: 'Quartile (~25%)' },
  { value: 'H', label: 'High (~30%)' },
] as const

type ProcessingAction = 'png' | 'svg' | 'pdf' | null
type PresetCategory = 'all' | 'restaurant' | 'cafe' | 'bar'

interface QrDesignPreset {
  id: string
  name: string
  description: string
  category: Exclude<PresetCategory, 'all'>
  fgColor: string
  bgColor: string
  backgroundMode: 'solid' | 'transparent'
  frameShape: QrFrameShape
  stylePreset: QrStylePreset
  errorLevel: 'L' | 'M' | 'Q' | 'H'
  includeMargin: boolean
}

const QR_DESIGN_PRESETS: QrDesignPreset[] = [
  {
    id: 'classic-bistro',
    name: 'Classic Bistro',
    description: 'Warm and elegant dining style.',
    category: 'restaurant',
    fgColor: '#2B2118',
    bgColor: '#FFF7EA',
    backgroundMode: 'solid',
    frameShape: 'square',
    stylePreset: 'standard',
    errorLevel: 'H',
    includeMargin: true,
  },
  {
    id: 'fine-dining-gold',
    name: 'Fine Dining Gold',
    description: 'Luxury look with premium contrast.',
    category: 'restaurant',
    fgColor: '#7A5B1B',
    bgColor: '#FFFDF7',
    backgroundMode: 'solid',
    frameShape: 'square',
    stylePreset: 'rounded',
    errorLevel: 'H',
    includeMargin: true,
  },
  {
    id: 'coffee-house',
    name: 'Coffee House',
    description: 'Cozy tones for specialty cafes.',
    category: 'cafe',
    fgColor: '#5B3A29',
    bgColor: '#FDF7F2',
    backgroundMode: 'solid',
    frameShape: 'circle',
    stylePreset: 'rounded',
    errorLevel: 'H',
    includeMargin: true,
  },
  {
    id: 'artisan-bakery',
    name: 'Artisan Bakery',
    description: 'Soft neutral palette for bakery counters.',
    category: 'cafe',
    fgColor: '#7C5E45',
    bgColor: '#FFF9F2',
    backgroundMode: 'solid',
    frameShape: 'square',
    stylePreset: 'dots',
    errorLevel: 'Q',
    includeMargin: true,
  },
  {
    id: 'cocktail-neon',
    name: 'Cocktail Neon',
    description: 'Night bar identity with bold colors.',
    category: 'bar',
    fgColor: '#DB2777',
    bgColor: '#111827',
    backgroundMode: 'solid',
    frameShape: 'diamond',
    stylePreset: 'dots',
    errorLevel: 'H',
    includeMargin: true,
  },
  {
    id: 'rooftop-lounge',
    name: 'Rooftop Lounge',
    description: 'Elegant dark style for lounge menus.',
    category: 'bar',
    fgColor: '#E5E7EB',
    bgColor: '#1F2937',
    backgroundMode: 'solid',
    frameShape: 'circle',
    stylePreset: 'rounded',
    errorLevel: 'Q',
    includeMargin: true,
  },
  {
    id: 'minimal-pro',
    name: 'Minimal Pro',
    description: 'Universal pro preset for print and screen.',
    category: 'restaurant',
    fgColor: '#000000',
    bgColor: '#FFFFFF',
    backgroundMode: 'solid',
    frameShape: 'square',
    stylePreset: 'standard',
    errorLevel: 'H',
    includeMargin: true,
  },
]

interface WebsiteQRStudioProps {
  locale: string
  websiteHost: string
  t: (path: string, fallback?: string) => string
  isRTL: boolean
}

export function WebsiteQRStudio({
  locale,
  websiteHost,
  t,
  isRTL,
}: WebsiteQRStudioProps) {
  const [target, setTarget] = useState<WebsiteQrTarget>('home')
  const [qrSize, setQrSize] = useState<number>(256)
  const [frameShape, setFrameShape] = useState<QrFrameShape>('square')
  const [stylePreset, setStylePreset] = useState<QrStylePreset>('standard')
  const [fgColor, setFgColor] = useState('#111827')
  const [bgColor, setBgColor] = useState('#FFFFFF')
  const [backgroundMode, setBackgroundMode] = useState<'solid' | 'transparent'>(
    'solid'
  )
  const [logoUrl, setLogoUrl] = useState<string | null>(null)
  const [processingAction, setProcessingAction] = useState<ProcessingAction>(null)
  const [errorLevel, setErrorLevel] = useState<'L' | 'M' | 'Q' | 'H'>('H')
  const [includeMargin, setIncludeMargin] = useState(true)
  const [presetCategory, setPresetCategory] = useState<PresetCategory>('all')
  const [selectedPresetId, setSelectedPresetId] = useState<string>('minimal-pro')
  const [labelMode, setLabelMode] = useState<'website' | 'custom' | 'none'>('website')
  const [customTitle, setCustomTitle] = useState('')
  const [customSubtitle, setCustomSubtitle] = useState('')

  const fileInputRef = useRef<HTMLInputElement | null>(null)
  const qrCanvasRef = useRef<HTMLCanvasElement | null>(null)
  const qrSvgRef = useRef<SVGSVGElement | null>(null)

  const targetUrl = useMemo(
    () => buildWebsiteQrTargetUrl(target, locale, websiteHost),
    [locale, target, websiteHost]
  )
  const baseFileName = useMemo(
    () => sanitizeQrFileName(`website-${target}-${websiteHost}-qr`),
    [target, websiteHost]
  )

  const visiblePresets = useMemo(() => {
    if (presetCategory === 'all') {
      return QR_DESIGN_PRESETS
    }
    return QR_DESIGN_PRESETS.filter((preset) => preset.category === presetCategory)
  }, [presetCategory])

  const selectedPreset = useMemo(
    () => QR_DESIGN_PRESETS.find((preset) => preset.id === selectedPresetId) || null,
    [selectedPresetId]
  )

  const defaultTitle = useMemo(() => {
    if (target === 'menu') {
      return t('themeBuilder.qrDefaultTitleMenu', 'Scan to open menu')
    }
    return t('themeBuilder.qrDefaultTitleHome', 'Scan to visit website')
  }, [target, t])

  const resolvedLabelTitle = useMemo(() => {
    if (labelMode === 'custom') {
      return customTitle.trim() || defaultTitle
    }
    return defaultTitle
  }, [customTitle, defaultTitle, labelMode])

  const resolvedLabelSubtitle = useMemo(() => {
    if (labelMode === 'custom') {
      return customSubtitle.trim() || websiteHost
    }
    return websiteHost
  }, [customSubtitle, labelMode, websiteHost])

  const applyPreset = (preset: QrDesignPreset) => {
    setSelectedPresetId(preset.id)
    setFgColor(preset.fgColor)
    setBgColor(preset.bgColor)
    setBackgroundMode(preset.backgroundMode)
    setFrameShape(preset.frameShape)
    setStylePreset(preset.stylePreset)
    setErrorLevel(preset.errorLevel)
    setIncludeMargin(preset.includeMargin)
  }

  const handleCopyLink = async () => {
    try {
      await navigator.clipboard.writeText(targetUrl)
      toast.success(t('tables.linkCopied', 'Link copied'))
    } catch (error) {
      toast.error(t('tables.copyFailed', 'Failed to copy link'))
      debugError('[WEBSITE_QR_STUDIO] Copy link failed:', error)
    }
  }

  const handleDownloadPng = async () => {
    const canvas = qrCanvasRef.current
    if (!canvas) {
      toast.error(t('tables.downloadFailed', 'Failed to download QR image'))
      return
    }

    setProcessingAction('png')
    try {
      const anchor = document.createElement('a')
      anchor.href = canvas.toDataURL('image/png')
      anchor.download = `${baseFileName}.png`
      document.body.appendChild(anchor)
      anchor.click()
      anchor.remove()
    } catch (error) {
      toast.error(t('tables.downloadFailed', 'Failed to download QR image'))
      debugError('[WEBSITE_QR_STUDIO] PNG download failed:', error)
    } finally {
      setProcessingAction(null)
    }
  }

  const handleDownloadSvg = async () => {
    const svgElement = qrSvgRef.current
    if (!svgElement) {
      toast.error(t('tables.downloadFailed', 'Failed to download QR image'))
      return
    }

    setProcessingAction('svg')
    try {
      const serializer = new XMLSerializer()
      const svgMarkup = serializer.serializeToString(svgElement)
      const blob = new Blob([svgMarkup], { type: 'image/svg+xml;charset=utf-8' })
      const blobUrl = URL.createObjectURL(blob)
      const anchor = document.createElement('a')
      anchor.href = blobUrl
      anchor.download = `${baseFileName}.svg`
      document.body.appendChild(anchor)
      anchor.click()
      anchor.remove()
      URL.revokeObjectURL(blobUrl)
    } catch (error) {
      toast.error(t('tables.downloadFailed', 'Failed to download QR image'))
      debugError('[WEBSITE_QR_STUDIO] SVG download failed:', error)
    } finally {
      setProcessingAction(null)
    }
  }

  const handlePrintPdf = async () => {
    const canvas = qrCanvasRef.current
    if (!canvas) {
      toast.error(t('tables.downloadFailed', 'Failed to download QR image'))
      return
    }

    setProcessingAction('pdf')
    try {
      const doc = new jsPDF({
        orientation: 'portrait',
        unit: 'mm',
        format: 'a4',
      })

      const pageWidth = doc.internal.pageSize.getWidth()
      const qrWidth = 90
      const qrHeight = 90
      const qrX = (pageWidth - qrWidth) / 2
      const qrY = 48

      doc.setFontSize(17)
      doc.text(
        `${t('themeBuilder.qrBuilderTitle', 'Website QR Builder')} - ${resolvedLabelTitle}`,
        pageWidth / 2,
        20,
        { align: 'center' }
      )

      doc.setFontSize(11)
      doc.text(resolvedLabelSubtitle, pageWidth / 2, 28, {
        align: 'center',
      })

      if (selectedPreset) {
        doc.setFontSize(10)
        doc.text(
          `${t('tables.qrPreset', 'Preset')}: ${selectedPreset.name}`,
          pageWidth / 2,
          34,
          { align: 'center' }
        )
      }

      const qrPng = canvas.toDataURL('image/png')
      doc.addImage(qrPng, 'PNG', qrX, qrY, qrWidth, qrHeight)

      doc.setFontSize(10)
      const wrappedUrl = doc.splitTextToSize(targetUrl, pageWidth - 24)
      doc.text(`${t('tables.qrUrl', 'URL')}:`, 12, 150)
      doc.text(wrappedUrl, 12, 156)

      doc.save(`${baseFileName}.pdf`)
    } catch (error) {
      toast.error(t('tables.downloadFailed', 'Failed to download QR image'))
      debugError('[WEBSITE_QR_STUDIO] PDF export failed:', error)
    } finally {
      setProcessingAction(null)
    }
  }

  const handleUploadLogo = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) {
      return
    }

    const reader = new FileReader()
    reader.onload = () => {
      if (typeof reader.result === 'string') {
        setLogoUrl(reader.result)
      }
    }
    reader.onerror = (error) => {
      toast.error(t('tables.uploadLogoFailed', 'Failed to load logo image.'))
      debugError('[WEBSITE_QR_STUDIO] Logo upload failed:', error)
    }
    reader.readAsDataURL(file)
  }

  return (
    <div className="space-y-6">
      <div className="rounded-lg border border-blue-200 bg-blue-50 p-4 dark:border-blue-900/40 dark:bg-blue-950/20">
        <div className="flex items-start gap-3">
          <div className="rounded-lg bg-blue-600 p-2 text-white">
            <Globe className="h-5 w-5" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              {t('themeBuilder.qrBuilderTitle', 'Website QR Builder')}
            </h1>
            <p className="text-sm text-gray-600 dark:text-gray-300">
              {t(
                'themeBuilder.qrBuilderSubtitle',
                'Create QR codes for your website home page or menu page.'
              )}
            </p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 xl:grid-cols-[1.6fr_1fr]">
        <div className="space-y-6">
          <section className="rounded-lg border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-700 dark:bg-gray-800">
            <h2 className="mb-4 text-base font-semibold text-gray-900 dark:text-white">
              {t('themeBuilder.qrTargetLabel', 'QR Target')}
            </h2>

            <div className="mb-3 flex flex-wrap items-center gap-4 text-sm text-gray-700 dark:text-gray-300">
              <label className="inline-flex items-center gap-2">
                <input
                  type="radio"
                  name="qr-target"
                  checked={target === 'home'}
                  onChange={() => setTarget('home')}
                />
                {t('themeBuilder.qrTargetHome', 'Home Page')}
              </label>
              <label className="inline-flex items-center gap-2">
                <input
                  type="radio"
                  name="qr-target"
                  checked={target === 'menu'}
                  onChange={() => setTarget('menu')}
                />
                {t('themeBuilder.qrTargetMenu', 'Menu Page')}
              </label>
            </div>

            <p className="mb-2 text-sm text-gray-700 dark:text-gray-300">
              {t('themeBuilder.qrCurrentHost', 'Current Host')}:{' '}
              <span className="font-semibold">{websiteHost}</span>
            </p>

            <div className="rounded-md border border-gray-200 bg-gray-50 p-3 text-xs text-gray-700 break-all dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200">
              {targetUrl}
            </div>
          </section>

          <section className="rounded-lg border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-700 dark:bg-gray-800">
            <div className="mb-4 flex items-center gap-2">
              <Wand2 className="h-4 w-4 text-primary" />
              <h2 className="text-base font-semibold text-gray-900 dark:text-white">
                {t('tables.qrPresetLibrary', 'Preset Library')}
              </h2>
            </div>

            <div className="mb-4 flex flex-wrap gap-2">
              {([
                { value: 'all', label: t('tables.qrPresetAll', 'All') },
                { value: 'restaurant', label: t('tables.qrPresetRestaurant', 'Restaurant') },
                { value: 'cafe', label: t('tables.qrPresetCafe', 'Cafe') },
                { value: 'bar', label: t('tables.qrPresetBar', 'Bar') },
              ] as const).map((option) => (
                <button
                  key={option.value}
                  type="button"
                  onClick={() => setPresetCategory(option.value)}
                  className={`rounded-full px-3 py-1 text-xs font-medium transition ${
                    presetCategory === option.value
                      ? 'bg-primary text-white'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600'
                  }`}
                >
                  {option.label}
                </button>
              ))}
            </div>

            <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
              {visiblePresets.map((preset) => {
                const active = selectedPresetId === preset.id
                return (
                  <button
                    key={preset.id}
                    type="button"
                    onClick={() => applyPreset(preset)}
                    className={`rounded-lg border p-3 text-left transition ${
                      active
                        ? 'border-primary bg-primary/5'
                        : 'border-gray-200 hover:border-gray-300 dark:border-gray-700 dark:hover:border-gray-500'
                    }`}
                  >
                    <div className="flex items-center justify-between">
                      <p className="text-sm font-semibold text-gray-900 dark:text-white">
                        {preset.name}
                      </p>
                      {active && <Check className="h-4 w-4 text-primary" />}
                    </div>
                    <p className="mt-1 text-xs text-gray-600 dark:text-gray-400">
                      {preset.description}
                    </p>
                  </button>
                )
              })}
            </div>
          </section>

          <section className="rounded-lg border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-700 dark:bg-gray-800">
            <h2 className="mb-4 text-base font-semibold text-gray-900 dark:text-white">
              {t('tables.qrBuilderControls', 'Builder Controls')}
            </h2>

            <div className="space-y-5">
              <div>
                <p className="mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                  {t('tables.qrSize', 'QR Size')}
                </p>
                <div className="mb-3 flex flex-wrap gap-2">
                  {QR_SIZE_OPTIONS.map((option) => (
                    <button
                      key={option.value}
                      type="button"
                      onClick={() => setQrSize(option.value)}
                      className={`rounded-lg px-3 py-2 text-sm font-medium transition ${
                        qrSize === option.value
                          ? 'bg-primary text-white'
                          : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600'
                      }`}
                    >
                      {t(`tables.qrSize${option.label}`, option.label)} ({option.value}px)
                    </button>
                  ))}
                </div>
                <div className="flex items-center gap-3">
                  <input
                    type="range"
                    min={128}
                    max={1024}
                    step={16}
                    value={qrSize}
                    onChange={(event) => setQrSize(Number(event.target.value))}
                    className="w-full"
                  />
                  <span className="w-20 text-right text-sm text-gray-700 dark:text-gray-300">
                    {qrSize}px
                  </span>
                </div>
              </div>

              <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                <div>
                  <p className="mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                    {t('tables.qrShape', 'Shape')}
                  </p>
                  <div className="flex flex-wrap gap-2">
                    {FRAME_SHAPES.map((shape) => (
                      <button
                        key={shape.value}
                        type="button"
                        onClick={() => setFrameShape(shape.value)}
                        className={`rounded-lg px-3 py-2 text-sm font-medium transition ${
                          frameShape === shape.value
                            ? 'bg-primary text-white'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600'
                        }`}
                      >
                        {t(`tables.qrShapes.${shape.value}`, shape.label)}
                      </button>
                    ))}
                  </div>
                </div>

                <div>
                  <p className="mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                    {t('tables.qrStyle', 'Style')}
                  </p>
                  <div className="flex flex-wrap gap-2">
                    {STYLE_PRESETS.map((style) => (
                      <button
                        key={style.value}
                        type="button"
                        onClick={() => setStylePreset(style.value)}
                        className={`rounded-lg px-3 py-2 text-sm font-medium transition ${
                          stylePreset === style.value
                            ? 'bg-primary text-white'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600'
                        }`}
                      >
                        {t(`tables.qrStyles.${style.value}`, style.label)}
                      </button>
                    ))}
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                <label className="block">
                  <span className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {t('tables.qrErrorLevel', 'Error Correction')}
                  </span>
                  <select
                    value={errorLevel}
                    onChange={(event) =>
                      setErrorLevel(event.target.value as 'L' | 'M' | 'Q' | 'H')
                    }
                    className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                  >
                    {ERROR_LEVELS.map((level) => (
                      <option key={level.value} value={level.value}>
                        {level.label}
                      </option>
                    ))}
                  </select>
                </label>
                <label className="mt-7 inline-flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
                  <input
                    type="checkbox"
                    checked={includeMargin}
                    onChange={(event) => setIncludeMargin(event.target.checked)}
                  />
                  {t('tables.qrIncludeMargin', 'Include quiet margin')}
                </label>
              </div>
            </div>
          </section>

          <section className="rounded-lg border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-700 dark:bg-gray-800">
            <h2 className="mb-4 text-base font-semibold text-gray-900 dark:text-white">
              {t('tables.qrBranding', 'Branding & Content')}
            </h2>

            <div className="space-y-5">
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <label className="block">
                  <span className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {t('tables.qrForegroundColor', 'Foreground')}
                  </span>
                  <input
                    type="color"
                    value={fgColor}
                    onChange={(event) => setFgColor(event.target.value)}
                    className="h-10 w-full cursor-pointer rounded border border-gray-200 bg-white p-1 dark:border-gray-700 dark:bg-gray-800"
                  />
                </label>
                <label className="block">
                  <span className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {t('tables.qrBackgroundColor', 'Background')}
                  </span>
                  <input
                    type="color"
                    value={bgColor}
                    onChange={(event) => setBgColor(event.target.value)}
                    disabled={backgroundMode === 'transparent'}
                    className="h-10 w-full cursor-pointer rounded border border-gray-200 bg-white p-1 disabled:cursor-not-allowed disabled:opacity-50 dark:border-gray-700 dark:bg-gray-800"
                  />
                </label>
              </div>

              <div className="flex flex-wrap items-center gap-4 text-sm text-gray-700 dark:text-gray-300">
                <label className="inline-flex items-center gap-2">
                  <input
                    type="radio"
                    name="bg-mode"
                    checked={backgroundMode === 'solid'}
                    onChange={() => setBackgroundMode('solid')}
                  />
                  {t('tables.qrBackgroundSolid', 'Solid Background')}
                </label>
                <label className="inline-flex items-center gap-2">
                  <input
                    type="radio"
                    name="bg-mode"
                    checked={backgroundMode === 'transparent'}
                    onChange={() => setBackgroundMode('transparent')}
                  />
                  {t('tables.qrBackgroundTransparent', 'Transparent Background')}
                </label>
              </div>

              <div className="space-y-2">
                <p className="text-sm font-medium text-gray-700 dark:text-gray-300">
                  {t('tables.qrLogoOverlay', 'Logo Overlay')}
                </p>
                <div className="flex flex-wrap gap-2">
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    onChange={handleUploadLogo}
                    className="block text-xs text-gray-600 file:mr-3 file:rounded file:border-0 file:bg-primary file:px-3 file:py-2 file:text-xs file:font-medium file:text-white hover:file:bg-primary/90 dark:text-gray-300"
                  />
                  {logoUrl && (
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => {
                        setLogoUrl(null)
                        if (fileInputRef.current) {
                          fileInputRef.current.value = ''
                        }
                      }}
                    >
                      {t('tables.removeLogo', 'Remove Logo')}
                    </Button>
                  )}
                </div>
              </div>

              <div>
                <p className="mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                  {t('tables.qrLabelMode', 'Label mode')}
                </p>
                <div className="flex flex-wrap gap-4 text-sm text-gray-700 dark:text-gray-300">
                  <label className="inline-flex items-center gap-2">
                    <input
                      type="radio"
                      name="label-mode"
                      checked={labelMode === 'website'}
                      onChange={() => setLabelMode('website')}
                    />
                    {t('themeBuilder.qrLabelWebsite', 'Use website info')}
                  </label>
                  <label className="inline-flex items-center gap-2">
                    <input
                      type="radio"
                      name="label-mode"
                      checked={labelMode === 'custom'}
                      onChange={() => setLabelMode('custom')}
                    />
                    {t('tables.qrLabelCustom', 'Custom label')}
                  </label>
                  <label className="inline-flex items-center gap-2">
                    <input
                      type="radio"
                      name="label-mode"
                      checked={labelMode === 'none'}
                      onChange={() => setLabelMode('none')}
                    />
                    {t('tables.qrLabelNone', 'No label')}
                  </label>
                </div>
              </div>

              {labelMode === 'custom' && (
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                  <label className="block">
                    <span className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                      {t('tables.qrCustomTitle', 'Title')}
                    </span>
                    <input
                      value={customTitle}
                      onChange={(event) => setCustomTitle(event.target.value)}
                      placeholder={defaultTitle}
                      className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                    />
                  </label>
                  <label className="block">
                    <span className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                      {t('tables.qrCustomSubtitle', 'Subtitle')}
                    </span>
                    <input
                      value={customSubtitle}
                      onChange={(event) => setCustomSubtitle(event.target.value)}
                      placeholder={websiteHost}
                      className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
                    />
                  </label>
                </div>
              )}
            </div>
          </section>

          <section className="rounded-lg border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-700 dark:bg-gray-800">
            <h2 className="mb-3 text-base font-semibold text-gray-900 dark:text-white">
              {t('tables.qrExport', 'Export')}
            </h2>
            <div className="mb-4 rounded-lg border border-gray-200 bg-gray-50 p-3 text-xs text-gray-700 break-all dark:border-gray-700 dark:bg-gray-900 dark:text-gray-200">
              {targetUrl}
            </div>
            <div className={`flex flex-wrap gap-2 ${isRTL ? 'justify-start' : ''}`}>
              <Button
                variant="outline"
                onClick={handleCopyLink}
                disabled={processingAction !== null}
              >
                <Copy className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
                {t('tables.copyLink', 'Copy Link')}
              </Button>
              <Button
                variant="outline"
                onClick={() => window.open(targetUrl, '_blank', 'noopener,noreferrer')}
                disabled={processingAction !== null}
              >
                <ExternalLink className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
                {t('tables.openLink', 'Open Link')}
              </Button>
              <Button
                variant="outline"
                onClick={handleDownloadPng}
                disabled={processingAction !== null}
              >
                {processingAction === 'png' ? (
                  <Loader2 className={isRTL ? 'ml-2 h-4 w-4 animate-spin' : 'mr-2 h-4 w-4 animate-spin'} />
                ) : (
                  <Download className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
                )}
                {t('tables.downloadPng', 'Download PNG')}
              </Button>
              <Button
                variant="outline"
                onClick={handleDownloadSvg}
                disabled={processingAction !== null}
              >
                {processingAction === 'svg' ? (
                  <Loader2 className={isRTL ? 'ml-2 h-4 w-4 animate-spin' : 'mr-2 h-4 w-4 animate-spin'} />
                ) : (
                  <QrCode className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
                )}
                {t('tables.downloadSvg', 'Download SVG')}
              </Button>
              <Button
                variant="outline"
                onClick={handlePrintPdf}
                disabled={processingAction !== null}
              >
                {processingAction === 'pdf' ? (
                  <Loader2 className={isRTL ? 'ml-2 h-4 w-4 animate-spin' : 'mr-2 h-4 w-4 animate-spin'} />
                ) : (
                  <FileText className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
                )}
                {t('tables.printPdf', 'Print PDF')}
              </Button>
            </div>
          </section>
        </div>

        <aside className="h-fit rounded-lg border border-gray-200 bg-white p-5 shadow-sm xl:sticky xl:top-6 dark:border-gray-700 dark:bg-gray-800">
          <div className="mb-4">
            <h2 className="text-base font-semibold text-gray-900 dark:text-white">
              {t('tables.qrLivePreview', 'Live Preview')}
            </h2>
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
              {selectedPreset
                ? `${t('tables.qrPreset', 'Preset')}: ${selectedPreset.name}`
                : t('themeBuilder.qrBuilderTitle', 'Website QR Builder')}
            </p>
          </div>

          <div className="flex items-center justify-center rounded-lg border border-dashed border-gray-300 bg-gray-50 p-4 dark:border-gray-700 dark:bg-gray-900">
            <QRPreview
              value={targetUrl}
              size={qrSize}
              fgColor={fgColor}
              bgColor={bgColor}
              backgroundMode={backgroundMode}
              frameShape={frameShape}
              stylePreset={stylePreset}
              logoUrl={logoUrl}
              tableNumber={1}
              tableName={resolvedLabelSubtitle}
              errorLevel={errorLevel}
              includeMargin={includeMargin}
              showLabel={labelMode !== 'none'}
              labelTitle={resolvedLabelTitle}
              labelSubtitle={resolvedLabelSubtitle}
              onCanvasReady={(canvas) => {
                qrCanvasRef.current = canvas
              }}
              onSvgReady={(svg) => {
                qrSvgRef.current = svg
              }}
              t={t}
            />
          </div>
        </aside>
      </div>
    </div>
  )
}
