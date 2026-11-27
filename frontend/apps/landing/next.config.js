/** @type {import('next').NextConfig} */
const isWindows = process.platform === 'win32'

const nextConfig = {
  reactStrictMode: true,
  transpilePackages: ['@pos-saas/ui', '@pos-saas/types', '@pos-saas/form-foundation', '@pos-saas/i18n'],
  images: {
    qualities: [50, 75],
    remotePatterns: [
      { protocol: 'http', hostname: 'localhost' },
      { protocol: 'https', hostname: 'images.unsplash.com' },
    ],
  },
  // Avoid Windows EPERM symlink failures during standalone trace copy.
  output: isWindows ? undefined : 'standalone',
  outputFileTracingRoot: require('path').join(__dirname, '../../'),
  async headers() {
    const htmlCache = [{ key: 'Cache-Control', value: 'public, max-age=0, s-maxage=300, stale-while-revalidate=60' }]
    const security = [
      { key: 'X-Frame-Options', value: 'DENY' },
      { key: 'X-Content-Type-Options', value: 'nosniff' },
      { key: 'Referrer-Policy', value: 'strict-origin-when-cross-origin' },
      { key: 'Permissions-Policy', value: 'camera=(), microphone=(), geolocation=()' },
    ]
    const localeRoutes = ['', '/features', '/features/:slug', '/pricing', '/how-it-works', '/demo', '/contact', '/privacy']
    return [
      ...localeRoutes.map(route => ({ source: `/:locale(en|ar)${route}`, headers: [...htmlCache, ...security] })),
      { source: '/:locale(en|ar)/auth/:path*', headers: [{ key: 'Cache-Control', value: 'private, no-store, max-age=0' }, ...security] },
      { source: '/_next/static/:path*', headers: [{ key: 'Cache-Control', value: 'public, max-age=31536000, immutable' }] },
      { source: '/images/:path*', headers: [{ key: 'Cache-Control', value: 'public, max-age=604800, stale-while-revalidate=86400' }] },
    ]
  },
}

module.exports = nextConfig
