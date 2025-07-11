// WebSocket Service for handling real-time communications
// This centralizes WebSocket connections and handling to avoid duplication

import { ref } from 'vue'

let socket = null
let reconnectTimeout = null
const RECONNECT_DELAY = 2000
const listeners = new Map()
let isConnecting = false

// Reactive connection status
const connectionStatus = ref(false)

const wsService = {
  // Connect to the WebSocket server
  connect() {
    if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
      console.log('WebSocket already connected or connecting')
      return
    }
    
    if (isConnecting) return
    isConnecting = true
    
    try {
      console.log('Connecting to WebSocket server...')
      // Use the same origin as the current page to inherit cookies/session
      const wsUrl = `ws://${window.location.hostname}:8081/ws`
      console.log('WebSocket URL:', wsUrl)
      socket = new WebSocket(wsUrl)
      
      socket.onopen = () => {
        console.log('✅ WebSocket connection established')
        isConnecting = false
        connectionStatus.value = true // Update reactive status
        if (reconnectTimeout) {
          clearTimeout(reconnectTimeout)
          reconnectTimeout = null
        }
        
        console.log('🎯 WebSocket ready to receive notifications')
        // Don't send any initial messages, just wait for notifications from backend
      }
      
      socket.onmessage = (event) => {
        try {
          console.log('📨 WebSocket message received:', event.data)
          const data = JSON.parse(event.data)
          console.log('📨 Parsed WebSocket data:', data)
          
          // Debug: Show notification structure
          if (data.action === 'notification') {
            console.log('🔔 Notification details:', {
              action: data.action,
              notification: data.notification,
              type: data.notification?.type,
              sender: data.notification?.sender,
              target: data.notification?.target
            })
          }
          
          // Notify all listeners about the message
          listeners.forEach((callback, key) => {
            console.log(`📤 Notifying listener: ${key}`)
            callback(data)
          })
        } catch (e) {
          console.error('❌ Error parsing WebSocket message:', e)
        }
      }
      
      socket.onclose = (event) => {
        console.log('❌ WebSocket connection closed:', {
          code: event.code,
          reason: event.reason,
          wasClean: event.wasClean
        })
        
        isConnecting = false
        connectionStatus.value = false // Update reactive status
        
        // Analyze close codes
        if (event.code === 1000) {
          console.log('ℹ️ Normal closure')
        } else if (event.code === 1001) {
          console.log('ℹ️ Going away')
        } else if (event.code === 1006) {
          console.log('⚠️ Abnormal closure (network issue or server crash)')
        } else if (event.code === 1008) {
          console.log('🚫 Policy violation (possibly authentication issue)')
        } else if (event.code === 1011) {
          console.log('💥 Server error')
        } else {
          console.log('❓ Unknown close code:', event.code)
        }
        
        // Only auto-reconnect for network issues
        if (event.code === 1006 || event.code === 1011) {
          console.log('🔄 Attempting to reconnect WebSocket...')
          // Schedule reconnection
          if (!reconnectTimeout) {
            reconnectTimeout = setTimeout(() => {
              this.connect()
            }, RECONNECT_DELAY)
          }
        } else {
          console.log('🚫 Not reconnecting due to close code:', event.code)
        }
      }
      
      socket.onerror = (error) => {
        console.error('❌ WebSocket error:', error)
        isConnecting = false
        connectionStatus.value = false // Update reactive status
        // Don't close here, let onclose handle it
      }
    } catch (e) {
      console.error('❌ Error setting up WebSocket:', e)
      isConnecting = false
    }
  },
  
  // Disconnect from the WebSocket server
  disconnect() {
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout)
      reconnectTimeout = null
    }
    
    if (socket) {
      socket.close()
      socket = null
    }
    
    connectionStatus.value = false // Update reactive status
    
    // Clear all listeners
    listeners.clear()
  },
  
  // Add a listener for WebSocket messages
  addListener(key, callback) {
    listeners.set(key, callback)
    
    // Connect if not already connected
    this.connect()
    
    // Return a function to remove the listener
    return () => {
      listeners.delete(key)
      
      // If no more listeners, disconnect
      if (listeners.size === 0) {
        this.disconnect()
      }
    }
  },
  
  // Send a message through the WebSocket
  send(message) {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(typeof message === 'string' ? message : JSON.stringify(message))
      return true
    }
    return false
  },
  
  // Check if WebSocket is connected
  isConnected() {
    return socket && socket.readyState === WebSocket.OPEN
  },
  
  // Get reactive connection status
  getConnectionStatus() {
    return connectionStatus
  },
  
  // Get listeners map (for debugging/testing)
  getListeners() {
    return listeners
  }
}

export default wsService
