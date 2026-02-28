import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
  },
}))

import { api } from '@/lib/api'
import { posApi } from './posApi'

const mockedApi = api as unknown as {
  get: ReturnType<typeof vi.fn>
  post: ReturnType<typeof vi.fn>
  put: ReturnType<typeof vi.fn>
}

describe('posApi', () => {
  beforeEach(() => {
    mockedApi.get.mockReset()
    mockedApi.post.mockReset()
    mockedApi.put.mockReset()
  })

  it('maps dashboard product payload fields to POS shape', async () => {
    mockedApi.get.mockResolvedValue({
      data: {
        products: [
          {
            id: 42,
            name_en: 'Grilled Chicken',
            name_ar: 'دجاج مشوي',
            main_image_url: 'https://cdn.local/chicken.png',
            category_id: 3,
            price: 14.99,
            is_available: true,
            quantity_in_stock: 5,
            track_inventory: true,
            barcode: 123456,
          },
        ],
      },
    })

    const products = await posApi.getProducts()

    expect(mockedApi.get).toHaveBeenCalledWith('/products', {
      params: { page_size: 500 },
    })
    expect(products).toHaveLength(1)
    expect(products[0]).toEqual(
      expect.objectContaining({
        id: 42,
        name: 'Grilled Chicken',
        name_ar: 'دجاج مشوي',
        image: 'https://cdn.local/chicken.png',
        category_id: 3,
        price: 14.99,
        barcode: '123456',
      })
    )
  })

  it('ignores invalid products missing both name and name_en', async () => {
    mockedApi.get.mockResolvedValue({
      data: {
        products: [
          { id: 1, price: 10 },
          { id: 2, name: 'Valid Product', price: 12 },
        ],
      },
    })

    const products = await posApi.getProducts()
    expect(products).toHaveLength(1)
    expect(products[0].id).toBe(2)
  })
})
