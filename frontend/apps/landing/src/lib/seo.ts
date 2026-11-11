const DEVELOPMENT_LANDING_ORIGIN = 'http://localhost:3001'

export function getLandingOrigin(): URL {
  const configuredOrigin = process.env.NEXT_PUBLIC_LANDING_URL

  if (!configuredOrigin) {
    if (process.env.NODE_ENV === 'production') {
      throw new Error('NEXT_PUBLIC_LANDING_URL is required for production SEO URLs')
    }
    return new URL(DEVELOPMENT_LANDING_ORIGIN)
  }

  const origin = new URL(configuredOrigin)
  if (process.env.NODE_ENV === 'production' && origin.protocol !== 'https:') {
    throw new Error('NEXT_PUBLIC_LANDING_URL must use HTTPS in production')
  }
  return origin
}
