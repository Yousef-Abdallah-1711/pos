/** @type {import('next').NextConfig} */
const createNextIntlPlugin = require('next-intl/plugin')
const withNextIntl = createNextIntlPlugin('./src/i18n/request.ts')
const isWindows = process.platform === 'win32'
const enableCustomSplitChunks = process.env.NEXT_CUSTOM_SPLIT_CHUNKS === 'true'
const isProd = process.env.NODE_ENV === 'production'
const paymentFrameSrc = [
  "'self'",
  'http://localhost:8080',
  'https://accept.paymob.com',
  'https://*.paymob.com',
].join(' ')

const nextConfig = {
  reactStrictMode: true,
  // The dashboard currently has pre-existing type errors; Docker prod builds can opt out of the type gate.
  typescript: { ignoreBuildErrors: process.env.NEXT_IGNORE_BUILD_ERRORS === 'true' },
  transpilePackages: ['@pos-saas/ui', '@pos-saas/types'],
  images: {
    domains: ['localhost', 'images.unsplash.com'],
    // Enable AVIF for better compression (25-30% smaller than WebP)
    formats: ['image/avif', 'image/webp'],
    // Image optimization settings
    deviceSizes: [320, 640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
    // Use modern formats and optimize by default
    minimumCacheTTL: 31536000, // 1 year for static images
  },
  // Windows local builds can fail with EPERM on standalone symlink tracing.
  output: isWindows ? undefined : 'standalone',
  outputFileTracingRoot: require('path').join(__dirname, '../../'),

  // ===== CODE SPLITTING CONFIGURATION =====
  // Optimized webpack configuration for better code splitting
  webpack: (config, { isServer, dev }) => {
    // Keep Next.js defaults during development and SSR to avoid runtime chunk
    // instability in App Router pages (e.g. products page module factory errors).
    if (dev || isServer || !enableCustomSplitChunks) {
      return config
    }

    config.optimization = {
      ...config.optimization,
      splitChunks: {
        chunks: 'all',
        cacheGroups: {
          // React and ReactDOM go to their own chunk
          react: {
            test: /[\\/]node_modules[\\/](react|react-dom)[\\/]/,
            name: 'react',
            priority: 50,
            reuseExistingChunk: true,
            enforce: true,
          },
          // UI library chunks
          ui: {
            test: /[\\/]node_modules[\\/](@pos-saas\/ui|@radix-ui)[\\/]/,
            name: 'ui-libs',
            priority: 40,
            reuseExistingChunk: true,
          },
          // Form and validation libraries
          forms: {
            test: /[\\/]node_modules[\\/](react-hook-form|zod|@hookform)[\\/]/,
            name: 'form-libs',
            priority: 35,
            reuseExistingChunk: true,
          },
          // Data fetching and state management
          data: {
            test: /[\\/]node_modules[\\/](@tanstack\/react-query|axios|swr)[\\/]/,
            name: 'data-libs',
            priority: 30,
            reuseExistingChunk: true,
          },
          // Utility libraries
          utils: {
            test: /[\\/]node_modules[\\/](lodash|date-fns|clsx)[\\/]/,
            name: 'utils-libs',
            priority: 20,
            reuseExistingChunk: true,
          },
          // All other node_modules
          vendors: {
            test: /[\\/]node_modules[\\/]/,
            name: 'vendors',
            priority: 10,
            reuseExistingChunk: true,
          },
          // Shared components used across multiple routes
          sharedComponents: {
            test: /[\\/]src[\\/]components[\\/]/,
            name: 'shared-components',
            minChunks: 2,
            priority: 5,
            reuseExistingChunk: true,
          },
        },
      },
    }
    return config
  },

  // ===== EXPERIMENTAL FEATURES FOR OPTIMIZATION =====
  experimental: {
    // Optimized package imports for better tree-shaking
    optimizePackageImports: [
      '@radix-ui/react-icons',
      '@radix-ui/react-dialog',
      'lucide-react',
      'date-fns',
      'recharts',
      '@tanstack/react-table',
      '@tanstack/react-query',
      'jspdf',
      'html2canvas',
    ],
  },

  // ===== PRODUCTION OPTIMIZATIONS =====
  compress: true,
  productionBrowserSourceMaps: false,
  poweredByHeader: false,

  // ===== HEADERS FOR BETTER CACHING =====
  async headers() {
    return [
      {
        source: '/static/:path*',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=31536000, immutable',
          },
        ],
      },
      {
        source: '/_next/static/:path*',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=31536000, immutable',
          },
        ],
      },
      {
        source: '/fonts/:path*',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=31536000, immutable',
          },
        ],
      },
      {
        source: '/:path*',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=3600, s-maxage=3600',
          },
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'Referrer-Policy',
            value: 'strict-origin-when-cross-origin',
          },
          {
            key: 'Permissions-Policy',
            value: 'camera=(), microphone=(), geolocation=()',
          },
          {
            key: 'Content-Security-Policy',
            value:
              `default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://va.vercel-scripts.com; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob: https:; connect-src 'self' https: http: ws: wss:; font-src 'self' data:; frame-src ${paymentFrameSrc}; frame-ancestors 'none'; base-uri 'self'; form-action 'self'`,
          },
          ...(isProd
            ? [
                {
                  key: 'Strict-Transport-Security',
                  value: 'max-age=31536000; includeSubDomains; preload',
                },
              ]
            : []),
        ],
      },
    ]
  },
}

module.exports = withNextIntl(nextConfig)
