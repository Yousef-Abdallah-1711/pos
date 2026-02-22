'use client'

import React from 'react'

export interface WhyChooseItem {
  id: string | number
  title: string
  icon: string
  description?: string
}

export interface WhyChooseUsProps {
  title: string
  items?: WhyChooseItem[]
  primaryColor: string
}

export function WhyChooseUs({
  title,
  items = [
    { id: 1, title: 'Quality', icon: 'âœ¨', description: 'Premium quality ingredients' },
    { id: 2, title: 'Fast Delivery', icon: 'âš¡', description: 'Quick delivery service' },
    { id: 3, title: 'Wide Variety', icon: 'ðŸŽ¯', description: 'Wide menu selection' },
    { id: 4, title: '24/7 Support', icon: 'ðŸ›Ÿ', description: 'Always available' },
  ],
  primaryColor,
}: WhyChooseUsProps) {
  const displayItems =
    items.length > 0
      ? items
      : [
          { id: 1, title: 'Quality', icon: 'âœ¨', description: 'Premium quality ingredients' },
          { id: 2, title: 'Fast Delivery', icon: 'âš¡', description: 'Quick delivery service' },
          { id: 3, title: 'Wide Variety', icon: 'ðŸŽ¯', description: 'Wide menu selection' },
          { id: 4, title: '24/7 Support', icon: 'ðŸ›Ÿ', description: 'Always available' },
        ]

  return (
    <div className="p-12 md:p-16 bg-white dark:bg-gray-900">
      <h2 className="text-3xl md:text-4xl font-bold mb-8 text-gray-900 dark:text-white text-center">{title}</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
        {displayItems.map((item) => (
          <div key={item.id} className="text-center">
            <div className="text-5xl mb-4">{item.icon}</div>
            <h3 className="why-us-item-title font-semibold text-lg mb-2">{item.title}</h3>
            {item.description && <p className="text-sm text-gray-600 dark:text-gray-400">{item.description}</p>}
          </div>
        ))}
      </div>

      <style jsx>{`
        .why-us-item-title {
          color: ${primaryColor};
        }
      `}</style>
    </div>
  )
}
