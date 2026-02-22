'use client'

import React from 'react'

export interface FeaturedItem {
  id: string | number
  title: string
  description: string
  price?: string
}

export interface FeaturedItemsProps {
  title: string
  description?: string
  items?: FeaturedItem[]
  primaryColor: string
  accentColor: string
}

export function FeaturedItems({
  title,
  description,
  items = [
    { id: 1, title: 'Dish Name', description: 'Delicious and authentic dish' },
    { id: 2, title: 'Dish Name', description: 'Delicious and authentic dish' },
    { id: 3, title: 'Dish Name', description: 'Delicious and authentic dish' },
  ],
  primaryColor,
  accentColor,
}: FeaturedItemsProps) {
  const displayItems =
    items.length > 0
      ? items
      : [
          { id: 1, title: 'Dish Name', description: 'Delicious and authentic dish' },
          { id: 2, title: 'Dish Name', description: 'Delicious and authentic dish' },
          { id: 3, title: 'Dish Name', description: 'Delicious and authentic dish' },
        ]

  return (
    <div className="p-12 md:p-16 bg-gray-50 dark:bg-gray-800">
      <h2 className="text-3xl md:text-4xl font-bold mb-8 text-gray-900 dark:text-white">{title}</h2>
      {description && <p className="text-gray-600 dark:text-gray-400 mb-6">{description}</p>}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {displayItems.map((item) => (
          <div
            key={item.id}
            className="bg-white dark:bg-gray-700 rounded-lg overflow-hidden shadow hover:shadow-lg transition-shadow"
          >
            <div className="featured-item-hero h-32 flex items-center justify-center text-white text-sm font-semibold p-4">
              Featured Product
            </div>
            <div className="p-4">
              <h3 className="font-semibold text-gray-900 dark:text-white mb-2">{item.title}</h3>
              <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">{item.description}</p>
              {item.price && <p className="text-lg font-bold text-gray-900 dark:text-white mb-4">{item.price}</p>}
              <button className="featured-item-button w-full py-2 rounded text-sm font-medium hover:opacity-90 transition-opacity">
                View Details
              </button>
            </div>
          </div>
        ))}
      </div>

      <style jsx>{`
        .featured-item-hero {
          background-color: ${primaryColor};
        }

        .featured-item-button {
          background-color: ${accentColor};
          color: #000000;
        }
      `}</style>
    </div>
  )
}
