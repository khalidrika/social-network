<template>
  <div class="modal-overlay" @click="$emit('cancel')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>Quitter le groupe</h3>
      </div>
      <div class="modal-body">
        <!-- Error Message -->
        <div v-if="errorMessage" class="error-message">
          <X :size="20" />
          {{ errorMessage }}
        </div>
        
        <!-- Success Message -->
        <div v-else-if="successMessage" class="success-message">
          <Check :size="20" />
          {{ successMessage }}
        </div>
        
        <!-- Normal Content -->
        <template v-else>
          <p>Êtes-vous sûr de vouloir quitter le groupe <strong>{{ group?.name }}</strong> ?</p>
          <div v-if="props.isAdmin" class="admin-warning-text">
            <AlertTriangle :size="20" />
            <strong>Attention :</strong> En tant qu'administrateur, vous ne pouvez pas quitter votre propre groupe.
            <br>
            Vous devez d'abord transférer votre rôle d'administrateur à un autre membre.
          </div>
          <p v-else class="warning-text">Cette action ne peut pas être annulée.</p>
        </template>
      </div>
      <div class="modal-actions">
        <button @click="handleCancel" class="btn-secondary">
          {{ errorMessage || successMessage ? 'Fermer' : 'Annuler' }}
        </button>
        <button 
          v-if="!props.isAdmin && !errorMessage && !successMessage"
          @click="handleConfirm" 
          class="btn-danger" 
          :disabled="isLeaving"
        >
          <component 
            :is="isLeaving ? Clock : LogOut"
            :size="18" 
            :class="{ 'animate-spin': isLeaving }"
          />
          {{ isLeaving ? 'Quitter...' : 'Quitter' }}
        </button>
        <button 
          v-else-if="props.isAdmin && !errorMessage && !successMessage"
          @click="handleCancel" 
          class="btn-disabled" 
          disabled
        >
          <span class="icon">🚫</span>
          Action impossible
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { X, Check, LogOut, Clock, AlertTriangle } from 'lucide-vue-next'

const props = defineProps({
  group: {
    type: Object,
    required: true
  },
  isAdmin: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: null
  }
})

const emit = defineEmits(['confirm', 'cancel'])

const isLeaving = ref(false)
const errorMessage = ref(props.error)
const successMessage = ref(null)

// Watch for error prop changes
watch(() => props.error, (newError) => {
  errorMessage.value = newError
  if (newError) {
    successMessage.value = null
  }
})

const handleConfirm = async () => {
  if (isLeaving.value) return
  
  isLeaving.value = true
  errorMessage.value = null
  
  try {
    // Show animation for minimum time
    await new Promise(resolve => setTimeout(resolve, 1200))
    
    // After animation, emit confirm event
    emit('confirm')
  } catch (error) {
    console.error('LeaveGroupModal: Error during confirm:', error)
    errorMessage.value = 'Une erreur est survenue lors de la tentative de quitter le groupe.'
  } finally {
    isLeaving.value = false
  }
}

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  background: rgba(15, 15, 23, 0.95);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 0;
  max-width: 400px;
  width: 90%;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

.modal-header {
  padding: 24px 24px 0 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  margin-bottom: 20px;
}

.modal-header h3 {
  color: white;
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.modal-body {
  padding: 0 24px 20px 24px;
}

.modal-body p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 16px 0;
  line-height: 1.5;
}

.modal-body .warning-text {
  color: rgba(255, 107, 107, 0.8);
  font-size: 0.875rem;
  margin: 0;
}

.modal-body .admin-warning-text {
  color: rgba(255, 193, 7, 0.95);
  font-size: 0.875rem;
  margin: 16px 0;
  padding: 16px;
  background: rgba(255, 193, 7, 0.15);
  border-radius: 12px;
  border: 1px solid rgba(255, 193, 7, 0.3);
  line-height: 1.4;
}

.error-message {
  background: rgba(255, 107, 107, 0.15);
  color: rgba(255, 107, 107, 0.9);
  padding: 16px;
  border-radius: 12px;
  border: 1px solid rgba(255, 107, 107, 0.3);
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.875rem;
  line-height: 1.4;
}

.success-message {
  background: rgba(46, 204, 113, 0.15);
  color: rgba(46, 204, 113, 0.9);
  padding: 16px;
  border-radius: 12px;
  border: 1px solid rgba(46, 204, 113, 0.3);
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.875rem;
  line-height: 1.4;
}

.modal-actions {
  padding: 20px 24px 24px 24px;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.btn-secondary,
.btn-danger,
.btn-disabled {
  padding: 10px 20px;
  font-size: 0.875rem;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
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

.btn-disabled {
  padding: 10px 20px;
  font-size: 0.875rem;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: not-allowed;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.btn-secondary:hover,
.btn-danger:hover {
  transform: translateY(-1px);
}

.btn-disabled:hover {
  transform: none;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.icon {
  font-size: 0.875rem;
}
</style>
