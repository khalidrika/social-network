<template>
  <div class="groups-grid">
    <div v-if="groups.length === 0" class="empty-state">
      <span class="empty-icon">{{ emptyIcon }}</span>
      <h3>{{ emptyTitle }}</h3>
      <p>{{ emptyMessage }}</p>
      <button 
        v-if="showCreateButton" 
        @click="$emit('create-group')" 
        class="btn-primary"
      >
        Créer votre premier groupe
      </button>
    </div>
    <GroupCard 
      v-for="group in groups" 
      :key="group.id"
      :group="group"
      :is-member="isMember(group)"
      :loading="group.loading || false"
      @view="$emit('view-group', group)"
      @chat="$emit('chat-group', group.id)"
      @join-request="$emit('join-group', group.id)"
    />
  </div>
</template>

<script setup>
import GroupCard from './GroupCard.vue'

const props = defineProps({
  groups: {
    type: Array,
    default: () => []
  },
  emptyIcon: {
    type: String,
    default: '👥'
  },
  emptyTitle: {
    type: String,
    default: 'Aucun groupe'
  },
  emptyMessage: {
    type: String,
    default: 'Aucun groupe disponible'
  },
  showCreateButton: {
    type: Boolean,
    default: false
  },
  type: {
    type: String,
    default: 'my' // 'my' or 'discover'
  }
})

defineEmits(['view-group', 'chat-group', 'join-group', 'create-group'])

const isMember = (group) => {
  return props.type === 'my' ? true : (group.member || false)
}
</script>

<style scoped>
.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  background: rgba(15, 15, 23, 0.4);
  border-radius: 20px;
  border: 2px dashed rgba(255, 255, 255, 0.1);
  grid-column: 1 / -1;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.6;
}

.empty-state h3 {
  color: rgba(255, 255, 255, 0.8);
  margin: 16px 0 8px 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.empty-state p {
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 24px 0;
}

.btn-primary {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
  border: none;
  border-radius: 12px;
  padding: 10px 20px;
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
  .groups-grid {
    grid-template-columns: 1fr;
  }
}
</style>
