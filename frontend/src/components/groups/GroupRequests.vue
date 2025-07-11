<template>
  <div class="group-requests">
    <div class="requests-header">
      <h2>Demandes d'adhésion</h2>
      <p>Gérez les demandes de rejoindre le groupe</p>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <Loader2 class="spinner-icon" />
      <p>Chargement des demandes...</p>
    </div>

    <!-- No Requests -->
    <div v-else-if="!requests.length" class="no-requests">
      <div class="no-requests-icon">
        <Mail :size="64" />
      </div>
      <h3>Aucune demande</h3>
      <p>Il n'y a actuellement aucune demande d'adhésion en attente.</p>
    </div>

    <!-- Requests List -->
    <div v-else class="requests-list">
      <div
        v-for="request in requests"
        :key="request.id"
        class="request-card"
      >
        <!-- User Info -->
        <div class="request-info">
          <div class="user-avatar">
            <img 
              v-if="request.user.avatar" 
              :src="getImageUrl(request.user.avatar)" 
              :alt="`${request.user.firstName} ${request.user.lastName}`"
              @error="handleImageError"
            />
            <div v-else class="default-avatar">
              <Users :size="24" />
            </div>
          </div>
          
          <div class="request-details">
            <h3 class="user-name">
              {{ request.user.firstName }} {{ request.user.lastName }}
            </h3>
            <p class="request-text">
              Demande à rejoindre ce groupe
            </p>
            <div class="request-meta">
              <span v-if="request.user.nickname" class="user-nickname">
                <AtSign :size="16" />
                {{ request.user.nickname }}
              </span>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="request-actions">
          <button 
            @click="handleRequestResponse(request.id, 'decline')"
            class="btn-decline"
            :disabled="processingRequests.includes(request.id)"
          >
            <UserX :size="18" />
            Refuser
          </button>
          <button 
            @click="handleRequestResponse(request.id, 'accept')"
            class="btn-accept"
            :disabled="processingRequests.includes(request.id)"
          >
            <UserCheck :size="18" />
            Accepter
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Users, UserCheck, UserX, Mail, AtSign, Loader2 } from 'lucide-vue-next'

const props = defineProps({
  groupId: {
    type: String,
    required: true
  }
})

const loading = ref(false)
const requests = ref([])
const processingRequests = ref([])

const fetchGroupRequests = async () => {
  loading.value = true
  try {
    const response = await fetch(`http://localhost:8081/groupRequests?groupId=${props.groupId}`, {
      method: 'GET',
      credentials: 'include',
    })

    if (response.ok) {
      const data = await response.json()
      const notifications = data.notifications || []
      requests.value = notifications
    } else {
      console.error('Failed to fetch group requests')
      requests.value = []
    }
  } catch (error) {
    console.error('Error fetching group requests:', error)
    requests.value = []
  } finally {
    loading.value = false
  }
}

const handleRequestResponse = async (requestId, action) => {
  try {
    // Add to processing list
    processingRequests.value.push(requestId)
    
    const payload = {
      groupId: props.groupId,
      requestId: requestId,
      response: action.toLowerCase() // Ensure lowercase
    }
    
    const response = await fetch('http://localhost:8081/responseGroupRequest', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(payload)
    })

    if (response.ok) {
      // Remove the request from the list
      requests.value = requests.value.filter(req => req.id !== requestId)
    } else {
      const errorText = await response.text()
      console.error(`Failed to ${action} request:`, errorText)
    }
  } catch (error) {
    console.error(`Error ${action}ing request:`, error)
  } finally {
    // Remove from processing list
    processingRequests.value = processingRequests.value.filter(id => id !== requestId)
  }
}

const getImageUrl = (imagePath) => {
  if (!imagePath) return null
  if (imagePath.startsWith('http')) return imagePath
  return `http://localhost:8081/${imagePath}`
}

const handleImageError = (event) => {
  event.target.style.display = 'none'
  if (event.target.nextElementSibling) {
    event.target.nextElementSibling.style.display = 'flex'
  }
}

onMounted(() => {
  fetchGroupRequests()
})
</script>

<style scoped>
.group-requests {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.requests-header {
  text-align: center;
  margin-bottom: 30px;
}

.requests-header h2 {
  color: #fff;
  margin-bottom: 8px;
  font-size: 2em;
  font-weight: 700;
}

.requests-header p {
  color: rgba(255, 255, 255, 0.7);
  font-size: 1.1em;
  margin: 0;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: rgba(255, 255, 255, 0.7);
}

.spinner-icon {
  color: #e879c6;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

.no-requests {
  text-align: center;
  padding: 60px 20px;
  color: rgba(255, 255, 255, 0.7);
}

.no-requests-icon {
  margin-bottom: 20px;
  color: rgba(255, 255, 255, 0.4);
}

.no-requests h3 {
  margin-bottom: 10px;
  color: #fff;
  font-size: 1.5em;
}

.no-requests p {
  font-size: 1.1em;
  line-height: 1.6;
  max-width: 400px;
  margin: 0 auto;
}

.requests-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.request-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.request-card:hover {
  transform: translateY(-2px);
  border-color: rgba(232, 121, 198, 0.3);
  box-shadow: 0 8px 32px rgba(232, 121, 198, 0.1);
}

.request-info {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
}

.user-avatar {
  width: 60px;
  height: 60px;
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

.request-details {
  flex: 1;
}

.user-name {
  color: #fff;
  font-size: 1.2em;
  font-weight: 600;
  margin: 0 0 8px 0;
}

.request-text {
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 8px 0;
  font-size: 0.95em;
}

.request-meta {
  display: flex;
  gap: 16px;
  align-items: center;
}

.user-nickname {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.9em;
  display: flex;
  align-items: center;
  gap: 4px;
}

.request-actions {
  display: flex;
  gap: 12px;
  flex-shrink: 0;
}

.btn-decline,
.btn-accept {
  padding: 10px 18px;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9em;
}

.btn-decline {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.btn-decline:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.25);
  border-color: rgba(239, 68, 68, 0.5);
  transform: translateY(-1px);
}

.btn-accept {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.btn-accept:hover:not(:disabled) {
  background: rgba(34, 197, 94, 0.25);
  border-color: rgba(34, 197, 94, 0.5);
  transform: translateY(-1px);
}

.btn-decline:disabled,
.btn-accept:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .request-card {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
  
  .request-info {
    flex-direction: column;
    text-align: center;
    gap: 12px;
  }
  
  .request-actions {
    width: 100%;
    justify-content: center;
  }
}
</style>
