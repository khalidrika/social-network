<template>
  <div class="group-header">
    <!-- Banner Background -->
    <div class="group-banner" :style="bannerStyle">
      <div class="banner-overlay"></div>
      <div class="banner-content">
        <div class="group-info">
          <div class="group-details">
            <h1>{{ group.name }}</h1>
            <p v-if="group.description" class="group-description">{{ group.description }}</p>
            <div class="group-meta">
              <span class="member-count">
                <Users :size="16" />
                {{ group.memberCount || 0 }} membres
              </span>
              <span class="group-privacy" :class="group.privacy">
                <component :is="group.privacy === 'public' ? Globe : Lock" :size="16" />
                {{ group.privacy === 'public' ? 'Public' : 'Privé' }}
              </span>
            </div>
          </div>
        </div>
        
        <GroupActions 
          :group="group"
          :is-member="isMember"
          :is-admin="isAdmin"
          @join="$emit('join')"
          @leave="$emit('leave')"
          @invite="$emit('invite')"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Users, Globe, Lock } from 'lucide-vue-next'
import GroupActions from './GroupActions.vue'

const props = defineProps({
  group: {
    type: Object,
    required: true
  },
  isMember: {
    type: Boolean,
    default: false
  },
  isAdmin: {
    type: Boolean,
    default: false
  }
})

defineEmits(['join', 'leave', 'invite'])

const bannerStyle = computed(() => {
  if (props.group.image) {
    // Convert backend file path to full URL
    const imageUrl = props.group.image.startsWith('http') 
      ? props.group.image 
      : `http://localhost:8081/${props.group.image}`
    
    return {
      '--banner-image': `url("${imageUrl}")`
    }
  }
  
  // Use random banner if no image uploaded
  const randomBanners = [
    '/src/assets/randoms/random1.gif',
    '/src/assets/randoms/random2.gif',
    '/src/assets/randoms/random3.gif',
    '/src/assets/randoms/random4.gif',
    '/src/assets/randoms/random5.gif',
    '/src/assets/randoms/random6.gif',
    '/src/assets/randoms/random7.gif',
    '/src/assets/randoms/random8.gif'
  ]
  
  // Generate consistent random index based on group name/id to avoid changing on each render
  const groupIdentifier = props.group.name || props.group.id || 'default'
  const hash = groupIdentifier.split('').reduce((a, b) => {
    a = ((a << 5) - a) + b.charCodeAt(0)
    return a & a
  }, 0)
  const randomIndex = Math.abs(hash) % randomBanners.length
  
  return {
    '--banner-image': `url("${randomBanners[randomIndex]}")`
  }
})
</script>

<style scoped>
.group-header {
  margin-bottom: 30px;
  width: 100%;
}

.group-banner {
  position: relative;
  width: 100%;
  height: 300px;
  border-radius: 20px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-sizing: border-box;
  background: var(--banner-image);
  background-size: cover;
  background-position: center center;
  background-repeat: no-repeat;
}

.banner-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    to bottom, 
    rgba(0, 0, 0, 0.3) 0%, 
    rgba(0, 0, 0, 0.5) 70%, 
    rgba(0, 0, 0, 0.8) 100%
  );
  backdrop-filter: blur(1px);
}

.banner-content {
  position: relative;
  z-index: 2;
  height: 100%;
  display: flex;
  align-items: end;
  justify-content: space-between;
  padding: 25px 30px;
  color: white;
}

.group-info {
  flex: 1;
}

.group-details h1 {
  font-size: 2.5rem;
  font-weight: 800;
  color: white;
  margin: 0 0 8px 0;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.7);
  line-height: 1.2;
}

.group-description {
  font-size: 1.1rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0 0 16px 0;
  line-height: 1.5;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.7);
  max-width: 600px;
}

.group-meta {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}

.member-count,
.group-privacy {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.95);
  text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.8);
  background: rgba(255, 255, 255, 0.15);
  padding: 10px 18px;
  border-radius: 30px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.group-privacy.public {
  background: rgba(120, 199, 255, 0.15);
  border-color: rgba(120, 199, 255, 0.3);
}

.group-privacy.private {
  background: rgba(232, 121, 198, 0.15);
  border-color: rgba(232, 121, 198, 0.3);
}

.icon {
  font-size: 1.1em;
}

/* Responsive Design */
@media (max-width: 768px) {
  .group-banner {
    height: 250px;
  }
  
  .banner-content {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
    padding: 14px 16px;
  }
  
  .group-details h1 {
    font-size: 2rem;
  }
  
  .group-description {
    font-size: 1rem;
  }
  
  .group-meta {
    gap: 12px;
  }
  
  .member-count,
  .group-privacy {
    font-size: 0.9rem;
    padding: 8px 14px;
  }
}

@media (max-width: 480px) {
  .group-banner {
    height: 200px;
    border-radius: 16px;
  }
  
  .banner-content {
    padding: 12px 14px;
  }
  
  .group-details h1 {
    font-size: 1.75rem;
  }
  
  .group-meta {
    flex-direction: column;
    gap: 8px;
  }
  
  .member-count,
  .group-privacy {
    padding: 6px 12px;
    font-size: 0.85rem;
  }
}
</style>
