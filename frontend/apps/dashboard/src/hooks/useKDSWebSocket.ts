'use client'

import { useEffect, useRef } from 'react'
import { getConfig } from '@/lib/config'
import { debugError, debugLog } from '@/lib/debug'
import { useAuthStore } from '@/stores/authStore'
import { KDSOrder, useKDSStore } from '@/stores/kdsStore'

interface KDSWebSocketMessage {
  type?: string
  data?: KDSOrder | KDSOrder[] | null
}

interface UseKDSWebSocketOptions {
  onOrderNew?: () => void | Promise<void>
}

const MAX_RECONNECT_DELAY_MS = 10_000

const buildWebSocketURL = () => {
  const { apiUrl } = getConfig()

  try {
    const parsed = new URL(apiUrl)
    const protocol = parsed.protocol === 'https:' ? 'wss:' : 'ws:'
    const normalizedPath = parsed.pathname.replace(/\/$/, '')
    return `${protocol}//${parsed.host}${normalizedPath}/kds/ws`
  } catch {
    if (typeof window !== 'undefined') {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      return `${protocol}//${window.location.host}/api/v1/kds/ws`
    }
    return ''
  }
}

const parseSocketMessage = (raw: string): KDSWebSocketMessage | null => {
  try {
    return JSON.parse(raw) as KDSWebSocketMessage
  } catch {
    return null
  }
}

export const useKDSWebSocket = ({ onOrderNew }: UseKDSWebSocketOptions = {}) => {
  const token = useAuthStore((state) => state.token)
  const setConnected = useKDSStore((state) => state.setConnected)
  const setOrders = useKDSStore((state) => state.setOrders)
  const upsertOrder = useKDSStore((state) => state.upsertOrder)

  const socketRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<number | null>(null)
  const reconnectAttemptRef = useRef(0)
  const shouldReconnectRef = useRef(true)
  const onOrderNewRef = useRef<(() => void | Promise<void>) | undefined>(onOrderNew)

  useEffect(() => {
    onOrderNewRef.current = onOrderNew
  }, [onOrderNew])

  useEffect(() => {
    if (!token) {
      setConnected(false)
      return
    }

    shouldReconnectRef.current = true
    reconnectAttemptRef.current = 0
    setConnected(false)
    const baseURL = buildWebSocketURL()
    if (!baseURL) {
      debugError('[KDS WS] Failed to resolve websocket URL')
      setConnected(false)
      return
    }

    const clearReconnectTimer = () => {
      if (reconnectTimeoutRef.current !== null) {
        window.clearTimeout(reconnectTimeoutRef.current)
        reconnectTimeoutRef.current = null
      }
    }

    const closeSocketSafely = (socket: WebSocket) => {
      if (socket.readyState === WebSocket.CONNECTING) {
        // In React Strict Mode development pass, cleanup can run before handshake completes.
        // Closing after open avoids browser noise: "WebSocket is closed before connection is established."
        socket.addEventListener(
          'open',
          () => {
            socket.close()
          },
          { once: true }
        )
        return
      }

      if (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CLOSING) {
        socket.close()
      }
    }

    const handleMessage = (message: KDSWebSocketMessage | null) => {
      if (!message?.type) {
        return
      }

      switch (message.type) {
        case 'order_new':
        case 'order_updated':
        case 'order_bumped':
        case 'order_recalled': {
          if (Array.isArray(message.data)) {
            setOrders(message.data)
          } else if (message.data) {
            upsertOrder(message.data)
          }
          if (message.type === 'order_new') {
            void onOrderNewRef.current?.()
          }
          break
        }
        default:
          break
      }
    }

    const connect = () => {
      clearReconnectTimer()

      const wsURL = `${baseURL}?token=${encodeURIComponent(token)}`
      debugLog('[KDS WS] Connecting:', wsURL)

      const socket = new WebSocket(wsURL)
      socketRef.current = socket

      socket.onopen = () => {
        reconnectAttemptRef.current = 0
        setConnected(true)
      }

      socket.onmessage = (event) => {
        const chunks = String(event.data || '').split('\n').filter(Boolean)
        for (const chunk of chunks) {
          handleMessage(parseSocketMessage(chunk))
        }
      }

      socket.onerror = () => {
        // Let onclose drive reconnect logic.
      }

      socket.onclose = () => {
        if (socketRef.current === socket) {
          socketRef.current = null
        }
        setConnected(false)

        if (!shouldReconnectRef.current) {
          return
        }

        const delay = Math.min(
          MAX_RECONNECT_DELAY_MS,
          1000 * Math.pow(2, reconnectAttemptRef.current)
        )
        reconnectAttemptRef.current += 1

        reconnectTimeoutRef.current = window.setTimeout(() => {
          connect()
        }, delay)
      }
    }

    connect()

    return () => {
      shouldReconnectRef.current = false
      clearReconnectTimer()
      setConnected(false)

      if (socketRef.current) {
        closeSocketSafely(socketRef.current)
        socketRef.current = null
      }
    }
  }, [token, setConnected, setOrders, upsertOrder])
}
