import { renderHook, act } from '@testing-library/react'
import { useThemeBuilderStore } from './themeBuilderStore'
import {
  ThemeData,
  GlobalColors,
  TypographySettings,
  WebsiteIdentity,
  HeaderConfig,
  FooterConfig,
} from '@/types/theme'

// Mock data
const mockTheme: ThemeData = {
  id: 'theme-1',
  restaurantId: 'rest-1',
  tenantId: 'tenant-1',
  name: 'Test Theme',
  slug: 'test-theme',
  description: 'Test theme description',
  version: 1,
  isActive: false,
  isPublished: false,
  colors: {
    primary: '#007bff',
    secondary: '#6c757d',
    accent: '#ff6b6b',
    background: '#ffffff',
    text: '#212529',
    border: '#dee2e6',
    shadow: '#00000020',
  },
  typography: {
    fontFamily: 'Arial, sans-serif',
    baseFontSize: 16,
    borderRadius: 4,
    lineHeight: 1.5,
  },
  identity: {
    siteTitle: 'Test Restaurant',
  },
  header: {
    showLogo: true,
    navigationItems: [
      { id: '1', label: 'Menu', href: '/menu', order: 1 },
    ],
    backgroundColor: '#ffffff',
    textColor: '#212529',
    height: 80,
  },
  footer: {
    companyName: 'Test Restaurant',
    companyDescription: 'Test description',
    copyrightText: '© 2024 Test Restaurant',
    socialLinks: [],
    backgroundColor: '#f8f9fa',
    textColor: '#212529',
    showLinks: true,
  },
  components: [],
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
}

describe('useThemeBuilderStore', () => {
  beforeEach(() => {
    const { result } = renderHook(() => useThemeBuilderStore())
    act(() => {
      result.current.resetTheme()
    })
  })

  describe('Initial State', () => {
    it('should initialize with default values', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.currentTheme).toBeNull()
      expect(result.current.isLoading).toBe(false)
      expect(result.current.isSaving).toBe(false)
      expect(result.current.error).toBeNull()
      expect(result.current.successMessage).toBeNull()
      expect(result.current.isDirty).toBe(false)
    })

    it('should have default colors', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.colors.primary).toBe('#007bff')
      expect(result.current.colors.secondary).toBe('#6c757d')
      expect(result.current.colors.background).toBe('#ffffff')
      expect(result.current.colors.text).toBe('#212529')
    })

    it('should have default typography', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.typography.baseFontSize).toBe(16)
      expect(result.current.typography.borderRadius).toBe(4)
      expect(result.current.typography.lineHeight).toBe(1.5)
      expect(result.current.typography.fontFamily).toContain('Segoe')
    })

    it('should have empty components array', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.components).toEqual([])
    })

    it('should have empty custom CSS', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.customCss).toBe('')
    })
  })

  describe('Color Updates', () => {
    it('should update primary color', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateColors({ primary: '#ff0000' })
      })

      expect(result.current.colors.primary).toBe('#ff0000')
      expect(result.current.isDirty).toBe(true)
    })

    it('should update multiple colors at once', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateColors({
          primary: '#ff0000',
          secondary: '#00ff00',
          accent: '#0000ff',
        })
      })

      expect(result.current.colors.primary).toBe('#ff0000')
      expect(result.current.colors.secondary).toBe('#00ff00')
      expect(result.current.colors.accent).toBe('#0000ff')
    })

    it('should preserve unchanged colors when updating', () => {
      const { result } = renderHook(() => useThemeBuilderStore())
      const originalBackground = result.current.colors.background

      act(() => {
        result.current.updateColors({ primary: '#ff0000' })
      })

      expect(result.current.colors.background).toBe(originalBackground)
    })

    it('should mark theme as dirty after color update', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateColors({ primary: '#ff0000' })
      })

      expect(result.current.isDirty).toBe(true)
    })
  })

  describe('Typography Updates', () => {
    it('should update font family', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateTypography({ fontFamily: 'Georgia, serif' })
      })

      expect(result.current.typography.fontFamily).toBe('Georgia, serif')
      expect(result.current.isDirty).toBe(true)
    })

    it('should update font size', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateTypography({ baseFontSize: 18 })
      })

      expect(result.current.typography.baseFontSize).toBe(18)
    })

    it('should update multiple typography values', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateTypography({
          baseFontSize: 18,
          lineHeight: 1.8,
          borderRadius: 8,
        })
      })

      expect(result.current.typography.baseFontSize).toBe(18)
      expect(result.current.typography.lineHeight).toBe(1.8)
      expect(result.current.typography.borderRadius).toBe(8)
    })
  })

  describe('Header Updates', () => {
    it('should update header background color', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateHeader({ backgroundColor: '#f0f0f0' })
      })

      expect(result.current.header.backgroundColor).toBe('#f0f0f0')
      expect(result.current.isDirty).toBe(true)
    })

    it('should update header navigation items', () => {
      const { result } = renderHook(() => useThemeBuilderStore())
      const newNavItems = [
        { id: '1', label: 'Home', href: '/', order: 1 },
        { id: '2', label: 'Products', href: '/products', order: 2 },
      ]

      act(() => {
        result.current.updateHeader({ navigationItems: newNavItems })
      })

      expect(result.current.header.navigationItems).toEqual(newNavItems)
      expect(result.current.header.navigationItems.length).toBe(2)
    })

    it('should update header height', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateHeader({ height: 100 })
      })

      expect(result.current.header.height).toBe(100)
    })

    it('should toggle sticky header', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateHeader({ stickyHeader: false })
      })

      expect(result.current.header.stickyHeader).toBe(false)
    })
  })

  describe('Footer Updates', () => {
    it('should update company name', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateFooter({ companyName: 'New Company' })
      })

      expect(result.current.footer.companyName).toBe('New Company')
      expect(result.current.isDirty).toBe(true)
    })

    it('should update footer colors', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateFooter({
          backgroundColor: '#2c3e50',
          textColor: '#ffffff',
        })
      })

      expect(result.current.footer.backgroundColor).toBe('#2c3e50')
      expect(result.current.footer.textColor).toBe('#ffffff')
    })

    it('should update social links', () => {
      const { result } = renderHook(() => useThemeBuilderStore())
      const socialLinks = [
        { id: '1', platform: 'facebook', url: 'https://facebook.com/test' },
        { id: '2', platform: 'twitter', url: 'https://twitter.com/test' },
      ]

      act(() => {
        result.current.updateFooter({ socialLinks })
      })

      expect(result.current.footer.socialLinks).toEqual(socialLinks)
    })
  })

  describe('Identity Updates', () => {
    it('should update site title', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateIdentity({ siteTitle: 'My New Restaurant' })
      })

      expect(result.current.identity.siteTitle).toBe('My New Restaurant')
      expect(result.current.isDirty).toBe(true)
    })

    it('should update logo URL', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateIdentity({ logoUrl: 'https://example.com/logo.png' })
      })

      expect(result.current.identity.logoUrl).toBe('https://example.com/logo.png')
    })
  })

  describe('Theme Loading', () => {
    it('should load theme data', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.loadTheme(mockTheme)
      })

      expect(result.current.currentTheme).toEqual(mockTheme)
      expect(result.current.colors).toEqual(mockTheme.colors)
      expect(result.current.isDirty).toBe(false)
    })

    it('should load theme and preserve websiteIdentity as identity', () => {
      const { result } = renderHook(() => useThemeBuilderStore())
      const themeWithWebsiteIdentity: ThemeData = {
        ...mockTheme,
        websiteIdentity: { siteTitle: 'Website Title' },
      }

      act(() => {
        result.current.loadTheme(themeWithWebsiteIdentity)
      })

      expect(result.current.identity.siteTitle).toBe('Website Title')
    })

    it('should clear error when loading theme', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.setError('Previous error')
        result.current.loadTheme(mockTheme)
      })

      expect(result.current.error).toBeNull()
    })

    it('should load theme components', () => {
      const { result } = renderHook(() => useThemeBuilderStore())
      const themeWithComponents: ThemeData = {
        ...mockTheme,
        components: [
          {
            id: 'hero-1',
            type: 'hero',
            title: 'Hero Banner',
            displayOrder: 1,
            enabled: true,
          },
        ],
      }

      act(() => {
        result.current.loadTheme(themeWithComponents)
      })

      expect(result.current.components.length).toBe(1)
      expect(result.current.components[0].type).toBe('hero')
    })
  })

  describe('Theme Reset', () => {
    it('should reset to initial state', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.loadTheme(mockTheme)
        result.current.updateColors({ primary: '#ff0000' })
        result.current.setError('Some error')
      })

      act(() => {
        result.current.resetTheme()
      })

      expect(result.current.currentTheme).toBeNull()
      expect(result.current.colors.primary).toBe('#007bff')
      expect(result.current.error).toBeNull()
      expect(result.current.isDirty).toBe(false)
    })
  })

  describe('Dirty State', () => {
    it('should mark as dirty on color update', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.isDirty).toBe(false)

      act(() => {
        result.current.updateColors({ primary: '#ff0000' })
      })

      expect(result.current.isDirty).toBe(true)
    })

    it('should mark as clean manually', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => result.current.updateColors({ primary: '#ff0000' }))
      expect(result.current.isDirty).toBe(true)

      act(() => result.current.markClean())

      expect(result.current.isDirty).toBe(false)
    })

    it('should mark as dirty manually', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.markDirty()
      })

      expect(result.current.isDirty).toBe(true)
    })
  })

  describe('Validation', () => {
    it('should validate colors correctly', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.validateColors()).toBe(true)

      act(() => {
        result.current.updateColors({ primary: 'invalid' })
      })

      expect(result.current.validateColors()).toBe(false)
    })

    it('should validate typography correctly', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.validateTypography()).toBe(true)

      act(() => {
        result.current.updateTypography({ baseFontSize: -5 })
      })

      expect(result.current.validateTypography()).toBe(false)
    })

    it('should validate identity correctly', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.validateIdentity()).toBe(true)

      act(() => {
        result.current.updateIdentity({ siteTitle: '' })
      })

      expect(result.current.validateIdentity()).toBe(false)
    })

    it('should validate header correctly', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.validateHeader()).toBe(true)

      act(() => {
        result.current.updateHeader({ height: 0 })
      })

      expect(result.current.validateHeader()).toBe(false)
    })

    it('should validate footer correctly', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.validateFooter()).toBe(true)

      act(() => {
        result.current.updateFooter({ companyName: '' })
      })

      expect(result.current.validateFooter()).toBe(false)
    })

    it('should validate entire theme', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.validateTheme()).toBe(true)

      act(() => {
        result.current.updateColors({ primary: 'invalid' })
      })

      expect(result.current.validateTheme()).toBe(false)
    })
  })

  describe('Selectors', () => {
    it('should get color value', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateColors({ primary: '#ff0000' })
      })

      expect(result.current.getColorValue('primary')).toBe('#ff0000')
    })

    it('should get color value with fallback to default', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.getColorValue('secondary')).toBe('#6c757d')
    })

    it('should get component by ID', () => {
      const { result } = renderHook(() => useThemeBuilderStore())
      const components = [
        {
          id: 'comp-1',
          type: 'hero' as const,
          title: 'Hero',
          displayOrder: 1,
          enabled: true,
        },
        {
          id: 'comp-2',
          type: 'products' as const,
          title: 'Products',
          displayOrder: 2,
          enabled: true,
        },
      ]

      act(() => {
        result.current.updateComponents(components)
      })

      const comp = result.current.getComponentById('comp-1')
      expect(comp).toBeDefined()
      expect(comp?.type).toBe('hero')
    })

    it('should return undefined for non-existent component', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      const comp = result.current.getComponentById('non-existent')
      expect(comp).toBeUndefined()
    })

    it('should get form data', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.loadTheme(mockTheme)
      })

      const formData = result.current.getFormData()
      expect(formData.name).toBe('Test Theme')
      expect(formData.colors).toBeDefined()
      expect(formData.typography).toBeDefined()
    })
  })

  describe('Status Management', () => {
    it('should set loading state', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.isLoading).toBe(false)

      act(() => {
        result.current.setLoading(true)
      })

      expect(result.current.isLoading).toBe(true)
    })

    it('should set saving state', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      expect(result.current.isSaving).toBe(false)

      act(() => {
        result.current.setSaving(true)
      })

      expect(result.current.isSaving).toBe(true)
    })

    it('should set error message', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.setError('Error loading theme')
      })

      expect(result.current.error).toBe('Error loading theme')
      expect(result.current.hasErrors()).toBe(true)
    })

    it('should set success message', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.setSuccess('Theme saved successfully')
      })

      expect(result.current.successMessage).toBe('Theme saved successfully')
    })

    it('should clear error', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.setError('Some error')
        result.current.setError(null)
      })

      expect(result.current.error).toBeNull()
      expect(result.current.hasErrors()).toBe(false)
    })
  })

  describe('Component Management', () => {
    it('should add components', () => {
      const { result } = renderHook(() => useThemeBuilderStore())
      const components = [
        {
          id: 'hero-1',
          type: 'hero' as const,
          title: 'Hero Banner',
          displayOrder: 1,
          enabled: true,
        },
      ]

      act(() => {
        result.current.updateComponents(components)
      })

      expect(result.current.components.length).toBe(1)
      expect(result.current.components[0].type).toBe('hero')
    })

    it('should replace components', () => {
      const { result } = renderHook(() => useThemeBuilderStore())
      const initialComponents = [
        {
          id: 'hero-1',
          type: 'hero' as const,
          title: 'Hero',
          displayOrder: 1,
          enabled: true,
        },
      ]

      act(() => {
        result.current.updateComponents(initialComponents)
      })

      const newComponents = [
        {
          id: 'products-1',
          type: 'products' as const,
          title: 'Products',
          displayOrder: 1,
          enabled: true,
        },
      ]

      act(() => {
        result.current.updateComponents(newComponents)
      })

      expect(result.current.components.length).toBe(1)
      expect(result.current.components[0].type).toBe('products')
    })
  })

  describe('Custom CSS', () => {
    it('should update custom CSS', () => {
      const { result } = renderHook(() => useThemeBuilderStore())
      const css = ':root { --custom-var: red; }'

      act(() => {
        result.current.updateCustomCss(css)
      })

      expect(result.current.customCss).toBe(css)
      expect(result.current.isDirty).toBe(true)
    })

    it('should handle empty CSS', () => {
      const { result } = renderHook(() => useThemeBuilderStore())

      act(() => {
        result.current.updateCustomCss('')
      })

      expect(result.current.customCss).toBe('')
    })
  })
})
