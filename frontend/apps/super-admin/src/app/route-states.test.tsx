import { fireEvent, render, screen } from '@testing-library/react';
import '@testing-library/jest-dom/vitest';
import { describe, expect, it, vi } from 'vitest';
import SuperAdminError from './error';
import SuperAdminLoading from './loading';

describe('super admin route states', () => {
  it('shows a skeleton while route content loads', () => {
    render(<SuperAdminLoading />);

    expect(screen.getByRole('main')).toHaveAttribute('aria-busy', 'true');
    expect(screen.getByRole('status')).toHaveTextContent('Loading...');
  });

  it('offers a retry action when route content fails', () => {
    const reset = vi.fn();
    render(<SuperAdminError error={new Error('route failure')} reset={reset} />);

    expect(screen.getByRole('alert')).toBeVisible();
    fireEvent.click(screen.getByRole('button', { name: 'Try again' }));
    expect(reset).toHaveBeenCalledOnce();
  });
});
