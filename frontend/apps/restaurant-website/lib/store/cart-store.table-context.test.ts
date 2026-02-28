import { act, renderHook } from '@testing-library/react'
import { useCartStore } from './cart-store'

describe('cart-store table context', () => {
  beforeEach(() => {
    localStorage.clear()
    useCartStore.setState({
      items: [],
      tableId: null,
      tableName: null,
    })
  })

  it('sets table context with table id and trimmed name', () => {
    const { result } = renderHook(() => useCartStore())

    act(() => {
      result.current.setTableContext(17, '  VIP Window  ')
    })

    expect(result.current.tableId).toBe(17)
    expect(result.current.tableName).toBe('VIP Window')
  })

  it('clears table context without removing cart items', () => {
    const { result } = renderHook(() => useCartStore())

    act(() => {
      result.current.addItem({
        id: 10,
        name: 'Burger',
        price: 12,
      })
      result.current.setTableContext(3, 'Table 3')
    })

    expect(result.current.items.length).toBe(1)
    expect(result.current.tableId).toBe(3)

    act(() => {
      result.current.clearTableContext()
    })

    expect(result.current.tableId).toBeNull()
    expect(result.current.tableName).toBeNull()
    expect(result.current.items.length).toBe(1)
  })
})
