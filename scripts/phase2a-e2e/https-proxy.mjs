import { readFileSync } from 'node:fs'
import http from 'node:http'
import https from 'node:https'

const [listenPortValue, targetPortValue, certificatePath, keyPath] = process.argv.slice(2)
const listenPort = Number(listenPortValue)
const targetPort = Number(targetPortValue)

if (!Number.isSafeInteger(listenPort) || listenPort < 1 || listenPort > 65535
  || !Number.isSafeInteger(targetPort) || targetPort < 1 || targetPort > 65535
  || !certificatePath || !keyPath) {
  throw new Error('Usage: https-proxy.mjs <listen-port> <target-port> <certificate> <key>')
}

const server = https.createServer({
  cert: readFileSync(certificatePath),
  key: readFileSync(keyPath),
}, (request, response) => {
  const headers = { ...request.headers, 'x-forwarded-proto': 'https' }
  const upstream = http.request({
    hostname: '127.0.0.1',
    port: targetPort,
    path: request.url,
    method: request.method,
    headers,
  }, (upstreamResponse) => {
    response.writeHead(upstreamResponse.statusCode ?? 502, upstreamResponse.headers)
    upstreamResponse.pipe(response)
  })

  upstream.on('error', () => {
    if (!response.headersSent) response.writeHead(502, { 'content-type': 'text/plain' })
    response.end('Isolated acceptance upstream unavailable')
  })
  request.pipe(upstream)
})

server.listen(listenPort, '127.0.0.1')
