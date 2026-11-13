
import type { JSX } from 'react'
﻿/**
 * TestimonialsEditor Component
 * Specialized editor for testimonials/reviews carousel sections
 */

'use client'

import React, { useState, useCallback } from 'react'
import type { ThemeComponent } from '@/types/theme'
import { previewText } from '@/lib/previewText'

interface TestimonialItem {
  id: string
  author: string
  authorTitle?: string
  content: string
  rating: number
  avatarUrl?: string
  date?: string
  order: number
}

interface TestimonialsConfig {
  layout?: 'carousel' | 'grid' | 'list'
  itemsToShow?: number
  autoRotate?: boolean
  autoRotateInterval?: number
  showRatings?: boolean
  showDates?: boolean
  showAvatars?: boolean
  showAuthorTitle?: boolean
  columns?: 1 | 2 | 3
  cardBackgroundColor?: string
  cardBorderColor?: string
  titleColor?: string
  contentColor?: string
  authorColor?: string
  ratingColor?: string
  carouselButtonColor?: string
  testimonials?: TestimonialItem[]
  borderRadius?: number
  spacing?: number
  cardStyle?: 'minimal' | 'shadow' | 'border'
  hoverEffect?: 'lift' | 'glow' | 'color-change' | 'none'
  showQuoteIcon?: boolean
  quoteIconSize?: number
}

/**
 * TestimonialsEditor Component
 */
export function TestimonialsEditor({
  component,
  onChange,
  onPreview,
  className = '',
}: {
  component: ThemeComponent
  onChange: (component: ThemeComponent) => void
  onPreview?: (component: ThemeComponent) => void
  className?: string
}): JSX.Element {
  const [config, setConfig] = useState<TestimonialsConfig>(
    (component.config as unknown as TestimonialsConfig) || {
      layout: 'carousel',
      testimonials: [],
    }
  )

  const [previewMode, setPreviewMode] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)

  const handleConfigChange = useCallback(
    (newConfig: Partial<TestimonialsConfig>) => {
      const updatedConfig = { ...config, ...newConfig }
      setConfig(updatedConfig)

      const updatedComponent: ThemeComponent = {
        ...component,
        config: updatedConfig,
      }

      onChange(updatedComponent)

      if (onPreview) {
        onPreview(updatedComponent)
      }
    },
    [config, component, onChange, onPreview]
  )

  const addTestimonial = () => {
    const newTestimonial: TestimonialItem = {
      id: `testimonial-${Date.now()}`,
      author: 'Customer Name',
      authorTitle: 'Job Title',
      content: 'Great experience! Highly recommended.',
      rating: 5,
      order: (config.testimonials || []).length,
    }

    const updated = [...(config.testimonials || []), newTestimonial]
    handleConfigChange({ testimonials: updated })
  }

  const updateTestimonial = (id: string, updates: Partial<TestimonialItem>) => {
    const updated = (config.testimonials || []).map((t) =>
      t.id === id ? { ...t, ...updates } : t
    )
    handleConfigChange({ testimonials: updated })
  }

  const deleteTestimonial = (id: string) => {
    const updated = (config.testimonials || []).filter((t) => t.id !== id)
    handleConfigChange({ testimonials: updated })
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Tabs */}
      <div className="flex gap-2 border-b border-gray-200 dark:border-gray-700">
        <button
          onClick={() => setPreviewMode(false)}
          className={`px-4 py-2 font-medium transition-colors ${
            !previewMode
              ? 'border-b-2 border-blue-600 text-blue-600'
              : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
          }`}
        >
          Settings
        </button>
        <button
          onClick={() => setPreviewMode(true)}
          className={`px-4 py-2 font-medium transition-colors ${
            previewMode
              ? 'border-b-2 border-blue-600 text-blue-600'
              : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
          }`}
        >
          Preview
        </button>
      </div>

      {!previewMode && (
        <div className="space-y-6">
          {/* Layout Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Layout Settings</h3>

            {/* Layout Type */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Layout Type
              </label>
              <select
                value={config.layout || 'carousel'}
                onChange={(e) =>
                  handleConfigChange({ layout: e.target.value as 'carousel' | 'grid' | 'list' })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="carousel">Carousel (One at a time)</option>
                <option value="grid">Grid</option>
                <option value="list">List</option>
              </select>
            </div>

            {/* Columns for Grid */}
            {config.layout === 'grid' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Columns
                </label>
                <div className="flex gap-2">
                  {[1, 2, 3].map((col) => (
                    <button
                      key={col}
                      onClick={() => handleConfigChange({ columns: col as 1 | 2 | 3 })}
                      className={`flex-1 px-4 py-2 rounded-lg font-medium transition-colors ${
                        (config.columns || 1) === col
                          ? 'bg-blue-600 text-white'
                          : 'bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600'
                      }`}
                    >
                      {col}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {/* Auto Rotate (Carousel only) */}
            {config.layout === 'carousel' && (
              <>
                <label className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={config.autoRotate ?? true}
                    onChange={(e) => handleConfigChange({ autoRotate: e.target.checked })}
                    className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
                  />
                  <span className="text-sm text-gray-700 dark:text-gray-300">
                    Auto-rotate testimonials
                  </span>
                </label>

                {config.autoRotate && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Rotation Speed (seconds)
                    </label>
                    <div className="flex items-center gap-3">
                      <input
                        type="range"
                        min="2"
                        max="10"
                        step="1"
                        value={(config.autoRotateInterval || 5) / 1000}
                        onChange={(e) =>
                          handleConfigChange({ autoRotateInterval: parseInt(e.target.value) * 1000 })
                        }
                        className="flex-1"
                      />
                      <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                        {(config.autoRotateInterval || 5000) / 1000}s
                      </span>
                    </div>
                  </div>
                )}
              </>
            )}

            {/* Spacing */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Spacing (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="0"
                  max="40"
                  step="4"
                  value={config.spacing || 16}
                  onChange={(e) => handleConfigChange({ spacing: parseInt(e.target.value) })}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {config.spacing || 16}px
                </span>
              </div>
            </div>
          </div>

          {/* Display Options */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-3">
            <h3 className="font-semibold text-gray-900 dark:text-white mb-3">Display Options</h3>

            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showRatings ?? true}
                onChange={(e) => handleConfigChange({ showRatings: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Star Ratings</span>
            </label>

            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showAvatars ?? true}
                onChange={(e) => handleConfigChange({ showAvatars: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Avatar Images</span>
            </label>

            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showAuthorTitle ?? true}
                onChange={(e) => handleConfigChange({ showAuthorTitle: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Author Title</span>
            </label>

            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showDates ?? false}
                onChange={(e) => handleConfigChange({ showDates: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Dates</span>
            </label>

            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={config.showQuoteIcon ?? true}
                onChange={(e) => handleConfigChange({ showQuoteIcon: e.target.checked })}
                className="w-4 h-4 rounded border-gray-300 dark:border-gray-600"
              />
              <span className="text-sm text-gray-700 dark:text-gray-300">Show Quote Icon</span>
            </label>
          </div>

          {/* Style Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Card Style</h3>

            {/* Card Style */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Style Type
              </label>
              <select
                value={config.cardStyle || 'shadow'}
                onChange={(e) =>
                  handleConfigChange({ cardStyle: e.target.value as 'minimal' | 'shadow' | 'border' })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="minimal">Minimal</option>
                <option value="shadow">Shadow</option>
                <option value="border">Border</option>
              </select>
            </div>

            {/* Hover Effect */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Hover Effect
              </label>
              <select
                value={config.hoverEffect || 'lift'}
                onChange={(e) =>
                  handleConfigChange({
                    hoverEffect: e.target.value as 'lift' | 'glow' | 'color-change' | 'none',
                  })
                }
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                  bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                  focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="none">None</option>
                <option value="lift">Lift</option>
                <option value="glow">Glow</option>
                <option value="color-change">Color Change</option>
              </select>
            </div>

            {/* Border Radius */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Border Radius (px)
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="range"
                  min="0"
                  max="20"
                  step="2"
                  value={config.borderRadius || 8}
                  onChange={(e) => handleConfigChange({ borderRadius: parseInt(e.target.value) })}
                  className="flex-1"
                />
                <span className="text-sm font-medium text-gray-600 dark:text-gray-400 w-12">
                  {config.borderRadius || 8}px
                </span>
              </div>
            </div>
          </div>

          {/* Colors Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <h3 className="font-semibold text-gray-900 dark:text-white">Colors</h3>

            {/* Card Background */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Card Background Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.cardBackgroundColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ cardBackgroundColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.cardBackgroundColor || '#ffffff'}
                  onChange={(e) => handleConfigChange({ cardBackgroundColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Content Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Content Text Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.contentColor || '#6b7280'}
                  onChange={(e) => handleConfigChange({ contentColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.contentColor || '#6b7280'}
                  onChange={(e) => handleConfigChange({ contentColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Author Color */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Author Name Color
              </label>
              <div className="flex items-center gap-3">
                <input
                  type="color"
                  value={config.authorColor || '#1f2937'}
                  onChange={(e) => handleConfigChange({ authorColor: e.target.value })}
                  className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                />
                <input
                  type="text"
                  value={config.authorColor || '#1f2937'}
                  onChange={(e) => handleConfigChange({ authorColor: e.target.value })}
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                    bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                    focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            {/* Rating Color */}
            {config.showRatings && (
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Rating Star Color
                </label>
                <div className="flex items-center gap-3">
                  <input
                    type="color"
                    value={config.ratingColor || '#fbbf24'}
                    onChange={(e) => handleConfigChange({ ratingColor: e.target.value })}
                    className="h-10 w-14 rounded-lg border border-gray-300 dark:border-gray-600 cursor-pointer"
                  />
                  <input
                    type="text"
                    value={config.ratingColor || '#fbbf24'}
                    onChange={(e) => handleConfigChange({ ratingColor: e.target.value })}
                    className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg
                      bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                      focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>
            )}
          </div>

          {/* Testimonials Section */}
          <div className="rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="font-semibold text-gray-900 dark:text-white">Testimonials</h3>
              <button
                onClick={addTestimonial}
                className="px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-medium transition-colors"
              >
                + Add Testimonial
              </button>
            </div>

            <div className="space-y-3 max-h-96 overflow-y-auto">
              {(config.testimonials || []).map((testimonial) => (
                <div
                  key={testimonial.id}
                  className="border border-gray-300 dark:border-gray-600 rounded-lg p-3 space-y-2"
                >
                  {editingId === testimonial.id ? (
                    <>
                      <input
                        type="text"
                        value={testimonial.author}
                        onChange={(e) => updateTestimonial(testimonial.id, { author: e.target.value })}
                        placeholder="Author name"
                        className="w-full px-2 py-1 border border-gray-300 dark:border-gray-600 rounded
                          bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm
                          focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <input
                        type="text"
                        value={testimonial.authorTitle || ''}
                        onChange={(e) =>
                          updateTestimonial(testimonial.id, { authorTitle: e.target.value })
                        }
                        placeholder="Job title (optional)"
                        className="w-full px-2 py-1 border border-gray-300 dark:border-gray-600 rounded
                          bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm
                          focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <textarea
                        value={testimonial.content}
                        onChange={(e) =>
                          updateTestimonial(testimonial.id, { content: e.target.value })
                        }
                        placeholder="Testimonial content"
                        rows={2}
                        className="w-full px-2 py-1 border border-gray-300 dark:border-gray-600 rounded
                          bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm
                          focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <div>
                        <label className="text-xs font-medium text-gray-600 dark:text-gray-400">
                          Rating
                        </label>
                        <div className="flex gap-1">
                          {[1, 2, 3, 4, 5].map((star) => (
                            <button
                              key={star}
                              onClick={() => updateTestimonial(testimonial.id, { rating: star })}
                              className={`text-xl ${
                                star <= testimonial.rating
                                  ? 'text-yellow-400'
                                  : 'text-gray-300 dark:text-gray-600'
                              }`}
                            >
                              Ã¢Ëœâ€¦
                            </button>
                          ))}
                        </div>
                      </div>
                      <input
                        type="url"
                        value={testimonial.avatarUrl || ''}
                        onChange={(e) =>
                          updateTestimonial(testimonial.id, { avatarUrl: e.target.value })
                        }
                        placeholder="Avatar image URL"
                        className="w-full px-2 py-1 border border-gray-300 dark:border-gray-600 rounded
                          bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm
                          focus:outline-none focus:ring-2 focus:ring-blue-500"
                      />
                      <div className="flex gap-2">
                        <button
                          onClick={() => setEditingId(null)}
                          className="flex-1 px-3 py-1 rounded text-sm bg-blue-600 hover:bg-blue-700 text-white font-medium transition-colors"
                        >
                          Done
                        </button>
                        <button
                          onClick={() => deleteTestimonial(testimonial.id)}
                          className="flex-1 px-3 py-1 rounded text-sm bg-red-600 hover:bg-red-700 text-white font-medium transition-colors"
                        >
                          Delete
                        </button>
                      </div>
                    </>
                  ) : (
                    <>
                      <div className="flex items-start gap-3">
                        {testimonial.avatarUrl && (
                          <img
                            src={testimonial.avatarUrl}
                            alt={testimonial.author}
                            className="w-8 h-8 rounded-full object-cover"
                            onError={(e) => {
                              (e.target as HTMLImageElement).style.display = 'none'
                            }}
                          />
                        )}
                        <div className="flex-1">
                          <p className="font-medium text-gray-900 dark:text-white">
                            {testimonial.author}
                          </p>
                          {testimonial.authorTitle && (
                            <p className="text-xs text-gray-600 dark:text-gray-400">
                              {testimonial.authorTitle}
                            </p>
                          )}
                          <p className="text-sm text-gray-700 dark:text-gray-300 mt-1">
                            {testimonial.content.substring(0, 50)}...
                          </p>
                          <div className="flex gap-0.5 mt-1">
                            {[1, 2, 3, 4, 5].map((star) => (
                              <span
                                key={star}
                                className={`text-xs ${
                                  star <= testimonial.rating ? 'text-yellow-400' : 'text-gray-300'
                                }`}
                              >
                                Ã¢Ëœâ€¦
                              </span>
                            ))}
                          </div>
                        </div>
                        <button
                          onClick={() => setEditingId(testimonial.id)}
                          className="px-3 py-1 rounded text-sm bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
                        >
                          Edit
                        </button>
                      </div>
                    </>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {previewMode && (
        <div className="rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden bg-gray-50 dark:bg-gray-900 p-8">
          <div className="max-w-6xl mx-auto testimonials-preview">
            <h2 className="text-3xl font-bold text-gray-900 dark:text-white text-center mb-12">
              {previewText(component.title, 'What Customers Say')}
            </h2>

            {config.layout === 'carousel' ? (
              <div className="bg-white dark:bg-gray-800 rounded-lg p-8 shadow-lg">
                {(config.testimonials || []).length > 0 && (
                  <>
                    {config.showQuoteIcon && (
                      <div className="text-5xl text-gray-300 dark:text-gray-600 mb-4">&quot;</div>
                    )}
                    <p className="testimonial-content text-lg mb-6 italic">
                      {(config.testimonials || [])[0]?.content}
                    </p>
                    <div className="flex items-center gap-4">
                      {config.showAvatars && (config.testimonials || [])[0]?.avatarUrl && (
                        <img
                          src={(config.testimonials || [])[0].avatarUrl}
                          alt={(config.testimonials || [])[0].author}
                          className="w-12 h-12 rounded-full object-cover"
                          onError={(e) => {
                            e.currentTarget.classList.add('hidden')
                          }}
                        />
                      )}
                      <div>
                        <p className="testimonial-author font-semibold">
                          {(config.testimonials || [])[0]?.author}
                        </p>
                        {config.showAuthorTitle && (config.testimonials || [])[0]?.authorTitle && (
                          <p className="text-sm text-gray-600 dark:text-gray-400">
                            {(config.testimonials || [])[0].authorTitle}
                          </p>
                        )}
                      </div>
                    </div>
                    {config.showRatings && (
                      <div className="flex gap-1 mt-4">
                        {[1, 2, 3, 4, 5].map((star) => (
                          <span
                            key={star}
                            className="testimonial-rating text-xl"
                          >
                            {star <= ((config.testimonials || [])[0]?.rating || 5) ? 'â˜…' : 'â˜†'}
                          </span>
                        ))}
                      </div>
                    )}
                  </>
                )}
              </div>
            ) : (
              <div className="testimonials-preview-grid">
                {(config.testimonials || []).map((testimonial) => (
                  <div
                    key={testimonial.id}
                    className="testimonials-preview-card p-6"
                  >
                    {config.showQuoteIcon && (
                      <div className="text-4xl text-gray-300 dark:text-gray-600 mb-3">&quot;</div>
                    )}
                    <p className="testimonial-content mb-4 italic">
                      {testimonial.content}
                    </p>
                    <div className="flex items-center gap-3">
                      {config.showAvatars && testimonial.avatarUrl && (
                        <img
                          src={testimonial.avatarUrl}
                          alt={testimonial.author}
                          className="w-10 h-10 rounded-full object-cover"
                          onError={(e) => {
                            e.currentTarget.classList.add('hidden')
                          }}
                        />
                      )}
                      <div>
                        <p className="testimonial-author font-semibold">
                          {testimonial.author}
                        </p>
                        {config.showAuthorTitle && testimonial.authorTitle && (
                          <p className="text-xs text-gray-600 dark:text-gray-400">
                            {testimonial.authorTitle}
                          </p>
                        )}
                      </div>
                    </div>
                    {config.showRatings && (
                      <div className="flex gap-1 mt-3">
                        {[1, 2, 3, 4, 5].map((star) => (
                          <span
                            key={star}
                            className="testimonial-rating text-lg"
                          >
                            {star <= testimonial.rating ? 'â˜…' : 'â˜†'}
                          </span>
                        ))}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}
            <style jsx>{`
              .testimonials-preview-grid {
                display: grid;
                grid-template-columns: repeat(${config.columns || 1}, 1fr);
                gap: ${config.spacing || 16}px;
              }

              .testimonials-preview-card {
                background-color: ${config.cardBackgroundColor || '#ffffff'};
                border-radius: ${config.borderRadius || 8}px;
                ${config.cardStyle === 'shadow' ? 'box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);' : ''}
                ${config.cardStyle === 'border' ? 'border: 1px solid #e5e7eb;' : ''}
              }

              .testimonial-content {
                color: ${config.contentColor || '#6b7280'};
              }

              .testimonial-author {
                color: ${config.authorColor || '#1f2937'};
              }

              .testimonial-rating {
                color: ${config.ratingColor || '#fbbf24'};
              }
            `}</style>
          </div>
        </div>
      )}
    </div>
  )
}

export default TestimonialsEditor
