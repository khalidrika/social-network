<template>
  <div class="register-bg">
    <div class="register-card">
      <div class="register-header">
        <div class="register-logo">
          <img src="/src/assets/logo.png" alt="Social Network Logo" class="logo-image" />
        </div>
        <h1>Créer un compte</h1>
        <p class="register-subtitle">Rejoignez notre communauté créative</p>
      </div>

      <form class="register-form" @submit.prevent="handleRegister">
        <div class="form-columns">
          <!-- Left Column -->
          <div class="form-column">
            <div class="input-group">
              <label for="firstname">Prénom</label>
              <input id="firstname" v-model="form.firstname" type="text" placeholder="Votre prénom" required />
              <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
            </div>
            <div class="input-group">
              <label for="lastname">Nom</label>
              <input id="lastname" v-model="form.lastname" type="text" placeholder="Votre nom" required />
              <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
            </div>
            <div class="input-group">
              <label for="email">Email</label>
              <input id="email" v-model="form.email" type="email" placeholder="Votre email" required />
              <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                <polyline points="22,6 12,13 2,6"/>
              </svg>
            </div>
            <div class="input-group">
              <label for="dob">Date de naissance</label>
              <input id="dob" v-model="form.dateofbirth" type="date" required />
              <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                <line x1="16" y1="2" x2="16" y2="6"/>
                <line x1="8" y1="2" x2="8" y2="6"/>
                <line x1="3" y1="10" x2="21" y2="10"/>
              </svg>
            </div>
          </div>

          <!-- Right Column -->
          <div class="form-column">
            <div class="input-group">
              <label for="password">Mot de passe</label>
              <input id="password" v-model="form.password" type="password" placeholder="Votre mot de passe" required />
              <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <circle cx="12" cy="16" r="1"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
            </div>
            <div class="input-group">
              <label for="nickname">Pseudo (optionnel)</label>
              <input id="nickname" v-model="form.nickname" type="text" placeholder="Votre pseudo" />
              <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
                <rect x="8" y="2" width="8" height="4" rx="1" ry="1"/>
              </svg>
            </div>
            <div class="input-group file-group">
              <label for="avatar">Avatar (optionnel)</label>
              <input id="avatar" type="file" @change="onFileChange" class="file-input" />
              <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M14.828 14.828a4 4 0 0 1-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0z"/>
              </svg>
            </div>
            <div class="input-group">
              <label for="about">À propos de moi (optionnel)</label>
              <textarea id="about" v-model="form.aboutme" placeholder="Quelques mots sur vous..." rows="2"></textarea>
              <svg class="input-icon textarea-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
                <polyline points="10,9 9,9 8,9"/>
              </svg>
            </div>
          </div>
        </div>

        <button type="submit" class="submit-btn">S'inscrire</button>
      </form>
      
      <div v-if="successMsg" class="success-msg">{{ successMsg }}</div>
      <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
      
      <div class="login-link">
        Déjà un compte ? <router-link to="/login">Se connecter</router-link>
      </div>
    </div>
  </div>
</template>


<script>
export default {
  name: "Register",
  data() {
    return {
      form: {
        email: '',
        password: '',
        firstname: '',
        lastname: '',
        dateofbirth: '',
        nickname: '',
        aboutme: '',
        avatar: null
      },
      successMsg: '',
      errorMsg: ''
    }
  },
  methods: {
    check(txt) {
alert(txt)
    },
    onFileChange(e) {
      this.form.avatar = e.target.files[0];
    },
    async handleRegister() {
      this.successMsg = '';
      this.errorMsg = '';
      try {
        const formData = new FormData();
        for (const key in this.form) {
          if (this.form[key]) {
            formData.append(key, this.form[key]);
          }
        }
        const res = await fetch('http://localhost:8081/register', {
          method: 'POST',
          body: formData
        });
        if (res.ok) {
          this.successMsg = 'Inscription réussie ! Vous pouvez vous connecter.';
          this.form = { email: '', password: '', firstname: '', lastname: '', dateofbirth: '', nickname: '', aboutme: '', avatar: null };
          this.$router.push('/login'); // Redirection vers la page de connexion
        } else {
          const data = await res.json();
          console.log(data);
          
          this.errorMsg = data.message || "Erreur lors de l'inscription.";
        }
      } catch (err) {
        this.errorMsg = "Erreur réseau ou serveur.";
      }
    }
  }
}
</script>

<style scoped>
.register-bg {
  height: 100vh;
  background: #0a0a0f;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  overflow: hidden;
}

.register-bg::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 25% 25%, rgba(232, 121, 198, 0.4) 0%, transparent 50%),
    radial-gradient(circle at 75% 75%, rgba(120, 199, 255, 0.4) 0%, transparent 50%),
    radial-gradient(circle at 50% 50%, rgba(147, 51, 234, 0.3) 0%, transparent 50%);
  /* animation: registerPulse 18s ease-in-out infinite; */
}

.register-bg::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="dots" width="25" height="25" patternUnits="userSpaceOnUse"><circle cx="12.5" cy="12.5" r="1.2" fill="rgba(232,121,198,0.06)"/><circle cx="6" cy="19" r="0.8" fill="rgba(120,199,255,0.04)"/></pattern></defs><rect width="100" height="100" fill="url(%23dots)"/></svg>');
  pointer-events: none;
  opacity: 0.8;
}

@keyframes registerPulse {
  0%, 100% { transform: translate(0px, 0px) scale(1); opacity: 0.8; }
  25% { transform: translate(15px, -25px) scale(1.1); opacity: 0.9; }
  50% { transform: translate(-20px, 10px) scale(0.9); opacity: 1; }
  75% { transform: translate(30px, 15px) scale(1.05); opacity: 0.85; }
}

@keyframes floating {
  0%, 100% { transform: translate(0px, 0px); }
  25% { transform: translate(20px, -20px); }
  50% { transform: translate(-15px, 15px); }
  75% { transform: translate(25px, 8px); }
}

.register-card {
  background: rgba(15, 15, 23, 0.85);
  backdrop-filter: blur(32px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 28px;
  box-shadow: 
    0 40px 80px rgba(0, 0, 0, 0.5),
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 0 0 1px rgba(232, 121, 198, 0.1);
  padding: 40px 35px;
  width: 100%;
  max-width: 820px;
  max-height: 100vh;
  position: relative;
  overflow-y: auto;
  z-index: 10;
  animation: expandIn s ease-out;
  scale: 0.95;
}

@keyframes expandIn {
  0% { 
    opacity: 0; 
    transform: scale(0.8) rotateY(10deg); 
    filter: blur(5px);
  }
  100% { 
    opacity: 1; 
    transform: scale(1) rotateY(0deg); 
    filter: blur(0px);
  }
}

.register-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(232, 121, 198, 0.08) 0%, rgba(120, 199, 255, 0.08) 50%, rgba(147, 51, 234, 0.08) 100%);
  border-radius: 28px;
  pointer-events: none;
}

.register-header {
  text-align: center;
  margin-bottom: 35px;
  position: relative;
  z-index: 1;
}

.register-logo {
  margin-bottom: 15px;
  animation: logoFloat 6s ease-in-out infinite;
  cursor: pointer;
  transition: all 0.4s ease;
}

.register-logo .logo-image {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  animation: registerLogoGlow 5s ease-in-out infinite;
  transition: all 0.4s ease;
  filter: drop-shadow(0 0 15px rgba(232, 121, 198, 0.4));
}

.register-logo:hover .logo-image {
  transform: scale(1.2);
  filter: drop-shadow(0 0 30px rgba(120, 199, 255, 0.8)) 
          drop-shadow(0 0 50px rgba(232, 121, 198, 0.6)) 
          drop-shadow(0 0 70px rgba(147, 51, 234, 0.4));
}

@keyframes registerLogoGlow {
  0%, 100% { 
    filter: drop-shadow(0 0 20px rgba(232, 121, 198, 0.5)) 
            drop-shadow(0 0 30px rgba(147, 51, 234, 0.3));
  }
  25% { 
    filter: drop-shadow(0 0 30px rgba(120, 199, 255, 0.6)) 
            drop-shadow(0 0 40px rgba(232, 121, 198, 0.4));
  }
  50% { 
    filter: drop-shadow(0 0 35px rgba(147, 51, 234, 0.6)) 
            drop-shadow(0 0 45px rgba(120, 199, 255, 0.4));
  }
  75% { 
    filter: drop-shadow(0 0 25px rgba(232, 121, 198, 0.7)) 
            drop-shadow(0 0 40px rgba(147, 51, 234, 0.3));
  }
}

@keyframes logoFloat {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-5px); }
}

h1 {
  background: linear-gradient(135deg, #ffffff 0%, #e879c6 30%, #78c7ff 70%, #9333ea 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-size: 2rem;
  font-weight: 900;
  margin: 0 0 10px 0;
  letter-spacing: -0.03em;
  position: relative;
  z-index: 1;
}

.register-subtitle {
  color: rgba(255, 255, 255, 0.75);
  font-size: 1.05rem;
  margin: 0;
  font-weight: 400;
  letter-spacing: 0.01em;
}

.register-form {
  width: 100%;
  position: relative;
  z-index: 1;
}

.form-columns {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
  margin-bottom: 30px;
}

.form-column {
  display: flex;
  flex-direction: column;
}

.input-group {
  margin-bottom: 20px;
  position: relative;
  animation: fadeInUp 0.7s ease-out forwards;
  opacity: 0;
}

.form-column:first-child .input-group:nth-child(1) { animation-delay: 0.2s; }
.form-column:first-child .input-group:nth-child(2) { animation-delay: 0.3s; }
.form-column:first-child .input-group:nth-child(3) { animation-delay: 0.4s; }
.form-column:first-child .input-group:nth-child(4) { animation-delay: 0.5s; }
.form-column:last-child .input-group:nth-child(1) { animation-delay: 0.35s; }
.form-column:last-child .input-group:nth-child(2) { animation-delay: 0.45s; }
.form-column:last-child .input-group:nth-child(3) { animation-delay: 0.55s; }
.form-column:last-child .input-group:nth-child(4) { animation-delay: 0.65s; }

.input-group label {
  display: block;
  margin-bottom: 8px;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 600;
  font-size: 0.95rem;
  letter-spacing: 0.01em;
}

.input-group input,
.input-group textarea {
  width: 100%;
  padding: 16px 14px 16px 48px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  font-size: 0.95rem;
  transition: all 0.4s ease;
  background: rgba(255, 255, 255, 0.04);
  backdrop-filter: blur(12px);
  color: #ffffff;
  box-sizing: border-box;
  font-family: inherit;
}

.input-group textarea {
  resize: vertical;
  min-height: 60px;
  padding-top: 12px;
}

.input-group input::placeholder,
.input-group textarea::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.input-group input:focus,
.input-group textarea:focus {
  outline: none;
  border-color: #e879c6;
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 
    0 0 0 3px rgba(232, 121, 198, 0.15),
    0 8px 32px rgba(232, 121, 198, 0.25);
  transform: translateY(-2px);
}

.input-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: rgba(255, 255, 255, 0.35);
  transition: all 0.3s ease;
  margin-top: 14px;
  pointer-events: none;
}

.textarea-icon {
  top: 45px;
  transform: translateY(0);
}

.input-group input:focus + .input-icon,
.input-group textarea:focus + .textarea-icon {
  color: #e879c6;
  transform: translateY(-50%) scale(1.1);
}

.input-group textarea:focus + .textarea-icon {
  transform: translateY(0) scale(1.1);
}

.file-group .file-input {
  padding: 18px 14px 18px 48px !important;
  border: 2px dashed rgba(255, 255, 255, 0.15) !important;
  background: rgba(255, 255, 255, 0.02) !important;
  cursor: pointer;
  text-align: left;
  color: rgba(255, 255, 255, 0.7);
  transition: all 0.4s ease;
}

.file-group .file-input:hover {
  border-color: rgba(232, 121, 198, 0.4) !important;
  background: rgba(232, 121, 198, 0.05) !important;
  transform: translateY(-1px);
}

.file-group .file-input:focus {
  border-color: #e879c6 !important;
  background: rgba(232, 121, 198, 0.1) !important;
}

.submit-btn {
  width: 100%;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 50%, #9333ea 100%);
  color: #000000;
  border: none;
  padding: 18px;
  border-radius: 18px;
  font-size: 1.05rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.4s ease;
  box-shadow: 
    0 12px 40px rgba(232, 121, 198, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
  position: relative;
  overflow: hidden;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-top: 10px;
  animation: slideInButton 0.3s ease-out 0.4s forwards;
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

@keyframes bounceIn {
  0% { 
    opacity: 0; 
    transform: scale(0.3) rotate(15deg); 
  }
  50% { 
    opacity: 1; 
    transform: scale(1.05) rotate(-5deg); 
  }
  70% { 
    transform: scale(0.95) rotate(2deg); 
  }
  100% { 
    opacity: 1; 
    transform: scale(1) rotate(0deg); 
  }
}

.submit-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.3), transparent);
  transition: left 0.6s;
}

.submit-btn:hover {
  transform: translateY(-4px);
  box-shadow: 
    0 20px 60px rgba(232, 121, 198, 0.6),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.submit-btn:hover::before {
  left: 100%;
}

.submit-btn:active {
  transform: translateY(-2px);
}

.success-msg {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2) 0%, rgba(16, 185, 129, 0.1) 100%);
  color: #10b981;
  padding: 16px 20px;
  border-radius: 16px;
  margin: 20px 0;
  text-align: center;
  border: 1px solid rgba(16, 185, 129, 0.3);
  backdrop-filter: blur(10px);
  font-weight: 500;
}

.error-msg {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.2) 0%, rgba(239, 68, 68, 0.1) 100%);
  color: #ef4444;
  padding: 16px 20px;
  border-radius: 16px;
  margin: 20px 0;
  text-align: center;
  border: 1px solid rgba(239, 68, 68, 0.3);
  backdrop-filter: blur(10px);
  font-weight: 500;
}

.login-link {
  text-align: center;
  margin-top: 25px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 1rem;
}

.login-link a {
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s ease;
  position: relative;
}

.login-link a::after {
  content: '';
  position: absolute;
  width: 0;
  height: 2px;
  bottom: -3px;
  left: 0;
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%);
  transition: width 0.3s ease;
}

.login-link a:hover::after {
  width: 100%;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .register-bg {
    padding: 15px;
  }
  
  .register-card {
    padding: 30px 20px;
    max-height: 98vh;
    border-radius: 20px;
  }
  
  .form-columns {
    grid-template-columns: 1fr;
    gap: 0;
  }
  
  .input-group {
    margin-bottom: 16px;
  }
  
  h1 {
    font-size: 1.6rem;
  }
  
  .register-header {
    margin-bottom: 25px;
  }
  
  .submit-btn {
    padding: 16px;
    font-size: 1rem;
  }
}

/* Very small screens */
@media (max-width: 480px) {
  .register-card {
    padding: 25px 15px;
    margin: 5px;
  }
  
  h1 {
    font-size: 1.4rem;
  }
  
  .register-subtitle {
    font-size: 0.95rem;
  }
}

@keyframes fadeInUp {
  0% { 
    opacity: 0; 
    transform: translateY(40px); 
  }
  100% { 
    opacity: 1; 
    transform: translateY(0); 
  }
}
</style>