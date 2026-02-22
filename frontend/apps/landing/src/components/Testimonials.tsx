'use client'

import { useLanguage } from '@/contexts/LanguageContext'
import { Star, Quote } from 'lucide-react'

const avatars = [
  'https://i.pravatar.cc/150?img=11',
  'https://i.pravatar.cc/150?img=5',
  'https://i.pravatar.cc/150?img=33',
]

export function Testimonials() {
  const { t, dir } = useLanguage()

  const testimonials = [
    {
      content: t.testimonials.content_1,
      name: t.testimonials.name_1,
      role: t.testimonials.role_1,
      avatar: avatars[0],
      rating: 5,
    },
    {
      content: t.testimonials.content_2,
      name: t.testimonials.name_2,
      role: t.testimonials.role_2,
      avatar: avatars[1],
      rating: 5,
    },
    {
      content: t.testimonials.content_3,
      name: t.testimonials.name_3,
      role: t.testimonials.role_3,
      avatar: avatars[2],
      rating: 5,
    },
  ]

  return (
    <section id="testimonials" className="py-24 px-4 sm:px-6 lg:px-8 bg-gray-50 dark:bg-gray-800/50" dir={dir}>
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="text-center mb-16">
          <span className="inline-block text-sm font-semibold text-orange-600 dark:text-orange-400 uppercase tracking-wider mb-3">
            Testimonials
          </span>
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-black text-gray-900 dark:text-white mb-4 tracking-tight">
            {t.testimonials.title}
          </h2>
          <p className="text-lg text-gray-500 dark:text-gray-400 max-w-2xl mx-auto">
            {t.testimonials.subtitle}
          </p>
        </div>

        {/* Testimonials Grid */}
        <div className="grid md:grid-cols-3 gap-6">
          {testimonials.map((testimonial, index) => (
            <div
              key={index}
              className="bg-white dark:bg-gray-800 rounded-2xl p-6 border border-gray-200 dark:border-gray-700 hover:shadow-xl transition-all duration-300 flex flex-col"
            >
              {/* Quote icon */}
              <div className="w-10 h-10 rounded-xl bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center mb-4">
                <Quote className="h-5 w-5 text-orange-500" />
              </div>

              {/* Stars */}
              <div className="flex gap-0.5 mb-4">
                {[...Array(testimonial.rating)].map((_, i) => (
                  <Star key={i} className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                ))}
              </div>

              {/* Content */}
              <p className="text-gray-700 dark:text-gray-300 text-sm leading-relaxed flex-1 mb-6">
                &ldquo;{testimonial.content}&rdquo;
              </p>

              {/* Author */}
              <div className="flex items-center gap-3 pt-4 border-t border-gray-100 dark:border-gray-700">
                <img
                  src={testimonial.avatar}
                  alt={testimonial.name}
                  className="w-10 h-10 rounded-full object-cover ring-2 ring-orange-100 dark:ring-orange-900/30"
                />
                <div>
                  <div className="font-bold text-sm text-gray-900 dark:text-white">
                    {testimonial.name}
                  </div>
                  <div className="text-xs text-gray-500 dark:text-gray-400">
                    {testimonial.role}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
