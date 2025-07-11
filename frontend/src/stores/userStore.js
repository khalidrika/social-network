import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    isAuthenticated: false,
    loading: false,
    error: null,
  }),
  getters: {
    isLoggedIn: (state) => state.isAuthenticated,
  },
  actions: {
    setUser(user) {
      this.user = user
      this.isAuthenticated = !!user
    },
    logout() {
      this.user = null
      this.isAuthenticated = false
    },
    async fetchUserData() {
      if (!this.user || !this.user.id) return
      this.loading = true
      this.error = null
      try {
        const res = await fetch(`http://localhost:8081/userData?userId=${this.user.id}`, { credentials: 'include' })
        const data = await res.json()
        if (res.ok && data.users && data.users.length > 0) {
          this.user = data.users[0]
        } else {
          this.error = data.message || "Impossible de charger le profil"
        }
      } catch (e) {
        this.error = "Erreur réseau"
      } finally {
        this.loading = false
      }
    },
    async reloadUserAfterRefresh() {
      this.loading = true
      this.error = null
      try {
        const res = await fetch('http://localhost:8081/currentUser', { credentials: 'include' })
        const data = await res.json()
        if (res.ok && data.users && data.users.length > 0) {
          this.user = data.users[0]
          this.isAuthenticated = true
          await this.fetchUserData()
        } else {
          this.user = null
          this.isAuthenticated = false
        }
      } catch (e) {
        this.user = null
        this.isAuthenticated = false
        this.error = "Erreur réseau"
      } finally {
        this.loading = false
      }
    },
    async updateProfile(profileData) {
      this.loading = true
      this.error = null
      try {
        const res = await fetch('http://localhost:8081/updateProfile', {
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(profileData)
        })
        const data = await res.json()
        if (res.ok && data.type === 'Success') {
          await this.fetchUserData()
        } else {
          this.error = data.message || "Erreur lors de la mise à jour"
        }
      } catch (e) {
        this.error = "Erreur réseau"
      } finally {
        this.loading = false
      }
    },
    async changeStatus(newStatus) {
      this.loading = true
      this.error = null
      try {
        const res = await fetch(`http://localhost:8081/changeStatus?status=${newStatus}`, {
          method: 'POST',
          credentials: 'include'
        })
        const data = await res.json()
        if (res.ok && data.type === 'Success') {
          this.user.status = newStatus
        } else {
          this.error = data.message || "Erreur lors du changement de statut"
        }
      } catch (e) {
        this.error = "Erreur réseau"
      } finally {
        this.loading = false
      }
    }
  },
})
