/** @type {import('next').NextConfig} */
const isWindows = process.platform === 'win32'

const nextConfig = {
  reactStrictMode: true,
  transpilePackages: ['@pos-saas/ui', '@pos-saas/types', '@pos-saas/form-foundation'],
  images: {
    domains: ['localhost', 'images.unsplash.com'],
  },
  // Avoid Windows EPERM symlink failures during standalone trace copy.
  output: isWindows ? undefined : 'standalone',
  outputFileTracingRoot: require('path').join(__dirname, '../../'),
}

module.exports = nextConfig
