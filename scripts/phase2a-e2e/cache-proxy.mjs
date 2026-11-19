import { createHash } from 'node:crypto'
import { readFileSync } from 'node:fs'
import https from 'node:https'

const [listenPortValue, upstreamPortValue, certificatePath, keyPath, serviceName] = process.argv.slice(2)
const listenPort = Number(listenPortValue)
const upstreamPort = Number(upstreamPortValue)
const maxEntries = 128
const maxEntryBytes = 4 * 1024 * 1024
const maxCacheBytes = 16 * 1024 * 1024
const cache = new Map()
const cacheKeysByBaseKey = new Map()
let cacheBytes = 0
// These query names can carry credentials or proofs, so matching requests bypass shared caching.
const sensitiveQueryName = /token|proof|password|secret|code|jwt|session|signature/i
// These directives prohibit shared storage even when an origin also advertises public caching.
const privateCacheDirective = /^(private|no-store|no-cache)\b/

if (!Number.isSafeInteger(listenPort) || listenPort < 1 || listenPort > 65535
  || !Number.isSafeInteger(upstreamPort) || upstreamPort < 1 || upstreamPort > 65535
  || !certificatePath || !keyPath || !serviceName) {
  throw new Error('Usage: cache-proxy.mjs <listen-port> <upstream-port> <certificate> <key> <service>')
}

function cacheBaseKey(request, requestUrl) {
  const source = [
    serviceName,
    request.headers.host ?? '',
    requestUrl.pathname,
    requestUrl.search,
  ].join('\n')
  return createHash('sha256').update(source).digest('hex')
}

function headerValue(request, name) {
  const value = request.headers[name.toLowerCase()]
  return Array.isArray(value) ? value.join(',') : value ?? ''
}

function cacheKey(request, requestUrl, varyHeaders) {
  const dimensions = [...new Set(varyHeaders.map((name) => name.toLowerCase()))].sort()
  const values = dimensions.map((name) => [name, headerValue(request, name)])
  const source = JSON.stringify([cacheBaseKey(request, requestUrl), values])
  return createHash('sha256').update(source).digest('hex')
}

function responseVaryHeaders(headers) {
  return [...new Set(String(headers.vary ?? '').split(',').map((name) => name.trim().toLowerCase()).filter(Boolean))].sort()
}

function requestMatchesVary(request, varyHeaders, varyValues) {
  return varyHeaders.every((name) => headerValue(request, name) === varyValues[name])
}

function indexEntry(baseKey, key) {
  const indexedKeys = cacheKeysByBaseKey.get(baseKey) ?? new Set()
  indexedKeys.add(key)
  cacheKeysByBaseKey.set(baseKey, indexedKeys)
}

function removeEntry(key) {
  const entry = cache.get(key)
  if (!entry) return
  cache.delete(key)
  cacheBytes -= entry.body.length
  const indexedKeys = cacheKeysByBaseKey.get(entry.baseKey)
  indexedKeys?.delete(key)
  if (indexedKeys?.size === 0) cacheKeysByBaseKey.delete(entry.baseKey)
}

function cacheMaxAge(cacheControl) {
  const directives = cacheControl.split(',').map((directive) => directive.trim().toLowerCase())
  if (!directives.includes('public') || directives.some((directive) => privateCacheDirective.test(directive))) return 0
  const age = directives.find((directive) => directive.startsWith('s-maxage='))
    ?? directives.find((directive) => directive.startsWith('max-age='))
  if (!age) return 0
  const seconds = Number(age.split('=', 2)[1])
  return Number.isSafeInteger(seconds) && seconds > 0 ? Math.min(seconds, 3600) : 0
}

function originAgeSeconds(headers) {
  const seconds = Number(headers.age ?? 0)
  return Number.isSafeInteger(seconds) && seconds > 0 ? seconds : 0
}

function requestMustBypass(request, requestUrl) {
  const isImmutableAsset = requestUrl.pathname.startsWith('/_next/static/')
  const requestCacheControl = String(request.headers['cache-control'] ?? '').toLowerCase()
  const requestPragma = String(request.headers.pragma ?? '').toLowerCase()
  const sensitiveQuery = [...requestUrl.searchParams.keys()].some((name) => sensitiveQueryName.test(name))
  // Requests opting out of a shared cache must never be answered from an existing entry.
  const asksToBypass = /(?:^|,)\s*(?:no-cache|no-store|max-age=0)\b/.test(requestCacheControl)
    || /(?:^|,)\s*no-cache\b/.test(requestPragma)

  return request.method !== 'GET'
    || Boolean(request.headers.authorization)
    || (Boolean(request.headers.cookie) && !isImmutableAsset)
    || sensitiveQuery
    || asksToBypass
}

function responseMaxAge(statusCode, headers, bodyBytes, requestUrl, request) {
  const varyHeaders = String(headers.vary ?? '').split(',').map((name) => name.trim().toLowerCase()).filter(Boolean)
  // Next.js App Router responses vary by these request headers; cache keys include each announced value.
  const supportedVaryHeaders = new Set([
    'accept', 'accept-encoding', 'accept-language', 'origin', 'rsc', 'next-router-state-tree',
    'next-router-prefetch', 'next-router-segment-prefetch', 'next-url',
  ])
  if (statusCode !== 200 || bodyBytes > maxEntryBytes || headers['set-cookie']
    || varyHeaders.includes('*') || varyHeaders.some((name) => !supportedVaryHeaders.has(name))) return 0
  if (requestUrl.searchParams.size > 0 && [...requestUrl.searchParams.keys()].some((name) => sensitiveQueryName.test(name))) return 0
  if (request.headers.cookie && !requestUrl.pathname.startsWith('/_next/static/')) return 0

  const advertisedAge = cacheMaxAge(String(headers['cache-control'] ?? ''))
  return Math.max(0, advertisedAge - originAgeSeconds(headers))
}

function forwardedHeaders(headers, decision, key, ageSeconds) {
  const output = { ...headers, connection: 'close' }
  delete output['transfer-encoding']
  delete output['keep-alive']
  delete output['proxy-authenticate']
  delete output['proxy-authorization']
  delete output['proxy-connection']
  delete output.te
  delete output.trailer
  delete output.upgrade
  delete output.age
  output['x-servio-cache-status'] = decision
  output['x-servio-cache-key'] = key
  if (decision === 'HIT') output.age = String(ageSeconds)
  return output
}

function send(response, statusCode, headers, body, decision, key, ageSeconds = 0) {
  response.writeHead(statusCode, forwardedHeaders(headers, decision, key, ageSeconds))
  response.end(body)
}

function storeEntry(key, entry) {
  const previous = cache.get(key)
  if (previous) removeEntry(key)

  const indexedKeys = cacheKeysByBaseKey.get(entry.baseKey)
  for (const indexedKey of [...(indexedKeys ?? [])]) {
    const existing = cache.get(indexedKey)
    if (existing && JSON.stringify(existing.varyHeaders) !== JSON.stringify(entry.varyHeaders)) removeEntry(indexedKey)
  }

  cache.set(key, entry)
  cacheBytes += entry.body.length
  indexEntry(entry.baseKey, key)

  while (cache.size > maxEntries || cacheBytes > maxCacheBytes) {
    const oldestKey = cache.keys().next().value
    if (!oldestKey) break
    removeEntry(oldestKey)
  }
}

function serveCachedResponse(request, response, baseKey, bypassRequest) {
  if (bypassRequest) return false
  const indexedKeys = cacheKeysByBaseKey.get(baseKey)
  if (!indexedKeys) return false

  for (const key of indexedKeys) {
    const cached = cache.get(key)
    if (!cached || !requestMatchesVary(request, cached.varyHeaders, cached.varyValues)) continue

    const ageSeconds = cached.originAge + Math.floor((Date.now() - cached.storedAt) / 1000)
    if (ageSeconds < cached.maxAge) {
      cache.delete(key)
      cache.set(key, cached)
      send(response, cached.statusCode, cached.headers, cached.body, 'HIT', key, ageSeconds)
      return true
    }
    removeEntry(key)
  }
  return false
}

function proxyToUpstream(request, response, requestUrl, baseKey, bypassRequest) {
  const upstream = https.request({
    hostname: '127.0.0.1',
    servername: 'localhost',
    port: upstreamPort,
    path: `${requestUrl.pathname}${requestUrl.search}`,
    method: request.method,
    headers: { ...request.headers, 'x-forwarded-proto': 'https', connection: 'close' },
  }, (upstreamResponse) => {
    const statusCode = upstreamResponse.statusCode ?? 502
    const headers = { ...upstreamResponse.headers }
    const chunks = []
    let bodyBytes = 0
    let streamed = false

    upstreamResponse.on('data', (chunk) => {
      bodyBytes += chunk.length
      if (streamed) {
        response.write(chunk)
        return
      }
      chunks.push(chunk)
      if (bodyBytes > maxEntryBytes) {
        streamed = true
        response.writeHead(statusCode, forwardedHeaders(headers, 'BYPASS', baseKey, 0))
        for (const bufferedChunk of chunks) response.write(bufferedChunk)
        chunks.length = 0
      }
    })

    upstreamResponse.on('end', () => {
      if (streamed) {
        response.end()
        return
      }
      const body = Buffer.concat(chunks)
      const originAge = originAgeSeconds(headers)
      const varyHeaders = responseVaryHeaders(headers)
      const key = cacheKey(request, requestUrl, varyHeaders)
      const maxAge = bypassRequest ? 0 : responseMaxAge(statusCode, headers, bodyBytes, requestUrl, request)
      if (maxAge > 0) {
        const varyValues = Object.fromEntries(varyHeaders.map((name) => [name, headerValue(request, name)]))
        storeEntry(key, {
          baseKey,
          statusCode,
          headers,
          body,
          varyHeaders,
          varyValues,
          maxAge: maxAge + originAge,
          originAge,
          storedAt: Date.now(),
        })
        send(response, statusCode, headers, body, 'MISS', key, originAge)
        return
      }
      send(response, statusCode, headers, body, 'BYPASS', key)
    })
  })

  upstream.on('error', () => {
    if (!response.headersSent) response.writeHead(502, { 'cache-control': 'private, no-store', connection: 'close' })
    response.end('Isolated cache acceptance upstream unavailable')
  })
  request.pipe(upstream)
}

function handleCacheRequest(request, response) {
  let requestUrl
  try {
    requestUrl = new URL(request.url ?? '/', `https://localhost:${listenPort}`)
  } catch {
    response.writeHead(400, { 'cache-control': 'no-store', connection: 'close' })
    response.end('Invalid isolated cache request')
    return
  }

  const baseKey = cacheBaseKey(request, requestUrl)
  const bypassRequest = requestMustBypass(request, requestUrl)
  if (serveCachedResponse(request, response, baseKey, bypassRequest)) return
  proxyToUpstream(request, response, requestUrl, baseKey, bypassRequest)
}

const server = https.createServer({
  cert: readFileSync(certificatePath),
  key: readFileSync(keyPath),
}, handleCacheRequest)

server.listen(listenPort, '127.0.0.1', () => {
  process.stdout.write(`isolated cache proxy listening service=${serviceName} port=${listenPort}\n`)
})
