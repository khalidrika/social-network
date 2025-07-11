<template>
  <div class="invitation-card">
    <!-- Group Info -->
    <div class="invitation-info">
      <div class="group-avatar">
        <img 
          v-if="invitation.group.image" 
          :src="invitation.group.image" 
          :alt="invitation.group.name"
          @error="handleImageError"
        />
        <div v-else class="default-avatar">
          <Users :size="24" />
        </div>
      </div>
      
      <div class="invitation-details">
        <h3 class="group-name">{{ invitation.group.name }}</h3>
        <p class="invitation-text">
          <strong>{{ invitation.inviter.firstName }} {{ invitation.inviter.lastName }}</strong>
          vous a invité à rejoindre ce groupe
        </p>
        <div class="invitation-meta">
          <span class="member-count">
            <Users :size="16" />
            {{ invitation.group.memberCount || 0 }} membres
          </span>
          <span class="invitation-date">
            <Calendar :size="16" />
            {{ formatDate(invitation.createdAt) }}
          </span>
        </div>
      </div>
    </div>

    <!-- Actions -->
    <div class="invitation-actions">
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
import { Users, Calendar, UserCheck, UserX } from 'lucide-vue-next'

const props = defineProps({
  invitation: {
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
    await emit('respond', props.invitation.id, action)
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
.invitation-card {
  background: rgba(15, 15, 23, 0.8);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.invitation-card:hover {
  border-color: rgba(232, 121, 198, 0.3);
  box-shadow: 0 4px 20px rgba(232, 121, 198, 0.1);
}

.invitation-info {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.group-avatar {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  overflow: hidden;
  flex-shrink: 0;
}

.group-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-avatar {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.avatar-icon {
  font-size: 1.5rem;
}

.invitation-details {
  flex: 1;
  min-width: 0;
}

.group-name {
  font-size: 1.125rem;
  font-weight: 700;
  color: white;
  margin: 0 0 8px 0;
  line-height: 1.3;
}

.invitation-text {
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 12px 0;
  line-height: 1.5;
}

.invitation-text strong {
  color: white;
}

.invitation-meta {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.member-count,
.invitation-date {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
}

.icon {
  font-size: 0.875rem;
}

.invitation-actions {
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
  .invitation-info {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  
  .invitation-actions {
    flex-direction: column;
  }
  
  .btn-accept,
  .btn-decline {
    flex: none;
  }
}
</style>
