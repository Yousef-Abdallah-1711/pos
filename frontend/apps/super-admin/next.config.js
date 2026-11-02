/** @type {import('next').NextConfig} */
const path = require('path');

const nextConfig = {
  reactStrictMode: true,
  // Windows local builds can fail with EPERM on standalone symlink tracing.
  output: process.platform === 'win32' ? undefined : 'standalone',
  outputFileTracingRoot: path.join(__dirname, '../..'),
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "images.unsplash.com",
      },
      {
        protocol: "https",
        hostname: "api.placeholder.com",
      },
    ],
  },
  transpilePackages: ["@pos-saas/ui-components", "@pos-saas/ui-themes"],
  experimental: {
    optimizePackageImports: ["@pos-saas/ui-components", "@pos-saas/ui-themes"],
  },
};

module.exports = nextConfig;
