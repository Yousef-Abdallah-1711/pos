'use client'

import {
  Component,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ChangeEvent,
  type ReactNode,
} from 'react'
import { Copy, Download, ExternalLink, FileText, Loader2, QrCode, X } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { type RestaurantTable } from '@/lib/api/tablesApi'
import { debugError } from '@/lib/debug'
import { buildTableMenuUrl, sanitizeQrFileName } from '@/lib/table-qr'
import {
  QRPreview,
  type QrFrameShape,
  type QrStylePreset,
} from './QRPreview'

interface QRBuilderModalProps {
  isOpen: boolean
  table: RestaurantTable | null
  locale: string
  onClose: () => void
  t: (path: string, fallback?: string) => string
  isRTL: boolean
}

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

type ProcessingAction = 'png' | 'svg' | 'pdf' | null

interface QRPreviewBoundaryProps {
  children: ReactNode
  fallback: ReactNode
}

interface QRPreviewBoundaryState {
  hasError: boolean
}

class QRPreviewBoundary extends Component<
  QRPreviewBoundaryProps,
  QRPreviewBoundaryState
> {
  state: QRPreviewBoundaryState = {
    hasError: false,
  }

  static getDerivedStateFromError(): QRPreviewBoundaryState {
    return {
      hasError: true,
    }
  }

  componentDidUpdate(prevProps: QRPreviewBoundaryProps) {
    if (this.state.hasError && prevProps.children !== this.props.children) {
      this.setState({ hasError: false })
    }
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback
    }

    return this.props.children
  }
}

export function QRBuilderModal({
  isOpen,
  table,
  locale,
  onClose,
  t,
  isRTL,
}: QRBuilderModalProps) {
  const [qrSize, setQrSize] = useState<number>(256)
  const [frameShape, setFrameShape] = useState<QrFrameShape>('square')
  const [stylePreset, setStylePreset] = useState<QrStylePreset>('standard')
  const [fgColor, setFgColor] = useState('#111827')
  const [bgColor, setBgColor] = useState('#FFFFFF')
  const [backgroundMode, setBackgroundMode] = useState<'solid' | 'transparent'>('solid')
  const [logoUrl, setLogoUrl] = useState<string | null>(null)
  const [processingAction, setProcessingAction] = useState<ProcessingAction>(null)
  const fileInputRef = useRef<HTMLInputElement | null>(null)
  const qrCanvasRef = useRef<HTMLCanvasElement | null>(null)
  const qrSvgRef = useRef<SVGSVGElement | null>(null)

  useEffect(() => {
    if (!isOpen || !table) {
      return
    }

    setQrSize(256)
    setFrameShape('square')
    setStylePreset('standard')
    setFgColor('#111827')
    setBgColor('#FFFFFF')
    setBackgroundMode('solid')
    setLogoUrl(null)
    setProcessingAction(null)
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }, [isOpen, table?.id])

  const tableUrl = useMemo(() => {
    if (!table) {
      return ''
    }
    return buildTableMenuUrl(table.id, locale)
  }, [table, locale])

  if (!isOpen || !table) {
    return null
  }

  const baseFileName = sanitizeQrFileName(`table-${table.table_number}-qr`)

  const handleCopyLink = async () => {
    if (!tableUrl) {
      return
    }

    try {
      await navigator.clipboard.writeText(tableUrl)
      toast.success(t('tables.linkCopied'))
    } catch (error) {
      toast.error(t('tables.copyFailed'))
      debugError('[QR_BUILDER] Copy link failed:', error)
    }
  }

  const handleDownloadPng = async () => {
    const canvas = qrCanvasRef.current
    if (!canvas) {
      toast.error(t('tables.downloadFailed'))
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
      toast.error(t('tables.downloadFailed'))
      debugError('[QR_BUILDER] PNG download failed:', error)
    } finally {
      setProcessingAction(null)
    }
  }

  const handleDownloadSvg = async () => {
    const svgElement = qrSvgRef.current
    if (!svgElement) {
      toast.error(t('tables.downloadFailed'))
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
      toast.error(t('tables.downloadFailed'))
      debugError('[QR_BUILDER] SVG download failed:', error)
    } finally {
      setProcessingAction(null)
    }
  }

  const handlePrintPdf = async () => {
    const canvas = qrCanvasRef.current
    if (!canvas) {
      toast.error(t('tables.downloadFailed'))
      return
    }

    setProcessingAction('pdf')
    try {
      const { jsPDF } = await import('jspdf')
      const doc = new jsPDF({
        orientation: 'portrait',
        unit: 'mm',
        format: 'a4',
      })

      const pageWidth = doc.internal.pageSize.getWidth()
      const qrWidth = 90
      const qrHeight = 90
      const qrX = (pageWidth - qrWidth) / 2
      const qrY = 45

      doc.setFontSize(16)
      doc.text(`${t('tables.qrBuilder')} - ${t('tables.tableNumber')} ${table.table_number}`, pageWidth / 2, 20, {
        align: 'center',
      })

      doc.setFontSize(11)
      doc.text(`${t('tables.tableName')}: ${table.name || t('common.notAvailable', 'N/A')}`, pageWidth / 2, 28, {
        align: 'center',
      })

      const qrPng = canvas.toDataURL('image/png')
      doc.addImage(qrPng, 'PNG', qrX, qrY, qrWidth, qrHeight)

      doc.setFontSize(10)
      const wrappedUrl = doc.splitTextToSize(tableUrl, pageWidth - 24)
      doc.text(`${t('tables.qrUrl')}:`, 12, 145)
      doc.text(wrappedUrl, 12, 151)

      doc.save(`${baseFileName}.pdf`)
    } catch (error) {
      toast.error(t('tables.downloadFailed'))
      debugError('[QR_BUILDER] PDF export failed:', error)
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
      debugError('[QR_BUILDER] Logo upload failed:', error)
    }
    reader.readAsDataURL(file)
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 px-4">
      <div className="w-full max-w-5xl rounded-lg bg-white p-6 shadow-xl dark:bg-gray-800">
        <div className="mb-4 flex items-start justify-between">
          <div>
            <h2 className="text-xl font-bold text-gray-900 dark:text-white">
              {t('tables.qrBuilder')}
            </h2>
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
              {t('tables.qrForTable')} #{table.table_number}
            </p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="rounded-md p-1 text-gray-500 transition hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-gray-700 dark:hover:text-gray-100"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="grid grid-cols-1 gap-6 lg:grid-cols-[1.3fr_1fr]">
          <div className="space-y-5 rounded-lg border border-gray-200 bg-gray-50 p-4 dark:border-gray-700 dark:bg-gray-900">
            <div>
              <p className="mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('tables.qrSize')}
              </p>
              <div className="flex flex-wrap gap-2">
                {QR_SIZE_OPTIONS.map((option) => (
                  <button
                    key={option.value}
                    type="button"
                    onClick={() => setQrSize(option.value)}
                    className={`rounded-lg px-3 py-2 text-sm font-medium transition ${
                      qrSize === option.value
                        ? 'bg-primary text-white'
                        : 'bg-white text-gray-700 hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-200 dark:hover:bg-gray-700'
                    }`}
                  >
                    {t(`tables.qrSize${option.label}`, option.label)} ({option.value}px)
                  </button>
                ))}
              </div>
            </div>

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
                        : 'bg-white text-gray-700 hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-200 dark:hover:bg-gray-700'
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
                        : 'bg-white text-gray-700 hover:bg-gray-100 dark:bg-gray-800 dark:text-gray-200 dark:hover:bg-gray-700'
                    }`}
                  >
                    {t(`tables.qrStyles.${style.value}`, style.label)}
                  </button>
                ))}
              </div>
            </div>

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
                    onClick={() => setLogoUrl(null)}
                  >
                    {t('tables.removeLogo', 'Remove Logo')}
                  </Button>
                )}
              </div>
            </div>

            <div>
              <p className="mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('tables.qrUrl')}
              </p>
              <div className="rounded-lg border border-gray-200 bg-white p-3 text-xs text-gray-700 break-all dark:border-gray-700 dark:bg-gray-800 dark:text-gray-200">
                {tableUrl}
              </div>
            </div>

            <div className={`flex flex-wrap gap-2 ${isRTL ? 'justify-start' : ''}`}>
              <Button variant="outline" onClick={handleCopyLink}>
                <Copy className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
                {t('tables.copyLink')}
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
                {t('tables.downloadPng')}
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
              <Button
                variant="outline"
                onClick={() => {
                  window.open(tableUrl, '_blank', 'noopener,noreferrer')
                }}
              >
                <ExternalLink className={isRTL ? 'ml-2 h-4 w-4' : 'mr-2 h-4 w-4'} />
                {t('tables.openLink')}
              </Button>
            </div>
          </div>

          <div className="flex items-center justify-center rounded-lg border border-dashed border-gray-300 bg-gray-50 p-6 dark:border-gray-700 dark:bg-gray-900">
            {tableUrl ? (
              <QRPreviewBoundary
                key={`${table.id}-${tableUrl}-${qrSize}-${stylePreset}-${frameShape}`}
                fallback={
                  <div className="text-center text-gray-500 dark:text-gray-400">
                    <QrCode className="mx-auto mb-2 h-10 w-10" />
                    {t('tables.qrUnavailable')}
                  </div>
                }
              >
                <QRPreview
                  value={tableUrl}
                  size={qrSize}
                  fgColor={fgColor}
                  bgColor={bgColor}
                  backgroundMode={backgroundMode}
                  frameShape={frameShape}
                  stylePreset={stylePreset}
                  logoUrl={logoUrl}
                  tableNumber={table.table_number}
                  tableName={table.name}
                  onCanvasReady={(canvas) => {
                    qrCanvasRef.current = canvas
                  }}
                  onSvgReady={(svg) => {
                    qrSvgRef.current = svg
                  }}
                  t={t}
                />
              </QRPreviewBoundary>
            ) : (
              <div className="text-center text-gray-500 dark:text-gray-400">
                <QrCode className="mx-auto mb-2 h-10 w-10" />
                {t('tables.qrUnavailable')}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
