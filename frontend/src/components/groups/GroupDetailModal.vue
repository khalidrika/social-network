<template>
  <div class="modal-overlay" @click="handleOverlayClick">
    <div class="modal-content" @click.stop>
      <!-- Header -->
      <div class="modal-header">
        <div class="group-header-info">
          <div class="group-avatar">
            <img 
              v-if="group.image" 
              :src="group.image" 
              :alt="group.name"
              @error="handleImageError"
            />
            <div v-else class="default-avatar">
              <Users :size="32" />
            </div>
          </div>
          <div>
            <h2>{{ group.name }}</h2>
            <div class="group-meta">
              <span class="member-count">
                <Users :size="16" />
                {{ group.memberCount || 0 }} membres
              </span>
              <span class="group-privacy" :class="group.privacy">
                <component :is="group.privacy === 'public' ? Globe : Lock" :size="16" />
                {{ privacyLabel }}
              </span>
            </div>
          </div>
        </div>
        <button @click="$emit('close')" class="close-btn">
          <X :size="20" />
        </button>
      </div>

      <!-- Content -->
      <div class="modal-body">
        <!-- Description -->
        <div v-if="group.description" class="section">
          <h3>Description</h3>
          <p class="description">{{ group.description }}</p>
        </div>

        <!-- Actions -->
        <div class="section">
          <h3>Actions</h3>
          <div class="action-buttons">
            <button 
              @click="$emit('chat', group.id)"
              class="btn-primary"
            >
              <MessageCircle :size="18" />
              Chat du groupe
            </button>
            
            <button 
              @click="inviteUsers"
              class="btn-primary"
            >
              <UserPlus :size="18" />
              Inviter des utilisateurs
            </button>
            
            <button 
              @click="viewMembers"
              class="btn-secondary"
            >
              <Users :size="18" />
              Voir les membres
            </button>
            
            <button 
              v-if="isAdmin"
              @click="manageGroup"
              class="btn-secondary"
            >
              <Settings :size="18" />
              Gérer le groupe
            </button>
          </div>
        </div>

        <!-- Recent Activity -->
        <div class="section">
          <h3>Activité récente</h3>
          <div v-if="recentActivity.length === 0" class="empty-activity">
            <div class="empty-icon">
              <Clock :size="48" />
            </div>
            <p>Aucune activité récente</p>
          </div>
          <div v-else class="activity-list">
            <div 
              v-for="activity in recentActivity" 
              :key="activity.id"
              class="activity-item"
            >
              <div class="activity-icon">
                <component :is="getActivityIcon(activity.type)" :size="16" />
              </div>
              <div class="activity-content">
                <p>{{ activity.description }}</p>
                <span class="activity-time">{{ formatDate(activity.createdAt) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Members Preview -->
        <div v-if="members.length > 0" class="section">
          <h3>Membres</h3>
          <div class="members-preview">
            <div 
              v-for="member in members.slice(0, 6)" 
              :key="member.id"
              class="member-avatar"
              :title="`${member.firstName} ${member.lastName}`"
            >
              <img 
                v-if="member.profilePic" 
                :src="getFullImageUrl(member.profilePic)" 
                :alt="`${member.firstName} ${member.lastName}`"
                @error="handleMemberImageError"
              />
              <div v-else class="default-member-avatar">
                <Users :size="16" />
              </div>
            </div>
            <div v-if="group.memberCount > 6" class="more-members">
              +{{ group.memberCount - 6 }}
            </div>
          </div>
        </div>
      </div>

      <!-- Invite Users Modal -->
      <InviteUsersModal 
        v-if="showInviteModal"
        :group="group"
        @close="showInviteModal = false"
        @invited="handleUsersInvited"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useMainStore } from '../../stores/postsStore'
import { Users, Globe, Lock, X, MessageCircle, Clock, Settings, FileText, Camera, UserPlus } from 'lucide-vue-next'
import InviteUsersModal from './InviteUsersModal.vue'

const mainStore = useMainStore()
const { getFullImageUrl } = mainStore

const props = defineProps({
  group: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'updated', 'chat', 'invite'])

const members = ref([])
const recentActivity = ref([])
const isAdmin = ref(false)
const showInviteModal = ref(false)

const privacyIcon = computed(() => {
  return props.group.privacy === 'public' ? '🌐' : '🔒'
})

const privacyLabel = computed(() => {
  return props.group.privacy === 'public' ? 'Public' : 'Privé'
})

const handleOverlayClick = (event) => {
  if (event.target === event.currentTarget) {
    emit('close')
  }
}

const handleImageError = (event) => {
  event.target.style.display = 'none'
  event.target.nextElementSibling.style.display = 'flex'
}

const handleMemberImageError = (event) => {
  event.target.style.display = 'none'
  event.target.nextElementSibling.style.display = 'flex'
}

const viewMembers = () => {
  // TODO: Implement members view
  console.log('View members')
}

const manageGroup = () => {
  // TODO: Implement group management
  console.log('Manage group')
}

const inviteUsers = () => {
  showInviteModal.value = true
}

const handleUsersInvited = (count) => {
  console.log(`Successfully invited ${count} users`)
  // Could show a success message or update the group data
}

const getActivityIcon = (type) => {
  const iconMap = {
    'post': FileText,
    'member_joined': UserPlus,
    'photo': Camera,
    'default': Clock
  }
  return iconMap[type] || iconMap.default
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

const loadGroupDetails = async () => {
  try {
    // Load group members
    const membersRes = await fetch(`http://localhost:8081/groupMembers?groupId=${props.group.id}`, {
      credentials: 'include'
    })
    if (membersRes.ok) {
      const data = await membersRes.json()
      members.value = data.users || []
      // Check if current user is admin
      isAdmin.value = props.group.administrator || false
    }

    // Load recent activity - using group posts as activity
    const activityRes = await fetch(`http://localhost:8081/groupPosts?groupId=${props.group.id}`, {
      credentials: 'include'
    })
    if (activityRes.ok) {
      const data = await activityRes.json()
      // Transform posts into activity items
      recentActivity.value = (data.posts || []).slice(0, 5).map(post => ({
        id: post.id,
        type: post.image ? 'photo' : 'post',
        description: `${post.author?.firstName || 'Quelqu\'un'} a publié: "${post.content.substring(0, 50)}${post.content.length > 50 ? '...' : ''}"`,
        createdAt: post.createdAt
      }))
    }
  } catch (error) {
    console.error('Error loading group details:', error)
  }
}

onMounted(() => {
  loadGroupDetails()
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
  backdrop-filter: blur(8px);
}

.modal-content {
  background: rgba(15, 15, 23, 0.95);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  backdrop-filter: blur(20px);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 24px 24px 0 24px;
  margin-bottom: 24px;
}

.group-header-info {
  display: flex;
  gap: 16px;
  flex: 1;
}

.group-avatar {
  width: 80px;
  height: 80px;
  border-radius: 16px;
  overflow: hidden;
  flex-shrink: 0;
}

.group-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-avatar {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.avatar-icon {
  font-size: 2rem;
}

.group-header-info h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: white;
  margin: 0 0 8px 0;
  line-height: 1.3;
}

.group-meta {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.member-count,
.group-privacy {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
}

.group-privacy.public {
  color: rgba(120, 199, 255, 0.8);
}

.group-privacy.private {
  color: rgba(232, 121, 198, 0.8);
}

.close-btn {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 8px;
  padding: 8px;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.7);
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.modal-body {
  padding: 0 24px 24px 24px;
}

.section {
  margin-bottom: 32px;
}

.section h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: white;
  margin: 0 0 16px 0;
}

.description {
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
  margin: 0;
}

.action-buttons {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.btn-primary,
.btn-secondary {
  padding: 10px 16px;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.875rem;
}

.btn-primary {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.empty-activity {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
  color: rgba(255, 255, 255, 0.6);
}

.empty-activity p {
  margin: 8px 0 0 0;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.activity-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.activity-icon {
  width: 32px;
  height: 32px;
  background: rgba(232, 121, 198, 0.2);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #e879c6;
  flex-shrink: 0;
}

.activity-content {
  flex: 1;
}

.activity-content p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0 0 4px 0;
  line-height: 1.4;
}

.activity-time {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.5);
}

.members-preview {
  display: flex;
  gap: 8px;
  align-items: center;
}

.member-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.member-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-member-avatar {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #78c7ff 0%, #e879c6 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.more-members {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.7);
}

.icon {
  font-size: 0.875rem;
}

.empty-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
  opacity: 0.6;
}

.member-icon {
  font-size: 0.75rem;
}

/* Responsive */
@media (max-width: 768px) {
  .modal-content {
    margin: 0;
    border-radius: 20px 20px 0 0;
    max-height: 100vh;
  }
  
  .group-header-info {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .btn-primary,
  .btn-secondary {
    justify-content: center;
  }
}
</style>
