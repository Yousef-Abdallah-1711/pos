import { fireEvent, render, waitFor } from '@testing-library/react'
import { CheckoutPageContent } from '@/app/[locale]/checkout/checkout-page-content'
import { useCartStore } from '@/lib/store/cart-store'

const mutateSpy = jest.fn()

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
})
