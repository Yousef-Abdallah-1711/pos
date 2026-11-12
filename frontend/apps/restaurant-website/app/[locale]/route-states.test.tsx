import { fireEvent, render, screen } from '@testing-library/react';
import { usePathname } from 'next/navigation';
import ErrorBoundary from './error';
import StorefrontLoading from './loading';

jest.mock('next/navigation', () => ({
  usePathname: jest.fn(),
}));

jest.mock('@/lib/debug', () => ({
  debugError: jest.fn(),
}));

describe('storefront route states', () => {
  beforeEach(() => {
    jest.mocked(usePathname).mockReturnValue('/ar/menu');
  });

  it('shows a direction-aware localized loading skeleton', () => {
    render(<StorefrontLoading />);

    expect(screen.getByRole('main')).toHaveAttribute('dir', 'rtl');
    expect(screen.getByRole('main')).toHaveAttribute('aria-busy', 'true');
    expect(screen.getByRole('status')).toHaveTextContent('جاري التحميل...');
  });

  it('keeps the failure and recovery action in the active locale', () => {
    const reset = jest.fn();
    render(<ErrorBoundary error={new Error('route failure')} reset={reset} />);

    expect(screen.getByRole('alert')).toHaveAttribute('dir', 'rtl');
    expect(screen.getByRole('heading', { name: 'عذراً! حدث خطأ ما' })).toBeVisible();
    fireEvent.click(screen.getByRole('button', { name: 'حاول مرة أخرى' }));
    expect(reset).toHaveBeenCalledTimes(1);
  });
});
