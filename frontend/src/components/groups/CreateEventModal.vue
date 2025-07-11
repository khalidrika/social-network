<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h3>
          <Calendar :size="24" />
          Créer un événement
        </h3>
        <button @click="$emit('close')" class="close-btn">
          <X :size="20" />
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="event-form">
        <div class="form-group">
          <label for="title">
            <Type :size="18" />
            Titre de l'événement
          </label>
          <input
            id="title"
            v-model="eventData.title"
            type="text"
            placeholder="Ex: Réunion d'équipe, Fête d'anniversaire..."
            required
            maxlength="100"
          />
        </div>

        <div class="form-group">
          <label for="description">
            <FileText :size="18" />
            Description
          </label>
          <textarea
            id="description"
            v-model="eventData.description"
            placeholder="Décrivez votre événement..."
            rows="4"
            maxlength="500"
          ></textarea>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="date">
              <CalendarDays :size="18" />
              Date
            </label>
            <input
              id="date"
              v-model="eventData.date"
              type="date"
              required
              :min="today"
            />
          </div>

          <div class="form-group">
            <label for="time">
              <Clock :size="18" />
              Heure
            </label>
            <input
              id="time"
              v-model="eventData.time"
              type="time"
              required
            />
          </div>
        </div>

        <div class="rsvp-info">
          <div class="rsvp-header">
            <Users :size="18" />
            <span>Options de réponse</span>
          </div>
          <div class="rsvp-options">
            <div class="rsvp-option">
              <CheckCircle :size="16" />
              <span>Je participe</span>
            </div>
            <div class="rsvp-option">
              <XCircle :size="16" />
              <span>Je ne participe pas</span>
            </div>
            <div class="rsvp-option">
              <HelpCircle :size="16" />
              <span>Peut-être</span>
            </div>
          </div>
        </div>

        <div class="form-actions">
          <button type="button" @click="$emit('close')" class="btn-secondary">
            Annuler
          </button>
          <button type="submit" class="btn-primary" :disabled="loading">
            <component 
              :is="loading ? Loader2 : Calendar" 
              :size="18"
              :class="{ 'animate-spin': loading }"
            />
            {{ loading ? 'Création...' : 'Créer l\'événement' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { 
  Calendar, 
  CalendarDays, 
  Clock, 
  Type, 
  FileText, 
  Users, 
  CheckCircle, 
  XCircle, 
  HelpCircle, 
  X, 
  Loader2 
} from 'lucide-vue-next'

const props = defineProps({
  groupId: {
    type: String,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'create'])

const today = computed(() => {
  return new Date().toISOString().split('T')[0]
})

const eventData = ref({
  title: '',
  description: '',
  date: '',
  time: ''
})

const handleSubmit = () => {
  if (!eventData.value.title.trim() || !eventData.value.date || !eventData.value.time) {
    return
  }

  // Combine date and time
  const datetime = new Date(`${eventData.value.date}T${eventData.value.time}`)
  
  const eventPayload = {
    title: eventData.value.title.trim(),
    content: eventData.value.description.trim(),
    date: datetime.toISOString(),
    groupId: props.groupId
  }

  emit('create', eventPayload)
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
  padding: 20px;
}

.modal-content {
  background: rgba(15, 15, 23, 0.95);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(20px);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h3 {
  display: flex;
  align-items: center;
  gap: 12px;
  color: white;
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
}

.close-btn {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 12px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.event-form {
  padding: 24px;
}

.form-group {
  margin-bottom: 24px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 600;
  margin-bottom: 8px;
  font-size: 0.95rem;
}

.form-group input,
.form-group textarea {
  width: 100%;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 12px 16px;
  color: white;
  font-size: 1rem;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: rgba(232, 121, 198, 0.5);
  box-shadow: 0 0 0 3px rgba(232, 121, 198, 0.1);
  background: rgba(255, 255, 255, 0.15);
}

.form-group input::placeholder,
.form-group textarea::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.rsvp-info {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 24px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.rsvp-header {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 600;
  margin-bottom: 16px;
}

.rsvp-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rsvp-option {
  display: flex;
  align-items: center;
  gap: 12px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.95rem;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 32px;
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
  font-size: 1rem;
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

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
  color: white;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
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
  .modal-overlay {
    padding: 10px;
  }
  
  .modal-content {
    max-width: 100%;
    border-radius: 20px;
  }
  
  .modal-header,
  .event-form {
    padding: 20px;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn-primary,
  .btn-secondary {
    width: 100%;
    justify-content: center;
  }
}
</style>
