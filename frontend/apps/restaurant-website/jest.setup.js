// Learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom'

// Note: MSW server setup is handled in individual test files
// to avoid module resolution issues in Jest environment

// Mock localStorage
const localStorageData = new Map()
const localStorageMock = {
  getItem: jest.fn(),
	setItem: jest.fn((key, value) => localStorageData.set(key, String(value))),
	removeItem: jest.fn((key) => localStorageData.delete(key)),
	clear: jest.fn(() => localStorageData.clear()),
}

localStorageMock.getItem.mockImplementation((key) => localStorageData.get(key) ?? null)

Object.defineProperty(global, 'localStorage', {
  value: localStorageMock,
  writable: true,
})

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(), // deprecated
    removeListener: jest.fn(), // deprecated
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
})

// Mock requestAnimationFrame
global.requestAnimationFrame = (cb) => {
  return setTimeout(cb, 0)
}

global.cancelAnimationFrame = (id) => {
  clearTimeout(id)
}
