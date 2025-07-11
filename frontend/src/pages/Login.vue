<template>
  <div class="login-bg">
    <div class="login-card">
      <div class="login-header">
        <div class="login-logo">
          <img src="/src/assets/logo.png" alt="Social Network Logo" class="logo-image" />
        </div>
        <h1>Bienvenue sur Social Network</h1>
        <p class="login-subtitle">Connectez-vous pour retrouver vos amis</p>
      </div>

      <form class="login-form" @submit.prevent="handleLogin">
        <div class="input-group">
          <label for="email">Email</label>
          <input id="email" v-model="form.email" type="email" placeholder="Votre email" required />
          <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
            <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
            <polyline points="22,6 12,13 2,6"/>
          </svg>
        </div>
        <div class="input-group">
          <label for="password">Mot de passe</label>
          <input id="password" v-model="form.password" type="password" placeholder="Votre mot de passe" required />
          <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
            <circle cx="12" cy="16" r="1"/>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
          </svg>
        </div>
        <button type="submit">Se connecter</button>
      </form>
      
      <div v-if="successMsg" class="success-msg">{{ successMsg }}</div>
      <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
      
      <div class="register-link">
        Pas de compte ? <router-link to="/register">Créer un compte</router-link>
      </div>
    </div>
  </div>
</template>


<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/userStore'

export default {
  name: "Login",
  setup() {
    const router = useRouter()
    const userStore = useUserStore()
    
    const form = ref({
      email: '',
      password: ''
    })
    const successMsg = ref('')
    const errorMsg = ref('')

    // Fonction appelée lors de la soumission du formulaire de connexion
    const handleLogin = async () => {
      successMsg.value = ''
      errorMsg.value = ''
      try {
        // On envoie la requête de connexion au backend avec credentials: 'include'
        const res = await fetch('http://localhost:8081/signin', {
          method: 'POST',
          credentials: 'include', // Important pour les cookies de session
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            email: form.value.email.trim(),
            password: form.value.password
          })
        })
        const data = await res.json()
        console.log('Réponse du serveur:', data)
        // Si la connexion est réussie
        if (res.ok && data.type === 'Success') {
          successMsg.value = 'Connexion réussie !'
          // Récupérer le profil utilisateur après login
          const userRes = await fetch('http://localhost:8081/currentUser', {
            method: 'GET',
            credentials: 'include'
          })
          const userData = await userRes.json()
          if (userRes.ok && userData.users && userData.users.length > 0) {
            userStore.setUser(userData.users[0])
          } else {
            userStore.setUser(null)
          }
          router.push('/')
        } else {
          errorMsg.value = data.message || 'Erreur de connexion'
        }
      } catch (err) {
        errorMsg.value = 'Erreur réseau ou serveur'
      }
    }

    return {
      form,
      successMsg,
      errorMsg,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-bg {
  height: 100vh;
  background: #0a0a0f;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  overflow: hidden;
}

.login-bg::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 20% 80%, rgba(120, 119, 198, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(255, 119, 198, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 40% 40%, rgba(120, 219, 255, 0.2) 0%, transparent 50%);
  /* animation: loginFloat 15s ease-in-out infinite; */
}

.login-bg::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="dots" width="20" height="20" patternUnits="userSpaceOnUse"><circle cx="10" cy="10" r="1" fill="rgba(255,255,255,0.05)"/></pattern></defs><rect width="100" height="100" fill="url(%23dots)"/></svg>');
  pointer-events: none;
  opacity: 0.7;
}

@keyframes loginFloat {
  0%, 100% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(25px, -15px) scale(1.05); }
  66% { transform: translate(-15px, 20px) scale(0.95); }
}

@keyframes floating {
  0%, 100% { transform: translate(0px, 0px); }
  33% { transform: translate(15px, -15px); }
  66% { transform: translate(-10px, 10px); }
}

.login-card {
  background: rgba(15, 15, 23, 0.8);
  backdrop-filter: blur(24px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 24px;
  box-shadow: 
    0 32px 64px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  padding: 48px 40px;
  width: 100%;
  max-width: 440px;
  max-height: 90vh;
  text-align: center;
  overflow-y: auto;
  position: relative;
  z-index: 10;
  animation: slideInLogin 0.4s ease-out;
}

@keyframes slideInLogin {
  0% { 
    opacity: 0; 
    transform: translateY(50px) scale(0.9); 
  }
  100% { 
    opacity: 1; 
    transform: translateY(0) scale(1); 
  }
}

.login-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(120, 119, 198, 0.1) 0%, rgba(255, 119, 198, 0.05) 100%);
  border-radius: 24px;
  pointer-events: none;
}

.login-header {
  margin-bottom: 30px;
}

.login-logo {
  margin-bottom: 15px;
  animation: logoFloat 6s ease-in-out infinite;
  cursor: pointer;
  transition: all 0.3s ease;
}

.login-logo .logo-image {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  animation: loginLogoGlow 4s ease-in-out infinite;
  transition: all 0.3s ease;
  filter: drop-shadow(0 0 15px rgba(232, 121, 198, 0.4));
}

.login-logo:hover .logo-image {
  transform: scale(1.15);
  filter: drop-shadow(0 0 25px rgba(120, 199, 255, 0.7)) 
          drop-shadow(0 0 40px rgba(232, 121, 198, 0.5)) 
          drop-shadow(0 0 60px rgba(168, 85, 247, 0.3));
}

@keyframes loginLogoGlow {
  0%, 100% { 
    filter: drop-shadow(0 0 20px rgba(232, 121, 198, 0.5)) 
            drop-shadow(0 0 30px rgba(168, 85, 247, 0.3));
  }
  33% { 
    filter: drop-shadow(0 0 30px rgba(120, 199, 255, 0.6)) 
            drop-shadow(0 0 40px rgba(232, 121, 198, 0.4));
  }
  66% { 
    filter: drop-shadow(0 0 25px rgba(168, 85, 247, 0.5)) 
            drop-shadow(0 0 35px rgba(120, 199, 255, 0.4));
  }
}

@keyframes logoFloat {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-5px); }
}

h1 {
  background: linear-gradient(135deg, #ffffff 0%, #e879c6 50%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-size: 1.8rem;
  font-weight: 900;
  margin: 0 0 8px 0;
  letter-spacing: -0.02em;
  position: relative;
  z-index: 1;
}

.login-subtitle {
  color: rgba(255, 255, 255, 0.7);
  font-size: 1rem;
  margin: 0;
  font-weight: 400;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
  margin-bottom: 25px;
}

.input-group {
  position: relative;
  text-align: left;
  animation: slideInElement 0.3s ease-out forwards;
  opacity: 0;
}

.input-group:nth-child(1) { animation-delay: 0.2s; }
.input-group:nth-child(2) { animation-delay: 0.3s; }

@keyframes slideInElement {
  0% { 
    opacity: 0; 
    transform: translateX(-30px); 
  }
  100% { 
    opacity: 1; 
    transform: translateX(0); 
  }
}

.input-group label {
  display: block;
  margin-bottom: 8px;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 600;
  font-size: 0.95rem;
}

.input-group input {
  width: 100%;
  padding: 16px 14px 16px 48px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  font-size: 0.95rem;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  color: #ffffff;
  box-sizing: border-box;
}

.input-group input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.input-group input:focus {
  outline: none;
  border-color: #e879c6;
  background: rgba(255, 255, 255, 0.1);
  box-shadow: 
    0 0 0 3px rgba(232, 121, 198, 0.2),
    0 8px 32px rgba(232, 121, 198, 0.3);
  transform: translateY(-2px);
}

.input-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: rgba(255, 255, 255, 0.4);
  transition: color 0.3s ease;
  margin-top: 10px;
}

.input-group input:focus + .input-icon {
  color: #e879c6;
}

button {
  width: 100%;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  color: #000000;
  border: none;
  padding: 16px;
  border-radius: 16px;
  font-size: 1rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.4s ease;
  box-shadow: 
    0 8px 32px rgba(232, 121, 198, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
  position: relative;
  overflow: hidden;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  animation: slideInButton 0.3s ease-out 0.3s forwards;
  opacity: 0;
}

@keyframes slideInButton {
  0% { 
    opacity: 0; 
    transform: translateY(30px) scale(0.9); 
  }
  100% { 
    opacity: 1; 
    transform: translateY(0) scale(1); 
  }
}

button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.3), transparent);
  transition: left 0.6s;
}

button:hover {
  transform: translateY(-3px);
  box-shadow: 
    0 16px 48px rgba(232, 121, 198, 0.6),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

button:hover::before {
  left: 100%;
}

button:active {
  transform: translateY(-1px);
}

.success-msg {
  background: #d4edda;
  color: #155724;
  padding: 12px 16px;
  border-radius: 8px;
  margin: 15px 0;
  text-align: center;
  border: 1px solid #c3e6cb;
}

.error-msg {
  background: #f8d7da;
  color: #721c24;
  padding: 12px 16px;
  border-radius: 8px;
  margin: 15px 0;
  text-align: center;
  border: 1px solid #f5c6cb;
}

.register-link {
  text-align: center;
  color: #6b7280;
  font-size: 1rem;
}

.register-link a {
  color: #4f46e5;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s ease;
  position: relative;
}

.register-link a::after {
  content: '';
  position: absolute;
  width: 0;
  height: 2px;
  bottom: -2px;
  left: 0;
  background: linear-gradient(135deg, #4f46e5 0%, #5b21b6 100%);
  transition: width 0.3s ease;
}

.register-link a:hover::after {
  width: 100%;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .login-card {
    padding: 30px 20px;
    margin: 10px;
  }
  
  h1 {
    font-size: 1.4rem;
  }
}
</style>
