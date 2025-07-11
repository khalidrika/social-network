import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useGroupStore = defineStore('group', () => {
  // State
  const groups = ref([])
  const userGroups = ref([])
  const publicGroups = ref([])
  const currentGroup = ref(null)
  const members = ref([])
  const posts = ref([])
  const events = ref([])
  const invitations = ref([])
  const requests = ref([])
  const loading = ref(false)
  const error = ref(null)
  const eventsUpdateCounter = ref(0) // Add counter to force reactivity

  // Getters
  const currentGroupMembers = computed(() => members.value)
  const currentGroupPosts = computed(() => posts.value)
  const currentGroupEvents = computed(() => events.value)
  const groupInvitations = computed(() => invitations.value)
  const groupRequests = computed(() => requests.value)
  const discoveryGroups = computed(() => publicGroups.value.filter(g => !g.member))
  const isCurrentUserMember = computed(() => {
    return currentGroup.value?.member || false
  })
  const isCurrentUserAdmin = computed(() => {
    return currentGroup.value?.admin || false
  })

  // Actions
  const fetchGroups = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await fetch('http://localhost:8081/allGroups', {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        groups.value = data.groups || []
      } else {
        throw new Error('Failed to fetch groups')
      }
    } catch (err) {
      error.value = err.message
      console.error('Error fetching groups:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchUserGroups = async () => {
    try {
      const response = await fetch('http://localhost:8081/userGroups', {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        userGroups.value = data.groups || []
      } else {
        throw new Error('Failed to fetch user groups')
      }
    } catch (err) {
      error.value = err.message
      console.error('Error fetching user groups:', err)
    }
  }

  const fetchPublicGroups = async () => {
    try {
      const response = await fetch('http://localhost:8081/allGroups', {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        publicGroups.value = data.groups || []
      } else {
        throw new Error('Failed to fetch public groups')
      }
    } catch (err) {
      error.value = err.message
      console.error('Error fetching public groups:', err)
    }
  }

  const fetchNotifications = async () => {
    try {
      const response = await fetch('http://localhost:8081/notifications', {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        // Filter for group invitations
        invitations.value = (data.notifications || [])
          .filter(notif => notif.type === 'GROUP_INVITE')
          .map(notif => ({
            id: notif.id,
            group: notif.group || { name: 'Groupe', image: null, memberCount: 0 },
            inviter: notif.sender || { firstName: 'Utilisateur', lastName: '' },
            createdAt: notif.createdAt
          }))
        
        // Filter for group requests (if user is admin)
        requests.value = (data.notifications || [])
          .filter(notif => notif.type === 'GROUP_REQUEST')
          .map(notif => ({
            id: notif.id,
            user: notif.user || { firstName: 'Utilisateur', lastName: '', profilePic: null },
            group: notif.group || { name: 'Groupe' },
            createdAt: notif.createdAt
          }))
      } else {
        console.error('Failed to load notifications')
        invitations.value = []
        requests.value = []
      }
    } catch (err) {
      console.error('Error loading notifications:', err)
      invitations.value = []
      requests.value = []
    }
  }

  const loadGroupsData = async () => {
    loading.value = true
    try {
      await Promise.all([
        fetchUserGroups(),
        fetchPublicGroups(),
        fetchNotifications()
      ])
    } catch (error) {
      console.error('Store: Error loading groups data:', error)
    } finally {
      loading.value = false
    }
  }

  const fetchGroupDetails = async (groupId) => {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(`http://localhost:8081/groupInfo?groupId=${groupId}`, {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        currentGroup.value = data.groups && data.groups.length > 0 ? data.groups[0] : null
      } else {
        throw new Error('Failed to fetch group details')
      }
    } catch (err) {
      error.value = err.message
      console.error('Error fetching group details:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchGroupMembers = async (groupId) => {
    try {
      const response = await fetch(`http://localhost:8081/groupMembers?groupId=${groupId}`, {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        members.value = data.users || []
      }
    } catch (err) {
      console.error('Error fetching group members:', err)
    }
  }

  const fetchGroupPosts = async (groupId) => {
    try {
      const response = await fetch(`http://localhost:8081/groupPosts?groupId=${groupId}`, {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        posts.value = data.posts || []
      }
    } catch (err) {
      console.error('Error fetching group posts:', err)
    }
  }

  let lastFetchTime = 0
  const FETCH_COOLDOWN = 1000 // 1 second cooldown

  const fetchGroupEvents = async (groupId) => {
    const now = Date.now()
    if (loading.value || (now - lastFetchTime < FETCH_COOLDOWN)) {
      return
    }
    
    lastFetchTime = now
    loading.value = true
    try {
      const response = await fetch(`http://localhost:8081/getGroupEvents?groupId=${groupId}`, {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        events.value = Array.isArray(data) ? data : []
      } else {
        const errorText = await response.text()
        console.error('Events fetch failed:', response.status, errorText)
        events.value = []
      }
    } catch (err) {
      console.error('Error fetching group events:', err)
      events.value = []
    } finally {
      loading.value = false
    }
  }

  const updateEventResponse = async (eventId, response) => {
    try {
      const res = await fetch('http://localhost:8081/updateEventResponse', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ eventId, response })
      })
      
      if (res.ok) {
        // Update the event in the local events array
        const eventIndex = events.value.findIndex(e => e.id === eventId)
        if (eventIndex !== -1) {
          const event = events.value[eventIndex]
          
          // Update user response
          event.userResponse = response
          
          // Update counts (simplified approach - refetch events for accuracy)
          if (currentGroup.value?.id) {
            await fetchGroupEvents(currentGroup.value.id)
          }
        }
        return true
      }
      return false
    } catch (err) {
      console.error('Error updating event response:', err)
      return false
    }
  }

  const createGroupPost = async (groupId, content, image = null) => {
    try {
      const formData = new FormData()
      formData.append('groupId', groupId)
      formData.append('body', content)  // Backend expects 'body', not 'content'
      
      if (image) {
        formData.append('image', image)
      }

      const response = await fetch('http://localhost:8081/newGroupPost', {
        method: 'POST',
        credentials: 'include',
        body: formData
      })

      if (response.ok) {
        // Refresh group posts to include the new post
        await fetchGroupPosts(groupId)
        return { success: true }
      } else {
        const errorData = await response.json().catch(() => ({ message: 'Unknown error' }))
        return { 
          success: false, 
          message: errorData.message || 'Erreur lors de la création du post'
        }
      }
    } catch (err) {
      console.error('Error creating group post:', err)
      return { 
        success: false, 
        message: 'Erreur de connexion'
      }
    }
  }

  const createGroupPostComment = async (postId, content, image = null) => {
    try {
      const formData = new FormData()
      formData.append('postid', postId)  // Backend expects 'postid', not 'postId'
      formData.append('body', content)   // Backend expects 'body', not 'content'
      
      if (image) {
        formData.append('image', image)
      }

      const response = await fetch('http://localhost:8081/newComment', {
        method: 'POST',
        credentials: 'include',
        body: formData
      })

      if (response.ok) {
        // Refresh group posts to include the new comment
        if (currentGroup.value?.id) {
          await fetchGroupPosts(currentGroup.value.id)
        }
        return { success: true }
      } else {
        const errorData = await response.json().catch(() => ({ message: 'Unknown error' }))
        return { 
          success: false, 
          message: errorData.message || 'Erreur lors de la création du commentaire'
        }
      }
    } catch (err) {
      console.error('Error creating group post comment:', err)
      return { 
        success: false, 
        message: 'Erreur de connexion'
      }
    }
  }

  const joinGroup = async (groupId, isPublic = true) => {
    const endpoint = isPublic 
      ? 'http://localhost:8081/joinPublicGroup'
      : `http://localhost:8081/newGroupRequest?groupId=${groupId}`
    
    try {
      // Set group loading state
      setGroupLoading(groupId, true)

      const requestOptions = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include'
      }

      // For public groups, send groupId in body, for private groups it's in URL
      if (isPublic) {
        requestOptions.body = JSON.stringify({ groupId })
      }

      const response = await fetch(endpoint, requestOptions)

      if (response.ok) {
        if (isPublic) {
          // Update current group membership status
          if (currentGroup.value && currentGroup.value.id === groupId) {
            currentGroup.value.member = true
          }
          // Update public groups list
          const group = publicGroups.value.find(g => g.id === groupId)
          if (group) {
            group.member = true
          }
          // Refresh group data
          await fetchGroupDetails(groupId)
          await fetchGroupMembers(groupId)
          await loadGroupsData() // Refresh all data
        } else {
          // For private groups, mark request as pending locally first
          const group = publicGroups.value.find(g => g.id === groupId)
          if (group) {
            group.requestPending = true
          }
          // Also update currentGroup if it's the same group
          if (currentGroup.value && currentGroup.value.id === groupId) {
            currentGroup.value.requestPending = true
          }
          // Then refresh all group data to get the updated state from backend
          await loadGroupsData()
        }
        return { success: true, isPublic }
      }
      return { success: false, isPublic }
    } catch (err) {
      console.error('Error joining group:', err)
      return { success: false, isPublic }
    } finally {
      // Reset group loading state
      setGroupLoading(groupId, false)
    }
  }

  const leaveGroup = async (groupId) => {
    try {
      const response = await fetch('http://localhost:8081/leaveGroup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ groupId })
      })

      if (response.ok) {
        // Update current group membership status
        if (currentGroup.value && currentGroup.value.id === groupId) {
          currentGroup.value.member = false
          // Clear any pending request state when leaving
          currentGroup.value.requestPending = false
        }
        // Also update in publicGroups array
        const group = publicGroups.value.find(g => g.id === groupId)
        if (group) {
          group.member = false
          group.requestPending = false
        }
        // Refresh group data
        await fetchGroupDetails(groupId)
        await fetchGroupMembers(groupId)
        return { success: true }
      } else {
        // Handle different error cases
        const errorData = await response.json().catch(() => ({ message: 'Unknown error' }))
        if (response.status === 403) {
          return { 
            success: false, 
            error: 'admin_cannot_leave',
            message: errorData.message || 'Les administrateurs ne peuvent pas quitter leur propre groupe'
          }
        }
        return { 
          success: false, 
          error: 'general_error',
          message: errorData.message || 'Erreur lors de la sortie du groupe'
        }
      }
    } catch (err) {
      console.error('Error leaving group:', err)
      return { 
        success: false, 
        error: 'network_error',
        message: 'Erreur de connexion'
      }
    }
  }

  const clearCurrentGroup = () => {
    currentGroup.value = null
    members.value = []
    posts.value = []
    events.value = []
  }

  const respondToInvitation = async (invitationId, action) => {
    try {
      const response = await fetch('http://localhost:8081/responseGroupInvite', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          requestId: invitationId, 
          response: action === 'accept' ? 'ACCEPT' : 'DECLINE' 
        })
      })
      
      if (response.ok) {
        // Remove invitation from list
        invitations.value = invitations.value.filter(inv => inv.id !== invitationId)
        if (action === 'accept') {
          // Reload groups data
          await loadGroupsData()
        }
        return true
      }
      return false
    } catch (err) {
      console.error('Error responding to invitation:', err)
      return false
    }
  }

  const respondToRequest = async (requestId, action, groupId) => {
    try {
      const response = await fetch('http://localhost:8081/responseGroupRequest', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          requestId: requestId, 
          groupId: groupId,
          response: action === 'accept' ? 'accept' : 'decline'
        })
      })
      
      if (response.ok) {
        // Remove request from list
        requests.value = requests.value.filter(req => req.id !== requestId)
        if (action === 'accept') {
          // Reload groups data to update member counts
          await loadGroupsData()
        }
        return true
      }
      return false
    } catch (err) {
      console.error('Error responding to request:', err)
      return false
    }
  }

  const setGroupLoading = (groupId, isLoading) => {
    const group = publicGroups.value.find(g => g.id === groupId)
    if (group) {
      group.loading = isLoading
    }
  }

  return {
    // State
    groups,
    userGroups,
    publicGroups,
    currentGroup,
    members,
    posts,
    events,
    invitations,
    requests,
    loading,
    error,
    eventsUpdateCounter,
    
    // Getters
    currentGroupMembers,
    currentGroupPosts,
    currentGroupEvents,
    groupInvitations,
    groupRequests,
    discoveryGroups,
    isCurrentUserMember,
    isCurrentUserAdmin,
    
    // Actions
    fetchGroups,
    fetchUserGroups,
    fetchPublicGroups,
    fetchNotifications,
    loadGroupsData,
    fetchGroupDetails,
    fetchGroupMembers,
    fetchGroupPosts,
    fetchGroupEvents,
    updateEventResponse,
    createGroupPost,
    createGroupPostComment,
    joinGroup,
    leaveGroup,
    respondToInvitation,
    respondToRequest,
    clearCurrentGroup,
    setGroupLoading
  }
})
