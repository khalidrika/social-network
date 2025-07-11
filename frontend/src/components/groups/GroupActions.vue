<template>
  <div class="group-actions">
    <!-- Join button for non-members -->
    <button 
      v-if="canJoinGroup && !group.requestPending"
      @click="handleJoin" 
      class="btn-primary"
      :disabled="isJoining"
    >
      <component 
        :is="isJoining ? Loader2 : Plus" 
        :class="{ 'animate-spin': isJoining }"
        :size="18"
      />
      {{ isJoining ? 'Rejoindre...' : (group.privacy === 'public' ? 'Rejoindre le groupe' : 'Demander à rejoindre') }}
    </button>
    
    <!-- Request pending for non-members -->
    <div 
      v-if="canJoinGroup && group.requestPending"
      class="request-pending-badge"
    >
      <Clock :size="18" />
      Demande envoyée
    </div>
    
    <!-- Invite button for members -->
    <button 
      v-if="isMember || isAdmin"
      @click="$emit('invite')" 
      class="action-badge invite-action"
    >
      <UserPlus :size="18" />
      Inviter
    </button>
    
    <!-- Leave button for members (not admins) -->
    <button 
      v-if="canLeaveGroup"
      @click="handleLeave" 
      class="action-badge leave-action"
      :disabled="isLeaving"
      title="Quitter le groupe"
    >
      <component 
        :is="isLeaving ? Loader2 : LogOut" 
        :class="{ 'animate-spin': isLeaving }"
        :size="18"
      />
      {{ isLeaving ? '' : 'Quitter' }}
    </button>
    
    <button @click="$router.go(-1)" class="btn-secondary">
      <ArrowLeft :size="18" />
      Retour
    </button>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Plus, 
  LogOut, 
  ArrowLeft, 
  Loader2,
  UserPlus,
  Clock
} from 'lucide-vue-next'

const router = useRouter()

const props = defineProps({
  group: {
    type: Object,
    required: true
  },
  isMember: {
    type: Boolean,
    default: false
  },
  isAdmin: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['join', 'leave', 'invite'])

const isJoining = ref(false)
const isLeaving = ref(false)

const canJoinGroup = computed(() => {
  return props.group && 
         (props.group.privacy === 'public' || props.group.privacy === 'private') &&
         !props.isMember && 
         !props.isAdmin
})

const canLeaveGroup = computed(() => {
  return props.group && 
         (props.isMember || props.isAdmin)
})

const handleJoin = async () => {
  if (isJoining.value) return
  
  isJoining.value = true
  
  try {
    // Show animation for minimum time
    await new Promise(resolve => setTimeout(resolve, 1200))
    
    // After animation, emit the join event
    emit('join')
  } finally {
    isJoining.value = false
  }
}

const handleLeave = async () => {
  if (isLeaving.value) return
  
  isLeaving.value = true
  
  try {
    // Show animation for minimum time
    await new Promise(resolve => setTimeout(resolve, 800))
    
    // After animation, emit the leave event
    emit('leave')
  } finally {
    isLeaving.value = false
  }
}
</script>

<style scoped>
.group-actions {
  display: flex;
  gap: 12px;
}

.btn-primary,
.btn-secondary {
  padding: 12px 20px;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.95);
  text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.8);
  padding: 10px 18px;
  border-radius: 30px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  cursor: pointer;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.15);
}

.chat-action {
  background: rgba(34, 197, 94, 0.15);
  border-color: rgba(34, 197, 94, 0.3);
  color: rgba(34, 197, 94, 1);
}

.invite-action {
  background: rgba(59, 130, 246, 0.15);
  border-color: rgba(59, 130, 246, 0.3);
  color: rgba(59, 130, 246, 1);
}

.leave-action {
  background: rgba(239, 68, 68, 0.15);
  border-color: rgba(239, 68, 68, 0.3);
  color: rgba(239, 68, 68, 1);
}

.action-badge:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
}

.chat-action:hover {
  background: rgba(34, 197, 94, 0.2);
  border-color: rgba(34, 197, 94, 0.4);
}

.invite-action:hover {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.4);
}

.leave-action:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.4);
}

.action-badge:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.request-pending-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1rem;
  font-weight: 600;
  color: rgba(255, 193, 7, 1);
  text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.8);
  padding: 10px 18px;
  border-radius: 30px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 193, 7, 0.3);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  background: rgba(255, 193, 7, 0.15);
}

.btn-icon {
  width: 48px;
  height: 48px;
  padding: 0;
  border: none;
  border-radius: 50%;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-primary {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-danger {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
  color: white;
}

.btn-primary:hover,
.btn-secondary:hover,
.btn-icon:hover {
  transform: translateY(-1px);
}

.btn-danger:hover {
  transform: translateY(-1px);
  background: linear-gradient(135deg, #ff5252 0%, #d63031 100%);
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.icon {
  font-size: 1.1em;
}
</style>
