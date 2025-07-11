import { defineStore } from 'pinia'
import { ref, computed, nextTick } from 'vue'
import { useGroupStore } from './groupStore'
import { useFollowStore } from './followStore'

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref([])
  const loading = ref(false)
  const hasNewNotifications = ref(false)
  const processingRequests = ref(new Set()) // Track which requests are being processed
  const error = ref(null)
  
  // Computed properties
  const unreadCount = computed(() => {
    return notifications.value.filter(n => !n.read).length
  })
  
  // For actionable notifications (follow requests, group invites/requests), show all regardless of read status
  // These need user action and shouldn't disappear when marked as read
  const followRequests = computed(() => {
    return notifications.value.filter(n => n.type === 'FOLLOW')
  })
  
  const groupInvites = computed(() => {
    return notifications.value.filter(n => n.type === 'GROUP_INVITE')
  })
  
  const groupRequests = computed(() => {
    return notifications.value.filter(n => n.type === 'GROUP_REQUEST')
  })
  
  // For informational notifications (events), show all recent ones
  const eventNotifications = computed(() => {
    return notifications.value.filter(n => n.type === 'EVENT')
  })
  
  // ONLY UNREAD event notifications for dropdown (Instagram/Facebook style)
  const unreadEventNotifications = computed(() => {
    return notifications.value
      .filter(n => n.type === 'EVENT' && !n.read)
      .sort((a, b) => new Date(b.createdAt || b.timestamp) - new Date(a.createdAt || a.timestamp))
      .slice(0, 3)
  })
  
  // All notifications for dropdown preview (limited) - show actionable + UNREAD events only
  const recentUnreadNotifications = computed(() => {
    const pendingFollowRequests = followRequests.value.slice(0, 2)
    const pendingGroupInvites = groupInvites.value.slice(0, 2)
    const pendingGroupRequests = groupRequests.value.slice(0, 2)
    const unreadEvents = unreadEventNotifications.value // Only unread events!
    
    return [
      ...pendingFollowRequests,
      ...pendingGroupInvites, 
      ...pendingGroupRequests,
      ...unreadEvents
    ].slice(0, 6) // Max 6 notifications in dropdown
  })
  
  // Actions
  async function fetchNotifications() {
    loading.value = true
    error.value = null
    try {
      const res = await fetch('http://localhost:8081/notifications', { credentials: 'include' })
      const data = await res.json()
      
      notifications.value = data.notifications || []
      hasNewNotifications.value = false
      
      // Also update notifications in other stores
      const groupStore = useGroupStore()
      const followStore = useFollowStore()
      
      // Run these in parallel to improve performance
      await Promise.all([
        groupStore.fetchNotifications(),
        followStore.fetchFollowRequests()
      ])
    } catch (e) {
      console.error('Error fetching notifications:', e)
      error.value = 'Erreur lors du chargement des notifications'
    } finally {
      loading.value = false
    }
  }
  
  async function respondToFollowRequest(requestId, response) {
    processingRequests.value.add(requestId)
    error.value = null
    
    try {
      const followStore = useFollowStore()
      const result = await followStore.respondToFollowRequest(requestId, response)
      if (!result.error) {
        await fetchNotifications()
        // Also refresh follow requests in the follow store
        await followStore.fetchFollowRequests()
      } else {
        error.value = result.error
      }
    } catch (e) {
      error.value = 'Erreur lors du traitement de la demande'
      console.error('Error responding to follow request:', e)
    } finally {
      processingRequests.value.delete(requestId)
    }
  }
   // Using the groupStore for group invites/requests to avoid duplication
  async function respondToGroupInvite(inviteId, response) {
    processingRequests.value.add(inviteId)
    error.value = null
    
    try {
      const groupStore = useGroupStore()
      const action = response === 'ACCEPT' ? 'accept' : 'decline'
      const result = await groupStore.respondToInvitation(inviteId, action)
      if (result) {
        await fetchNotifications()
      } else {
        error.value = 'Erreur lors du traitement de l\'invitation'
      }
    } catch (e) {
      error.value = 'Erreur lors du traitement de l\'invitation'
      console.error('Error responding to group invite:', e)
    } finally {
      processingRequests.value.delete(inviteId)
    }
  }

  async function respondToGroupRequest(requestId, response, groupId) {
    processingRequests.value.add(requestId)
    error.value = null
    
    try {
      const groupStore = useGroupStore()
      const action = response === 'ACCEPT' ? 'accept' : 'decline'
      const result = await groupStore.respondToRequest(requestId, action, groupId)
      if (result) {
        await fetchNotifications()
      } else {
        error.value = 'Erreur lors du traitement de la demande'
      }
    } catch (e) {
      error.value = 'Erreur lors du traitement de la demande'
      console.error('Error responding to group request:', e)
    } finally {
      processingRequests.value.delete(requestId)
    }
  }
  
  async function handleNewNotification(notification) {
    // Add notification to the list if it's not already there
    const exists = notifications.value.some(n => n.id === notification.id)
    if (!exists) {
      // Create a new array to trigger reactivity
      const newNotification = {
        ...notification,
        read: false, // Ensure new notifications are marked as unread
        createdAt: notification.createdAt || new Date().toISOString()
      }
      
      // Use a new array reference to ensure reactivity
      notifications.value = [newNotification, ...notifications.value]
      
      // Set the flag for new notifications
      hasNewNotifications.value = true
      
      // Force Vue to update the DOM on the next tick
      await nextTick()
    }
  }
  
  async function markAllAsRead() {
    try {
      const response = await fetch('http://localhost:8081/notifications/markAllAsRead', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        // Update local state
        notifications.value = notifications.value.map(n => ({...n, read: true}))
        hasNewNotifications.value = false
      } else {
        console.error('Failed to mark notifications as read:', response.status)
        error.value = 'Erreur lors du marquage des notifications'
      }
    } catch (e) {
      console.error('Error marking notifications as read:', e)
      error.value = 'Erreur réseau'
    }
  }

  async function markAsRead(notificationId) {
    try {
      const response = await fetch('http://localhost:8081/notifications/markAsRead', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ notificationId })
      })
      
      if (response.ok) {
        // Update local state
        const notification = notifications.value.find(n => n.id === notificationId)
        if (notification) {
          notification.read = true
        }
      } else {
        console.error('Failed to mark notification as read:', response.status)
        error.value = 'Erreur lors du marquage de la notification'
      }
    } catch (e) {
      console.error('Error marking notification as read:', e)
      error.value = 'Erreur réseau'
    }
  }

  // Mark notification as read when viewed (for event notifications)
  async function markAsReadOnView(notificationId) {
    const notification = notifications.value.find(n => n.id === notificationId)
    if (notification && !notification.read && notification.type === 'EVENT') {
      await markAsRead(notificationId)
    }
  }

  async function dismissNotification(notificationId) {
    try {
      const response = await fetch(`http://localhost:8081/dismissNotification`, {
        method: 'DELETE',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ notificationId })
      })

      if (response.ok) {
        // Remove from local state
        notifications.value = notifications.value.filter(n => n.id !== notificationId)
      } else {
        const errorText = await response.text()
        console.error('Failed to dismiss notification:', response.status, errorText)
        error.value = 'Erreur lors de la suppression de la notification'
      }
    } catch (e) {
      console.error('Error dismissing notification:', e)
      error.value = 'Erreur réseau'
    }
  }

  function dismissError() {
    error.value = null
  }

  function isProcessing(requestId) {
    return processingRequests.value.has(requestId)
  }

  function formatNotificationTime(timestamp) {
    if (!timestamp) return ''
    
    const now = new Date()
    const notifTime = new Date(timestamp)
    const diffMs = now - notifTime
    const diffMins = Math.floor(diffMs / (1000 * 60))
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
    
    if (diffMins < 1) return 'À l\'instant'
    if (diffMins < 60) return `Il y a ${diffMins} min`
    if (diffHours < 24) return `Il y a ${diffHours}h`
    if (diffDays < 7) return `Il y a ${diffDays} jour${diffDays > 1 ? 's' : ''}`
    
    return notifTime.toLocaleDateString('fr-FR')
  }
  
  return {
    notifications,
    loading,
    hasNewNotifications,
    processingRequests,
    error,
    unreadCount,
    followRequests,
    groupInvites,
    groupRequests,
    eventNotifications,
    unreadEventNotifications,
    recentUnreadNotifications,
    fetchNotifications,
    respondToFollowRequest,
    respondToGroupInvite,
    respondToGroupRequest,
    handleNewNotification,
    markAllAsRead,
    markAsRead,
    markAsReadOnView,
    dismissNotification,
    dismissError,
    isProcessing,
    formatNotificationTime
  }
})
