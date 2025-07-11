<template>
  <div class="main-layout" @click="handleClickOutside">
    <nav class="navbar">
      <div class="navbar-left">
        <div class="navbar-logo">
          <img src="/src/assets/logo.png" alt="Social Network Logo" class="logo-image" />
          <span class="navbar-title">Social Network</span>
        </div>
      </div>
      
      <div class="navbar-center">
        <div class="search-bar">
          <MagnifyingGlassIcon class="search-icon" />
          <input
            type="text"
            :placeholder="getSearchPlaceholder()"
            v-model="searchQuery"
            @keyup.enter="searchProfiles"
            @input="onSearchInput"
          />
          <div v-if="showResults && searchResults.length > 0" class="search-results-dropdown">
            <div
              v-for="user in searchResults"
              :key="user.id"
              class="search-result-item"
              @click="goToProfile(user.id)"
            >
              <img :src="getFullImageUrl(user.avatar || user.imagePath)" class="search-result-avatar" />
              <span class="search-result-name">{{ user.nickname || (user.firstName + ' ' + user.lastName) }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="navbar-right">
        <div class="notification-wrapper">
          <button class="nav-icon-btn" @click.stop="toggleNotificationDropdown">
            <BellIcon class="icon" :class="{ 'has-new': notificationStore.hasNewNotifications }" />
            <span v-if="notificationStore.unreadCount > 0" class="notification-badge">{{ notificationStore.unreadCount }}</span>
          </button>
          
          <div v-if="showNotificationDropdown" class="notification-dropdown">
            <div class="notification-header">
              <h3>Notifications</h3>
              <button 
                @click="notificationStore.markAllAsRead()" 
                v-if="notificationStore.unreadCount > 0" 
                class="mark-all-read"
              >
                Marquer tout lu ({{ notificationStore.unreadCount }})
              </button>
            </div>
            <div v-if="notificationStore.loading" class="notification-loading">
              <div class="loading-spinner-sm"></div>
              <span>Chargement...</span>
            </div>
            <div v-else-if="notificationStore.recentUnreadNotifications.length === 0" class="notification-empty">
              Aucune nouvelle notification
            </div>
            <div v-else class="notification-list">
              <template v-for="notif in notificationStore.recentUnreadNotifications" :key="notif.id">
                <div v-if="notif.type === 'FOLLOW'" class="notification-item follow-request">
                  <div class="notification-avatar">
                    <img v-if="notif.user?.avatar || notif.user?.imagePath" :src="getFullImageUrl(notif.user.avatar || notif.user.imagePath)" alt="Avatar" />
                    <div v-else class="default-avatar">{{ notif.user?.firstName?.[0] || 'U' }}</div>
                  </div>
                  <div class="notification-content">
                    <div class="notification-text">
                      <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong> souhaite vous suivre
                    </div>
                    <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
                    <div class="notification-actions">
                      <button 
                        @click="notificationStore.respondToFollowRequest(notif.id, 'ACCEPT')" 
                        class="accept-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-xs"></span>
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Accepter' }}
                      </button>
                      <button 
                        @click="notificationStore.respondToFollowRequest(notif.id, 'DECLINE')" 
                        class="decline-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Refuser' }}
                      </button>
                    </div>
                  </div>
                </div>
                <div v-else-if="notif.type === 'GROUP_INVITE'" class="notification-item group-invite">
                  <div class="notification-avatar">
                    <UserGroupIcon class="icon-sm" />
                  </div>
                  <div class="notification-content">
                    <div class="notification-text">
                      <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong> vous invite à rejoindre le groupe <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                    </div>
                    <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
                    <div class="notification-actions">
                      <button 
                        @click="notificationStore.respondToGroupInvite(notif.id, 'ACCEPT')" 
                        class="accept-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-xs"></span>
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Rejoindre' }}
                      </button>
                      <button 
                        @click="notificationStore.respondToGroupInvite(notif.id, 'DECLINE')" 
                        class="decline-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Refuser' }}
                      </button>
                    </div>
                  </div>
                </div>
                <div v-else-if="notif.type === 'GROUP_REQUEST'" class="notification-item group-request">
                  <div class="notification-avatar">
                    <UserGroupIcon class="icon-sm" />
                  </div>
                  <div class="notification-content">
                    <div class="notification-text">
                      <strong>{{ notif.user?.nickname || notif.user?.firstName || 'Utilisateur' }}</strong> demande à rejoindre votre groupe <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                    </div>
                    <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
                    <div class="notification-actions">
                      <button 
                        @click="notificationStore.respondToGroupRequest(notif.id, 'ACCEPT', notif.group?.id)" 
                        class="accept-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        <span v-if="notificationStore.isProcessing(notif.id)" class="loading-spinner-xs"></span>
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Accepter' }}
                      </button>
                      <button 
                        @click="notificationStore.respondToGroupRequest(notif.id, 'DECLINE', notif.group?.id)" 
                        class="decline-btn"
                        :disabled="notificationStore.isProcessing(notif.id)"
                      >
                        {{ notificationStore.isProcessing(notif.id) ? 'En cours...' : 'Refuser' }}
                      </button>
                    </div>
                  </div>
                </div>
                <div v-else-if="notif.type === 'EVENT'"
                      class="notification-item event-notification"
                     @click="notificationStore.markAsReadOnView(notif.id)">
                  <div class="notification-avatar">
                    <CalendarIcon class="icon-sm" />
                  </div>
                  <div class="notification-content">
                    <div class="notification-text">
                      Nouvel événement <strong>{{ notif.event?.title || 'Événement' }}</strong> dans le groupe <strong>{{ notif.group?.name || 'Groupe' }}</strong>
                    </div>
                    <div class="notification-time">{{ notificationStore.formatNotificationTime(notif.createdAt || notif.timestamp) }}</div>
                  </div>
                  <div v-if="!notif.read" class="unread-indicator"></div>
                </div>
              </template>
            </div>
            <div class="notification-footer">
              <button @click="navigateTo('/notifications'); showNotificationDropdown = false">Voir toutes les notifications</button>
            </div>
          </div>
        </div>
        
        <button class="nav-icon-btn nav-chat-icon" @click="navigateTo('/chat')">
          <ChatBubbleLeftRightIcon class="icon" />
        </button>
        
        <div class="navbar-profile">
          <div class="profile-avatar">
            <img v-if="currentUser?.avatar || currentUser?.imagePath" 
                 :src="getFullImageUrl(currentUser.avatar || currentUser.imagePath)" 
                 :alt="currentUser?.firstName || 'User'"
                 class="profile-avatar-img" />
            <span v-else class="profile-avatar-initial">{{ currentUser?.firstName?.[0] || 'U' }}</span>
          </div>
          <ChevronDownIcon class="dropdown-icon" />
        </div>
        
        <button class="nav-icon-btn logout-btn desktop-logout" @click="logout">
          <ArrowRightOnRectangleIcon class="icon" />
        </button>
      </div>
    </nav>

    <aside class="sidebar">
      <div class="sidebar-section">
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/' }"
          @click="navigateTo('/')"
        >
          <HomeIcon />
          <span>Accueil</span>
        </button>
        
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/profile' }"
          @click="navigateTo('/profile')"
        >
          <UserIcon />
          <span>Profil</span>
        </button>
        
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/groups' }"
          @click="navigateTo('/groups')"
        >
          <UserGroupIcon />
          <span>Groupes</span>
        </button>
        
        <button 
          class="sidebar-item" 
          :class="{ active: currentRoute === '/chat' }"
          @click="navigateTo('/chat')"
        >
          <ChatBubbleLeftRightIcon />
          <span>Messages</span>
        </button>
      </div>

      <div class="sidebar-footer">
        <button class="sidebar-item logout-btn" @click="logout">
          <ArrowRightOnRectangleIcon />
          <span>Déconnexion</span>
        </button>
      </div>
    </aside>

    <nav class="bottom-nav">
       <button class="bottom-nav-item" :class="{ active: currentRoute === '/' }" @click="navigateTo('/')">
        <HomeIcon />
        <span>Accueil</span>
      </button>
      <button class="bottom-nav-item" :class="{ active: currentRoute === '/groups' }" @click="navigateTo('/groups')">
        <UserGroupIcon />
        <span>Groupes</span>
      </button>
      <button class="bottom-nav-item" :class="{ active: currentRoute.startsWith('/chat') }" @click="navigateTo('/chat')">
        <ChatBubbleLeftRightIcon />
        <span>Chat</span>
      </button>
      <button class="bottom-nav-item" :class="{ active: currentRoute.startsWith('/profile') }" @click="navigateTo('/profile')">
        <UserIcon />
        <span>Profil</span>
      </button>
    </nav>

    <main class="main-content">
      <div class="content-wrapper">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script>
import {
  MagnifyingGlassIcon,
  BellIcon,
  ChatBubbleLeftRightIcon,
  ChevronDownIcon,
  ArrowRightOnRectangleIcon,
  HomeIcon,
  UserIcon,
  UserGroupIcon,
  HeartIcon,
  DocumentDuplicateIcon,
  CalendarIcon
} from '@heroicons/vue/24/outline'
import { useUserStore } from '../stores/userStore'
import { useNotificationStore } from '../stores/notificationStore'
import { useMainStore } from '../stores/postsStore'
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import wsService from '../services/websocketService'

export default {
  name: 'MainLayout',
  components: {
    MagnifyingGlassIcon,
    BellIcon,
    ChatBubbleLeftRightIcon,
    ChevronDownIcon,
    ArrowRightOnRectangleIcon,
    HomeIcon,
    UserIcon,
    UserGroupIcon,
    HeartIcon,
    DocumentDuplicateIcon,
    CalendarIcon
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const userStore = useUserStore()
    const notificationStore = useNotificationStore()
    const mainStore = useMainStore()
    
    const { getFullImageUrl } = mainStore
    
    const showResults = ref(false)
    const searchQuery = ref('')
    const searchResults = ref([])
    const showUserMenu = ref(false)
    const showNotificationDropdown = ref(false)
    
    let searchTimeout = null
    let removeWebSocketListener = null

    const isLoggedIn = computed(() => userStore.isLoggedIn)
    const currentUser = computed(() => userStore.user)
    const currentRoute = computed(() => route.path)
    
    function navigateTo(path) {
        router.push(path)
    }
    
    async function toggleNotificationDropdown() {
      showNotificationDropdown.value = !showNotificationDropdown.value
      if (showNotificationDropdown.value) {
        await notificationStore.fetchNotifications()
      }
    }
    
    function handleClickOutside(event) {
      const notificationWrapper = document.querySelector('.notification-wrapper');
      if (notificationWrapper && !notificationWrapper.contains(event.target)) {
        showNotificationDropdown.value = false;
      }
    }

    async function logout() {
      try {
        await fetch('http://localhost:8081/logout', {
          method: 'POST',
          credentials: 'include',
        });
      } catch (error) {
        console.error('Error during logout:', error);
      } finally {
        userStore.logout()
        router.push('/login');
      }
    }

    async function searchProfiles() {
      if (!searchQuery.value.trim()) {
        searchResults.value = [];
        showResults.value = false;
        return;
      }
      try {
        const response = await fetch('http://localhost:8081/allUsers', {
          credentials: 'include'
        });
        if (response.ok) {
          const data = await response.json();
          const users = data.users || [];
          const query = searchQuery.value.trim().toLowerCase();
          searchResults.value = users.filter(user => 
            (user.nickname && user.nickname.toLowerCase().includes(query)) ||
            (user.firstName && user.firstName.toLowerCase().includes(query)) ||
            (user.lastName && user.lastName.toLowerCase().includes(query))
          );
          showResults.value = true;
        }
      } catch (e) {
        searchResults.value = [];
        showResults.value = false;
      }
    }

    function onSearchInput() {
      clearTimeout(searchTimeout);
      if (!searchQuery.value.trim()) {
        searchResults.value = [];
        showResults.value = false;
        return;
      }
      searchTimeout = setTimeout(() => {
        searchProfiles();
      }, 300);
    }

    function goToProfile(userId) {
      router.push(`/profile/${userId}`);
      showResults.value = false;
      searchQuery.value = '';
    }

    function getSearchPlaceholder() {
      if (window.innerWidth <= 480) {
        return "Recherche...";
      } else if (window.innerWidth <= 768) {
        return "Rechercher...";
      } else {
        return "Rechercher des utilisateurs...";
      }
    }

    onMounted(async () => {
      try {
        await userStore.reloadUserAfterRefresh()
        if (userStore.isLoggedIn) {
          await notificationStore.fetchNotifications()
          removeWebSocketListener = wsService.addListener('mainLayout', (data) => {
            if (data.action === 'notification') {
              notificationStore.handleNewNotification(data.notification)
            }
          });
        }
      } catch (error) {
        console.error('Error in MainLayout mount:', error)
      }
      document.addEventListener('click', handleClickOutside);
    })

    onUnmounted(() => {
      if (removeWebSocketListener) removeWebSocketListener();
      clearTimeout(searchTimeout);
      document.removeEventListener('click', handleClickOutside);
    })

    return {
      router,
      currentRoute,
      searchQuery,
      searchResults,
      showResults,
      showUserMenu,
      showNotificationDropdown,
      userStore,
      notificationStore,
      navigateTo,
      toggleNotificationDropdown,
      logout,
      searchProfiles,
      onSearchInput,
      goToProfile,
      getFullImageUrl,
      getSearchPlaceholder,
      isLoggedIn,
      currentUser,
      getSearchPlaceholder
    }
  }
}
</script>

<style scoped>
.main-layout {
  display: flex;
  min-height: 100vh;
  background: #0a0a0f;
}

/* --- Top Navbar --- */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: rgba(15, 15, 23, 0.9);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  z-index: 1000;
}
.navbar-left, .navbar-right {
  display: flex;
  align-items: center;
  gap: 15px;
}
.navbar-logo { display: flex; align-items: center; gap: 10px; }
.logo-image { width: 36px; height: 36px; }
.navbar-title {
  font-size: 1.2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
.navbar-center { flex: 1; max-width: 400px; margin: 0 40px; }
.search-bar { position: relative; width: 100%; }
.search-bar input {
  width: 100%;
  padding: 8px 15px 8px 45px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 25px;
  color: #ffffff;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}
.search-bar input:focus {
  outline: none;
  border-color: #e879c6;
  box-shadow: 0 0 0 3px rgba(232, 121, 198, 0.2);
}
.search-icon { position: absolute; left: 15px; top: 50%; transform: translateY(-50%); color: rgba(255, 255, 255, 0.4); width: 20px; height: 20px;}
.nav-icon-btn {
  position: relative;
  background: none;
  border: none;
  border-radius: 50%;
  padding: 8px;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.3s ease;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.nav-icon-btn:hover { background: rgba(255, 255, 255, 0.1); color: #e879c6; }
.nav-icon-btn .icon { width: 22px; height: 22px; }
.notification-badge {
  position: absolute;
  top: 4px;
  right: 4px;
  background: #e879c6;
  color: black;
  font-size: 10px;
  font-weight: bold;
  min-width: 16px;
  height: 16px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.navbar-profile { display: flex; align-items: center; cursor: pointer; }
.profile-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.profile-avatar-img { width: 100%; height: 100%; object-fit: cover; }
.profile-avatar-initial { font-weight: bold; font-size: 14px; color: black; }
.dropdown-icon { display: none; } /* Hidden for simplicity */

/* --- Notification Dropdown --- */
.notification-wrapper {
  position: relative;
}

.notification-dropdown {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  width: 380px;
  max-height: 450px;
  background: rgba(25, 25, 35, 0.98);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(20px);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  z-index: 1100;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.notification-header h3 {
  color: #fff;
  margin: 0;
  font-size: 1rem;
}

.mark-all-read {
  background: none;
  border: none;
  color: #78c7ff;
  font-size: 0.8rem;
  cursor: pointer;
}

.notification-loading,
.notification-empty {
  text-align: center;
  padding: 20px;
  color: rgba(255, 255, 255, 0.6);
}

.notification-list {
  flex-grow: 1;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  transition: background 0.2s;
}

.notification-item:hover {
  background: rgba(232, 121, 198, 0.1);
}

.notification-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  flex-shrink: 0;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e879c6, #78c7ff);
}
.notification-avatar img { width: 100%; height: 100%; object-fit: cover; }
.notification-avatar .icon-sm { width: 20px; height: 20px; color: #fff; }

.notification-content {
  flex-grow: 1;
}

.notification-text {
  font-size: 0.9rem;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 4px;
}

.notification-time {
  font-size: 0.75rem;
  color: #78c7ff;
}

.notification-actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}
.notification-actions .accept-btn,
.notification-actions .decline-btn {
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: bold;
  cursor: pointer;
  border: 1px solid transparent;
}
.notification-actions .accept-btn { background: #10b981; color: #fff; }
.notification-actions .decline-btn { background: #374151; color: #d1d5db; }

.notification-footer {
  padding: 12px 16px;
  text-align: center;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}
.notification-footer button {
  background: none;
  border: none;
  color: #78c7ff;
  font-weight: bold;
  cursor: pointer;
}

/* --- Desktop Sidebar --- */
.sidebar {
  position: fixed;
  top: 60px;
  left: 0;
  width: 250px;
  height: calc(100vh - 60px);
  background: #181824;
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px 0;
  display: flex;
  flex-direction: column;
  z-index: 999;
}
.sidebar-section { margin-bottom: 25px; }
.sidebar-section-title {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.8rem;
  text-transform: uppercase;
  margin: 0 20px 10px 20px;
}
.sidebar-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 12px 20px;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: left;
}
.sidebar-item:hover { background: rgba(255, 255, 255, 0.05); }
.sidebar-item.active {
  background: linear-gradient(90deg, rgba(232, 121, 198, 0.2) 0%, transparent 100%);
  color: #e879c6;
  border-right: 3px solid #e879c6;
}
.sidebar-item > svg { width: 22px; height: 22px; }
.sidebar-footer { margin-top: auto; }

/* --- Main Content --- */
.main-content {
  margin-left: 250px;
  width: calc(100% - 250px);
  padding-top: 60px;
  position: relative;
  z-index: 1;
}
.content-wrapper { padding: 30px; max-width: 1200px; margin: 0 auto; }

/* --- Mobile Bottom Nav --- */
.bottom-nav {
  display: none; /* Hidden by default */
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: rgba(15, 15, 23, 0.9);
  backdrop-filter: blur(20px);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  justify-content: space-around;
  align-items: center;
  z-index: 1000;
}
.bottom-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  font-size: 10px;
  transition: color 0.3s ease;
}
.bottom-nav-item.active { color: #e879c6; }
.bottom-nav-item > svg { width: 24px; height: 24px; }

/* --- Responsive Adjustments --- */
@media (max-width: 1024px) {
  .sidebar { width: 200px; }
  .main-content { margin-left: 200px; width: calc(100% - 200px); }
  .navbar-center { margin: 0 20px; }
}

@media (max-width: 768px) {
  .sidebar { display: none; }
  .main-content { margin-left: 0; width: 100%; padding-bottom: 60px; /* Space for bottom nav */ }
  
  /* Keep search bar visible but make it more compact */
  .navbar {
    padding: 0 10px;
  }
  
  .navbar-center { 
    flex: 1;
    margin: 0 10px;
    max-width: none;
  }
  
  .navbar-title { display: none; }
  .logo-image { width: 28px; height: 28px; }
  
  .search-bar input {
    padding: 6px 12px 6px 35px;
    font-size: 0.85rem;
    border-radius: 20px;
    min-height: 32px;
  }
  
  .search-icon {
    left: 12px;
    width: 16px;
    height: 16px;
  }
  
  .desktop-logout { display: none; }
  .nav-chat-icon { display: none; }
  .bottom-nav { display: flex; }
}

@media (max-width: 480px) {
  .navbar {
    padding: 0 8px;
  }
  
  .navbar-center {
    margin: 0 8px;
  }
  
  .search-bar input {
    padding: 5px 10px 5px 30px;
    font-size: 0.8rem;
  }
  
  .search-icon {
    left: 10px;
    width: 14px;
    height: 14px;
  }
  
  .nav-icon-btn {
    width: 36px;
    height: 36px;
    padding: 6px;
  }
  
  .nav-icon-btn .icon {
    width: 18px;
    height: 18px;
  }
}

.search-results-dropdown {
  position: absolute;
  top: 110%;
  left: 0;
  width: 100%;
  background: rgba(25, 25, 35, 0.98);
  border-radius: 10px;
  box-shadow: 0 4px 24px rgba(0,0,0,0.25);
  z-index: 2000;
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid rgba(255,255,255,0.1);
  backdrop-filter: blur(20px);
}

@media (max-width: 768px) {
  .search-results-dropdown {
    max-height: 250px;
    border-radius: 8px;
    box-shadow: 0 2px 16px rgba(0,0,0,0.3);
  }
}

@media (max-width: 480px) {
  .search-results-dropdown {
    max-height: 200px;
    border-radius: 6px;
  }
  
  .search-result-item {
    padding: 8px;
    gap: 8px;
  }
  
  .search-result-avatar {
    width: 28px;
    height: 28px;
  }
  
  .search-result-name {
    font-size: 0.85rem;
  }
}

.search-result-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  cursor: pointer;
  transition: background 0.2s;
}
.search-result-item:hover { background: rgba(232, 121, 198, 0.1); }
.search-result-avatar { width: 32px; height: 32px; border-radius: 50%; object-fit: cover; }
.search-result-name { color: #fff; }
</style>