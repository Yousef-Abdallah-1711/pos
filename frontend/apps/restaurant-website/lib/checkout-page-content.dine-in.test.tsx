import type { OrderQuote } from './order-quote'
import { fireEvent, render, waitFor } from '@testing-library/react'
import { CheckoutPageContent } from '@/app/[locale]/checkout/checkout-page-content'
import { useCartStore } from '@/lib/store/cart-store'

const mutateSpy = jest.fn()
const mockQuoteState: { data: OrderQuote | undefined; isError: boolean; isFetching: boolean; refetch: jest.Mock } = { data: undefined, isError: false, isFetching: false, refetch: jest.fn() }
jest.mock('@/lib/hooks/use-order-quote', () => ({ useOrderQuote: () => mockQuoteState }))

jest.mock('@pos-saas/ui', () => ({
  Button: ({ children, ...props }: any) => <button {...props}>{children}</button>,
}))

jest.mock('@/lib/hooks/use-api-queries', () => ({
  useCreateOrder: () => ({
    mutate: mutateSpy,
    isPending: false,
  }),
}))

jest.mock('@/lib/translations', () => ({
  createTranslator: () => (key: string) => key,
}))

describe('CheckoutPageContent dine-in flow', () => {
  beforeEach(() => {
    mutateSpy.mockReset()
    mockQuoteState.data = { subtotal: 24, tax_amount: 2.4, delivery_fee: 0, total_amount: 26.4, currency: 'EGP', items: [{ product_id: 101, quantity: 2, unit_price: 12, total_price: 24 }] }
    mockQuoteState.isError = false
    mockQuoteState.isFetching = false
    localStorage.clear()
    useCartStore.setState({
      items: [
        {
          id: 1,
          productId: 101,
          name: 'Burger',
          price: 12,
          quantity: 2,
          trackInventory: false,
          quantityInStock: 100,
          total: 24,
        },
      ],
      tableId: 8,
      tableName: 'Garden 8',
    })
  })

  it('hides delivery fields and submits dine-in payload with table context', async () => {
    const { container, queryByText } = render(
      <CheckoutPageContent
        locale="en"
        themeData={{ colors: {} } as any}
      />
    )

    // Dine-in checkout should not ask for delivery address.
    expect(container.querySelector('textarea[name="deliveryAddress"]')).toBeNull()
    expect(queryByText('checkout.deliveryInformation')).toBeNull()

    const phoneInput = container.querySelector('input[name="customerPhone"]') as HTMLInputElement | null
    expect(phoneInput).not.toBeNull()
    fireEvent.change(phoneInput!, { target: { value: '0100000000' } })

    const form = container.querySelector('form') as HTMLFormElement | null
    expect(form).not.toBeNull()
    fireEvent.submit(form!)

    await waitFor(() => {
      expect(mutateSpy).toHaveBeenCalledTimes(1)
    })

    const [payload] = mutateSpy.mock.calls[0]
    expect(payload).toMatchObject({
      expected_total: 26.4,
      customer_name: 'Garden 8',
      customer_phone: '0100000000',
      delivery_address: '',
      table_id: 8,
      order_source: 'dine_in',
      payment_method: 'cash',
    })

    expect(payload.items).toEqual([
      {
        product_id: 101,
        quantity: 2,
        special_instructions: '',
      },
    ])
  })
  it('does not place an order while authoritative pricing is unavailable', () => {
    mockQuoteState.data = undefined
    mockQuoteState.isError = true
    const { container, getByRole } = render(<CheckoutPageContent locale="en" themeData={{ colors: {} }} />)
    expect(getByRole('alert')).toBeVisible()
    expect(container.querySelector('button[type="submit"]')).toBeDisabled()
    expect(mutateSpy).not.toHaveBeenCalled()
  })

})
