<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
import { onMounted } from 'vue'
import { useUserStore } from './stores/userStore'

export default {
  name: 'App',
  setup() {
    const userStore = useUserStore()
    
    onMounted(async () => {
      console.log('🚀 App mounted, initializing user authentication...')
      // Check if user is already authenticated when app loads
      try {
        await userStore.reloadUserAfterRefresh()
        console.log('✅ User authentication initialized:', {
          isLoggedIn: userStore.isLoggedIn,
          user: userStore.user
        })
      } catch (error) {
        console.error('❌ Error initializing user authentication:', error)
      }
    })
    
    return {}
  }
}
</script>

<style>
#app {
  min-height: 100vh;
  background: #0a0a0f;
  position: relative;
}

#app::before {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 20% 80%, rgba(120, 119, 198, 0.2) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(255, 119, 198, 0.2) 0%, transparent 50%),
    radial-gradient(circle at 40% 40%, rgba(120, 219, 255, 0.15) 0%, transparent 50%);
  /* animation: globalFloat 20s ease-in-out infinite; */
  pointer-events: none;
  z-index: 0;
}

#app::after {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="dots" width="20" height="20" patternUnits="userSpaceOnUse"><circle cx="10" cy="10" r="1" fill="rgba(255,255,255,0.03)"/></pattern></defs><rect width="100" height="100" fill="url(%23dots)"/></svg>');
  pointer-events: none;
  z-index: 0;
}

@keyframes globalFloat {
  0%, 100% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(20px, -20px) scale(1.05); }
  66% { transform: translate(-15px, 15px) scale(0.95); }
}

body {
  margin: 0;
  font-family: system-ui, Avenir, Helvetica, Arial, sans-serif;
  background: transparent;
  overflow-x: hidden;
}
</style>
