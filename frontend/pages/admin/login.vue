<template>
  <div class="admin-login-page">
    <div class="login-box">
      <div class="login-header">
        <span class="logo-icon">🍵</span>
        <h1 class="logo-text">Admin<span class="logo-accent">Panel</span></h1>
        <p class="login-desc">Đăng nhập để quản trị hệ thống</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label>Tên đăng nhập</label>
          <input type="text" v-model="username" placeholder="Nhập tên đăng nhập..." required />
        </div>
        <div class="form-group">
          <label>Mật khẩu</label>
          <div class="password-input-wrapper">
            <input :type="showPassword ? 'text' : 'password'" v-model="password" placeholder="Nhập mật khẩu..." required />
            <button type="button" class="btn-toggle-password" @click="showPassword = !showPassword">
              {{ showPassword ? 'Ẩn' : 'Hiện' }}
            </button>
          </div>
        </div>
        
        <p v-if="error" class="error-msg">{{ error }}</p>

        <button type="submit" class="submit-btn" :disabled="loading">
          <span v-if="!loading">Đăng nhập</span>
          <span v-else class="loader-spinner-small"></span>
        </button>
      </form>

      <div class="login-footer">
        <NuxtLink to="/" class="back-site">← Về trang chủ Matcha Comic</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

definePageMeta({
  layout: 'admin',
  title: 'Đăng nhập Admin'
})

useHead({ title: 'Đăng nhập Admin - Matcha Comic' })

const router = useRouter()
const authStore = useAuthStore()
const { post } = useApi()

const username = ref('superadmin')
const password = ref('password')
const showPassword = ref(false)
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await post<{ token: string, user: any }>('/auth/login', {
      username: username.value,
      password: password.value
    })
    
    if (response.token) {
      authStore.setAuth(response.token, response.user)
      router.push('/admin')
    }
  } catch (err: any) {
    error.value = 'Sai tên đăng nhập hoặc mật khẩu. Hãy chắc chắn bạn đã chạy Backend và Seed dữ liệu.'
    console.error('Login error:', err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.admin-login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: radial-gradient(circle at center, #141720 0%, #0d0f14 100%);
  position: relative;
  overflow: hidden;
}

.admin-login-page::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle at center, rgba(156,167,100,0.05) 0%, transparent 40%);
  pointer-events: none;
}

.login-box {
  width: 100%;
  max-width: 420px;
  background: rgba(20, 23, 32, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255,255,255,0.05);
  border-top: 3px solid #9CA764;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 25px 50px -12px rgba(0,0,0,0.8), 0 0 40px rgba(156,167,100,0.1);
  position: relative;
  z-index: 10;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.logo-icon {
  font-size: 2.5rem;
  display: block;
  margin-bottom: 10px;
}

.logo-text {
  font-family: 'Montserrat', sans-serif;
  font-size: 1.8rem;
  font-weight: 800;
  color: #F1E8C7;
  margin: 0 0 8px 0;
}
.logo-accent { color: #9CA764; }

.login-desc {
  color: #A8A8B3;
  font-size: 0.95rem;
  margin: 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #A8A8B3;
}

.form-group input {
  background: rgba(0,0,0,0.2);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 8px;
  padding: 12px 16px;
  color: #fff;
  font-size: 1rem;
  transition: all 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: #9CA764;
  background: rgba(0,0,0,0.4);
  box-shadow: 0 0 0 3px rgba(156,167,100,0.2);
}

.password-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.password-input-wrapper input {
  width: 100%;
  padding-right: 60px; /* Space for the button */
}

.btn-toggle-password {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  color: #A8A8B3;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  padding: 4px;
}
.btn-toggle-password:hover {
  color: #9CA764;
}

.error-msg {
  color: #ef4444;
  font-size: 0.85rem;
  text-align: center;
  margin: 0;
}

.submit-btn {
  background: linear-gradient(135deg, #9CA764, #6B7445);
  color: #000;
  border: none;
  border-radius: 8px;
  padding: 14px;
  font-weight: 700;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 10px;
}

.submit-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #B8C878, #9CA764);
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(156,167,100,0.3);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.loader-spinner-small {
  width: 20px; height: 20px;
  border: 2px solid rgba(0,0,0,0.2);
  border-top-color: #000;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.login-footer {
  margin-top: 30px;
  text-align: center;
  border-top: 1px solid rgba(255,255,255,0.05);
  padding-top: 20px;
}

.back-site {
  color: #5C5C6B;
  font-size: 0.9rem;
  text-decoration: none;
  transition: color 0.2s;
}
.back-site:hover { color: #A8A8B3; }
</style>
