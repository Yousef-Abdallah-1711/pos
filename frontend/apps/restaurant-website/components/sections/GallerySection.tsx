'use client'

import { useMemo, useState } from 'react'
import Image from 'next/image'
import { Dialog, DialogContent, DialogTitle } from '@/components/ui/dialog'
import type { ThemeLocale, ThemeSectionConfig } from '@/types/theme'
import { readLocalized, toArray, toRecord } from '@/components/sections/section-utils'

interface GallerySectionProps {
  locale: ThemeLocale
  sectionConfig?: ThemeSectionConfig
}

export function GallerySection({ locale, sectionConfig }: GallerySectionProps) {
  const settings = toRecord(sectionConfig?.settings)
  const title = readLocalized(sectionConfig?.title ?? settings.title, locale, locale === 'ar' ? 'معرض الصور' : 'Gallery')
  const images = useMemo(() => {
    const raw = toArray<string>(settings.images)
    return raw.filter((image) => typeof image === 'string' && image.length > 0).slice(0, 6)
  }, [settings])
  const [selectedImage, setSelectedImage] = useState<string | null>(null)

  if (images.length === 0) {
    return null
  }

  return (
    <section className="px-4 py-16 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-6xl">
        <h2 className={`mb-8 text-3xl font-bold text-[hsl(var(--theme-text))] ${locale === 'ar' ? 'text-right' : 'text-left'}`}>{title}</h2>
        <div className="grid grid-cols-2 gap-3 md:grid-cols-3">
          {images.map((image, index) => (
            <button
              key={`${image}-${index}`}
              type="button"
              className="group relative aspect-[4/3] overflow-hidden rounded-lg"
              onClick={() => setSelectedImage(image)}
            >
              <Image
                src={image}
                alt={`${title} ${index + 1}`}
                fill
                className="object-cover transition-transform duration-300 group-hover:scale-105"
                sizes="(max-width: 768px) 50vw, 33vw"
              />
              <span className="absolute inset-0 bg-black/0 transition-colors group-hover:bg-black/20" />
            </button>
          ))}
        </div>
      </div>

      <Dialog open={Boolean(selectedImage)} onOpenChange={(open) => !open && setSelectedImage(null)}>
        <DialogContent className="max-w-4xl border-none bg-transparent p-0 shadow-none">
          <DialogTitle className="sr-only">{title}</DialogTitle>
          {selectedImage ? (
            <div className="relative h-[70vh] w-full">
              <Image src={selectedImage} alt={title} fill className="object-contain" sizes="100vw" />
            </div>
          ) : null}
        </DialogContent>
      </Dialog>
    </section>
  )
}

