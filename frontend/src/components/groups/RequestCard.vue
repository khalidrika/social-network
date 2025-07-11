<template>
  <div class="request-card">
    <!-- User Info -->
    <div class="request-info">
      <div class="user-avatar">
        <img 
          v-if="request.user.profilePic" 
          :src="getFullImageUrl(request.user.profilePic)" 
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
          Demande à rejoindre <strong>{{ request.group.name }}</strong>
        </p>
        <div class="request-meta">
          <span class="request-date">
            <Calendar :size="16" />
            {{ formatDate(request.createdAt) }}
          </span>
          <span v-if="request.user.mutualFriends" class="mutual-friends">
            <Users :size="16" />
            {{ request.user.mutualFriends }} amis en commun
          </span>
        </div>
      </div>
    </div>

    <!-- Actions -->
    <div class="request-actions">
      <button 
        @click="handleResponse('decline')"
        class="btn-decline"
        :disabled="loading"
      >
        <UserX :size="18" />
        Refuser
      </button>
      <button 
        @click="handleResponse('accept')"
        class="btn-accept"
        :disabled="loading"
      >
        <UserCheck :size="18" />
        Accepter
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useMainStore } from '../../stores/postsStore'
import { Users, Calendar, UserCheck, UserX } from 'lucide-vue-next'

const mainStore = useMainStore()
const { getFullImageUrl } = mainStore

const props = defineProps({
  request: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['respond'])

const loading = ref(false)

const handleImageError = (event) => {
  event.target.style.display = 'none'
  event.target.nextElementSibling.style.display = 'flex'
}

const handleResponse = async (action) => {
  loading.value = true
  try {
    await emit('respond', props.request.id, action)
  } finally {
    setTimeout(() => {
      loading.value = false
    }, 500)
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInHours = Math.floor((now - date) / (1000 * 60 * 60))
  
  if (diffInHours < 1) {
    return 'Il y a quelques minutes'
  } else if (diffInHours < 24) {
    return `Il y a ${diffInHours}h`
  } else {
    const diffInDays = Math.floor(diffInHours / 24)
    return `Il y a ${diffInDays}j`
  }
}
</script>

<style scoped>
.request-card {
  background: rgba(15, 15, 23, 0.8);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.request-card:hover {
  border-color: rgba(120, 199, 255, 0.3);
  box-shadow: 0 4px 20px rgba(120, 199, 255, 0.1);
}

.request-info {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
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
  background: linear-gradient(135deg, #78c7ff 0%, #e879c6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.avatar-icon {
  font-size: 1.5rem;
  opacity: 0.8;
}

.request-details {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 1.125rem;
  font-weight: 700;
  color: white;
  margin: 0 0 8px 0;
  line-height: 1.3;
}

.request-text {
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 12px 0;
  line-height: 1.5;
}

.request-text strong {
  color: white;
}

.request-meta {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.request-date,
.mutual-friends {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
}

.icon {
  font-size: 0.875rem;
}

.request-actions {
  display: flex;
  gap: 12px;
}

.btn-accept,
.btn-decline {
  flex: 1;
  padding: 10px 16px;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 0.875rem;
}

.btn-accept {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.btn-accept:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(16, 185, 129, 0.3);
}

.btn-decline {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.btn-decline:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.5);
}

.btn-accept:disabled,
.btn-decline:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* Responsive */
@media (max-width: 768px) {
  .request-info {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  
  .request-actions {
    flex-direction: column;
  }
  
  .btn-accept,
  .btn-decline {
    flex: none;
  }
}
</style>
