<template>
  <div class="profile-page">
    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <span>Chargement du profil...</span>
    </div>

    <!-- Profile Content -->
    <div v-else-if="user" class="profile-container">
      <!-- Profile Header Section -->
      <div class="profile-header-section">
        <!-- Cover Image Area -->
        <div class="profile-cover">
          <div class="profile-cover-gradient"></div>
        </div>

        <!-- Main Profile Info -->
        <div class="profile-main-info">
          <!-- Avatar Section -->
          <div class="avatar-section">
            <div class="avatar-container">
              <img v-if="user.avatar || user.imagePath"
                   :src="getFullImageUrl(user.avatar || user.imagePath)"
                   alt="Avatar" 
                   class="profile-avatar" />
              <div v-else class="default-avatar">
                {{ (user.firstName?.[0] || 'U').toUpperCase() }}
              </div>
              <div v-if="isOnline" class="online-indicator"></div>
            </div>
          </div>

          <!-- User Info Section -->
          <div class="user-info-section">
            <div class="user-names-header">
              <div class="user-names">
                <h1 class="display-name">{{ user.nickname || (user.firstName + ' ' + user.lastName) }}</h1>
                <p v-if="user.nickname" class="real-name">{{ user.firstName }} {{ user.lastName }}</p>
              </div>
              
              <!-- Privacy Status Switch & Edit Button (Own Profile) -->
              <div v-if="isOwnProfile" class="profile-controls">
                <div class="privacy-toggle">
                  <div class="privacy-switch-container">
                    <span class="privacy-label">{{ isPrivate ? 'Private' : 'Public' }}</span>
                    <div class="privacy-switch" @click="askPrivacyChange(isPrivate ? 'public' : 'private')">
                      <div class="switch-track" :class="{ active: isPrivate }">
                        <div class="switch-thumb" :class="{ active: isPrivate }">
                          <Lock v-if="isPrivate" class="switch-icon" />
                          <Globe v-else class="switch-icon" />
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                
              </div>
            </div>

            <!-- Stats Row -->
            <div class="profile-stats">
              <div class="stat-item">
                <span class="stat-number">{{ visibleFollowersCount }}</span>
                <span class="stat-label">Followers</span>
              </div>
              <div class="stat-item">
                <span class="stat-number">{{ visibleFollowingCount }}</span>
                <span class="stat-label">Following</span>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="action-buttons">
              <!-- For other users -->
              <div v-if="!isOwnProfile" class="external-actions">
                <button v-if="isFollowing" @click="unfollow" class="btn-unfollow">
                  <UserCheck class="btn-icon" />
                  Following
                </button>
                <button v-else-if="isPrivate && !isFollowRequested" @click="follow" class="btn-follow">
                  <UserPlus class="btn-icon" />
                  Request
                </button>
                <button v-else-if="!isFollowing && isPublic && !isFollowRequested" @click="follow" class="btn-follow">
                  <UserPlus class="btn-icon" />
                  Follow
                </button>
                <span v-if="isFollowRequested" class="request-pending">
                  <Clock class="btn-icon" />
                  Request Sent
                </span>
                
              </div>
              
              <!-- For own profile -->
              <div v-else class="own-actions">
                <!-- Removed edit and settings buttons - using header icon instead -->
              </div>
            </div>
          </div>
        </div>

        <!-- Bio Section -->
        <div v-if="user.about" class="bio-section">
          <p class="bio-text">{{ user.about }}</p>
        </div>

        <!-- Additional Info Chips -->
        <div class="info-chips">
          <div v-if="user.dateOfBirth" class="info-chip">
            <Calendar class="chip-icon" />
            Born {{ new Date(user.dateOfBirth).toLocaleDateString() }}
          </div>
          <div v-if="user.email && isOwnProfile" class="info-chip">
            <Mail class="chip-icon" />
            {{ user.email }}
          </div>
          <div class="info-chip">
            <Lock v-if="isPrivate" class="chip-icon" />
            <Globe v-else class="chip-icon" />
            {{ isPrivate ? 'Private' : 'Public' }} Profile
          </div>
        </div>
      </div>

      <!-- Content Sections -->
      <div class="profile-content">
        <!-- Navigation Tabs -->
        <div class="content-nav">
          <button 
            :class="['nav-tab', { active: activeTab === 'posts' }]"
            @click="activeTab = 'posts'"
          >
            <Image class="tab-icon" />
            Posts
            <span class="tab-count">{{ visiblePostsCount }}</span>
          </button>
          <button 
            :class="['nav-tab', { active: activeTab === 'followers' }]"
            @click="activeTab = 'followers'"
          >
            <Users class="tab-icon" />
            Followers
            <span class="tab-count">{{ visibleFollowersCount }}</span>
          </button>
          <button 
            :class="['nav-tab', { active: activeTab === 'following' }]"
            @click="activeTab = 'following'"
          >
            <UserPlus class="tab-icon" />
            Following
            <span class="tab-count">{{ visibleFollowingCount }}</span>
          </button>
        </div>

        <!-- Tab Content -->
        <div class="tab-content">
          <!-- Posts Tab -->
          <div v-if="activeTab === 'posts'" class="posts-feed">
            <!-- Private Profile - Hide posts from non-followers -->
            <div v-if="!isOwnProfile && isPrivate && !isFollowing" class="private-posts-message">
              <Lock class="private-icon" />
              <h3>Posts are private</h3>
              <p>Follow this account to see their posts.</p>
            </div>
            <!-- Show posts for own profile, public profiles, or if following private profile -->
            <div v-else-if="formattedUserPosts.length === 0" class="empty-state">
              <Image class="empty-icon" />
              <h3>No posts yet</h3>
              <p>{{ isOwnProfile ? 'Start sharing your thoughts!' : 'This user hasn\'t posted anything yet.' }}</p>
            </div>
            <!-- Show posts only if allowed to see them -->
            <template v-else-if="isOwnProfile || isPublic || isFollowing" v-for="post in formattedUserPosts" :key="post && post.id">
              <div v-if="post">
                <div class="post-card">
                  <div class="post-header">
                    <div class="user-avatar">
                      <img v-if="post.userAvatar" 
                           :src="getFullImageUrl(post.userAvatar)" 
                           :alt="post.userName" />
                      <span v-else>{{ post.userInitials }}</span>
                    </div>
                    <div class="post-meta">
                      <h3 class="user-name">{{ post.userName }}</h3>
                      <div class="post-time-container">
                        <component :is="getPostPrivacyIconComponent(post.visibility)" :size="14" class="privacy-icon-component" />
                        <span class="post-time">{{ post.timeAgo }}</span>
                      </div>
                    </div>
                  
                  </div>
                  <div class="post-content">
                    <p>{{ post.content }}</p>
                    <div v-if="post.hasImage" class="post-image-container">
                      <img v-if="post.imagePath" :src="getImageUrl(post.imagePath)" class="post-image" alt="Post image" @error="handleImageError" />
                      <div v-else class="post-image-placeholder">
                        <Image :size="60" />
                        <span>Image</span>
                      </div>
                    </div>
                  </div>
                  <div class="post-actions">
                    <button class="action-btn" @click="toggleComments(post.id)">
                      <MessageCircle :size="18" />
                      {{ post.commentsCount || 0 }} Commentaire{{ (post.commentsCount || 0) > 1 ? 's' : '' }}
                    </button>
                  </div>
                  
                  <!-- Comments Section -->
                  <div v-if="activeCommentPostId === post.id" class="comments-section">
                    <div v-if="post.comments && post.comments.length > 0" class="comments-list">
                      <div v-for="comment in post.comments" :key="comment.id" class="comment">
                        <div class="comment-header">
                          <div class="user-avatar small">
                            <img v-if="comment.author?.avatar || comment.author?.imagePath || comment.author?.image" 
                                 :src="getFullImageUrl(comment.author.avatar || comment.author.imagePath || comment.author.image)" 
                                 :alt="comment.userName || comment.author?.nickname || 'Utilisateur'" />
                            <span v-else>{{ comment.userInitials || (comment.author?.nickname ? comment.author.nickname.slice(0,2).toUpperCase() : 'U') }}</span>
                          </div>
                          <div class="comment-meta">
                            <h4 class="comment-user">
                              {{ comment.userName || comment.author?.nickname || 'Utilisateur' }}
                            </h4>
                            <span class="comment-time">
                              {{ isValidDate(comment.createdAt) ? formatDate(comment.createdAt) : '' }}
                            </span>
                          </div>
                        </div>
                        <div class="comment-content">
                          <p>{{ comment.content }}</p>
                          <img v-if="getCommentImage(comment)" :src="getImageUrl(getCommentImage(comment))" class="comment-image" @error="handleImageError" />
                        </div>
                      </div>
                    </div>
                    <div class="comment-form">
                      <div class="user-avatar small">
                        <img v-if="userStore.user?.avatar || userStore.user?.imagePath" 
                             :src="getFullImageUrl(userStore.user.avatar || userStore.user.imagePath)" 
                             alt="Avatar" />
                        <span v-else>{{ currentUserInitials }}</span>
                      </div>
                      <input 
                        type="text" 
                        v-model="newCommentText" 
                        placeholder="Écrire un commentaire..." 
                        class="comment-input"
                        @keyup.enter="submitProfileComment(post.id)"
                      />
                      <label class="comment-attach-btn">
                        <input type="file" @change="handleCommentImageChange" accept="image/*" style="display: none" />
                        <Image :size="18" />
                      </label>
                      <button class="comment-send-btn" @click="submitProfileComment(post.id)">
                        <Send :size="18" />
                      </button>
                    </div>
                    <div v-if="commentImagePreview" class="comment-image-preview-container">
                      <img :src="commentImagePreview" class="comment-image-preview" />
                      <button class="remove-image-btn" @click="removeCommentImage">×</button>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </div>

          <!-- Followers Tab -->
          <div v-if="activeTab === 'followers'" class="followers-grid">
            <!-- Private Profile - Hide followers from non-followers -->
            <div v-if="!isOwnProfile && isPrivate && !isFollowing" class="private-followers-message">
              <Lock class="private-icon" />
              <h3>Followers are private</h3>
              <p>Follow this account to see their followers.</p>
            </div>
            <!-- Show followers for own profile, public profiles, or if following private profile -->
            <div v-else-if="profileFollowers.length === 0" class="empty-state">
              <Users class="empty-icon" />
              <h3>No followers yet</h3>
              <p>{{ isOwnProfile ? 'People who follow you will appear here.' : 'This user has no followers yet.' }}</p>
            </div>
            <!-- Show followers only if allowed to see them -->
            <div v-else-if="isOwnProfile || isPublic || isFollowing" v-for="follower in profileFollowers" :key="follower.id" class="user-card">
              <img v-if="follower.imagePath || follower.image || follower.avatar" 
                   :src="getFullImageUrl(follower.imagePath || follower.image || follower.avatar)" 
                   :alt="follower.nickname || follower.firstName" 
                   class="user-avatar" />
              <div v-else class="user-avatar default-avatar">
                {{ (follower.firstName?.[0] || follower.nickname?.[0] || 'U').toUpperCase() }}
              </div>
              <div class="user-info">
                <h4>{{ follower.nickname || follower.firstName }}</h4>
                <p>{{ follower.firstName }} {{ follower.lastName }}</p>
              </div>
              
            </div>
          </div>

          <!-- Following Tab -->
          <div v-if="activeTab === 'following'" class="following-grid">
            <!-- Private Profile - Hide following from non-followers -->
            <div v-if="!isOwnProfile && isPrivate && !isFollowing" class="private-following-message">
              <Lock class="private-icon" />
              <h3>Following list is private</h3>
              <p>Follow this account to see who they follow.</p>
            </div>
            <!-- Show following for own profile, public profiles, or if following private profile -->
            <div v-else-if="profileFollowing.length === 0" class="empty-state">
              <UserPlus class="empty-icon" />
              <h3>Not following anyone</h3>
              <p>{{ isOwnProfile ? 'People you follow will appear here.' : 'This user isn\'t following anyone yet.' }}</p>
            </div>
            <!-- Show following only if allowed to see them -->
            <div v-else-if="isOwnProfile || isPublic || isFollowing" v-for="following in profileFollowing" :key="following.id" class="user-card">
              <img v-if="following.imagePath || following.image || following.avatar" 
                   :src="getFullImageUrl(following.imagePath || following.image || following.avatar)" 
                   :alt="following.nickname || following.firstName" 
                   class="user-avatar" />
              <div v-else class="user-avatar default-avatar">
                {{ (following.firstName?.[0] || following.nickname?.[0] || 'U').toUpperCase() }}
              </div>
              <div class="user-info">
                <h4>{{ following.nickname || following.firstName }}</h4>
                <p>{{ following.firstName }} {{ following.lastName }}</p>
              </div>
             
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Private Profile Message -->
    <div v-else-if="!isOwnProfile && !isFollowing && isPrivate" class="private-profile">
      <div class="private-content">
        <div class="private-avatar">
          <img v-if="user?.avatar || user?.imagePath"
               :src="getFullImageUrl(user?.avatar || user?.imagePath)"
               alt="Avatar" />
          <div v-else class="default-avatar">
            {{ (user?.firstName?.[0] || 'U').toUpperCase() }}
          </div>
        </div>
        <h2>{{ user?.nickname || (user?.firstName + ' ' + user?.lastName) || 'User Profile' }}</h2>
        <div class="private-stats">
          <span>{{ user?.followersCount ?? 0 }} followers</span>
          <span>{{ user?.followingCount ?? 0 }} following</span>
        </div>
        <div class="private-message">
          <Lock class="lock-icon" />
          <h3>This account is private</h3>
          <p>Follow this account to see their posts and activity.</p>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else class="error-state">
      <AlertCircle class="error-icon" />
      <h3>Profile not found</h3>
      <p>This user doesn't exist or has been removed.</p>
    </div>

    <!-- Privacy Confirmation Modal -->
    <div v-if="showPrivacyConfirm" class="modal-overlay">
      <div class="modal-content">
        <div class="modal-header">
          <h3>Change Profile Privacy</h3>
        </div>
        <div class="modal-body">
          <p>Are you sure you want to make your profile {{ pendingPrivacyStatus }}?</p>
          <div v-if="pendingPrivacyStatus === 'public'" class="privacy-info">
            <p class="info-text">Your posts and information will be visible to everyone.</p>
          </div>
          <div v-else class="privacy-info">
            <p class="info-text">Only approved followers will be able to see your posts and information.</p>
          </div>
        </div>
        <div class="modal-actions">
          <button @click="confirmPrivacyChange" class="confirm-btn">Confirm</button>
          <button @click="cancelPrivacyChange" class="cancel-btn">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '../stores/userStore'
import { useFollowStore } from '../stores/followStore'
import { useMainStore } from '../stores/postsStore'
import { storeToRefs } from 'pinia'
import { 
  Lock, 
  Globe, 
  Settings, 
  Edit3, 
  Heart, 
  MessageCircle, 
  UserPlus, 
  UserCheck, 
  Mail, 
  Calendar, 
  Users, 
  Image, 
  Clock,
  CheckCircle,
  Send,
  AlertCircle,
  MoreHorizontal,
  X,
  Info,
  Loader
} from 'lucide-vue-next'

const route = useRoute()
const userStore = useUserStore()
const followStore = useFollowStore()
const mainStore = useMainStore()

// Get everything from store - simple!
const {
  activeCommentPostId,
  newCommentText,
  commentSelectedImage,
  commentImagePreview
} = storeToRefs(mainStore)

// All methods from store - simple!
const {
  handleCommentImageChange,
  removeCommentImage,
  toggleComments,
  submitCommentAction,
  getImageUrl,
  getFullImageUrl,
  formatDate,
  handleImageError,
  getCommentImage,
  isValidDate,
  getPostPrivacyIconComponent,
  formatPostsData,
  fetchUserPosts
} = mainStore

const user = ref(null)
const profileFollowers = ref([])
const profileFollowing = ref([])
const following = followStore.following // Utilise directement le ref du store Pinia
const posts = ref([])
const loading = ref(true)
const isOwnProfile = computed(() => !route.params.id)
const isPublic = computed(() => user.value?.status?.toLowerCase() !== 'private')
const isPrivate = computed(() => user.value?.status?.toLowerCase() === 'private')
const showPrivacyConfirm = ref(false)
const pendingPrivacy = ref(null)
const pendingPrivacyStatus = ref(null)
const isFollower = computed(() => !!user.value?.follower)
const isFollowing = computed(() => !!user.value?.following)
const isFollowRequested = computed(() => {
  // Check both user data and follow store for pending requests
  return !!user.value?.requestPending || followStore.isFollowRequested(user.value?.id)
})
const isOnline = ref(true) // You can implement actual online status later
const activeTab = ref('posts') // For tab navigation

// Computed properties for tab counts that respect privacy
const visiblePostsCount = computed(() => {
  if (isOwnProfile.value || isPublic.value || isFollowing.value) {
    return posts.value.length;
  }
  return 0; // Hide count for private profiles when not following
});

const visibleFollowersCount = computed(() => {
  if (isOwnProfile.value || isPublic.value || isFollowing.value) {
    return profileFollowers.value.length;
  }
  return 0; // Hide count for private profiles when not following
});

const visibleFollowingCount = computed(() => {
  if (isOwnProfile.value || isPublic.value || isFollowing.value) {
    return profileFollowing.value.length;
  }
  return 0; // Hide count for private profiles when not following
});

const fetchUserPostsAndSet = async (userId) => {
  const userPosts = await fetchUserPosts(userId);
  posts.value = userPosts;
}

const askPrivacyChange = (newPrivacy) => {
  pendingPrivacy.value = newPrivacy
  pendingPrivacyStatus.value = newPrivacy
  showPrivacyConfirm.value = true
}
const confirmPrivacyChange = async () => {
  if (!user.value) return
  const newStatus = pendingPrivacy.value === 'public' ? 'public' : 'private'
  await userStore.changeStatus(newStatus)
  showPrivacyConfirm.value = false
  pendingPrivacy.value = null
  pendingPrivacyStatus.value = null
  await loadProfile()
}

const cancelPrivacyChange = () => {
  showPrivacyConfirm.value = false
  pendingPrivacy.value = null
  pendingPrivacyStatus.value = null
}
function follow() {
  if (!user.value?.id) return
  if (isPublic.value) {
    followStore.follow(user.value.id)
      .then(async () => {
        await loadProfile() // Recharge le profil pour mettre à jour le bouton
      })
      .catch(() => {
        // Optionnel : afficher une erreur
      })
  } else {
    followStore.requestFollow(user.value.id)
      .then(async () => {
        await loadProfile() // Recharge le profil pour mettre à jour le bouton
      })
      .catch(() => {
        // Optionnel : afficher une erreur
      })
  }
}

function unfollow() {
  if (!user.value?.id) return
  followStore.unfollow(user.value.id)
    .then(async () => {
      await loadProfile() // Recharge le profil pour mettre à jour le bouton
    })
}

async function loadProfile() {
  loading.value = true
  let userId
  if (route.params.id) {
    // Profil d'un autre utilisateur
    userId = route.params.id
    const res = await fetch(`http://localhost:8081/userData?userId=${userId}`, { credentials: 'include' })
    const data = await res.json()
    user.value = (data.users && data.users.length > 0) ? data.users[0] : null
    // Affichage : followers/following du profil visité
    await followStore.fetchFollowers(userId)
    profileFollowers.value = [...followStore.followers]
    await followStore.fetchFollowing(userId)
    profileFollowing.value = [...followStore.following]
    // Pour le bouton : always fetch following du user courant
    if (userStore.user?.id) {
      await followStore.fetchFollowing(userStore.user.id)
      console.log('following (utilisateur courant):', following.value)
    }
  } else {
    // Mon propre profil
    await userStore.reloadUserAfterRefresh()
    await userStore.fetchUserData()
    user.value = userStore.user
    userId = user.value?.id
    await followStore.fetchFollowers(userId)
    profileFollowers.value = [...followStore.followers]
    await followStore.fetchFollowing(userId)
    profileFollowing.value = [...followStore.following]
    // Recharge aussi la liste des following du user courant (lui-même)
    await followStore.fetchFollowing(userId)
    console.log('following (utilisateur courant):', following.value)
  }
  if (userId) {
    await fetchUserPostsAndSet(userId)
  } else {
    profileFollowers.value = []
    profileFollowing.value = []
    posts.value = []
  }
  loading.value = false
}

onMounted(loadProfile)
watch(() => route.params.id, loadProfile)

function formatPostTime(timestamp) {
  if (!timestamp) return ''
  const now = new Date()
  const postTime = new Date(timestamp)
  const diffMs = now - postTime
  const diffMins = Math.floor(diffMs / (1000 * 60))
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  
  if (diffMins < 1) return 'now'
  if (diffMins < 60) return `${diffMins}m`
  if (diffHours < 24) return `${diffHours}h`
  if (diffDays < 7) return `${diffDays}d`
  
  return postTime.toLocaleDateString()
}

// Computed properties specific to Profile
const currentUserInitials = computed(() => {
  if (userStore.user?.nickname) {
    return userStore.user.nickname.slice(0, 2).toUpperCase();
  }
  return 'VU';
});

const formattedUserPosts = computed(() => {
  return formatPostsData(posts.value, user.value);
});

// Custom submitComment for profile that refreshes user posts
const submitProfileComment = async (postId) => {
  await submitCommentAction(postId);
  
  // Refresh user posts after comment submission
  if (user.value?.id) {
    await fetchUserPostsAndSet(user.value.id);
  }
};
</script>

<style scoped>
/* Base Profile Page */
.profile-page {
  padding: 1rem;
}

/* Loading State */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 50vh;
  color: #ffffff;
  font-size: 1.1rem;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(232, 121, 198, 0.3);
  border-top: 3px solid #e879c6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Profile Container */
.profile-container {
  max-width: 700px;
  margin: 0 auto;
}

/* Profile Header Section */
.profile-header-section {
  background: #181824;
  border-radius: 18px;
  padding: 0;
  margin-bottom: 1.5rem;
  box-shadow: 0 4px 24px rgba(120, 119, 198, 0.10);
  overflow: hidden;
}

/* Cover Area */
.profile-cover {
  height: 120px;
  background: linear-gradient(135deg, #0f0f14 0%, #16161d 50%, #1a1a22 100%);
  position: relative;
  display: flex;
  align-items: flex-end;
  justify-content: flex-end;
  padding: 0.75rem;
}

.profile-cover-gradient {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(15, 15, 20, 0.95) 0%, rgba(22, 22, 29, 0.9) 50%, rgba(26, 26, 34, 0.95) 100%);
}

.privacy-toggle {
  flex-shrink: 0;
  margin-left: 1rem;
}

.privacy-switch-container {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(8px);
  padding: 0.4rem 0.75rem;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.privacy-label {
  color: white;
  font-size: 0.8rem;
  font-weight: 600;
  white-space: nowrap;
}

.privacy-switch {
  cursor: pointer;
  user-select: none;
}

.switch-track {
  width: 42px;
  height: 22px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 11px;
  position: relative;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.switch-track.active {
  background: linear-gradient(135deg, #e879c6 0%, #78dbff 100%);
  border-color: transparent;
  box-shadow: 0 0 20px rgba(232, 121, 198, 0.4);
}

.switch-thumb {
  width: 18px;
  height: 18px;
  background: white;
  border-radius: 50%;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.switch-thumb.active {
  transform: translateX(20px);
  background: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.switch-icon {
  width: 10px;
  height: 10px;
  color: #6b7280;
  transition: all 0.3s ease;
}

.switch-thumb.active .switch-icon {
  color: #e879c6;
}

.privacy-switch:hover .switch-track {
  transform: scale(1.05);
}

.privacy-switch:hover .switch-track.active {
  box-shadow: 0 0 25px rgba(232, 121, 198, 0.6);
}

.profile-controls {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.btn-edit-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  color: rgba(255, 255, 255, 0.8);
}

.btn-edit-icon:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  color: white;
}

.btn-edit-icon .btn-icon {
  width: 16px;
  height: 16px;
}

/* Main Profile Info */
.profile-main-info {
  padding: 1.5rem;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 2rem;
  position: relative;
}

/* Avatar Section */
.avatar-section {
  flex-shrink: 0;
  position: relative;
  top: -30px;
}

.avatar-container {
  position: relative;
}

.profile-avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
  border: 4px solid #181824;
  box-shadow: 0 4px 16px rgba(232, 121, 198, 0.2);
}

.default-avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e879c6 0%, #78dbff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.5rem;
  font-weight: bold;
  color: white;
  border: 4px solid #181824;
  box-shadow: 0 4px 16px rgba(232, 121, 198, 0.2);
}

.online-indicator {
  position: absolute;
  bottom: 5px;
  right: 5px;
  width: 18px;
  height: 18px;
  background: #10b981;
  border: 2px solid #181824;
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

/* User Info Section */
.user-info-section {
  flex: 1;
  text-align: left;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 1.25rem;
  margin-top: 0;
  padding-top: 0;
}

.user-names-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  width: 100%;
}

.user-names {
  margin: 0;
  flex: 1;
}

.display-name {
  font-size: 1.75rem;
  font-weight: 600;
  color: #ffffff;
  margin: 0 0 0.5rem 0;
}

.real-name {
  font-size: 0.95rem;
  color: #9ca3af;
  margin: 0;
}

/* Profile Stats */
.profile-stats {
  display: flex;
  gap: 2rem;
  margin: 0;
}

.stat-item {
  display: flex;
  flex-direction: row;
  align-items: baseline;
  gap: 0.5rem;
}

.stat-number {
  font-size: 1.25rem;
  font-weight: 700;
  color: #e879c6;
}

.stat-label {
  font-size: 0.9rem;
  color: #9ca3af;
  text-transform: lowercase;
}

/* Action Buttons */
.action-buttons {
  display: flex;
  justify-content: flex-start;
}

.external-actions,
.own-actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.btn-follow,
.btn-unfollow,
.btn-message,
.btn-edit,
.btn-settings {
  padding: 0.5rem 1rem;
  border-radius: 10px;
  font-weight: 600;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.4rem;
  border: none;
}

.btn-follow {
  background: linear-gradient(135deg, #e879c6 0%, #78dbff 100%);
  color: white;
}

.btn-follow:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(232, 121, 198, 0.3);
}

.btn-unfollow {
  background: #23233a;
  color: #9ca3af;
  border: 1px solid #374151;
}

.btn-unfollow:hover {
  background: #374151;
  color: #ef4444;
  border-color: #ef4444;
}

.btn-message {
  background: #23233a;
  color: #9ca3af;
  border: 1px solid #374151;
}

.btn-message:hover {
  background: #374151;
  color: #78dbff;
  border-color: #78dbff;
}

.btn-edit,
.btn-settings {
  background: #23233a;
  color: #9ca3af;
  border: 1px solid #374151;
}

.btn-edit:hover,
.btn-settings:hover {
  background: #374151;
  color: #ffffff;
  transform: translateY(-1px);
}

.request-pending {
  padding: 0.5rem 1rem;
  background: #23233a;
  color: #fbbf24;
  border-radius: 10px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 0.4rem;
  border: 1px solid #374151;
}

.btn-icon {
  width: 16px;
  height: 16px;
}

/* Bio Section */
.bio-section {
  padding: 0 1.5rem 0.75rem;
  text-align: center;
}

.bio-text {
  font-size: 0.95rem;
  color: #d1d5db;
  line-height: 1.5;
  margin: 0;
}

/* Info Chips */
.info-chips {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 0.5rem;
  padding: 0 1.5rem 1.5rem;
}

.info-chip {
  background: #23233a;
  color: #9ca3af;
  padding: 0.4rem 0.8rem;
  border-radius: 16px;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  gap: 0.4rem;
  border: 1px solid #374151;
}

.chip-icon {
  width: 14px;
  height: 14px;
}

/* Profile Content */
.profile-content {
  background: #181824;
  border-radius: 18px;
  overflow: hidden;
  box-shadow: 0 4px 24px rgba(120, 119, 198, 0.10);
}

/* Content Navigation */
.content-nav {
  display: flex;
  border-bottom: 1px solid #374151;
  background: #181824;
}

.nav-tab {
  flex: 1;
  padding: 1rem 0.75rem;
  background: none;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  font-weight: 600;
  color: #9ca3af;
  transition: all 0.3s ease;
  position: relative;
}

.nav-tab:hover {
  background: #23233a;
  color: #d1d5db;
}

.nav-tab.active {
  color: #e879c6;
  background: #23233a;
}

.nav-tab.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(135deg, #e879c6 0%, #78dbff 100%);
}

.tab-icon {
  width: 18px;
  height: 18px;
}

.tab-count {
  background: #374151;
  color: #9ca3af;
  font-size: 0.8rem;
  padding: 0.2rem 0.5rem;
  border-radius: 10px;
  min-width: 20px;
  text-align: center;
}

.nav-tab.active .tab-count {
  background: #e879c6;
  color: white;
}

/* Tab Content */
.tab-content {
  padding: 1.5rem;
  min-height: 300px;
  background: #181824;
}

/* Followers/Following Grids */
.followers-grid,
.following-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 0.4rem;
  min-height: 200px;
  grid-auto-rows: min-content;
}

/* Center empty states in followers/following grids */
.followers-grid .empty-state,
.following-grid .empty-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 180px;
}

.user-card {
  background: #23233a;
  border-radius: 6px;
  padding: 0.4rem;
  box-shadow: 0 1px 4px rgba(120, 119, 198, 0.05);
  border: 1px solid #374151;
  display: flex;
  align-items: center;
  gap: 0.4rem;
  transition: all 0.3s ease;
  min-height: auto;
}

.user-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(120, 119, 198, 0.1);
  border-color: #e879c6;
}

.user-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  object-fit: cover;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-info h4 {
  margin: 0 0 0.1rem 0;
  font-size: 0.8rem;
  font-weight: 600;
  color: #ffffff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-info p {
  margin: 0;
  color: #9ca3af;
  font-size: 0.7rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.follow-btn,
.following-btn {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-weight: 600;
  font-size: 0.7rem;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
  flex-shrink: 0;
}

.follow-btn {
  background: linear-gradient(135deg, #e879c6 0%, #78dbff 100%);
  color: white;
}

.follow-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(232, 121, 198, 0.3);
}

.following-btn {
  background: #23233a;
  color: #9ca3af;
  border: 1px solid #374151;
}

.following-btn:hover {
  background: #374151;
  color: #ef4444;
  border-color: #ef4444;
}

/* Empty States */
.empty-state {
  text-align: center;
  padding: 2rem 1rem;
  color: #9ca3af;
}

.empty-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 0.75rem;
  opacity: 0.5;
}

.empty-state h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
  color: #d1d5db;
}

.empty-state p {
  margin: 0;
  font-size: 0.9rem;
}

/* Private Posts Message */
.private-posts-message {
  text-align: center;
  padding: 3rem 2rem;
  color: #9ca3af;
  background: #23233a;
  border-radius: 16px;
  border: 1px solid #374151;
}

/* Private Followers Message */
.private-followers-message {
  text-align: center;
  padding: 3rem 2rem;
  color: #9ca3af;
  background: #23233a;
  border-radius: 16px;
  border: 1px solid #374151;
}

/* Private Following Message */
.private-following-message {
  text-align: center;
  padding: 3rem 2rem;
  color: #9ca3af;
  background: #23233a;
  border-radius: 16px;
  border: 1px solid #374151;
}

.private-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 1rem;
  color: #9ca3af;
}

.private-posts-message h3,
.private-followers-message h3,
.private-following-message h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.2rem;
  color: #d1d5db;
}

.private-posts-message p,
.private-followers-message p,
.private-following-message p {
  margin: 0;
  font-size: 0.95rem;
  color: #9ca3af;
}

/* Private Profile */
.private-profile {
  max-width: 600px;
  margin: 0 auto;
  text-align: center;
  color: white;
}

.private-content {
  background: #181824;
  border-radius: 18px;
  padding: 3rem 2rem;
  color: #ffffff;
  box-shadow: 0 4px 24px rgba(120, 119, 198, 0.10);
}

.private-avatar {
  margin: 0 auto 1.5rem;
  width: 120px;
  height: 120px;
}

.private-avatar img,
.private-avatar .default-avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.private-content h2 {
  margin: 0 0 1rem 0;
  font-size: 1.8rem;
  font-weight: 700;
  color: #ffffff;
}

.private-stats {
  display: flex;
  justify-content: center;
  gap: 2rem;
  margin-bottom: 2rem;
  font-size: 1.1rem;
  color: #9ca3af;
}

.private-message {
  background: #23233a;
  border-radius: 16px;
  padding: 2rem;
  border: 1px solid #374151;
}

.lock-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 1rem;
  color: #9ca3af;
}

.private-message h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
  color: #d1d5db;
}

.private-message p {
  margin: 0;
  color: #9ca3af;
}

/* Error State */
.error-state {
  max-width: 600px;
  margin: 0 auto;
  text-align: center;
  color: white;
  background: #181824;
  border-radius: 18px;
  padding: 3rem 2rem;
  color: #ffffff;
  box-shadow: 0 4px 24px rgba(120, 119, 198, 0.10);
}

.error-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 1rem;
  color: #ef4444;
}

.error-state h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.5rem;
  color: #d1d5db;
}

.error-state p {
  margin: 0;
  color: #9ca3af;
}

/* Privacy Confirmation Modal */
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
  z-index: 1001;
}

.modal-content {
  background: #181824;
  border-radius: 16px;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  border: 1px solid #374151;
}

.modal-body {
  padding: 1.5rem;
}

.modal-body p {
  margin: 0 0 1rem 0;
  font-size: 1.1rem;
  color: #d1d5db;
}

.privacy-info {
  background: #23233a;
  border-radius: 8px;
  padding: 1rem;
  border-left: 4px solid #e879c6;
}

.info-text {
  margin: 0 !important;
  font-size: 0.95rem !important;
  color: #9ca3af !important;
}

.modal-actions {
  padding: 1.5rem;
  border-top: 1px solid #374151;
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.confirm-btn,
.cancel-btn {
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: all 0.3s ease;
}

.confirm-btn {
  background: linear-gradient(135deg, #e879c6 0%, #78dbff 100%);
  color: white;
}

.confirm-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(232, 121, 198, 0.3);
}

.cancel-btn {
  background: #23233a;
  color: #9ca3af;
  border: 1px solid #374151;
}

.cancel-btn:hover {
  background: #374151;
  color: #ffffff;
}

/* Posts Feed Styles */
.posts-feed {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.posts-feed .post-card {
  background: rgba(15, 15, 23, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
}

.posts-feed .post-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.4);
}

.posts-feed .post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 15px;
}

.posts-feed .user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 2px solid rgba(232, 121, 198, 0.3);
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #000;
  font-weight: 600;
  font-size: 14px;
  flex-shrink: 0;
  overflow: hidden;
}

.posts-feed .user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.posts-feed .user-avatar.small {
  width: 30px;
  height: 30px;
  font-size: 12px;
}

.posts-feed .user-avatar.small img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.posts-feed .post-meta {
  flex: 1;
}

.posts-feed .user-name {
  color: rgba(255, 255, 255, 0.9);
  font-size: 0.95rem;
  font-weight: 600;
  margin: 0 0 2px 0;
}

.posts-feed .post-time-container {
  display: flex;
  align-items: center;
  gap: 6px;
}

.posts-feed .privacy-icon-component {
  opacity: 1;
  color: rgba(255, 255, 255, 0.8);
}

.posts-feed .post-time {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.8rem;
}

.posts-feed .post-menu-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  padding: 5px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.posts-feed .post-menu-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
}

.posts-feed .post-content p {
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
  margin: 0 0 15px 0;
}

.posts-feed .post-image-container {
  margin-bottom: 15px;
}

.posts-feed .post-image {
  width: 100%;
  border-radius: 12px;
  object-fit: cover;
  max-height: 400px;
}

.posts-feed .post-image-placeholder {
  width: 100%;
  height: 200px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.4);
  gap: 10px;
}

.posts-feed .post-image-placeholder span {
  font-size: 0.9rem;
  font-weight: 500;
}

.posts-feed .post-actions {
  display: flex;
  align-items: center;
  gap: 20px;
  padding-top: 15px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.posts-feed .post-actions .action-btn {
  background: none;
  border: none;
  padding: 8px 12px;
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.posts-feed .post-actions .action-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #e879c6;
}

.comments-section {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.comments-list {
  margin-bottom: 15px;
}

.comment {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
  padding: 12px;
  margin-bottom: 10px;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.comment-meta {
  flex: 1;
}

.comment-user {
  font-size: 0.85rem;
  margin: 0 0 2px 0;
  color: rgba(255, 255, 255, 0.9);
}

.comment-time {
  font-size: 0.7rem;
  color: rgba(255, 255, 255, 0.5);
}

.comment-content {
  color: #fff;
}

.comment-content p {
  font-size: 0.9rem;
  margin: 0 0 8px 0;
}

.comment-image {
  max-width: 100%;
  max-height: 150px;
  border-radius: 8px;
  margin-top: 8px;
}

.comment-form {
  display: flex;
  align-items: center;
  gap: 10px;
}

.comment-input {
  flex: 1;
  padding: 8px 15px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: white;
  font-size: 0.9rem;
}

.comment-input:focus {
  outline: none;
  border-color: #e879c6;
  box-shadow: 0 0 0 2px rgba(232, 121, 198, 0.2);
}

.comment-input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.comment-attach-btn,
.comment-send-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.7);
  transition: all 0.2s ease;
}

.comment-attach-btn:hover,
.comment-send-btn:hover {
  background: rgba(232, 121, 198, 0.1);
  border-color: rgba(232, 121, 198, 0.3);
  color: #e879c6;
}

.comment-send-btn {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  border: none;
  color: #000000;
}

.comment-image-preview-container {
  margin-top: 10px;
  position: relative;
  width: 150px;
}

.comment-image-preview {
  width: 100%;
  max-height: 100px;
  border-radius: 8px;
  object-fit: cover;
}

.remove-image-btn {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: none;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.remove-image-btn:hover {
  background: rgba(255, 0, 0, 0.7);
  transform: scale(1.1);
}

/* Responsive Design */
@media (max-width: 768px) {
  .profile-page {
    padding: 0.75rem 0.5rem;
  }
  
  .profile-main-info {
    padding: 1.25rem;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
  }
  
  .avatar-section {
    position: relative;
    top: -25px;
    flex-shrink: 0;
  }
  
  .user-info-section {
    text-align: center;
    align-items: center;
    margin-top: 0;
    padding-top: 0;
  }
  
  .user-names-header {
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }
  
  .privacy-toggle {
    margin-left: 0;
  }
  
  .profile-stats {
    justify-content: center;
    gap: 1.5rem;
  }
  
  .stat-item {
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .action-buttons {
    justify-content: center;
  }
  
  .external-actions,
  .own-actions {
    justify-content: center;
    flex-wrap: wrap;
  }
  
  .content-nav {
    flex-direction: column;
  }
  
  .nav-tab {
    border-bottom: 1px solid #374151;
    padding: 0.75rem;
  }
  
  .nav-tab.active::after {
    height: 2px;
    background: linear-gradient(90deg, #e879c6 0%, #78dbff 100%);
  }
  
  .followers-grid,
  .following-grid {
    grid-template-columns: 1fr;
  }
  
  .tab-content {
    padding: 1rem;
  }
}
</style>
