<template>
  <div class="events-container">
    <!-- Header with Create Button -->
    <div class="events-header">
      <div class="header-content">
        <Calendar :size="24" />
        <h2>Événements du groupe</h2>
      </div>
      <button 
        v-if="canCreateEvent"
        @click="showCreateModal = true" 
        class="create-event-btn"
      >
        <Plus :size="18" />
        Nouvel événement
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement des événements...</p>
    </div>

    <!-- Events List -->
    <div v-else-if="events && events.length > 0" class="events-list">
      <EventCard
        v-for="event in events"
        :key="event.id"
        :event="event"
        :loading="false"
      />
    </div>

    <!-- Empty State -->
    <div v-else class="empty-state">
      <Calendar :size="48" />
      <h3>Aucun événement</h3>
      <p v-if="canCreateEvent">Créez le premier événement de ce groupe !</p>
      <p v-else>Il n'y a pas encore d'événements dans ce groupe.</p>
      <button 
        v-if="canCreateEvent"
        @click="showCreateModal = true" 
        class="btn-primary"
      >
        <Plus :size="18" />
        Créer un événement
      </button>
    </div>

    <!-- Create Event Modal -->
    <CreateEventModal
      v-if="showCreateModal"
      :group-id="groupId"
      :loading="creatingEvent"
      @close="showCreateModal = false"
      @create="handleCreateEvent"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Calendar, Plus } from 'lucide-vue-next'
import EventCard from './EventCard.vue'
import CreateEventModal from './CreateEventModal.vue'
import { useGroupStore } from '../../stores/groupStore'

const props = defineProps({
  groupId: {
    type: String,
    required: true
  },
  canCreateEvent: {
    type: Boolean,
    default: false
  }
})

const groupStore = useGroupStore()
const showCreateModal = ref(false)
const creatingEvent = ref(false)

// Use groupStore events
const events = computed(() => groupStore.events || [])
const loading = computed(() => groupStore.loading)

const loadEvents = async () => {
  if (!props.groupId) {
    console.warn('No groupId provided for loadEvents')
    return
  }
  await groupStore.fetchGroupEvents(props.groupId)
}

const handleCreateEvent = async (eventData) => {
  creatingEvent.value = true
  try {
    console.log('Creating event with data:', eventData)
    const response = await fetch('http://localhost:8081/newEvent', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify(eventData)
    })

    console.log('Create event response status:', response.status)
    if (response.ok) {
      const responseData = await response.json()
      console.log('Create event response data:', responseData)
      showCreateModal.value = false
      await loadEvents() // Reload events
    } else {
      const error = await response.json()
      console.error('Create event error:', error)
      alert(error.message || 'Erreur lors de la création de l\'événement')
    }
  } catch (error) {
    console.error('Error creating event:', error)
    alert('Erreur lors de la création de l\'événement')
  } finally {
    creatingEvent.value = false
  }
}

onMounted(() => {
  // Don't load events here - let GroupView handle the initial loading
  console.log('EventsList mounted, events already managed by GroupView')
})
</script>

<style scoped>
.events-container {
  max-width: 800px;
  margin: 0 auto;
}

.events-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 30px;
  padding: 0 4px;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-content h2 {
  color: white;
  margin: 0;
  font-size: 1.75rem;
  font-weight: 700;
}

.create-event-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.create-event-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-top: 3px solid #e879c6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-state p {
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
}

.events-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
  background: rgba(15, 15, 23, 0.4);
  border-radius: 20px;
  border: 2px dashed rgba(255, 255, 255, 0.1);
}

.empty-state svg {
  color: rgba(255, 255, 255, 0.4);
  margin-bottom: 20px;
}

.empty-state h3 {
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 12px 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.empty-state p {
  color: rgba(255, 255, 255, 0.6);
  margin: 0 0 24px 0;
  line-height: 1.5;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

/* Responsive */
@media (max-width: 768px) {
  .events-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .header-content {
    justify-content: center;
  }
  
  .create-event-btn {
    justify-content: center;
    width: 100%;
  }
  
  .empty-state {
    padding: 60px 20px;
  }
}
</style>
