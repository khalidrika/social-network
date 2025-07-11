<template>
  <div class="modal-overlay" @click="handleOverlayClick">
    <div class="modal-content" @click.stop>
      <!-- Header -->
      <div class="modal-header">
        <h2>Créer un nouveau groupe</h2>
        <button @click="$emit('close')" class="close-btn">
          <X :size="20" />
        </button>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="group-form">
        <!-- Group Image -->
        <div class="form-group">
          <label class="form-label">Image du groupe</label>
          <div class="image-upload">
            <div class="image-preview" @click="triggerFileInput">
              <img 
                v-if="form.imagePreview" 
                :src="form.imagePreview" 
                alt="Aperçu"
              />
              <div v-else class="upload-placeholder">
                <Camera :size="32" />
                <span>Cliquez pour ajouter une image</span>
              </div>
            </div>
            <input 
              ref="fileInput"
              type="file" 
              accept="image/*" 
              @change="handleImageUpload"
              hidden
            />
          </div>
        </div>

        <!-- Group Name -->
        <div class="form-group">
          <label for="name" class="form-label">Nom du groupe *</label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            class="form-input"
            placeholder="Entrez le nom du groupe"
            required
            maxlength="100"
          />
          <span class="char-count">{{ form.name.length }}/100</span>
        </div>

        <!-- Description -->
        <div class="form-group">
          <label for="description" class="form-label">Description</label>
          <textarea
            id="description"
            v-model="form.description"
            class="form-textarea"
            placeholder="Décrivez votre groupe..."
            rows="4"
            maxlength="500"
          ></textarea>
          <span class="char-count">{{ form.description.length }}/500</span>
        </div>

        <!-- Privacy -->
        <div class="form-group">
          <label class="form-label">Confidentialité</label>
          <div class="privacy-options">
            <label class="privacy-option">
              <input 
                v-model="form.privacy" 
                type="radio" 
                value="public"
              />
              <div class="option-content">
                <Globe :size="24" />
                <div>
                  <div class="option-title">Public</div>
                  <div class="option-desc">Tout le monde peut voir et rejoindre</div>
                </div>
              </div>
            </label>
            
            <label class="privacy-option">
              <input 
                v-model="form.privacy" 
                type="radio" 
                value="private"
              />
              <div class="option-content">
                <Lock :size="24" />
                <div>
                  <div class="option-title">Privé</div>
                  <div class="option-desc">Seuls les membres peuvent voir le contenu</div>
                </div>
              </div>
            </label>
          </div>
        </div>

        <!-- Actions -->
        <div class="form-actions">
          <button 
            type="button" 
            @click="$emit('close')" 
            class="btn-secondary"
          >
            Annuler
          </button>
          <button 
            type="submit" 
            class="btn-primary"
            :disabled="!isFormValid || loading"
          >
            <template v-if="loading">
              <div class="spinner"></div>
              Création...
            </template>
            <template v-else>
              <Plus :size="18" />
              Créer le groupe
            </template>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive } from 'vue'
import { X, Camera, Globe, Lock, Plus } from 'lucide-vue-next'

const emit = defineEmits(['close', 'created'])

const fileInput = ref(null)
const loading = ref(false)

const form = reactive({
  name: '',
  description: '',
  privacy: 'public',
  image: null,
  imagePreview: null
})

const isFormValid = computed(() => {
  return form.name.trim().length >= 3
})

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleImageUpload = (event) => {
  const file = event.target.files[0]
  if (file) {
    form.image = file
    const reader = new FileReader()
    reader.onload = (e) => {
      form.imagePreview = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

const handleOverlayClick = (event) => {
  if (event.target === event.currentTarget) {
    emit('close')
  }
}

const handleSubmit = async () => {
  if (!isFormValid.value || loading.value) return

  loading.value = true
  
  try {
    // Create FormData for file upload
    const formData = new FormData()
    formData.append('name', form.name.trim())
    formData.append('description', form.description.trim())
    formData.append('privacy', form.privacy)
    
    // Add image if selected
    if (form.image) {
      formData.append('image', form.image)
    }
    
    // Add invitations as JSON string
    formData.append('invitations', JSON.stringify([]))

    console.log('Creating group with FormData')

    const response = await fetch('http://localhost:8081/newGroup', {
      method: 'POST',
      credentials: 'include',
      body: formData // Don't set Content-Type header, let browser set it with boundary
    })

    console.log('Response status:', response.status)
    
    const responseData = await response.json()
    console.log('Response data:', responseData)

    if (response.ok && responseData.type === 'Success') {
      console.log('Group created successfully, emitting created event')
      emit('created')
    } else {
      throw new Error(responseData.message || 'Erreur lors de la création du groupe')
    }
  } catch (error) {
    console.error('Error creating group:', error)
    alert('Erreur lors de la création du groupe: ' + error.message)
  } finally {
    loading.value = false
  }
}
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
  margin-top: auto;
  background: rgba(15, 15, 23, 0.95);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  backdrop-filter: blur(20px);
  
  /* Hide scrollbar for all browsers */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* Internet Explorer and Edge Legacy */
  scrollbar-color: transparent transparent; /* Firefox fallback */
}

.modal-content::-webkit-scrollbar {
  display: none; /* Chromium browsers (Chrome, Safari, Edge Chromium) */
  width: 0;
  height: 0;
}

.modal-content::-webkit-scrollbar-track {
  background: transparent;
}

.modal-content::-webkit-scrollbar-thumb {
  background: transparent;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 24px 0 24px;
  margin-bottom: 24px;
}

.modal-header h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: white;
  margin: 0;
}

.close-btn {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 8px;
  padding: 8px;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.7);
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.group-form {
  padding: 0 24px 24px 24px;
}

.form-group {
  margin-bottom: 24px;
}

.form-label {
  display: block;
  color: white;
  font-weight: 600;
  margin-bottom: 8px;
}

.image-upload {
  display: flex;
  justify-content: center;
}

.image-preview {
  width: 120px;
  height: 120px;
  border-radius: 16px;
  border: 2px dashed rgba(255, 255, 255, 0.3);
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
}

.image-preview:hover {
  border-color: rgba(232, 121, 198, 0.5);
}

.image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 120px;
  border: 2px dashed rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  transition: all 0.3s ease;
}



.form-input,
.form-textarea {
  width: 100%;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 12px 16px;
  color: white;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: rgba(232, 121, 198, 0.5);
  background: rgba(255, 255, 255, 0.15);
}

.form-input::placeholder,
.form-textarea::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.char-count {
  display: block;
  text-align: right;
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 4px;
}

.privacy-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.privacy-option {
  display: block;
  cursor: pointer;
}

.privacy-option input {
  display: none;
}

.option-content {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.privacy-option input:checked + .option-content {
  background: rgba(232, 121, 198, 0.2);
  border-color: rgba(232, 121, 198, 0.5);
}

.option-content:hover {
  background: rgba(255, 255, 255, 0.1);
}

.option-title {
  font-weight: 600;
  color: white;
  margin-bottom: 4px;
}

.option-desc {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 32px;
}

.btn-primary,
.btn-secondary {
  flex: 1;
  padding: 12px 20px;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn-primary {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(232, 121, 198, 0.3);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
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

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Responsive */
@media (max-width: 768px) {
  .modal-content {
    margin: 0;
    border-radius: 20px 20px 0 0;
    max-height: 100vh;
  }
}
</style>
