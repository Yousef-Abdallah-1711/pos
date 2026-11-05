import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./vitest.setup.ts'],
    // tests/visual-regression holds Playwright specs (run with Playwright, not Vitest).
    // components/ and lib/ at the app root are unreferenced duplicates of src/ code (see PHASE_2B_HANDOFF.md).
    exclude: [
      '**/node_modules/**',
      'tests/visual-regression/**',
      'tests/phase2a-e2e/**',
      'components/**',
      'lib/**',
    ],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html', 'lcov'],
      exclude: [
        'node_modules/',
        'dist/',
        '**/*.test.ts',
        '**/*.test.tsx',
        '**/types/',
      ],
      lines: 70,
      functions: 70,
      branches: 70,
      statements: 70,
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
})
