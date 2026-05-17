import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null as string | null,
    user: null as any | null,
    initialized: false
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin',
  },
  
  actions: {
    init() {
      if (typeof window !== 'undefined') {
        this.token = localStorage.getItem('token')
        this.user = JSON.parse(localStorage.getItem('user') || 'null')
        this.initialized = true
      }
    },
    setAuth(token: string, user: any) {
      this.token = token
      this.user = user
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(user))
    },
    
    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      localStorage.removeItem('admin_logged_in') // Clean up old mock key
    }
  }
})
