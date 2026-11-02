const path = require('path');

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Windows local builds can fail with EPERM on standalone symlink tracing.
  output: process.platform === 'win32' ? undefined : 'standalone',
  outputFileTracingRoot: path.join(__dirname, '../..'),
  experimental: {
    optimizePackageImports: ["@radix-ui/react-icons"],
  },
  // NEXT_PUBLIC_* variables are exposed automatically by Next.js.
  // Do not inject localhost defaults here; runtime config should fail in production when missing.
  // Configure external image sources
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'images.unsplash.com',
      },
      {
        protocol: 'https',
        hostname: 'cdn.example.com',
      },
      {
        protocol: 'http',
        hostname: 'localhost',
      },
    ],
  },
  webpack: (config) => {
    config.resolve.fallback = {
      ...config.resolve.fallback,
      fs: false,
      path: false,
      crypto: false,
    };
    return config;
  },
};

module.exports = nextConfig;
