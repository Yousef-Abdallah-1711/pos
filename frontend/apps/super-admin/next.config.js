/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
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
