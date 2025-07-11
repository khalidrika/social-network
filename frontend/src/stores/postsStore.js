import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useMainStore = defineStore('main', () => {
  // Global state
  const count = ref(0);
  const posts = ref([]);
  
  // Comment UI state
  const activeCommentPostId = ref(null);
  const newCommentText = ref('');
  const commentSelectedImage = ref(null);
  const commentImagePreview = ref('');
  
  // Post creation UI state
  const newPostContent = ref('');
  const showCreatePost = ref(false);
  const selectedImage = ref(null);
  const imagePreview = ref('');
  const selectedPrivacy = ref('public');
  const selectedFollowers = ref([]);
  
  // Basic actions
  const increment = () => {
    count.value++;
  };

  const fetchPosts = async () => {
    try {
      const res = await fetch('http://localhost:8081/allPosts', { credentials: 'include' });
      if (res.ok) {
        const data = await res.json();
        posts.value = Array.isArray(data.posts) ? data.posts : [];
      } else {
        console.error('Failed to fetch posts:', await res.text());
      }
    } catch (error) {
      console.error('Error fetching posts:', error);
    }
  };

  const fetchUserPosts = async (userId) => {
    try {
      const res = await fetch(`http://localhost:8081/userPosts?id=${userId}`, { credentials: 'include' });
      if (res.ok) {
        const data = await res.json();
        return Array.isArray(data.posts) ? data.posts : [];
      } else {
        console.error('Failed to fetch user posts:', await res.text());
        return [];
      }
    } catch (error) {
      console.error('Error fetching user posts:', error);
      return [];
    }
  };

  const submitComment = async (postId, comment) => {
    let commentData = new FormData();
    commentData.set('postid', postId);
    commentData.set('body', comment.body);
    if (comment.image) {
      commentData.append('image', comment.image);
    }
    await fetch('http://localhost:8081/newComment', {
      method: 'POST',
      credentials: 'include',
      body: commentData
    });
  };

  // Image handling methods
  const handleCommentImageChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      commentSelectedImage.value = file;
      commentImagePreview.value = URL.createObjectURL(file);
    }
  };

  const removeCommentImage = () => {
    commentSelectedImage.value = null;
    commentImagePreview.value = '';
  };

  const handleImageChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      selectedImage.value = file;
      imagePreview.value = URL.createObjectURL(file);
    }
  };

  const removeImage = () => {
    selectedImage.value = null;
    imagePreview.value = '';
  };

  // Privacy methods
  const togglePrivacy = () => {
    const privacyOptions = ['public', 'almost_private', 'private'];
    const currentIndex = privacyOptions.indexOf(selectedPrivacy.value);
    const nextIndex = (currentIndex + 1) % privacyOptions.length;
    selectedPrivacy.value = privacyOptions[nextIndex];
  };

  const getPrivacyIconComponent = () => {
    switch (selectedPrivacy.value) {
      case 'public': return 'Globe';
      case 'almost_private': return 'Users';
      case 'private': return 'Lock';
      default: return 'Globe';
    }
  };

  const getPrivacyText = () => {
    switch (selectedPrivacy.value) {
      case 'public': return 'Public';
      case 'almost_private': return 'Presque privé';
      case 'private': return 'Privé';
      default: return 'Public';
    }
  };

  const getPrivacyTooltip = () => {
    switch (selectedPrivacy.value) {
      case 'public': return 'Cliquez pour changer en Presque privé';
      case 'almost_private': return 'Cliquez pour changer en Privé';
      case 'private': return 'Cliquez pour changer en Public';
      default: return 'Cliquez pour changer la confidentialité';
    }
  };

  // Comment actions
  const toggleComments = (postId) => {
    activeCommentPostId.value = activeCommentPostId.value === postId ? null : postId;
    newCommentText.value = '';
    removeCommentImage();
  };

  const submitCommentAction = async (postId) => {
    if (!newCommentText.value && !commentSelectedImage.value) return;
    
    await submitComment(postId, {
      body: newCommentText.value,
      image: commentSelectedImage.value
    });
    
    newCommentText.value = '';
    removeCommentImage();
    await fetchPosts(); // Refresh posts
  };

  // Post creation
  const publishPost = async () => {
    if (!newPostContent.value && !selectedImage.value) return;
    
    const formData = new FormData();
    formData.set('body', newPostContent.value);
    formData.set('privacy', selectedPrivacy.value.toUpperCase().replace('_', '_'));
    
    if (selectedPrivacy.value === 'private') {
      formData.set('checkedfollowers', selectedFollowers.value.join(','));
    }
    
    if (selectedImage.value) {
      formData.append('image', selectedImage.value);
    }
    
    try {
      const response = await fetch('http://localhost:8081/newPost', {
        method: 'POST',
        credentials: 'include',
        body: formData
      });
      
      const result = await response.json();
      
      if (!response.ok) {
        console.error("Error creating post:", result);
        alert("Erreur lors de la création du post: " + (result.message || "Erreur inconnue"));
        return false;
      }
      
      // Reset form
      newPostContent.value = '';
      removeImage();
      showCreatePost.value = false;
      selectedPrivacy.value = 'public';
      selectedFollowers.value = [];
      
      await fetchPosts();
      return true;
    } catch (error) {
      console.error("Post creation failed:", error);
      alert("Erreur de connexion lors de la création du post");
      return false;
    }
  };

  // Utility methods
  const getImageUrl = (path) => {
    if (!path) return '';
    if (path.startsWith('http://') || path.startsWith('https://')) {
      return path;
    }
    const cleanPath = path.startsWith('/') ? path.substring(1) : path;
    return `http://localhost:8081/${cleanPath}`;
  };

  // Global user avatar/image URL handler - works for all user images regardless of privacy
  const getFullImageUrl = (path) => {
    if (!path) return '';
    if (path.startsWith('http')) return path;
    return `http://localhost:8081/${path.replace(/^\/+/, '')}`;
  };

  const formatDate = (dateString) => {
    if (!dateString || dateString === "" || dateString === null || dateString === undefined) {
      return 'Récemment';
    }
    
    try {
      let date;
      
      if (typeof dateString === 'string' && /^\d+$/.test(dateString)) {
        const timestamp = parseInt(dateString);
        date = timestamp > 10000000000 ? new Date(timestamp) : new Date(timestamp * 1000);
      } else if (typeof dateString === 'number') {
        date = dateString < 10000000000 ? new Date(dateString * 1000) : new Date(dateString);
      } else {
        date = new Date(dateString);
      }
      
      if (isNaN(date.getTime())) {
        return 'Récemment';
      }
      
      const now = new Date();
      const diffMs = now - date;
      const diffMins = Math.floor(diffMs / 60000);
      const diffHours = Math.floor(diffMins / 60);
      const diffDays = Math.floor(diffHours / 24);
      
      if (diffMins < 1) return 'À l\'instant';
      if (diffMins < 60) return `Il y a ${diffMins} min`;
      if (diffHours < 24) return `Il y a ${diffHours}h`;
      if (diffDays < 7) return `Il y a ${diffDays} jour${diffDays > 1 ? 's' : ''}`;
      
      return date.toLocaleDateString('fr-FR', {
        day: 'numeric',
        month: 'short',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch (error) {
      return 'Récemment';
    }
  };

  const handleImageError = (event) => {
    console.error('Image failed to load:', event.target.src);
    event.target.style.display = 'none';
    const placeholder = document.createElement('div');
    placeholder.className = 'post-image-placeholder';
    placeholder.innerHTML = `
      <svg width="60" height="60" viewBox="0 0 24 24" fill="currentColor">
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
        <circle cx="8.5" cy="8.5" r="1.5"/>
        <polyline points="21,15 16,10 5,21"/>
      </svg>
      <span>Image non disponible</span>
    `;
    event.target.parentElement.appendChild(placeholder);
  };

  const getCommentImage = (comment) => {
    return comment.imagePath || comment.image || (comment.images && comment.images.length > 0 ? comment.images[0] : null);
  };

  const isValidDate = (date) => {
    return !isNaN(new Date(date).getTime());
  };

  const getPostPrivacyIconComponent = (visibility) => {
    switch (visibility) {
      case 'PUBLIC': return 'Globe';
      case 'ALMOST_PRIVATE': return 'Users';
      case 'PRIVATE': return 'Lock';
      default: return 'Globe';
    }
  };

  const formatPostsData = (postsArray, currentUser = null) => {
    const arr = Array.isArray(postsArray) ? postsArray : [];
    
    return arr
      .filter(post => post)
      .map(post => {
        const imagePath = post.imagePath || post.image || (post.images && post.images.length > 0 ? post.images[0] : null);
        
        // Get user avatar from author data - backend provides this as post.author.avatar (JSON tag from User model)
        const userAvatar = post.author?.avatar || post.author?.imagePath || post.author?.image || null;
        
        return {
          id: post.id,
          userName: post.author?.nickname || post.author?.firstName || currentUser?.nickname || 'Utilisateur',
          userInitials: post.author?.nickname ? post.author.nickname.slice(0,2).toUpperCase() : 
                       (post.author?.firstName ? post.author.firstName.slice(0,2).toUpperCase() : 
                       (currentUser?.nickname ? currentUser.nickname.slice(0,2).toUpperCase() : '??')),
          userAvatar: userAvatar,
          timeAgo: post.createdAt ? formatDate(post.createdAt) : '',
          content: post.content,
          hasImage: !!imagePath,
          imagePath: imagePath,
          likes: post.likes || 0,
          comments: post.comments || [],
          commentsCount: post.comments ? post.comments.length : 0,
          visibility: post.visibility || 'PUBLIC',
        };
      });
  };

  // Computed
  const formattedPosts = computed(() => formatPostsData(posts.value));

  return {
    // State
    count,
    posts,
    
    // Comment UI state
    activeCommentPostId,
    newCommentText,
    commentSelectedImage,
    commentImagePreview,
    
    // Post creation UI state
    newPostContent,
    showCreatePost,
    selectedImage,
    imagePreview,
    selectedPrivacy,
    selectedFollowers,
    
    // Computed
    formattedPosts,
    
    // Basic actions
    increment,
    fetchPosts,
    fetchUserPosts,
    submitComment,
    
    // Image handling
    handleCommentImageChange,
    removeCommentImage,
    handleImageChange,
    removeImage,
    
    // Privacy
    togglePrivacy,
    getPrivacyIconComponent,
    getPrivacyText,
    getPrivacyTooltip,
    
    // Comments
    toggleComments,
    submitCommentAction,
    
    // Post creation
    publishPost,
    
    // Utilities
    getImageUrl,
    getFullImageUrl,
    formatDate,
    handleImageError,
    getCommentImage,
    isValidDate,
    getPostPrivacyIconComponent,
    formatPostsData
  };
});
