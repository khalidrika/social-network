Logo
Issues
Pull Requests
Milestones
Explore
Intra
abelbach
/
social-network
Private
Code
Issues
Pull Requests
Actions
Packages
Projects
Releases
Wiki

social-network
/
frontend
/
src
/
pages
/
Notifications.vue

rachid ichouane
richouan
7733af02c6
feat : dismiss notifs in notif page
last week
637 lines
17 KiB
Vue

This file contains ambiguous Unicode characters
This file contains Unicode characters that might be confused with other characters. If you think that this is intentional, you can safely ignore this warning. Use the Escape button to reveal them.

<template>
  <div class="notifications-page">
    <div class="notifications-header">
      <h1>Notifications</h1>
      <div class="header-actions">
        <button 
          v-if="notificationStore.unreadCount > 0" 
          @click="notificationStore.markAllAsRead()" 
          class="mark-all-read-btn"
        >
          Tout marquer comme lu ({{ notificationStore.unreadCount }})
        </button>
      </div>
    </div>

    <!-- Error Display -->
    <div v-if="notificationStore.error" class="error-banner">
      <div class="error-content">
        <span>{{ notificationStore.error }}</span>
	        <button @click="notificationStore.dismissError()" class="dismiss-btn">×</button>
      </div>
    </div>

    <div v-if="notificationStore.loading" class="loading">
      <div class="loading-spinner"></div>
      <span>Chargement des notifications...</span>
    </div>
    <div v-else-if="notificationStore.notifications.length === 0" class="empty-state">
      Vous n'avez aucune notification.
    </div>
    <div v-else>
      <!-- Follow Requests Section -->
      <div class="notification-section" v-if="notificationStore.followRequests.length > 0">
        <h2>Demandes de suivi</h2>
        <div class="notification-list">
          <div v-for="notif in notificationStore.followRequests" :key="notif.id" class="notification-card follow-request">
            <div class="notification-header">
              <div class="notification-avatar">
                <img v-if="notif.user?.avatar" :src="notif.user.avatar" alt="Avatar" />
                <div v-else class="default-avatar">{{ notif.user?.firstName?.[0] || 'U' }}</div>
              </div>
              <div class="notification-info">
                <div class="notification-title">
                  <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong> souhaite vous suivre
                </div>
                <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
              </div>
            </div>
            <div class="notification-actions">
              <button 
                @click="notificationStore.respondToFollowRequest(notif.id, 'ACCEPT')" 
                class="accept-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-sm"></span>
                {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Accepter' }}
              </button>
              <button 
                @click="notificationStore.respondToFollowRequest(notif.id, 'DECLINE')" 
                class="decline-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Refuser' }}
              </button>
              <button 
                @click="dismissWithConfirmation(notif.id, 'follow')" 
                class="dismiss-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                Supprimer
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Group Invites Section -->
      <!-- Note: Group invites cannot be dismissed as they are required for the invite system to function -->
      <div class="notification-section" v-if="notificationStore.groupInvites.length > 0">
        <h2>Invitations de groupe</h2>
        <p class="section-info">Utilisez "Rejoindre" ou "Refuser" pour répondre aux invitations.</p>
        <div class="notification-list">
          <div v-for="notif in notificationStore.groupInvites" :key="notif.id" class="notification-card group-invite">
            <div class="notification-header">
              <div class="notification-avatar group">
                <UserGroupIcon class="icon-sm" />
              </div>
              <div class="notification-info">
                <div class="notification-title">
                  <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong> vous invite à rejoindre le groupe <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                </div>
                <div class="notification-subtitle" v-if="notif.group?.description">
                  {{ notif.group.description }}
                </div>
                <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
              </div>
            </div>
            <div class="notification-actions">
              <button 
                @click="notificationStore.respondToGroupInvite(notif.id, 'ACCEPT')" 
                class="accept-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-sm"></span>
                {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Rejoindre' }}
              </button>
              <button 
                @click="notificationStore.respondToGroupInvite(notif.id, 'DECLINE')" 
                class="decline-btn"
                :disabled="notificationStore.isProcessing(notif.id)"
              >
                {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Refuser' }}
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Group Requests Section -->
      <!-- Note: Group requests cannot be dismissed as they are required for the request system to function -->
      <div class="notification-section" v-if="notificationStore.groupRequests.length > 0">
        <h2>Demandes d'adhésion aux groupes</h2>
        <p class="section-info">Utilisez "Accepter" ou "Refuser" pour répondre aux demandes.</p>
        <div class="notification-list">
          <div v-for="notif in notificationStore.groupRequests" :key="notif.id" class="notification-card group-request">
            <div class="notification-header">
              <div class="notification-avatar group">
                <UserGroupIcon class="icon-sm" />
              </div>
              <div class="notification-info">
                <div class="notification-title">
                  <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong> demande à rejoindre votre groupe <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                </div>
                <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
              </div>
            </div>              <div class="notification-actions">
                <button 
                  @click="notificationStore.respondToGroupRequest(notif.id, 'ACCEPT', notif.group?.id)" 
                  class="accept-btn"
                  :disabled="notificationStore.isProcessing(notif.id)"
                >
                  <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-sm"></span>
                  {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Accepter' }}
                </button>
                <button 
                  @click="notificationStore.respondToGroupRequest(notif.id, 'DECLINE', notif.group?.id)" 
                  class="decline-btn"
                  :disabled="notificationStore.isProcessing(notif.id)"
                >
                  {{ notificationStore.isProcessing(notif.id) ? 'Traitement...' : 'Refuser' }}
                </button>
              </div>
          </div>
        </div>
      </div>
      
      <!-- Event Notifications Section -->
      <div class="notification-section" v-if="notificationStore.eventNotifications.length > 0">
        <h2>Événements de groupe</h2>
        <div class="notification-list">
          <div 
            v-for="notif in notificationStore.eventNotifications" 
            :key="notif.id" 
            class="notification-card event-notification"
            :class="{ 'unread': !notif.read }"
            @click="notificationStore.markAsReadOnView(notif.id)"
          >
            <div class="notification-header">
              <div class="notification-avatar event">
                <CalendarIcon class="icon-sm" />
              </div>
              <div class="notification-info">
                <div class="notification-title">
                  Nouvel événement <strong>{{ notif.event?.title || 'Événement' }}</strong> dans le groupe <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                  <span v-if="!notif.read" class="unread-dot"></span>
                </div>
                <div class="notification-subtitle" v-if="notif.event?.description">
                  {{ notif.event.description }}
                </div>
                <div class="notification-date" v-if="notif.event?.date">
                  <CalendarDaysIcon class="icon-xs" /> {{ new Date(notif.event.date).toLocaleDateString() }}
                </div>
                <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
              </div>
            </div>
            <div class="notification-actions">
              <button 
                @click.stop="dismissWithConfirmation(notif.id, 'event')" 
                class="dismiss-btn"
              >
                Supprimer
              </button>
              <button 
                @click.stop="notificationStore.markAsRead(notif.id)" 
                class="mark-read-btn" 
                v-if="!notif.read"
              >
                Marquer comme lu
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '../stores/notificationStore'
import { UserGroupIcon, CalendarIcon, CalendarDaysIcon } from '@heroicons/vue/24/outline'
import wsService from '../services/websocketService'

const router = useRouter()
const notificationStore = useNotificationStore()
let removeWebSocketListener = null

function dismissWithConfirmation(notificationId, type) {
  let message = ''
  
  switch(type) {
    case 'follow':
      message = 'Supprimer cette demande de suivi ? La personne pourra soumettre une nouvelle demande.'
      break
    case 'event':
      message = 'Supprimer cette notification d\'événement ?'
      break
    default:
      message = 'Supprimer cette notification ?'
  }
  
  if (confirm(message)) {
    notificationStore.dismissNotification(notificationId)
  }
}

onMounted(() => {
  notificationStore.fetchNotifications()
  
  // Setup WebSocket listener for notifications
  removeWebSocketListener = wsService.addListener('notificationsPage', (data) => {
    // If we receive a notification, refresh the list
    if (data.action === 'notification') {
      notificationStore.fetchNotifications()
    }
  })
})

onUnmounted(() => {
  // Clean up WebSocket listener
  if (removeWebSocketListener) {
    removeWebSocketListener()
  }
})
</script>

<style scoped>
.notifications-page {
  max-width: 800px;
  margin: 2rem auto;
  padding: 2rem;
}

.notifications-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.notifications-header h1 {
  color: #e879c6;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 1rem;
}

.mark-all-read-btn {
  background: #23233a;
  color: #78dbff;
  border: 1px solid rgba(120, 219, 255, 0.3);
  border-radius: 8px;
  padding: 0.5rem 1rem;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.mark-all-read-btn:hover {
  background: #282846;
  border-color: rgba(120, 219, 255, 0.5);
}

.error-banner {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.error-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #ef4444;
}

.error-content .dismiss-btn {
  background: none;
  border: none;
  color: #ef4444;
  font-size: 20px;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dismiss-btn {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3) !important;
  border-radius: 5px;
  padding: 0.5rem 1rem;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dismiss-btn:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.5) !important;
  color: #ff6b6b;
}

.dismiss-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading {
  text-align: center;
  padding: 2rem;
  color: #888;
  font-style: italic;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}

.loading-spinner {
  border: 3px solid rgba(120, 219, 255, 0.1);
  border-radius: 50%;
  border-top: 3px solid #78dbff;
  width: 24px;
  height: 24px;
  animation: spin 1s linear infinite;
}

.loading-spinner-sm {
  border: 2px solid rgba(120, 219, 255, 0.1);
  border-radius: 50%;
  border-top: 2px solid #78dbff;
  width: 16px;
  height: 16px;
  animation: spin 1s linear infinite;
  display: inline-block;
  margin-right: 0.5rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #888;
  font-style: italic;
}

.notification-section {
  margin-bottom: 2rem;
}

.notification-section h2 {
  color: #78dbff;
  margin-bottom: 1rem;
  font-size: 1.3rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid rgba(120, 119, 198, 0.1);
}

.section-info {
  color: #888;
  font-size: 14px;
  font-style: italic;
  margin-bottom: 1rem;
  margin-top: -0.5rem;
}

.notification-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.notification-card {
  background: #181824;
  border-radius: 10px;
  padding: 1.2rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.notification-header {
  display: flex;
  gap: 1rem;
}

.notification-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  overflow: hidden;
  background: #23233a;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.notification-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.notification-avatar .default-avatar {
  color: white;
  font-size: 20px;
}

.icon-sm {
  width: 24px;
  height: 24px;
  color: white;
}

.icon-xs {
  width: 16px;
  height: 16px;
  display: inline-block;
  vertical-align: middle;
  margin-right: 4px;
}

.notification-info {
  flex: 1;
}

.notification-title {
  font-size: 16px;
  color: #d1d5db;
  margin-bottom: 0.5rem;
}

.notification-title strong {
  color: white;
}

.notification-subtitle {
  font-size: 14px;
  color: #888;
  margin-bottom: 0.5rem;
}

.notification-time {
  font-size: 12px;
  color: #78dbff;
  margin-top: 0.5rem;
}

.notification-date {
  font-size: 14px;
  color: #78dbff;
  margin-bottom: 0.5rem;
}

.notification-actions {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}

.notification-actions button {
  padding: 0.5rem 1rem;
  border-radius: 5px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.notification-actions button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.accept-btn {
  background: linear-gradient(90deg, #78dbff 0%, #7879c6 100%);
  color: white;
}

.accept-btn:hover:not(:disabled) {
  opacity: 0.9;
}

.decline-btn {
  background: #23233a;
  color: #d1d5db;
  border: 1px solid rgba(120, 119, 198, 0.3) !important;
}

.decline-btn:hover:not(:disabled) {
  background: #282846;
}

.view-btn {
  background: #23233a;
  color: #78dbff;
  border: 1px solid rgba(120, 119, 198, 0.3) !important;
}

.view-btn:hover:not(:disabled) {
  background: #282846;
}

.mark-read-btn {
  background: transparent;
  color: #888;
  border: 1px solid rgba(136, 136, 136, 0.3) !important;
  font-size: 12px;
  padding: 0.25rem 0.5rem;
}

.mark-read-btn:hover:not(:disabled) {
  color: #78dbff;
  border-color: rgba(120, 219, 255, 0.3) !important;
}

.dismiss-btn {
  background: transparent;
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3) !important;
  font-size: 12px;
  padding: 0.25rem 0.5rem;
}

.dismiss-btn:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.5) !important;
}

/* Different notification types */
.follow-request {
  border-left: 4px solid #e879c6;
}

.group-invite {
  border-left: 4px solid #78dbff;
}

.group-request {
  border-left: 4px solid #7879c6;
}

.event-notification {
  border-left: 4px solid #84cc16;
  position: relative;
  cursor: pointer;
  transition: all 0.2s ease;
}

.event-notification:hover {
  background: rgba(132, 204, 22, 0.05);
}

.event-notification.unread {
  background: rgba(132, 204, 22, 0.08);
  border-left-color: #a3e635;
}

.notification-avatar.group {
  background: linear-gradient(135deg, #7879c6, #78dbff);
}

.notification-avatar.event {
  background: linear-gradient(135deg, #84cc16, #78dbff);
}

.unread-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  background: #e879c6;
  border-radius: 50%;
  margin-left: 8px;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.5; }
  100% { opacity: 1; }
}
</style>
Powered by Gitea
Version: 1.23.7 Page:
94ms
Template:
10ms
Licenses
API