<template>
  <div class="page-header">
    <div class="header-content">
      <h1>Groupes</h1>
      <p>Découvrez et rejoignez des communautés</p>
    </div>
    <div class="header-actions">
      <button @click="$emit('create')" class="create-btn">
        <Plus :size="18" />
        Créer un groupe
      </button>
      <button @click="$emit('refresh')" class="refresh-btn" :disabled="loading">
        <RefreshCw :size="18" :class="{ 'animate-spin': loading }" />
        Actualiser
      </button>
    </div>
  </div>
</template>

<script setup>
import { Plus, RefreshCw } from 'lucide-vue-next'

defineProps({
  loading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['create', 'refresh'])
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 30px;
  background: rgba(15, 15, 23, 0.8);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.header-actions {
  display: flex;
  gap: 12px;
}

.header-content h1 {
  font-size: 2.5rem;
  font-weight: 700;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.header-content p {
  color: rgba(255, 255, 255, 0.7);
  margin: 8px 0 0 0;
  font-size: 1.1rem;
}

.create-btn {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
  border: none;
  border-radius: 15px;
  padding: 12px 24px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

.create-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(232, 121, 198, 0.4);
}

.refresh-btn {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 15px;
  padding: 12px 20px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.refresh-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  transform: translateY(-1px);
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.icon {
  font-size: 1rem;
  font-weight: bold;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }
}
</style>
