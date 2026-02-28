'use client'

import { useEffect, useMemo, useRef } from 'react'
import { QRCodeCanvas, QRCodeSVG } from 'qrcode.react'

type QrFrameShape = 'square' | 'circle' | 'diamond'
type QrStylePreset = 'standard' | 'rounded' | 'dots'

interface QRPreviewProps {
  value: string
  size: number
  fgColor: string
  bgColor: string
  backgroundMode: 'solid' | 'transparent'
  frameShape: QrFrameShape
  stylePreset: QrStylePreset
  logoUrl: string | null
  tableNumber: number
  tableName?: string
  errorLevel?: 'L' | 'M' | 'Q' | 'H'
  includeMargin?: boolean
  showLabel?: boolean
  labelTitle?: string
  labelSubtitle?: string
  onCanvasReady?: (canvas: HTMLCanvasElement | null) => void
  onSvgReady?: (svg: SVGSVGElement | null) => void
  t: (path: string, fallback?: string) => string
}

function getFrameClass(shape: QrFrameShape) {
  if (shape === 'circle') {
    return 'rounded-full'
  }
  if (shape === 'diamond') {
    return 'rotate-45 rounded-2xl'
  }
  return 'rounded-xl'
}

function getQrBorderRadius(stylePreset: QrStylePreset) {
  if (stylePreset === 'rounded') {
    return 16
  }
  if (stylePreset === 'dots') {
    return 999
  }
  return 0
}

export function QRPreview({
  value,
  size,
  fgColor,
  bgColor,
  backgroundMode,
  frameShape,
  stylePreset,
  logoUrl,
  tableNumber,
  tableName,
  errorLevel = 'H',
  includeMargin = true,
  showLabel = true,
  labelTitle,
  labelSubtitle,
  onCanvasReady,
  onSvgReady,
  t,
}: QRPreviewProps) {
  const canvasContainerRef = useRef<HTMLDivElement | null>(null)
  const svgContainerRef = useRef<HTMLDivElement | null>(null)

  const effectiveBgColor = backgroundMode === 'transparent' ? 'transparent' : bgColor
  const qrBorderRadius = getQrBorderRadius(stylePreset)
  const frameClass = getFrameClass(frameShape)
  const innerRotationClass = frameShape === 'diamond' ? '-rotate-45' : ''
  const resolvedTitle = labelTitle || `${t('tables.tableNumber')} ${tableNumber}`
  const resolvedSubtitle = labelSubtitle || tableName || t('tables.notAvailable', 'N/A')

  const imageSettings = useMemo(() => {
    if (!logoUrl) {
      return undefined
    }

    const logoSize = Math.max(20, Math.floor(size * 0.2))
    return {
      src: logoUrl,
      height: logoSize,
      width: logoSize,
      excavate: true,
    }
  }, [logoUrl, size])

  useEffect(() => {
    if (!onCanvasReady) {
      return
    }
    const canvas = canvasContainerRef.current?.querySelector('canvas') || null
    onCanvasReady(canvas)
  }, [
    backgroundMode,
    bgColor,
    fgColor,
    frameShape,
    imageSettings,
    onCanvasReady,
    size,
    stylePreset,
    value,
  ])

  useEffect(() => {
    if (!onSvgReady) {
      return
    }
    const svg = svgContainerRef.current?.querySelector('svg') || null
    onSvgReady(svg)
  }, [
    backgroundMode,
    bgColor,
    fgColor,
    frameShape,
    imageSettings,
    onSvgReady,
    size,
    stylePreset,
    value,
  ])

  return (
    <div className="flex flex-col items-center">
      <div
        className={`relative flex h-[300px] w-[300px] items-center justify-center border border-gray-200 bg-gray-100 p-4 shadow-sm dark:border-gray-700 dark:bg-gray-900 ${frameClass}`}
      >
        <div className={`${innerRotationClass} flex items-center justify-center`}>
          <div ref={canvasContainerRef}>
            <QRCodeCanvas
              value={value}
              size={size}
              fgColor={fgColor}
              bgColor={effectiveBgColor}
              level={errorLevel}
              includeMargin={includeMargin}
              imageSettings={imageSettings}
              style={{
                width: 240,
                height: 240,
                borderRadius: qrBorderRadius,
                backgroundColor: backgroundMode === 'transparent' ? 'transparent' : bgColor,
              }}
            />
          </div>
        </div>
      </div>

      {showLabel && (
        <>
          <p className="mt-3 text-sm font-semibold text-gray-900 dark:text-white">
            {resolvedTitle}
          </p>
          <p className="text-xs text-gray-600 dark:text-gray-400">
            {resolvedSubtitle}
          </p>
        </>
      )}

      <div ref={svgContainerRef} className="hidden">
        <QRCodeSVG
          value={value}
          size={size}
          fgColor={fgColor}
          bgColor={effectiveBgColor}
          level={errorLevel}
          includeMargin={includeMargin}
          imageSettings={imageSettings}
          style={{
            borderRadius: qrBorderRadius,
            backgroundColor: backgroundMode === 'transparent' ? 'transparent' : bgColor,
          }}
        />
      </div>
    </div>
  )
}

export type { QrFrameShape, QrStylePreset }
