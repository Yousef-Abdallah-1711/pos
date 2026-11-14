import 'vitest'
import type { TestingLibraryMatchers } from '@testing-library/jest-dom/matchers'

// Vitest 5 extends Matchers<R, T>; jest-dom's legacy Assertion augmentation
// does not reach this interface. Keep the real DOM matcher argument/return types.
declare module 'vitest' {
  interface Matchers<R extends void | Promise<void> = void | Promise<void>, T = unknown>
    extends TestingLibraryMatchers<T, R> {}
}
