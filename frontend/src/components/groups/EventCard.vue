<template>
  <div class="event-card">
    <!-- Event Header -->
    <div class="event-header">
      <div class="event-date">
        <div class="day">{{ formatDay(event.date) }}</div>
        <div class="month">{{ formatMonth(event.date) }}</div>
      </div>
      
      <div class="event-info">
        <h3 class="event-title">{{ event.title }}</h3>
        <div class="event-meta">
          <div class="event-time">
            <Clock :size="16" />
            {{ formatTime(event.date) }}
          </div>
          <div class="event-creator">
            <User :size="16" />
            {{ event.author?.firstName }} {{ event.author?.lastName }}
          </div>
        </div>
      </div>
    </div>

    <!-- Event Description -->
    <div v-if="event.content" class="event-description">
      <p>{{ event.content }}</p>
    </div>

    <!-- RSVP Section -->
    <div class="rsvp-section">
      <div class="rsvp-stats">
        <div class="rsvp-stat going">
          <CheckCircle :size="16" />
          <span>{{ event.goingCount || 0 }}</span>
          <span class="label">Participent</span>
        </div>
        <div class="rsvp-stat not-going">
          <XCircle :size="16" />
          <span>{{ event.notGoingCount || 0 }}</span>
          <span class="label">Absents</span>
        </div>
        <div class="rsvp-stat maybe">
          <HelpCircle :size="16" />
          <span>{{ event.maybeCount || 0 }}</span>
          <span class="label">Peut-être</span>
        </div>
      </div>

      <!-- User Response Buttons -->
      <div class="rsvp-actions">
        <button
          @click="updateResponse('going')"
          class="rsvp-btn going"
          :class="{ active: event.userResponse === 'going' }"
          :disabled="loading"
        >
          <CheckCircle :size="18" />
          <span>J'y vais</span>
        </button>
        
        <button
          @click="updateResponse('maybe')"
          class="rsvp-btn maybe"
          :class="{ active: event.userResponse === 'maybe' }"
          :disabled="loading"
        >
          <HelpCircle :size="18" />
          <span>Peut-être</span>
        </button>
        
        <button
          @click="updateResponse('not_going')"
          class="rsvp-btn not-going"
          :class="{ active: event.userResponse === 'not_going' }"
          :disabled="loading"
        >
          <XCircle :size="18" />
          <span>Non</span>
        </button>
      </div>
    </div>

    <!-- Participants List (collapsed/expanded) -->
    <div v-if="event.responses && event.responses.length > 0" class="participants-section">
      <button @click="showParticipants = !showParticipants" class="participants-toggle">
        <Users :size="16" />
        <span>{{ event.responses.length }} participant{{ event.responses.length > 1 ? 's' : '' }}</span>
        <ChevronDown :size="16" :class="{ 'rotated': showParticipants }" />
      </button>
      
      <div v-if="showParticipants" class="participants-list">
        <div 
          v-for="response in event.responses" 
          :key="response.userId"
          class="participant"
          :class="response.response"
        >
          <div class="participant-avatar">
  <img v-if="response.avatar || response.user?.avatar || response.profilePic || response.user?.profilePic" 
       :src="getFullImageUrl(response.avatar || response.user?.avatar || response.profilePic || response.user?.profilePic)" 
       :alt="response.userName" />
  <span v-else>{{ getInitials(response.userName) }}</span>
</div>
          <span class="participant-name">{{ response.userName }}</span>
          <div class="participant-status">
            <component 
              :is="getResponseIcon(response.response)" 
              :size="14" 
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { 
  Calendar, 
  Clock, 
  User, 
  Users, 
  CheckCircle, 
  XCircle, 
  HelpCircle,
  ChevronDown
} from 'lucide-vue-next'
import { useGroupStore } from '../../stores/groupStore'
import { useUserStore } from '../../stores/userStore'
import { useMainStore } from '../../stores/postsStore'

const props = defineProps({
  event: {
    type: Object,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update-response'])

const groupStore = useGroupStore()
const userStore = useUserStore()
const mainStore = useMainStore()
const { getFullImageUrl } = mainStore

const showParticipants = ref(false)

const formatDay = (dateString) => {
  const date = new Date(dateString)
  return date.getDate().toString().padStart(2, '0')
}

const formatMonth = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('fr-FR', { month: 'short' }).toUpperCase()
}

const formatTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleTimeString('fr-FR', { 
    hour: '2-digit', 
    minute: '2-digit'
  })
}

const getInitials = (name) => {
  return name.split(' ').map(n => n[0]).join('').toUpperCase()
}

const getResponseIcon = (response) => {
  switch (response) {
    case 'going': return CheckCircle
    case 'not_going': return XCircle
    case 'maybe': return HelpCircle
    default: return HelpCircle
  }
}

const updateResponse = async (response) => {
  if (props.loading) return
  const success = await groupStore.updateEventResponse(props.event.id, response)
  if (!success) {
    console.error('Failed to update event response')
  }
}
</script>

<style scoped>
.event-card {
  background: rgba(15, 15, 23, 0.8);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  overflow: hidden;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.event-card:hover {
  border-color: rgba(232, 121, 198, 0.3);
  box-shadow: 0 8px 32px rgba(232, 121, 198, 0.15);
}

.event-header {
  display: flex;
  gap: 20px;
  padding: 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  position: relative;
}

.event-date {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  border-radius: 16px;
  padding: 16px;
  text-align: center;
  color: white;
  flex-shrink: 0;
  min-width: 80px;
}

.event-date .day {
  display: block;
  font-size: 2rem;
  font-weight: 800;
  line-height: 1;
}

.event-date .month {
  display: block;
  font-size: 0.875rem;
  opacity: 0.9;
  margin-top: 4px;
  font-weight: 600;
}

.event-info {
  flex: 1;
}

.event-title {
  color: white;
  margin: 0 0 12px 0;
  font-size: 1.5rem;
  font-weight: 700;
  line-height: 1.3;
}

.event-meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.event-time,
.event-creator {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.95rem;
}

.event-description {
  padding: 0 24px 20px 24px;
}

.event-description p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
  line-height: 1.6;
}

.rsvp-section {
  padding: 20px 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.rsvp-stats {
  display: flex;
  gap: 24px;
  margin-bottom: 20px;
  justify-content: center;
}

.rsvp-stat {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.9rem;
  font-weight: 600;
}

.rsvp-stat.going {
  color: rgba(34, 197, 94, 1);
}

.rsvp-stat.not-going {
  color: rgba(239, 68, 68, 1);
}

.rsvp-stat.maybe {
  color: rgba(251, 191, 36, 1);
}

.rsvp-stat .label {
  font-size: 0.8rem;
  opacity: 0.8;
}

.rsvp-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.rsvp-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 25px;
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  font-size: 0.9rem;
}

.rsvp-btn:hover {
  transform: translateY(-1px);
  background: rgba(255, 255, 255, 0.15);
}

.rsvp-btn.going.active {
  background: rgba(34, 197, 94, 0.2);
  border-color: rgba(34, 197, 94, 0.4);
  color: rgba(34, 197, 94, 1);
}

.rsvp-btn.not-going.active {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.4);
  color: rgba(239, 68, 68, 1);
}

.rsvp-btn.maybe.active {
  background: rgba(251, 191, 36, 0.2);
  border-color: rgba(251, 191, 36, 0.4);
  color: rgba(251, 191, 36, 1);
}

.rsvp-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.participants-section {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px 24px;
}

.participants-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  padding: 8px 0;
  font-weight: 600;
  transition: all 0.3s ease;
}

.participants-toggle:hover {
  color: white;
}

.participants-toggle .rotated {
  transform: rotate(180deg);
}

.participants-list {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.participant {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
}

.participant-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #78c7ff 0%, #e879c6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 0.9rem;
  overflow: hidden;
}

.participant-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.participant-name {
  flex: 1;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
}

.participant-status {
  opacity: 0.8;
}

.participant.going .participant-status {
  color: rgba(34, 197, 94, 1);
}

.participant.not-going .participant-status {
  color: rgba(239, 68, 68, 1);
}

.participant.maybe .participant-status {
  color: rgba(251, 191, 36, 1);
}

/* Responsive */
@media (max-width: 768px) {
  .event-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
    padding-top: 50px; /* Make room for delete button */
  }
  
  .event-actions {
    top: 12px;
    right: 12px;
  }
  
  .event-date {
    align-self: center;
  }
  
  .event-meta {
    align-items: center;
  }
  
  .rsvp-stats {
    flex-direction: column;
    gap: 12px;
    align-items: center;
  }
  
  .rsvp-actions {
    flex-direction: column;
  }
  
  .rsvp-btn {
    justify-content: center;
    width: 100%;
  }
}
</style>
