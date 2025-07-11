<template>
  <div class="members-section">
    <div v-if="members.length === 0" class="empty-state">
      <div class="empty-icon">
        <Users :size="64" />
      </div>
      <h3>Aucun membre</h3>
      <p>Ce groupe n'a pas encore de membres</p>
    </div>
    <div v-else class="members-grid">
      <div v-for="member in members" :key="member.id" class="member-card">
        <div class="member-avatar">
          <img v-if="member.avatar || member.profilePic" :src="getFullImageUrl(member.avatar || member.profilePic)" :alt="member.name" />
          <div v-else class="default-member">
            <Users :size="24" />
          </div>
        </div>
        <div class="member-info">
          <h4>{{ member.firstName }} {{ member.lastName }} 
            <Crown v-if="member.following" class="admin-crown" :size="16" />
          </h4>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Users, Crown } from 'lucide-vue-next'
import { useMainStore } from '../../stores/postsStore'

const mainStore = useMainStore()
const { getFullImageUrl } = mainStore

defineProps({
  members: {
    type: Array,
    default: () => []
  }
})
</script>

<style scoped>
.members-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
}

.member-card {
  background: rgba(15, 15, 23, 0.6);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 16px;
  display: flex;
  gap: 12px;
  align-items: center;
}

.member-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.member-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-member {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #78c7ff 0%, #e879c6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
}

.member-info h4 {
  color: white;
  margin: 0 0 4px 0;
  font-size: 1rem;
  font-weight: 600;
}

.admin-crown {
  color: #ffd700;
  font-size: 0.9rem;
  margin-left: 4px;
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
  margin: 0;
}

/* Responsive */
@media (max-width: 768px) {
  .members-grid {
    grid-template-columns: 1fr;
  }
}
</style>
