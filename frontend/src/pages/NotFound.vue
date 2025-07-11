<template>
  <div class="not-found-page">
    <div class="content-card">
      <div class="error-code">404</div>
      <h1>Page introuvable</h1>
      <p>La page que vous recherchez n'existe pas.</p>
      
      <div class="action-buttons">
        <button @click="goHome" class="btn btn-primary">
          <HomeIcon :size="16" />
          Accueil
        </button>
        
        <button @click="goBack" class="btn btn-secondary">
          <ArrowLeftIcon :size="16" />
          Retour
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import {
  HomeIcon,
  ArrowLeftIcon
} from '@heroicons/vue/24/outline'

export default {
  name: 'NotFound',
  components: {
    HomeIcon,
    ArrowLeftIcon
  },
  setup() {
    const router = useRouter()

    const goHome = () => {
      router.push('/')
    }

    const goBack = () => {
      if (window.history.length > 1) {
        router.go(-1)
      } else {
        router.push('/')
      }
    }

    return {
      goHome,
      goBack
    }
  }
}
</script>

<style scoped>
.not-found-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  padding: 20px;
}

.content-card {
  text-align: center;
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 40px;
  max-width: 400px;
  width: 100%;
}

.error-code {
  font-size: 4rem;
  font-weight: 700;
  margin: 0 0 20px 0;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.content-card h1 {
  color: #ffffff;
  font-size: 1.5rem;
  margin: 0 0 10px 0;
  font-weight: 600;
}

.content-card p {
  color: rgba(255, 255, 255, 0.7);
  font-size: 1rem;
  margin: 0 0 30px 0;
  line-height: 1.5;
}

.action-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
}

.btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
  background: none;
}

.btn-primary {
  background: rgba(232, 121, 198, 0.1);
  color: #ffffff;
  border-color: rgba(232, 121, 198, 0.3);
}

.btn-primary:hover {
  background: rgba(232, 121, 198, 0.2);
  border-color: rgba(232, 121, 198, 0.4);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.8);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
}

@media (max-width: 768px) {
  .content-card {
    padding: 30px 20px;
  }
  
  .error-code {
    font-size: 3rem;
  }
  
  .action-buttons {
    flex-direction: column;
  }
}
</style>
