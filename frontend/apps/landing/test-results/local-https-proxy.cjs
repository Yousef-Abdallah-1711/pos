const fs = require('node:fs')
const http = require('node:http')
const https = require('node:https')
const path = require('node:path')

https.createServer({
  key: fs.readFileSync(path.join(__dirname, 'local-proxy-key.pem')),
  cert: fs.readFileSync(path.join(__dirname, 'local-proxy-cert.pem')),
}, (request, response) => {
  const upstream = http.request({
    hostname: '127.0.0.1', port: 3004, path: request.url,
    method: request.method, headers: { ...request.headers, host: 'localhost:3004' },
  }, upstreamResponse => {
    response.writeHead(upstreamResponse.statusCode, upstreamResponse.headers)
    upstreamResponse.pipe(response)
  })
  upstream.on('error', error => { response.writeHead(502); response.end(error.message) })
  request.pipe(upstream)
}).listen(3001, '127.0.0.1', () => console.log('HTTPS proxy on 3001'))
