<template>
  <div class="modal-overlay" @click="handleOverlayClick">
    <div class="modal-content" @click.stop>
      <!-- Header -->
      <div class="modal-header">
        <h2>Inviter des abonnés</h2>
        <p>Invitez vos abonnés à rejoindre {{ group.name }}</p>
        <button @click="$emit('close')" class="close-btn">
          <X :size="20" />
        </button>
      </div>

      <!-- Search Users -->
      <div class="modal-body">
        <div class="search-section">
          <div class="search-input-wrapper">
            <Search :size="20" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Rechercher parmi vos abonnés..."
              class="search-input"
            />
          </div>
        </div>

        <!-- Followers List -->
        <div v-if="filteredFollowers.length > 0" class="users-list">
          <h3>{{ searchQuery ? 'Résultats de recherche' : 'Vos abonnés' }} ({{ filteredFollowers.length }})</h3>
          <div
            v-for="user in filteredFollowers"
            :key="user.id"
            class="user-item"
            :class="{ 'invited': invitedUsers.includes(user.id) }"
          >
            <div class="user-info">
              <div class="user-avatar">
                <img 
                  v-if="user.avatar" 
                  :src="user.avatar" 
                  :alt="`${user.firstName || user.first_name} ${user.lastName || user.last_name}`"
                  @error="handleImageError"
                />
                <div v-else class="default-avatar">
                  <Users :size="24" />
                </div>
              </div>
              <div class="user-details">
                <h4>{{ user.firstName || user.first_name }} {{ user.lastName || user.last_name }}</h4>
                <p v-if="user.nickname" class="user-nickname">@{{ user.nickname }}</p>
              </div>
            </div>
            <button
              @click="toggleInvitation(user)"
              class="select-btn"
              :class="{ 
                'selected': invitedUsers.includes(user.id),
                'loading': loadingUsers.includes(user.id)
              }"
              :disabled="loadingUsers.includes(user.id)"
            >
              <component 
                :is="loadingUsers.includes(user.id) ? Loader2 : (invitedUsers.includes(user.id) ? UserMinus : UserPlus)" 
                :size="18"
                :class="{ 'animate-spin': loadingUsers.includes(user.id) }"
              />
              {{ 
                loadingUsers.includes(user.id) ? 'Envoi...' : 
                (invitedUsers.includes(user.id) ? 'Invité' : 'Inviter') 
              }}
            </button>
          </div>
        </div>

        <!-- No Results -->
        <div v-if="!followersLoading && allFollowers.length === 0" class="no-results">
          <div class="no-results-icon">
            <Users :size="48" />
          </div>
          <h3>Aucun abonné trouvé</h3>
          <p>Vous n'avez pas encore d'abonnés à inviter</p>
        </div>

        <div v-if="!followersLoading && allFollowers.length > 0 && filteredFollowers.length === 0" class="no-results">
          <div class="no-results-icon">
            <Users :size="48" />
          </div>
          <h3>Aucun abonné trouvé</h3>
          <p>Aucun abonné ne correspond à votre recherche</p>
        </div>

        <!-- Loading -->
        <div v-if="followersLoading" class="loading-state">
          <Loader2 class="spinner" />
          <p>Chargement de vos abonnés...</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { 
  X, 
  Search, 
  Users, 
  UserPlus, 
  UserMinus, 
  Loader2 
} from 'lucide-vue-next'
import { useFollowStore } from '../../stores/followStore'
import { useUserStore } from '../../stores/userStore'

const props = defineProps({
  group: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'invited'])

const followStore = useFollowStore()
const userStore = useUserStore()

const searchQuery = ref('')
const allFollowers = ref([])
const invitedUsers = ref([])
const loadingUsers = ref([])
const followersLoading = ref(false)

// Computed property to filter followers based on search query
const filteredFollowers = computed(() => {
  if (!searchQuery.value.trim()) {
    return allFollowers.value
  }
  
  const query = searchQuery.value.trim().toLowerCase()
  return allFollowers.value.filter(user => {
    return (
      (user.nickname && user.nickname.toLowerCase().includes(query)) ||
      (user.first_name && user.first_name.toLowerCase().includes(query)) ||
      (user.last_name && user.last_name.toLowerCase().includes(query)) ||
      (user.firstName && user.firstName.toLowerCase().includes(query)) ||
      (user.lastName && user.lastName.toLowerCase().includes(query))
    )
  })
})

// Load followers when component mounts
onMounted(async () => {
  followersLoading.value = true
  try {
    // Fetch followers and existing invitations in parallel
    const [_, existingInvitationsResponse] = await Promise.all([
      followStore.fetchFollowers(userStore.user?.id),
      fetch(`http://localhost:8081/checkGroupInvitations?groupId=${props.group.id}`, {
        credentials: 'include'
      })
    ])
    
    // Process existing invitations
    if (existingInvitationsResponse.ok) {
      const invitationData = await existingInvitationsResponse.json()
      invitedUsers.value = invitationData.invitedUsers || []
    }
    
    // Map the followers to ensure consistent property names
    allFollowers.value = (followStore.followers || []).map(user => ({
      ...user,
      // Ensure we have both camelCase and snake_case for compatibility
      firstName: user.firstName || user.first_name,
      lastName: user.lastName || user.last_name,
      first_name: user.first_name || user.firstName,
      last_name: user.last_name || user.lastName,
      avatar: user.image ? `/${user.image}` : null
    })).filter(user => 
      // Filter out the group admin and any existing members
      user.id !== props.group.adminId
    )
  } catch (err) {
    console.error('Error loading followers:', err)
    allFollowers.value = []
  } finally {
    followersLoading.value = false
  }
})

const handleOverlayClick = (event) => {
  if (event.target === event.currentTarget) {
    emit('close')
  }
}

const handleImageError = (event) => {
  event.target.style.display = 'none'
  event.target.nextElementSibling.style.display = 'flex'
}

const toggleInvitation = async (user) => {
  const isInvited = invitedUsers.value.includes(user.id)
  
  // Add user to loading state
  loadingUsers.value.push(user.id)
  
  try {
    if (isInvited) {
      // Uninvite the user by calling the new endpoint
      const response = await fetch('http://localhost:8081/cancelGroupInvite', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          groupId: props.group.id,
          targetId: user.id
        })
      })

      if (response.ok) {
        const index = invitedUsers.value.indexOf(user.id)
        if (index > -1) {
          invitedUsers.value.splice(index, 1)
        }
      } else {
        throw new Error('Erreur lors de l\'annulation de l\'invitation')
      }
    } else {
      // Send invitation immediately
      const response = await fetch('http://localhost:8081/newGroupInvite', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          id: props.group.id,
          invitations: [user.id]
        })
      })

      if (response.ok) {
        // Add to invited users list
        invitedUsers.value.push(user.id)
        emit('invited', 1)
      } else {
        throw new Error('Erreur lors de l\'envoi de l\'invitation')
      }
    }
  } catch (error) {
    console.error('Error toggling invitation:', error)
    // Silently fail - no alert message
  } finally {
    // Remove user from loading state
    const loadingIndex = loadingUsers.value.indexOf(user.id)
    if (loadingIndex > -1) {
      loadingUsers.value.splice(loadingIndex, 1)
    }
  }
}

const getUserName = (userId) => {
  const user = filteredFollowers.value.find(u => u.id === userId)
  return user ? `${user.firstName || user.first_name} ${user.lastName || user.last_name}` : 'Utilisateur'
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
  backdrop-filter: blur(8px);
}

.modal-content {
  background: rgba(15, 15, 23, 0.95);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow: hidden;
  backdrop-filter: blur(20px);
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  position: relative;
}

.modal-header h2 {
  color: white;
  margin: 0 0 8px 0;
  font-size: 1.5rem;
  font-weight: 700;
}

.modal-header p {
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
  font-size: 0.95rem;
}

.close-btn {
  position: absolute;
  top: 20px;
  right: 20px;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.modal-body {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.search-section {
  margin-bottom: 24px;
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-input-wrapper svg {
  position: absolute;
  left: 16px;
  color: rgba(255, 255, 255, 0.5);
  z-index: 1;
}

.search-input {
  width: 100%;
  padding: 12px 16px 12px 48px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.05);
  color: white;
  font-size: 0.95rem;
  transition: all 0.3s ease;
}

.search-input:focus {
  outline: none;
  border-color: rgba(232, 121, 198, 0.5);
  background: rgba(255, 255, 255, 0.1);
}

.search-input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.users-list {
  margin-bottom: 24px;
}

.users-list h3 {
  color: white;
  margin: 0 0 16px 0;
  font-size: 1.1rem;
  font-weight: 600;
}

.user-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  margin-bottom: 12px;
  transition: all 0.3s ease;
}

.user-item:hover {
  border-color: rgba(232, 121, 198, 0.3);
  background: rgba(255, 255, 255, 0.05);
}

.user-item.invited {
  border-color: rgba(34, 197, 94, 0.5);
  background: rgba(34, 197, 94, 0.1);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-avatar {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.user-details h4 {
  color: white;
  margin: 0 0 4px 0;
  font-size: 1rem;
  font-weight: 600;
}

.user-nickname {
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
  font-size: 0.875rem;
}

.select-btn {
  padding: 8px 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  font-weight: 500;
}

.select-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.select-btn.selected {
  background: rgba(34, 197, 94, 0.2);
  border-color: rgba(34, 197, 94, 0.5);
  color: #22c55e;
}

.select-btn.loading {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.5);
  color: #3b82f6;
  cursor: not-allowed;
}

.no-results,
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  text-align: center;
}

.no-results-icon {
  color: rgba(255, 255, 255, 0.3);
  margin-bottom: 16px;
}

.no-results h3 {
  color: white;
  margin: 0 0 8px 0;
  font-size: 1.1rem;
}

.no-results p {
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
}

.spinner {
  animation: spin 1s linear infinite;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 16px;
}

.loading-state p {
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
}

.modal-actions {
  display: flex;
  gap: 12px;
  padding: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  justify-content: flex-end;
}

.btn-primary,
.btn-secondary {
  padding: 12px 24px;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.95rem;
}

.btn-primary {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Responsive */
@media (max-width: 768px) {
  .modal-content {
    margin: 0;
    border-radius: 20px 20px 0 0;
    max-height: 100vh;
  }
  
  .user-item {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .user-info {
    justify-content: center;
  }
}
</style>
