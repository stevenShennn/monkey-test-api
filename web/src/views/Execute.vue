<template>
  <div class="execute-page">
    <div class="emoji-background">
      <span v-for="(emoji, index) in emojis" 
            :key="index" 
            class="emoji"
            :style="{
              animationDelay: `${index * 0.2}s`,
              transform: `rotate(${Math.random() * 20 - 10}deg)`
            }">
        {{ emoji }}
      </span>
    </div>

    <header class="execute-header">
      <h1>🚀 测试执行</h1>
      <div class="status-panel">
        <el-row :gutter="24">
          <el-col :span="6">
            <div class="status-card">
              <span class="label">总请求数</span>
              <span class="value">{{ stats.totalCount }}</span>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="status-card">
              <span class="label">成功数量</span>
              <span class="value success">{{ stats.successCount }}</span>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="status-card">
              <span class="label">失败数量</span>
              <span class="value error">{{ stats.failureCount }}</span>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="status-card">
              <span class="label">平均响应时间</span>
              <span class="value">{{ stats.avgResponseTime }}ms</span>
            </div>
          </el-col>
        </el-row>
      </div>
    </header>

    <main class="execute-content">
      <section class="progress-section">
        <el-card>
          <template #header>
            <div class="progress-header">
              <h2>执行进度</h2>
              <span class="status" :class="{ running: isRunning }">
                {{ isRunning ? '正在执行...' : '测试完成' }}
              </span>
            </div>
          </template>
          <el-progress 
            :percentage="progress" 
            :status="progressStatus"
          />
          <div class="progress-details">
            {{ stats.completedCount }} / {{ stats.totalRequests }} 个请求已完成
          </div>
        </el-card>
      </section>

      <section class="results-section">
        <el-card>
          <template #header>
            <div class="results-header">
              <h2>测试结果</h2>
              <div class="results-actions">
                <el-button @click="exportResults">导出结果</el-button>
                <el-button type="primary" @click="retestAll">重新测试</el-button>
              </div>
            </div>
          </template>

          <div class="results-filters">
            <el-select v-model="filters.status" placeholder="状态筛选">
              <el-option label="全部状态" value="all" />
              <el-option label="成功" value="success" />
              <el-option label="失败" value="failure" />
            </el-select>
            <el-input 
              v-model="filters.search" 
              placeholder="搜索 URL..."
              clearable
            />
          </div>

          <el-table :data="filteredResults" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="url" label="URL" />
            <el-table-column prop="method" label="方法" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
                  {{ row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="responseTime" label="响应时间" width="120">
              <template #default="{ row }">
                {{ row.responseTime }}ms
              </template>
            </el-table-column>
            <el-table-column prop="error" label="错误信息" />
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="{ row }">
                <el-button 
                  link 
                  type="primary" 
                  @click="retestSingle(row.id)"
                >
                  重测
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { createWebSocket } from '@/utils/websocket'

// 状态数据
const stats = ref({
  totalCount: 0,
  successCount: 0,
  failureCount: 0,
  avgResponseTime: 0,
  completedCount: 0,
  totalRequests: 0
})

const isRunning = ref(true)
const progress = computed(() => {
  if (stats.value.totalRequests === 0) return 0
  return Math.round((stats.value.completedCount / stats.value.totalRequests) * 100)
})

const progressStatus = computed(() => {
  if (!isRunning.value) return 'success'
  return ''
})

// 测试结果
const testResults = ref([])
const filters = ref({
  status: 'all',
  search: ''
})

const filteredResults = computed(() => {
  return testResults.value.filter(result => {
    const statusMatch = filters.value.status === 'all' || result.status === filters.value.status
    const searchMatch = result.url.toLowerCase().includes(filters.value.search.toLowerCase())
    return statusMatch && searchMatch
  })
})

// WebSocket 处理
const ws = ref(null)

onMounted(() => {
  ws.value = createWebSocket('ws://localhost:8080/ws/execute')
  
  ws.value.on('progress', updateProgress)
  ws.value.on('result', updateResult)
  ws.value.on('complete', handleComplete)
  
  ws.value.connect()
})

onUnmounted(() => {
  if (ws.value) {
    ws.value.close()
  }
})

const updateProgress = (data) => {
  stats.value = {
    ...stats.value,
    ...data
  }
}

const updateResult = (data) => {
  testResults.value.push(data)
}

const handleComplete = () => {
  isRunning.value = false
  ElMessage.success('测试执行完成')
}

// 功能方法
const exportResults = () => {
  const csv = convertToCSV(testResults.value)
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'test_results.csv'
  a.click()
  window.URL.revokeObjectURL(url)
}

const retestAll = () => {
  ElMessage.confirm('确定要重新测试所有请求吗？')
    .then(() => {
      testResults.value = []
      isRunning.value = true
      ws.value.send({ type: 'retest_all' })
    })
    .catch(() => {})
}

const retestSingle = (id) => {
  isRunning.value = true
  ws.value.send({ type: 'retest_single', id })
}

const convertToCSV = (data) => {
  const headers = ['ID', 'URL', '方法', '状态', '响应时间', '错误信息']
  const rows = data.map(item => [
    item.id,
    item.url,
    item.method,
    item.status,
    item.responseTime,
    item.error || ''
  ])
  
  return [headers, ...rows]
    .map(row => row.map(cell => `"${cell}"`).join(','))
    .join('\n')
}

// 添加 emoji 列表
const emojis = [
  '🚀', '⚡️', '📊', '🎯', '🔍', '💫',
  '📈', '🎮', '🔄', '✨', '💡', '🌟',
  '🔮', '📡', '🎲', '🔋', '⚙️', '🎪'
];
</script>

<style lang="scss" scoped>
.execute-page {
  padding: 40px;
  max-width: 1280px;
  margin: 0 auto;
  min-height: 100vh;
  position: relative;
}

.emoji-background {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-wrap: wrap;
  gap: 40px;
  padding: 40px;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;

  .emoji {
    font-size: 24px;
    opacity: 0.08;
    animation: floating 8s ease-in-out infinite;
    transform-origin: center;

    &:nth-child(2n) { animation-duration: 12s; }
    &:nth-child(3n) { animation-duration: 10s; }
    &:nth-child(4n) { animation-duration: 14s; }
  }
}

.execute-header {
  margin-bottom: 32px;
  position: relative;
  z-index: 1;

  h1 {
    font-size: 36px;
    margin: 0 0 24px 0;
    font-weight: 600;
    background: linear-gradient(120deg, #24292f 30%, #57606a);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    letter-spacing: -0.5px;
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

.status-panel {
  .status-card {
    background: rgba(255, 255, 255, 0.7);
    backdrop-filter: blur(20px);
    border-radius: 16px;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    border: 1px solid var(--color-border-default);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
    }

    .label {
      font-size: 14px;
      color: #57606a;
    }
    
    .value {
      font-size: 32px;
      font-weight: 600;
      color: #24292f;
      
      &.success { color: var(--color-success-fg); }
      &.error { color: var(--color-danger-fg); }
    }
  }
}

.progress-section, .results-section {
  position: relative;
  z-index: 1;

  .el-card {
    border-radius: 16px;
    border: 1px solid var(--color-border-default);
    background: rgba(255, 255, 255, 0.7);
    backdrop-filter: blur(20px);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
    margin-bottom: 24px;
    transition: transform 0.3s ease, box-shadow 0.3s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
    }
  }
}

.progress-header, .results-header {
  h2 {
    font-size: 18px;
    font-weight: 600;
    color: #24292f;
    margin: 0;
    letter-spacing: -0.3px;
  }

  .status {
    font-size: 14px;
    
    &.running {
      color: var(--color-accent-fg);
    }
  }
}

.results-filters {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;

  .el-select, .el-input {
    width: 200px;
  }
}

:deep(.el-table) {
  background: transparent;
  
  th {
    background: var(--color-canvas-subtle);
    font-weight: 600;
    color: #24292f;
  }

  .el-button {
    padding: 6px 14px;
    font-size: 13px;
    border-radius: 8px;
    transition: all 0.2s ease;
    
    &:hover {
      color: var(--color-accent-fg);
      background: rgba(9, 105, 218, 0.06);
    }
  }
}

@keyframes floating {
  0%, 100% {
    transform: translateY(0) rotate(0);
  }
  50% {
    transform: translateY(15px) rotate(5deg);
  }
}
</style> 