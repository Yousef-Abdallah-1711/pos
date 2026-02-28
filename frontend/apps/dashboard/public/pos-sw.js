/* eslint-disable no-restricted-globals */

const CACHE_NAME = 'pos-cache-v1'
const POS_ORDER_SYNC_TAG = 'pos-order-sync'
const OFFLINE_DB_NAME = 'pos-offline-db'
const OFFLINE_DB_VERSION = 1
const OFFLINE_STORE_NAME = 'pending-orders'
const API_PREFIX = '/api/v1/'

const PRECACHE_URLS = [
  '/sounds/pos-add.mp3',
  '/sounds/pos-remove.mp3',
  '/sounds/pos-success.mp3',
  '/sounds/pos-error.mp3',
]

const STATIC_ASSET_PATH_REGEX =
  /\.(?:js|mjs|css|png|jpg|jpeg|svg|gif|webp|ico|woff2?|ttf|otf|mp3)$/i

const isGetRequest = (request) => request.method === 'GET'

const isApiRequest = (url) => url.pathname.startsWith(API_PREFIX)

const isNavigationRequest = (request) => request.mode === 'navigate'

const isStaticAssetRequest = (request, url) => {
  if (url.origin !== self.location.origin) {
    return false
  }

  if (url.pathname.startsWith('/_next/')) {
    return true
  }

  if (STATIC_ASSET_PATH_REGEX.test(url.pathname)) {
    return true
  }

  return (
    request.destination === 'script' ||
    request.destination === 'style' ||
    request.destination === 'image' ||
    request.destination === 'font' ||
    request.destination === 'audio'
  )
}

const cacheFirst = async (request) => {
  const cache = await caches.open(CACHE_NAME)
  const cached = await cache.match(request)
  if (cached) {
    return cached
  }

  const response = await fetch(request)
  if (response && response.ok) {
    await cache.put(request, response.clone())
  }
  return response
}

const networkFirst = async (request) => {
  const cache = await caches.open(CACHE_NAME)

  try {
    const response = await fetch(request)
    if (response && response.ok && isGetRequest(request)) {
      await cache.put(request, response.clone())
    }
    return response
  } catch (error) {
    const cached = await cache.match(request)
    if (cached) {
      return cached
    }
    throw error
  }
}

const openOfflineQueueDB = () =>
  new Promise((resolve, reject) => {
    try {
      const request = indexedDB.open(OFFLINE_DB_NAME, OFFLINE_DB_VERSION)

      request.onupgradeneeded = () => {
        const db = request.result
        if (!db.objectStoreNames.contains(OFFLINE_STORE_NAME)) {
          db.createObjectStore(OFFLINE_STORE_NAME, { keyPath: 'id' })
        }
      }

      request.onsuccess = () => resolve(request.result)
      request.onerror = () => reject(request.error || new Error('Failed to open offline queue DB'))
    } catch (error) {
      reject(error)
    }
  })

const readQueuedOrders = async () => {
  const db = await openOfflineQueueDB()

  return new Promise((resolve, reject) => {
    const transaction = db.transaction(OFFLINE_STORE_NAME, 'readonly')
    const store = transaction.objectStore(OFFLINE_STORE_NAME)
    const request = store.getAll()

    request.onsuccess = () => {
      const rows = Array.isArray(request.result) ? request.result : []
      rows.sort((left, right) => {
        if (left.queuedAt === right.queuedAt) {
          return String(left.id).localeCompare(String(right.id))
        }
        return String(left.queuedAt).localeCompare(String(right.queuedAt))
      })
      resolve(rows)
    }
    request.onerror = () => reject(request.error || new Error('Failed to read queued orders'))
    transaction.oncomplete = () => db.close()
    transaction.onerror = () => db.close()
  })
}

const deleteQueuedOrder = async (id) => {
  const db = await openOfflineQueueDB()

  return new Promise((resolve, reject) => {
    const transaction = db.transaction(OFFLINE_STORE_NAME, 'readwrite')
    const store = transaction.objectStore(OFFLINE_STORE_NAME)
    const request = store.delete(id)

    request.onsuccess = () => resolve(true)
    request.onerror = () => reject(request.error || new Error('Failed to delete queued order'))
    transaction.oncomplete = () => db.close()
    transaction.onerror = () => db.close()
  })
}

const notifyClients = async (message) => {
  const clients = await self.clients.matchAll({
    type: 'window',
    includeUncontrolled: true,
  })

  for (const client of clients) {
    client.postMessage(message)
  }
}

const replayQueuedOrdersFromSync = async () => {
  const queuedOrders = await readQueuedOrders()
  let processed = 0
  let failed = 0

  for (const queuedOrder of queuedOrders) {
    const response = await fetch(`${API_PREFIX}public/orders`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(queuedOrder.payload),
    })

    if (!response.ok) {
      failed += 1
      throw new Error(`Sync request failed with status ${response.status}`)
    }

    await deleteQueuedOrder(queuedOrder.id)
    processed += 1
  }

  return { processed, failed }
}

self.addEventListener('install', (event) => {
  event.waitUntil(
    (async () => {
      const cache = await caches.open(CACHE_NAME)
      await Promise.allSettled(
        PRECACHE_URLS.map(async (url) => {
          const response = await fetch(url, { cache: 'reload' })
          if (response && response.ok) {
            await cache.put(url, response.clone())
          }
        })
      )
      await self.skipWaiting()
    })()
  )
})

self.addEventListener('activate', (event) => {
  event.waitUntil(
    (async () => {
      const cacheNames = await caches.keys()
      await Promise.all(
        cacheNames
          .filter((cacheName) => cacheName !== CACHE_NAME)
          .map((cacheName) => caches.delete(cacheName))
      )
      await self.clients.claim()
    })()
  )
})

self.addEventListener('fetch', (event) => {
  const request = event.request
  if (!isGetRequest(request)) {
    return
  }

  const url = new URL(request.url)

  if (isApiRequest(url)) {
    event.respondWith(networkFirst(request))
    return
  }

  if (isNavigationRequest(request)) {
    event.respondWith(networkFirst(request))
    return
  }

  if (isStaticAssetRequest(request, url)) {
    event.respondWith(cacheFirst(request))
  }
})

self.addEventListener('sync', (event) => {
  if (event.tag !== POS_ORDER_SYNC_TAG) {
    return
  }

  event.waitUntil(
    (async () => {
      try {
        const result = await replayQueuedOrdersFromSync()
        await notifyClients({
          type: 'pos-sync-complete',
          processed: result.processed,
          failed: result.failed,
        })
      } catch (error) {
        await notifyClients({
          type: 'pos-sync-error',
          reason: error instanceof Error ? error.message : 'pos-sync-failed',
        })
        throw error
      }
    })()
  )
})
