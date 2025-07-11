<template>
  <div class="group-view">
    <!-- Loading State -->
    <div v-if="groupStore.loading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement du groupe...</p>
    </div>

    <!-- Group Content -->
    <div v-else-if="groupStore.currentGroup" class="group-content">
      <!-- Header -->
      <GroupHeader 
        :group="groupStore.currentGroup"
        :is-member="groupStore.isCurrentUserMember"
        :is-admin="groupStore.isCurrentUserAdmin"
        @join="handleJoinGroup"
        @leave="showLeaveConfirmation"
        @invite="showInviteModal = true"
      />

      <!-- Navigation Tabs -->
      <GroupTabs 
        :active-tab="activeTab"
        :tabs="tabs"
        @tab-change="activeTab = $event"
      />

      <!-- Tab Content -->
      <div class="tab-content">
        <!-- Posts Tab -->
        <GroupPosts 
          v-if="activeTab === 'posts'"
          :posts="groupStore.currentGroupPosts"
          :group-id="route.params.id"
          :can-post="groupStore.isCurrentUserMember || groupStore.isCurrentUserAdmin"
        />

        <!-- Members Tab -->
        <GroupMembers 
          v-if="activeTab === 'members'"
          :members="groupStore.currentGroupMembers"
        />

        <!-- Events Tab -->
        <GroupEvents 
          v-if="activeTab === 'events'"
          :group-id="route.params.id"
          :can-create-event="groupStore.isCurrentUserMember || groupStore.isCurrentUserAdmin"
        />
      </div>
    </div>

    <!-- Error State -->
    <div v-else class="error-state">
      <div class="error-icon">
        <X :size="48" />
      </div>
      <h3>Groupe introuvable</h3>
      <p>Ce groupe n'existe pas ou vous n'avez pas accès</p>
      <button @click="$router.push('/groups')" class="btn-primary">
        Retour aux groupes
      </button>
    </div>

    <!-- Leave Group Confirmation Modal -->
    <LeaveGroupModal 
      v-if="showLeaveModal"
      :group="groupStore.currentGroup"
      :is-admin="groupStore.isCurrentUserAdmin"
      :error="leaveError"
      @confirm="handleLeaveGroup"
      @cancel="showLeaveModal = false"
    />

    <!-- Invite Users Modal -->
    <InviteUsersModal 
      v-if="showInviteModal"
      :group="groupStore.currentGroup"
      @close="showInviteModal = false"
      @invited="handleUsersInvited"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useGroupStore } from '../stores/groupStore'
import { FileText, Users, Calendar, X } from 'lucide-vue-next'

// Components
import GroupHeader from '../components/groups/GroupHeader.vue'
import GroupTabs from '../components/groups/GroupTabs.vue'
import GroupPosts from '../components/groups/GroupPosts.vue'
import GroupMembers from '../components/groups/GroupMembers.vue'
import GroupEvents from '../components/groups/GroupEvents.vue'
import LeaveGroupModal from '../components/groups/LeaveGroupModal.vue'
import InviteUsersModal from '../components/groups/InviteUsersModal.vue'

const route = useRoute()
const router = useRouter()
const groupStore = useGroupStore()

// State
const activeTab = ref('posts')
const showLeaveModal = ref(false)
const showInviteModal = ref(false)
const leaveError = ref(null)

// Computed
const tabs = computed(() => [
  { id: 'posts', label: 'Posts', icon: FileText },
  { id: 'members', label: 'Membres', icon: Users },
  { id: 'events', label: 'Événements', icon: Calendar }
])

// Methods
const loadGroupData = async () => {
  const groupId = route.params.id
  
  // Load group details first
  await groupStore.fetchGroupDetails(groupId)
  
  if (groupStore.currentGroup) {
    // Load additional data in parallel
    await Promise.all([
      groupStore.fetchGroupPosts(groupId),
      groupStore.fetchGroupMembers(groupId),
      groupStore.fetchGroupEvents(groupId)
    ])
  }
}

const handleJoinGroup = async () => {
  const groupId = route.params.id
  const isPublic = groupStore.currentGroup.privacy === 'public'
  
  // The GroupActions component has already shown its animation and is now emitting this event
  // So we can proceed directly with the API call
  const success = await groupStore.joinGroup(groupId, isPublic)
  
  if (success && isPublic) {
    // Small delay to ensure state updates complete
    await new Promise(resolve => setTimeout(resolve, 200))
    router.push(`/group/${groupId}`)
  }
}

const showLeaveConfirmation = () => {
  leaveError.value = null // Clear any previous errors
  showLeaveModal.value = true
}

const handleLeaveGroup = async () => {
  const groupId = route.params.id
  
  // The LeaveGroupModal component has already shown its animation and is now emitting this event
  // So we can proceed directly with the API call
  const result = await groupStore.leaveGroup(groupId)
  
  if (result.success) {
    showLeaveModal.value = false
    leaveError.value = null
    // Small delay to ensure state updates complete
    await new Promise(resolve => setTimeout(resolve, 200))
  } else {
    // Show error in the modal instead of an alert
    if (result.error === 'admin_cannot_leave') {
      leaveError.value = `En tant qu'administrateur, vous ne pouvez pas quitter votre propre groupe. Vous devez d'abord transférer votre rôle à un autre membre.`
    } else {
      leaveError.value = result.message || 'Une erreur est survenue lors de la tentative de quitter le groupe.'
    }
    
    // Don't close the modal, keep it open to show the error
  }
}

const handleUsersInvited = (count) => {
  console.log(`Successfully invited ${count} users`)
  // Silently handle the invitation success
}

// Lifecycle
onMounted(() => {
  loadGroupData()
})

onUnmounted(() => {
  groupStore.clearCurrentGroup()
})
</script>

<style scoped>
.group-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.loading-state,
.error-state {
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
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.6;
}

.error-state h3 {
  color: rgba(255, 255, 255, 0.8);
  margin: 16px 0 8px 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.error-state p {
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

.tab-content {
  min-height: 400px;
}

/* Responsive */
@media (max-width: 768px) {
  .group-view {
    padding: 15px;
  }
}
</style>

<style scoped>
.group-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.loading-state,
.error-state {
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
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.animate-spin {
  animation: spin 1s linear infinite !important;
}

.error-icon,
.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.6;
}

.group-header {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  /* padding: 30px; */
  margin-bottom: 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.group-info {
  display: flex;
  gap: 20px;
  align-items: center;
}

.group-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.group-avatar img {
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
}

.group-icon {
  font-size: 2.5rem;
  opacity: 0.8;
}

.group-details h1 {
  font-size: 2rem;
  font-weight: 700;
  color: white;
  margin: 0 0 8px 0;
}

.group-details p {
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 12px 0;
}

.group-meta {
  display: flex;
  gap: 16px;
}

.member-count,
.group-privacy {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
}

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
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a52 100%);
  color: white;
  border: 1px solid rgba(255, 107, 107, 0.3);
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

.btn-primary:hover,
.btn-secondary:hover,
.btn-danger:hover,
.btn-icon:hover {
  transform: translateY(-1px);
}

.tabs-container {
  display: flex;
  gap: 4px;
  margin-bottom: 30px;
  background: rgba(15, 15, 23, 0.6);
  padding: 8px;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.tab-btn {
  flex: 1;
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  padding: 12px 20px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-weight: 500;
}

.tab-btn.active {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
  transform: translateY(-1px);
}

.tab-content {
  min-height: 400px;
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

.posts-list,
.events-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.post-card,
.event-card {
  background: rgba(15, 15, 23, 0.6);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px;
}

.post-header {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-author {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #78c7ff 0%, #e879c6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
}

.post-info h4 {
  color: white;
  margin: 0 0 4px 0;
  font-size: 0.875rem;
  font-weight: 600;
}

.post-date {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.75rem;
}

.post-content p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 16px 0;
  line-height: 1.6;
}

.post-image {
  width: 100%;
  max-height: 300px;
  object-fit: cover;
  border-radius: 12px;
}

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

.member-role {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
}

.event-card {
  display: flex;
  gap: 20px;
  align-items: center;
}

.event-date {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  border-radius: 12px;
  padding: 12px;
  text-align: center;
  color: white;
  flex-shrink: 0;
}

.event-date .day {
  display: block;
  font-size: 1.5rem;
  font-weight: 700;
}

.event-date .month {
  display: block;
  font-size: 0.875rem;
  opacity: 0.8;
}

.event-info h4 {
  color: white;
  margin: 0 0 8px 0;
  font-size: 1.125rem;
  font-weight: 600;
}

.event-info p {
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 8px 0;
}

.event-time {
  display: flex;
  align-items: center;
  gap: 6px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
}

.icon {
  font-size: 0.875rem;
}

/* Responsive */
@media (max-width: 768px) {
  .group-view {
    padding: 15px;
  }
  
  .group-header {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }
  
  .group-info {
    flex-direction: column;
    text-align: center;
  }
  
  .tabs-container {
    flex-direction: column;
  }
  
  .members-grid {
    grid-template-columns: 1fr;
  }
  
  .event-card {
    flex-direction: column;
    text-align: center;
  }
}

/* Modal Styles */
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

.modal-actions {
  padding: 20px 24px 24px 24px;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-actions .btn-secondary,
.modal-actions .btn-danger {
  padding: 10px 20px;
  font-size: 0.875rem;
}
</style>
