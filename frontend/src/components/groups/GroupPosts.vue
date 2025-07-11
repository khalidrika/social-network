<template>
  <div class="posts-section">
    <!-- Create New Post -->
    <div v-if="canPost" class="create-post-section">
      <div class="create-post-card" :class="{ 'expanded': isCreatePostExpanded }">
        <div class="create-post-header">
          <div class="user-avatar">
            <img v-if="userStore.user?.avatar || userStore.user?.imagePath" 
                 :src="getFullImageUrl(userStore.user.avatar || userStore.user.imagePath)" 
                 alt="Avatar" />
            <span v-else>{{ userStore.user?.nickname?.slice(0, 2).toUpperCase() || 'VU' }}</span>
          </div>
          <textarea
            v-model="newPost.content"
            placeholder="Partagez quelque chose avec le groupe..."
            class="create-post-input"
            @focus="isCreatePostExpanded = true"
            @input="handleInputChange"
            rows="1"
            @keydown="handleKeydown"
          ></textarea>
        </div>
        
        <div v-if="isCreatePostExpanded" class="create-post-actions">
          <div class="action-buttons-left">
            <button class="action-btn" @click="showEmojiPicker = !showEmojiPicker">
              😊 Emoji
            </button>
            <input
              ref="imageInput"
              type="file"
              accept="image/*"
              @change="handleImageSelect"
              style="display: none"
            />
            <label class="action-btn" @click="$refs.imageInput.click()">
              <Image :size="16" />
              Photo
            </label>
          </div>
          <button 
            class="publish-btn" 
            @click="createPost"
            :disabled="(!newPost.content.trim() && !newPost.image) || creating"
            :class="{ 'loading': creating }"
          >
            <div v-if="creating" class="loading-spinner"></div>
            {{ creating ? 'Publication...' : 'Publier' }}
          </button>
        </div>
        
        <!-- Emoji Picker -->
        <div v-if="isCreatePostExpanded && showEmojiPicker" class="emoji-picker">
          <div class="emoji-grid">
            <span v-for="emoji in commonEmojis" :key="emoji" @click="addEmoji(emoji)" class="emoji-item">
              {{ emoji }}
            </span>
          </div>
        </div>
        
        <div v-if="isCreatePostExpanded && newPost.imagePreview" class="image-preview-container">
          <img :src="newPost.imagePreview" class="image-preview" alt="Preview" />
          <button class="remove-image-btn" @click="removeImage">×</button>
        </div>
      </div>
    </div>

    <!-- Posts List -->
    <div v-if="displayPosts.length === 0" class="empty-state">
      <div class="empty-icon">
        <FileText :size="64" />
      </div>
      <h3>Aucun post</h3>
      <p>Soyez le premier à publier dans ce groupe</p>
    </div>
    
    <div v-else class="posts-list">
      <div v-for="post in displayPosts" :key="post.id" class="post-card">
        <div class="post-header">
          <div class="author-avatar">
            <img v-if="post.userAvatar" 
                 :src="getFullImageUrl(post.userAvatar)" 
                 :alt="post.userName" />
            <div v-else class="default-author">
              <span>{{ post.userInitials || 'U' }}</span>
            </div>
          </div>
          <div class="post-info">
            <h4>{{ post.userName || 'Utilisateur inconnu' }}</h4>
            <span class="post-date">{{ formatDate(post.createdAt) }}</span>
          </div>
        </div>
        
        <div class="post-content">
          <p v-if="post.content">{{ post.content }}</p>
          <p v-else-if="!post.hasImage" class="no-content">Contenu non disponible</p>
          <img v-if="post.hasImage && post.imagePath" :src="getImageUrl(post.imagePath)" :alt="post.content || 'Image du post'" class="post-image" />
        </div>
        
        <!-- Comments Section -->
        <div class="comments-section">
          <div class="comments-header">
            <button 
              @click="toggleComments(post.id)"
              class="comments-toggle"
            >
              <MessageCircle :size="16" />
              {{ post.comments?.length || 0 }} commentaire{{ (post.comments?.length || 0) > 1 ? 's' : '' }}
            </button>
          </div>
          
          <div v-if="showComments[post.id]" class="comments-content">
            <!-- Existing Comments -->
            <div v-if="post.comments && post.comments.length > 0" class="comments-list">
              <div v-for="comment in post.comments" :key="comment.id" class="comment">
                <div class="comment-header">
                  <div class="comment-avatar">
                    <img v-if="comment.author && (comment.author.avatar || comment.author.imagePath || comment.author.image)" 
                         :src="getFullImageUrl(comment.author.avatar || comment.author.imagePath || comment.author.image)" 
                         :alt="`${comment.author.firstName} ${comment.author.lastName}`" />
                    <div v-else class="default-comment-author">
                      <Users :size="16" />
                    </div>
                  </div>
                  <div class="comment-author">
                    <span class="author-name">{{ comment.author.firstName }} {{ comment.author.lastName }}</span>
                    <span class="comment-date">{{ formatDate(comment.createdAt) }}</span>
                  </div>
                </div>
                <div class="comment-content">
                  <p>{{ comment.content }}</p>
                  <img v-if="comment.image" :src="getImageUrl(comment.image)" :alt="comment.content" class="comment-image" />
                </div>
              </div>
            </div>
            
            <!-- Add New Comment -->
            <div class="add-comment">
              <form @submit.prevent="addComment(post.id)">
                <div class="comment-form">
                  <div class="current-user-avatar">
                    <img v-if="userStore.user?.avatar || userStore.user?.imagePath" 
                         :src="getFullImageUrl(userStore.user.avatar || userStore.user.imagePath)" 
                         alt="Avatar" />
                    <span v-else>{{ userStore.user?.nickname?.slice(0, 2).toUpperCase() || 'VU' }}</span>
                  </div>
                  <textarea
                    v-model="getCommentRef(post.id).content"
                    placeholder="Ajouter un commentaire..."
                    rows="2"
                  ></textarea>
                  
                  <div class="comment-actions">
                    <input
                      :ref="`commentImage${post.id}`"
                      type="file"
                      accept="image/*"
                      @change="handleCommentImageSelect(post.id, $event)"
                      style="display: none"
                    />
                    <button 
                      type="button" 
                      @click="$refs[`commentImage${post.id}`][0].click()" 
                      class="btn-icon"
                    >
                      <Image :size="16" />
                    </button>
                    <button 
                      type="submit" 
                      :disabled="(!getCommentRef(post.id).content?.trim() && !getCommentRef(post.id).image) || commentLoading[post.id]"
                      class="btn-primary"
                    >
                      <Send :size="16" />
                    </button>
                  </div>
                </div>
                
                <div v-if="getCommentRef(post.id).imagePreview" class="comment-image-preview">
                  <img :src="getCommentRef(post.id).imagePreview" alt="Preview" />
                  <button type="button" @click="removeCommentImage(post.id)" class="remove-image">
                    <X :size="16" />
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useGroupStore } from '../../stores/groupStore'
import { useMainStore } from '../../stores/postsStore'
import { useUserStore } from '../../stores/userStore'
import { Image, Send, X, FileText, Users, MessageCircle } from 'lucide-vue-next'

const mainStore = useMainStore()
const userStore = useUserStore()
const { getFullImageUrl, formatPostsData } = mainStore

// Props
const props = defineProps({
  posts: {
    type: Array,
    default: () => []
  },
  groupId: {
    type: [String, Number],
    required: true
  },
  canPost: {
    type: Boolean,
    default: false
  }
})

// Store
const groupStore = useGroupStore()

// Reactive data
const creating = ref(false)
const showComments = ref({})
const commentLoading = ref({})
const newComment = ref({})
const isCreatePostExpanded = ref(false)
const showEmojiPicker = ref(false)

const newPost = ref({
  content: '',
  image: null,
  imagePreview: null
})

// Common emojis for the picker
const commonEmojis = [
  '😊', '😂', '🥰', '😍', '🤔', '😎', '😢', '😡', '👍', '👎', 
  '❤️', '💙', '💚', '💛', '🔥', '⭐', '🎉', '🚀', '💪', '👏'
]

// Computed
const displayPosts = computed(() => {
  const posts = props.posts.length > 0 ? props.posts : groupStore.currentGroupPosts
  
  // Format the posts using the same function as Home.vue for consistency
  const formattedPosts = formatPostsData(posts, userStore.user)
  
  return formattedPosts
})

const getCommentRef = (postId) => {
  if (!newComment.value[postId]) {
    newComment.value[postId] = { content: '', image: null, imagePreview: null }
  }
  return newComment.value[postId]
}

// Methods
const handleInputChange = () => {
  if (!newPost.value.content.trim()) {
    isCreatePostExpanded.value = false
    showEmojiPicker.value = false
  }
}

const handleKeydown = (event) => {
  // Allow Enter key for new lines
  if (event.key === 'Enter' && !event.shiftKey) {
    // Don't submit on Enter, just add new line
    return
  }
}

const addEmoji = (emoji) => {
  newPost.value.content += emoji
  showEmojiPicker.value = false
}

// Click outside to close emoji picker
const handleClickOutside = (event) => {
  if (showEmojiPicker.value && !event.target.closest('.emoji-picker') && !event.target.closest('.action-btn')) {
    showEmojiPicker.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

const createPost = async () => {
  if (!newPost.value.content.trim() && !newPost.value.image) {
    return
  }
  
  creating.value = true
  try {
    const result = await groupStore.createGroupPost(props.groupId, newPost.value.content, newPost.value.image)
    
    if (result.success) {
      // Reset form
      newPost.value = { content: '', image: null, imagePreview: null }
      isCreatePostExpanded.value = false
      showEmojiPicker.value = false
    } else {
      console.error('Post creation failed:', result.message)
      alert(result.message || 'Erreur lors de la création du post')
    }
  } catch (error) {
    console.error('Error creating post:', error)
    alert('Erreur lors de la création du post')
  } finally {
    creating.value = false
  }
}

const handleImageSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    newPost.value.image = file
    newPost.value.imagePreview = URL.createObjectURL(file)
  }
}

const removeImage = () => {
  if (newPost.value.imagePreview) {
    URL.revokeObjectURL(newPost.value.imagePreview)
  }
  newPost.value.image = null
  newPost.value.imagePreview = null
}

const toggleComments = (postId) => {
  showComments.value[postId] = !showComments.value[postId]
}

const handleCommentImageSelect = (postId, event) => {
  const file = event.target.files[0]
  if (file) {
    const commentRef = getCommentRef(postId)
    commentRef.image = file
    commentRef.imagePreview = URL.createObjectURL(file)
  }
}

const removeCommentImage = (postId) => {
  const commentRef = getCommentRef(postId)
  if (commentRef.imagePreview) {
    URL.revokeObjectURL(commentRef.imagePreview)
  }
  commentRef.image = null
  commentRef.imagePreview = null
}

const addComment = async (postId) => {
  commentLoading.value[postId] = true
  try {
    const commentRef = getCommentRef(postId)
    const result = await groupStore.createGroupPostComment(postId, commentRef.content, commentRef.image)
    
    if (result.success) {
      // Reset comment form
      newComment.value[postId] = { content: '', image: null, imagePreview: null }
    } else {
      console.error('Comment creation failed:', result.message)
      alert(result.message || 'Erreur lors de la création du commentaire')
    }
  } catch (error) {
    console.error('Error creating comment:', error)
    alert('Erreur lors de la création du commentaire')
  } finally {
    commentLoading.value[postId] = false
  }
}

const getImageUrl = (imagePath) => {
  if (!imagePath) return ''
  if (imagePath.startsWith('http')) return imagePath
  return `http://localhost:8081/${imagePath}`
}

const formatDate = (dateString) => {
  // Handle completely missing date - provide a friendly fallback
  if (!dateString || dateString === "" || dateString === null || dateString === undefined) {
    return 'Récemment'  // More user-friendly than "Date non disponible"
  }
  
  try {
    let date;
    
    // Handle Unix timestamp (both string and number)
    if (typeof dateString === 'string' && /^\d+$/.test(dateString)) {
      const timestamp = parseInt(dateString)
      // Try both seconds and milliseconds
      date = timestamp > 10000000000 ? new Date(timestamp) : new Date(timestamp * 1000)
    } else if (typeof dateString === 'number') {
      // Check if it's Unix timestamp (seconds) or JavaScript timestamp (milliseconds)
      date = dateString < 10000000000 ? new Date(dateString * 1000) : new Date(dateString)
    } else if (typeof dateString === 'string') {
      // Try parsing as ISO string
      date = new Date(dateString)
    } else {
      date = new Date(dateString)
    }
    
    if (isNaN(date.getTime())) {
      return 'Récemment'  // Fallback for invalid dates
    }
    
    const now = new Date()
    const diffMs = now - date
    const diffMins = Math.floor(diffMs / 60000)
    const diffHours = Math.floor(diffMins / 60)
    const diffDays = Math.floor(diffHours / 24)
    
    if (diffMins < 1) return 'À l\'instant'
    if (diffMins < 60) return `Il y a ${diffMins} min`
    if (diffHours < 24) return `Il y a ${diffHours}h`
    if (diffDays < 7) return `Il y a ${diffDays} jour${diffDays > 1 ? 's' : ''}`
    
    return date.toLocaleDateString('fr-FR', {
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (error) {
    return 'Récemment'  // Fallback for any errors
  }
}
</script>

<style scoped>
.posts-section {
  max-width: 600px;
  margin: 0 auto;
  padding: 0 20px;
}

@media (max-width: 768px) {
  .posts-section {
    padding: 16px;
  }
}

/* Create Post Section */
.create-post-section {
  margin-bottom: 30px;
}

.create-post-card {
  background: rgba(15, 15, 23, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.create-post-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, 
    transparent 0%, 
    rgba(232, 121, 198, 0.5) 25%, 
    rgba(120, 199, 255, 0.5) 75%, 
    transparent 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.create-post-card.expanded {
  border-color: rgba(232, 121, 198, 0.2);
}

.create-post-card.expanded::before {
  opacity: 1;
}

.create-post-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.user-avatar {
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

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.create-post-input {
  flex: 1;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 25px;
  color: #ffffff;
  font-size: 0.95rem;
  transition: all 0.3s ease;
  resize: none;
  min-height: 44px;
  max-height: 120px;
  font-family: inherit;
  line-height: 1.4;
}

.create-post-input:focus {
  outline: none;
  border-color: #e879c6;
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 3px rgba(232, 121, 198, 0.2);
}

.create-post-input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.create-post-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 15px;
}

.action-buttons-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 15px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
  text-decoration: none;
  white-space: nowrap;
}

.action-btn:hover {
  background: rgba(232, 121, 198, 0.1);
  border-color: rgba(232, 121, 198, 0.3);
  color: #e879c6;
  transform: translateY(-1px);
}

.publish-btn {
  padding: 10px 20px;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  border: none;
  border-radius: 12px;
  color: #000000;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 120px;
  justify-content: center;
}

.publish-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(232, 121, 198, 0.4);
}

.publish-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.publish-btn.loading {
  pointer-events: none;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Emoji Picker */
.emoji-picker {
  position: relative;
  margin-top: 10px;
  background: rgba(15, 15, 23, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 15px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  z-index: 10;
}

.emoji-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(32px, 1fr));
  gap: 8px;
  max-height: 160px;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: rgba(232, 121, 198, 0.5) transparent;
}

.emoji-grid::-webkit-scrollbar {
  width: 6px;
}

.emoji-grid::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}

.emoji-grid::-webkit-scrollbar-thumb {
  background: rgba(232, 121, 198, 0.5);
  border-radius: 3px;
}

.emoji-grid::-webkit-scrollbar-thumb:hover {
  background: rgba(232, 121, 198, 0.7);
}

.emoji-item {
  font-size: 1.4rem;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  aspect-ratio: 1;
}

.emoji-item:hover {
  background: rgba(232, 121, 198, 0.3);
  transform: scale(1.2);
}

.image-preview-container {
  position: relative;
  margin-top: 15px;
  max-width: 200px;
}

.image-preview {
  width: 100%;
  border-radius: 12px;
  max-height: 150px;
  object-fit: cover;
}

.remove-image-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 24px;
  height: 24px;
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

/* Empty State */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: rgba(255, 255, 255, 0.6);
  background: rgba(15, 15, 23, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  margin-top: 30px;
}

.empty-icon {
  margin-bottom: 20px;
  opacity: 0.6;
  color: rgba(232, 121, 198, 0.5);
}

.empty-state h3 {
  margin: 0 0 10px 0;
  font-size: 1.2rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
}

.empty-state p {
  font-size: 0.9rem;
  margin: 0;
  color: rgba(255, 255, 255, 0.5);
}

/* Posts List */
.posts-list {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.post-card {
  background: rgba(15, 15, 23, 0.8);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.post-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background: linear-gradient(180deg, 
    rgba(232, 121, 198, 0.6) 0%, 
    rgba(120, 199, 255, 0.6) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.post-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.4);
  border-color: rgba(232, 121, 198, 0.2);
}

.post-card:hover::before {
  opacity: 1;
}

.post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 15px;
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  border: 2px solid rgba(232, 121, 198, 0.3);
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.author-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-author {
  color: #000;
  font-weight: 600;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

.post-info {
  flex: 1;
}

.post-info h4 {
  color: rgba(255, 255, 255, 0.9);
  font-size: 0.95rem;
  font-weight: 600;
  margin: 0 0 2px 0;
  background: linear-gradient(135deg, #ffffff 0%, #e879c6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.post-date {
  color: rgba(255, 255, 255, 0.5);
  font-size: 0.8rem;
}

.post-content {
  margin-bottom: 15px;
}

.post-content p {
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
  margin: 0 0 15px 0;
}

.no-content {
  color: rgba(255, 255, 255, 0.5);
  font-style: italic;
  margin: 0 0 15px 0;
}

.post-image {
  width: 100%;
  border-radius: 12px;
  margin-top: 10px;
  object-fit: cover;
  max-height: 300px;
}

/* Comments Section */
.comments-section {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding-top: 15px;
  margin-top: 15px;
}

.comments-header {
  margin-bottom: 15px;
}

.comments-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.85rem;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 12px;
  transition: all 0.3s ease;
}

.comments-toggle:hover {
  background: rgba(232, 121, 198, 0.1);
  border-color: rgba(232, 121, 198, 0.3);
  color: #e879c6;
}

.comments-content {
  margin-top: 15px;
}

.comments-list {
  margin-bottom: 15px;
}

.comment {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 12px;
  padding: 12px;
  margin-bottom: 10px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.comment-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  overflow: hidden;
  border: 1px solid rgba(232, 121, 198, 0.3);
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.comment-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-comment-author {
  color: #000;
  font-weight: 600;
  font-size: 10px;
}

.comment-author {
  flex: 1;
}

.author-name {
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 600;
  background: linear-gradient(135deg, #ffffff 0%, #e879c6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.comment-date {
  font-size: 0.7rem;
  color: rgba(255, 255, 255, 0.5);
  margin-left: 8px;
}

.comment-content p {
  font-size: 0.9rem;
  margin: 0 0 8px 0;
  color: rgba(255, 255, 255, 0.8);
}

.comment-image {
  max-width: 100%;
  max-height: 150px;
  border-radius: 8px;
  margin-top: 8px;
}

.add-comment {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 15px;
}

.comment-form {
  display: flex;
  align-items: flex-end;
  gap: 10px;
}

.current-user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 2px solid rgba(232, 121, 198, 0.3);
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #000;
  font-weight: 600;
  font-size: 12px;
  flex-shrink: 0;
  overflow: hidden;
}

.current-user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.comment-form textarea {
  flex: 1;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 8px 12px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 0.9rem;
  resize: none;
  min-height: 36px;
  transition: all 0.3s ease;
}

.comment-form textarea:focus {
  outline: none;
  border-color: #e879c6;
  box-shadow: 0 0 0 2px rgba(232, 121, 198, 0.2);
}

.comment-form textarea::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.comment-actions {
  display: flex;
  gap: 6px;
  align-items: center;
}

.btn-icon {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  padding: 8px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-icon:hover {
  background: rgba(232, 121, 198, 0.1);
  border-color: rgba(232, 121, 198, 0.3);
  color: #e879c6;
}

.comment-actions .btn-primary {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  border: none;
  color: #000000;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 0.85rem;
  min-width: auto;
}

.comment-actions .btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 3px 10px rgba(232, 121, 198, 0.3);
}

.comment-image-preview {
  margin-top: 10px;
  position: relative;
  max-width: 120px;
}

.comment-image-preview img {
  width: 100%;
  border-radius: 8px;
  max-height: 80px;
  object-fit: cover;
}

.remove-image {
  position: absolute;
  top: 4px;
  right: 4px;
  background: rgba(0, 0, 0, 0.6);
  border: none;
  border-radius: 50%;
  color: white;
  cursor: pointer;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

.remove-image:hover {
  background: rgba(255, 0, 0, 0.7);
}
</style>
