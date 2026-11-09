import { chromium } from '@playwright/test'

const BASE_URL = process.env.POS_QA_BASE_URL
if (!BASE_URL) {
  throw new Error('POS_QA_BASE_URL must be set to the dashboard URL under test')
}
const POS_PATH = '/en/dashboard/pos'
const SW_SCOPE = '/en/dashboard/pos'
const SW_URL = '/pos-sw.js'
const SYNC_TAG = 'pos-order-sync'
const API_PREFIX = '/api/v1/'
const SW_ACTIVE_TIMEOUT_MS = 7000

const result = {
  baseUrl: BASE_URL,
  steps: [],
  summary: {
    swRegistered: false,
    cacheInstalled: false,
    offlineReloadWorked: false,
    backgroundSyncSupported: false,
    backgroundSyncTriggered: false,
  },
}

const addStep = (name, status, details = {}) => {
  result.steps.push({
    name,
    status,
    timestamp: new Date().toISOString(),
    details,
  })
}

const wait = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

const isExecutionContextDestroyed = (error) =>
  error instanceof Error && error.message.includes('Execution context was destroyed')

const evaluateWithRetry = async (page, fn, arg, attempts = 3) => {
  let lastError = null
  for (let attempt = 1; attempt <= attempts; attempt += 1) {
    try {
      return await page.evaluate(fn, arg)
    } catch (error) {
      lastError = error
      if (!isExecutionContextDestroyed(error) || attempt === attempts) {
        throw error
      }

      await page.waitForLoadState('domcontentloaded', { timeout: 8000 }).catch(() => {})
      await wait(350)
    }
  }

  throw lastError || new Error('evaluateWithRetry failed')
}

const ensureOnPOSPage = async (page) => {
  const currentPath = new URL(page.url()).pathname
  if (currentPath.includes('/dashboard/pos')) {
    return
  }

  await page.goto(`${BASE_URL}${POS_PATH}`, { waitUntil: 'domcontentloaded' })
  await page.waitForLoadState('networkidle', { timeout: 12000 }).catch(() => {})
}

const queueOfflineOrder = async (page, orderId) => {
  return evaluateWithRetry(page, async ({ id }) => {
    const open = () =>
      new Promise((resolve, reject) => {
        const request = indexedDB.open('pos-offline-db', 1)
        request.onupgradeneeded = () => {
          const db = request.result
          if (!db.objectStoreNames.contains('pending-orders')) {
            db.createObjectStore('pending-orders', { keyPath: 'id' })
          }
        }
        request.onsuccess = () => resolve(request.result)
        request.onerror = () => reject(request.error || new Error('open failed'))
      })

    const db = await open()
    await new Promise((resolve, reject) => {
      const tx = db.transaction('pending-orders', 'readwrite')
      const store = tx.objectStore('pending-orders')
      store.put({
        id,
        payload: {
          order_source: 'pos',
          customer_name: 'Phase4 QA',
          payment_method: 'cash',
          items: [{ product_id: 1, quantity: 1 }],
        },
        queuedAt: new Date().toISOString(),
        retryCount: 0,
      })
      tx.oncomplete = () => resolve(true)
      tx.onerror = () => reject(tx.error || new Error('put failed'))
    })
    db.close()
    return true
  }, { id: orderId })
}

const readOfflineQueueCount = async (page) => {
  return evaluateWithRetry(page, async () => {
    const open = () =>
      new Promise((resolve, reject) => {
        const request = indexedDB.open('pos-offline-db', 1)
        request.onupgradeneeded = () => {
          const db = request.result
          if (!db.objectStoreNames.contains('pending-orders')) {
            db.createObjectStore('pending-orders', { keyPath: 'id' })
          }
        }
        request.onsuccess = () => resolve(request.result)
        request.onerror = () => reject(request.error || new Error('open failed'))
      })

    const db = await open()
    const count = await new Promise((resolve, reject) => {
      const tx = db.transaction('pending-orders', 'readonly')
      const store = tx.objectStore('pending-orders')
      const req = store.count()
      req.onsuccess = () => resolve(req.result || 0)
      req.onerror = () => reject(req.error || new Error('count failed'))
    })
    db.close()
    return count
  })
}

const run = async () => {
  const browser = await chromium.launch({
    headless: true,
    args: ['--enable-experimental-web-platform-features'],
  })
  const context = await browser.newContext()

  await context.addInitScript(() => {
    const mockUser = {
      id: 999,
      email: 'phase4.qa@local.dev',
      name: 'Phase4 QA',
      role: 'admin',
      tenant_id: 1,
      restaurant_id: 1,
    }
    const mockToken = 'phase4-qa-token'
    const authPayload = {
      state: {
        user: mockUser,
        token: mockToken,
      },
      version: 0,
    }
    window.localStorage.setItem('auth-storage', JSON.stringify(authPayload))
    window.localStorage.setItem('token', mockToken)
    window.localStorage.setItem('user', JSON.stringify(mockUser))
  })

  const page = await context.newPage()
  page.on('console', (message) => {
    if (['error', 'warning'].includes(message.type())) {
      addStep('browser-console', 'info', {
        type: message.type(),
        text: message.text(),
      })
    }
  })

  let syncReplayRequestCount = 0
  await context.route(`**${API_PREFIX}**`, async (route) => {
    const url = new URL(route.request().url())
    const path = url.pathname
    const method = route.request().method().toUpperCase()

    const reply = async (payload, status = 200) =>
      route.fulfill({
        status,
        contentType: 'application/json',
        body: JSON.stringify(payload),
      })

    if (path.endsWith('/rbac/check-permission')) {
      await reply({
        has_permission: true,
        module_id: 10,
        permission_level: 'WRITE',
      })
      return
    }

    if (path.endsWith('/tenant/status')) {
      await reply({
        success: true,
        data: {
          can_access: true,
          block_reason: null,
          status: 'active',
          days_remaining: 30,
        },
      })
      return
    }

    if (path.endsWith('/products')) {
      await reply({
        success: true,
        data: {
          products: [
            {
              id: 1,
              name: 'QA Burger',
              name_ar: 'برجر',
              price: 25,
              category_id: 2,
              is_available: true,
              barcode: '123456',
            },
          ],
        },
      })
      return
    }

    if (path.endsWith('/categories')) {
      await reply({
        success: true,
        data: {
          categories: [
            {
              id: 2,
              name: 'Main',
              name_ar: 'وجبات',
            },
          ],
        },
      })
      return
    }

    if (path.endsWith('/tables')) {
      await reply({
        success: true,
        data: [
          {
            id: 1,
            tenant_id: 1,
            restaurant_id: 1,
            table_number: 1,
            name: 'T1',
            capacity: 4,
            location: 'indoor',
            status: 'available',
            position_x: 0,
            position_y: 0,
            shape: 'square',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          },
        ],
      })
      return
    }

    if (path.endsWith('/public/orders') && method === 'POST') {
      syncReplayRequestCount += 1
      await reply({
        success: true,
        data: {
          id: 9090 + syncReplayRequestCount,
          order_number: `SW-QA-${9090 + syncReplayRequestCount}`,
          status: 'pending',
          total_amount: 25,
          created_at: new Date().toISOString(),
        },
      })
      return
    }

    if (/\/admin\/orders\/\d+\/payment-status$/.test(path) && method === 'PUT') {
      await reply({ success: true })
      return
    }

    await reply({ success: true, data: {} })
  })

  try {
    await page.goto(`${BASE_URL}${POS_PATH}`, { waitUntil: 'domcontentloaded' })
    await page.waitForLoadState('networkidle', { timeout: 12000 }).catch(() => {})

    const currentPath = new URL(page.url()).pathname
    if (!currentPath.includes('/dashboard/pos')) {
      addStep('open-pos-page', 'failed', {
        expectedPath: POS_PATH,
        actualUrl: page.url(),
      })
      throw new Error(`Unexpected redirect to ${page.url()}`)
    }

    addStep('open-pos-page', 'passed')

    const swInstall = await evaluateWithRetry(page, async ({ swUrl, scope, activeTimeoutMs }) => {
      if (!('serviceWorker' in navigator)) {
        return { supported: false }
      }

      const waitForActivation = async (registration, timeoutMs) =>
        new Promise((resolve) => {
          if (registration.active) {
            resolve(registration.active)
            return
          }

          const candidate = registration.installing || registration.waiting
          if (!candidate) {
            resolve(null)
            return
          }

          const timer = window.setTimeout(() => resolve(null), timeoutMs)
          const onStateChange = () => {
            if (candidate.state === 'activated') {
              window.clearTimeout(timer)
              resolve(candidate)
            }
          }
          candidate.addEventListener('statechange', onStateChange)
        })

      const old = await navigator.serviceWorker.getRegistration(scope)
      if (old) {
        await old.unregister()
      }

      const reg = await navigator.serviceWorker.register(swUrl, { scope })
      const active = await waitForActivation(reg, activeTimeoutMs)

      return {
        supported: true,
        scope: reg.scope,
        scriptURL: active?.scriptURL || reg.active?.scriptURL || '',
        hasController: Boolean(navigator.serviceWorker.controller),
        activeResolved: Boolean(active || reg.active),
      }
    }, { swUrl: SW_URL, scope: SW_SCOPE, activeTimeoutMs: SW_ACTIVE_TIMEOUT_MS })

    if (swInstall.supported && swInstall.scriptURL.includes('pos-sw.js')) {
      result.summary.swRegistered = true
      addStep('service-worker-registration', 'passed', swInstall)
    } else {
      addStep('service-worker-registration', 'failed', swInstall)
    }

    const controlInfoBeforeReload = await evaluateWithRetry(page, async () => ({
      hasController: Boolean(navigator.serviceWorker.controller),
      url: window.location.pathname,
    }))
    addStep('service-worker-control-before-reload', 'info', controlInfoBeforeReload)

    await page.reload({ waitUntil: 'domcontentloaded', timeout: 15000 })
    await page.waitForLoadState('networkidle', { timeout: 12000 }).catch(() => {})

    const controlInfoAfterReload = await evaluateWithRetry(page, async () => ({
      hasController: Boolean(navigator.serviceWorker.controller),
      url: window.location.pathname,
    }))
    addStep('service-worker-control-after-reload', 'info', controlInfoAfterReload)

    const cacheInfo = await evaluateWithRetry(page, async () => {
      const keys = await caches.keys()
      const hasPosCache = keys.includes('pos-cache-v1')
      const cache = hasPosCache ? await caches.open('pos-cache-v1') : null
      const requests = cache ? await cache.keys() : []
      const paths = requests.map((request) => new URL(request.url).pathname)

      return {
        keys,
        hasPosCache,
        cachedCount: paths.length,
        requiredAssets: {
          add: paths.includes('/sounds/pos-add.mp3'),
          remove: paths.includes('/sounds/pos-remove.mp3'),
          success: paths.includes('/sounds/pos-success.mp3'),
          error: paths.includes('/sounds/pos-error.mp3'),
        },
      }
    })

    if (
      cacheInfo.hasPosCache &&
      cacheInfo.requiredAssets.add &&
      cacheInfo.requiredAssets.remove &&
      cacheInfo.requiredAssets.success &&
      cacheInfo.requiredAssets.error
    ) {
      result.summary.cacheInstalled = true
      addStep('cache-inspection', 'passed', cacheInfo)
    } else {
      addStep('cache-inspection', 'failed', cacheInfo)
    }

    await context.setOffline(true)
    let offlineReloadOk = false
    try {
      await page.reload({ waitUntil: 'domcontentloaded', timeout: 12000 })
      offlineReloadOk = true
    } catch {
      offlineReloadOk = false
    }
    await context.setOffline(false)

    if (offlineReloadOk) {
      result.summary.offlineReloadWorked = true
      addStep('offline-reload', 'passed')
    } else {
      addStep('offline-reload', 'failed')
      await ensureOnPOSPage(page)
      addStep('post-offline-recovery', 'info', {
        url: page.url(),
      })
    }

    await ensureOnPOSPage(page)

    await evaluateWithRetry(page, () => {
      window.__posSyncMessages = []
      if ('serviceWorker' in navigator) {
        navigator.serviceWorker.addEventListener('message', (event) => {
          window.__posSyncMessages.push(event.data)
        })
      }
    })

    const queueId = `qa-sync-${Date.now()}`
    let queuedBefore = 0
    try {
      await queueOfflineOrder(page, queueId)
      queuedBefore = await readOfflineQueueCount(page)
    } catch (error) {
      addStep('queue-offline-order', 'failed', {
        message: error instanceof Error ? error.message : String(error),
      })
      return
    }

    const syncRegistration = await evaluateWithRetry(page, async ({ tag, scope }) => {
      if (!('serviceWorker' in navigator)) {
        return { supported: false, registered: false }
      }

      const reg = await navigator.serviceWorker.getRegistration(scope)
      if (!reg) {
        return { supported: false, registered: false, reason: 'registration-not-found' }
      }
      if (!('sync' in reg) || !reg.sync || typeof reg.sync.register !== 'function') {
        return { supported: false, registered: false }
      }

      try {
        await reg.sync.register(tag)
        return { supported: true, registered: true }
      } catch (error) {
        return {
          supported: true,
          registered: false,
          reason: error instanceof Error ? error.message : String(error),
        }
      }
    }, { tag: SYNC_TAG, scope: SW_SCOPE })

    result.summary.backgroundSyncSupported = Boolean(
      syncRegistration.supported && syncRegistration.registered
    )

    if (syncRegistration.supported && syncRegistration.registered) {
      await wait(5000)

      const syncMessages = await evaluateWithRetry(page, () => window.__posSyncMessages || [])
      const queuedAfter = await readOfflineQueueCount(page)

      const completedMessage = syncMessages.find((message) => message?.type === 'pos-sync-complete')
      if (completedMessage && queuedBefore > queuedAfter) {
        result.summary.backgroundSyncTriggered = true
        addStep('background-sync-replay', 'passed', {
          queuedBefore,
          queuedAfter,
          syncMessages,
          syncReplayRequestCount,
        })
      } else {
        addStep('background-sync-replay', 'failed', {
          queuedBefore,
          queuedAfter,
          syncMessages,
          syncReplayRequestCount,
        })
      }
    } else {
      addStep('background-sync-replay', 'skipped', {
        reason: syncRegistration.reason || 'sync-manager-not-supported',
        queuedBefore,
        syncReplayRequestCount,
      })

      await context.setOffline(true)
      await wait(700)
      await context.setOffline(false)
      await wait(2200)

      let queuedAfterFallback = await readOfflineQueueCount(page)
      if (queuedAfterFallback >= queuedBefore) {
        await evaluateWithRetry(page, () => {
          window.dispatchEvent(new Event('online'))
        })
        await wait(1500)
        queuedAfterFallback = await readOfflineQueueCount(page)
      }

      if (queuedBefore > queuedAfterFallback) {
        addStep('manual-online-replay-fallback', 'passed', {
          queuedBefore,
          queuedAfter: queuedAfterFallback,
          replayRequests: syncReplayRequestCount,
        })
      } else {
        addStep('manual-online-replay-fallback', 'failed', {
          queuedBefore,
          queuedAfter: queuedAfterFallback,
          replayRequests: syncReplayRequestCount,
        })
      }
    }
  } finally {
    await context.close()
    await browser.close()
  }
}

run()
  .then(() => {
    // eslint-disable-next-line no-console
    console.log(JSON.stringify(result, null, 2))
    if (!result.summary.swRegistered || !result.summary.cacheInstalled || !result.summary.offlineReloadWorked) {
      process.exitCode = 1
    }
  })
  .catch((error) => {
    addStep('fatal', 'failed', {
      message: error instanceof Error ? error.message : String(error),
    })
    // eslint-disable-next-line no-console
    console.log(JSON.stringify(result, null, 2))
    process.exitCode = 1
  })
