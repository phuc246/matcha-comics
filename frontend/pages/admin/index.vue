<template>
  <div class="admin-dashboard">
    <div class="stats-grid">
      <div class="stat-card" v-for="s in stats" :key="s.label">
        <div class="stat-icon" :style="{ backgroundColor: s.color + '15', color: s.color }">{{ s.icon }}</div>
        <div class="stat-info">
          <span class="stat-label">{{ s.label }}</span>
          <span class="stat-value">{{ s.value }}</span>
          <span v-if="s.trend !== undefined" class="stat-trend" :class="s.trend > 0 ? 'up' : 'down'">
            {{ s.trend > 0 ? '↑' : '↓' }} {{ Math.abs(s.trend) }}% <span class="trend-text">vs tuần trước</span>
          </span>
          <div v-if="s.isStorage" class="storage-mini-bar">
            <div class="progress" :style="{ width: s.percent + '%', backgroundColor: s.color }"></div>
          </div>
        </div>
      </div>
    </div>

    <div class="dashboard-charts">
      <div class="chart-box main-chart">
        <div class="chart-header">
          <h3>Thống kê lượt xem (7 ngày qua)</h3>
          <div class="chart-actions">
            <button class="active">Lượt xem</button>
            <button>Bình luận</button>
          </div>
        </div>
        <div class="chart-container">
          <ClientOnly>
            <Line :data="lineData" :options="lineOptions" />
          </ClientOnly>
        </div>
      </div>

      <div class="chart-box side-chart">
        <div class="chart-header">
          <h3>Phân loại truyện</h3>
        </div>
        <div class="chart-container pie-wrap">
          <ClientOnly>
            <Doughnut :data="doughnutData" :options="doughnutOptions" />
          </ClientOnly>
        </div>
        <div class="chart-legend">
          <div class="legend-item"><span class="dot" style="background:#9CA764"></span> Truyện tranh</div>
          <div class="legend-item"><span class="dot" style="background:#F1E8C7"></span> Truyện chữ</div>
        </div>
      </div>
    </div>

    <div class="recent-activity">
      <div class="activity-header">
        <h3>Hoạt động gần đây</h3>
        <NuxtLink to="/admin/comics" class="view-all">Xem tất cả</NuxtLink>
      </div>
      <div class="activity-table-wrap">
        <table class="activity-table">
          <thead>
            <tr>
              <th>Truyện</th>
              <th>Loại</th>
              <th>Hành động</th>
              <th>Thời gian</th>
              <th>Người thực hiện</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="a in activities" :key="a.id">
              <td>
                <div class="story-cell">
                  <img :src="a.cover" alt="" class="mini-cover" />
                  <span>{{ a.title }}</span>
                </div>
              </td>
              <td><span class="type-badge" :class="a.type">{{ a.type === 'comic' ? 'Tranh' : 'Chữ' }}</span></td>
              <td>{{ a.action }}</td>
              <td>{{ a.time }}</td>
              <td><span class="user-tag">{{ a.user }}</span></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Line, Doughnut } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  LinearScale,
  CategoryScale,
  PointElement,
  ArcElement
} from 'chart.js'

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  LinearScale,
  CategoryScale,
  PointElement,
  ArcElement
)

definePageMeta({
  layout: 'admin',
  title: 'Dashboard'
})

useHead({ title: 'Admin Dashboard - Matcha Comic' })

const stats = ref<any[]>([
  { label: 'Tổng lượt xem', value: '0', trend: 0, icon: '👁', color: '#9CA764' },
  { label: 'Tổng số truyện', value: '0', trend: 0, icon: '📚', color: '#F1E8C7' },
  { label: 'Lưu trữ (R2 Free)', value: '0/10GB', isStorage: true, percent: 0, icon: '☁️', color: '#60a5fa' },
  { label: 'Doanh thu', value: '0đ', trend: 0, icon: '💰', color: '#E2D5A2' },
])

const { get } = useApi()
const authStore = useAuthStore()

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

onMounted(async () => {
  try {
    const s = await get<any>('/admin/storage-stats', {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    if (s) {
      const storageIdx = stats.value.findIndex(st => st.isStorage)
      if (storageIdx !== -1) {
        stats.value[storageIdx].value = `${formatSize(s.totalSize)} / 10GB`
        stats.value[storageIdx].percent = Math.min(100, (s.totalSize / s.limitBytes) * 100)
        if (stats.value[storageIdx].percent > 90) stats.value[storageIdx].color = '#ef4444'
      }
    }
    
    // Fetch other real stats if needed...
    const storyData = await get<any[]>('/stories')
    if (storyData) {
      stats.value[1].value = storyData.length.toString()
    }
  } catch (err) {
    console.error(err)
  }
})

const lineData = {
  labels: ['Thứ 2', 'Thứ 3', 'Thứ 4', 'Thứ 5', 'Thứ 6', 'Thứ 7', 'CN'],
  datasets: [
    {
      label: 'Lượt xem',
      backgroundColor: 'rgba(156, 167, 100, 0.2)',
      borderColor: '#9CA764',
      pointBackgroundColor: '#9CA764',
      data: [40, 39, 10, 40, 39, 80, 40],
      tension: 0.4,
      fill: true
    }
  ]
}

const lineOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false }
  },
  scales: {
    y: { grid: { color: 'rgba(255,255,255,0.05)' }, ticks: { color: '#5C5C6B' } },
    x: { grid: { display: false }, ticks: { color: '#5C5C6B' } }
  }
}

const doughnutData = {
  labels: ['Truyện tranh', 'Truyện chữ'],
  datasets: [
    {
      backgroundColor: ['#9CA764', '#F1E8C7'],
      data: [65, 35],
      borderWidth: 0
    }
  ]
}

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false }
  },
  cutout: '70%'
}

const activities = ref([
  { id: 1, title: 'Đấu Phá Thương Khung', cover: 'https://images.unsplash.com/photo-1541963463532-d68292c34b19?w=50&h=70&fit=crop', type: 'comic', action: 'Cập nhật Chương 420', time: '5 phút trước', user: 'Admin' },
  { id: 2, title: 'Võ Luyện Đỉnh Phong', cover: 'https://images.unsplash.com/photo-1512820790803-83ca734da794?w=50&h=70&fit=crop', type: 'comic', action: 'Thêm mới', time: '1 giờ trước', user: 'Admin' },
  { id: 3, title: 'Kiếm Đạo Độc Tôn', cover: 'https://images.unsplash.com/photo-1589998059171-988d887df646?w=50&h=70&fit=crop', type: 'novel', action: 'Cập nhật Chương 12', time: '3 giờ trước', user: 'Editor_01' },
  { id: 4, title: 'Toàn Chức Pháp Sư', cover: 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=50&h=70&fit=crop', type: 'comic', action: 'Sửa thông tin', time: '6 giờ trước', user: 'Admin' },
])
</script>

<style scoped>
.admin-dashboard { display: flex; flex-direction: column; gap: 30px; }

/* Stats Cards */
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 20px; }
.stat-card { background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 16px; padding: 24px; display: flex; align-items: center; gap: 20px; transition: transform 0.3s; }
.stat-card:hover { transform: translateY(-5px); border-color: rgba(156,167,100,0.2); }
.stat-icon { width: 56px; height: 56px; border-radius: 14px; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; }
.stat-info { display: flex; flex-direction: column; gap: 4px; }
.stat-label { color: #A8A8B3; font-size: 0.9rem; font-weight: 500; }
.stat-value { color: #fff; font-size: 1.6rem; font-weight: 800; }
.stat-trend { font-size: 0.8rem; font-weight: 600; display: flex; align-items: center; gap: 4px; }
.stat-trend.up { color: #4ade80; }
.stat-trend.down { color: #ef4444; }
.trend-text { color: #5C5C6B; font-weight: 400; }
.storage-mini-bar { width: 100%; height: 4px; background: rgba(255,255,255,0.05); border-radius: 2px; margin-top: 8px; overflow: hidden; }
.storage-mini-bar .progress { height: 100%; transition: width 0.5s; }

/* Charts */
.dashboard-charts { display: grid; grid-template-columns: 2fr 1fr; gap: 24px; }
@media (max-width: 1024px) { .dashboard-charts { grid-template-columns: 1fr; } }
.chart-box { background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 16px; padding: 24px; }
.chart-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.chart-header h3 { font-size: 1.1rem; font-weight: 700; color: #F1E8C7; margin: 0; }
.chart-actions { display: flex; background: rgba(0,0,0,0.2); padding: 4px; border-radius: 8px; }
.chart-actions button { background: none; border: none; padding: 6px 12px; font-size: 0.8rem; color: #A8A8B3; cursor: pointer; border-radius: 6px; }
.chart-actions button.active { background: #9CA764; color: #000; font-weight: 600; }

.chart-container { height: 300px; position: relative; }
.pie-wrap { height: 240px; }

.chart-legend { display: flex; flex-direction: column; gap: 10px; margin-top: 20px; }
.legend-item { display: flex; align-items: center; gap: 10px; font-size: 0.9rem; color: #A8A8B3; }
.dot { width: 10px; height: 10px; border-radius: 50%; }

/* Activity Table */
.recent-activity { background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 16px; padding: 24px; }
.activity-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.activity-header h3 { font-size: 1.1rem; color: #F1E8C7; margin: 0; }
.view-all { font-size: 0.85rem; color: #9CA764; text-decoration: none; font-weight: 600; }

.activity-table-wrap { overflow-x: auto; }
.activity-table { width: 100%; border-collapse: collapse; text-align: left; }
.activity-table th { color: #5C5C6B; font-weight: 600; font-size: 0.85rem; padding: 12px; border-bottom: 1px solid rgba(255,255,255,0.05); }
.activity-table td { padding: 16px 12px; border-bottom: 1px solid rgba(255,255,255,0.02); color: #A8A8B3; font-size: 0.9rem; }
.story-cell { display: flex; align-items: center; gap: 12px; color: #fff; font-weight: 500; }
.mini-cover { width: 32px; height: 44px; border-radius: 4px; object-fit: cover; }
.type-badge { padding: 4px 8px; border-radius: 4px; font-size: 0.75rem; font-weight: 700; }
.type-badge.comic { background: rgba(156,167,100,0.1); color: #9CA764; }
.type-badge.novel { background: rgba(241,232,199,0.1); color: #F1E8C7; }
.user-tag { color: #9CA764; background: rgba(156,167,100,0.05); padding: 4px 10px; border-radius: 100px; font-size: 0.8rem; }
</style>
