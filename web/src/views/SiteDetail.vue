<template>
  <div class="site-detail-page">
    <!-- 动态 emoji 背景 -->
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

    <header class="site-header">
      <h1>{{ siteInfo.domain }}</h1>
      <p class="description">最近更新时间: {{ formatDate(siteInfo.lastTestTime) }}</p>
    </header>

    <main class="site-content">
      <!-- 评分卡片 -->
      <el-row :gutter="24" class="score-section">
        <el-col :span="8">
          <div class="score-card">
            <div class="score">{{ siteInfo.performanceScore }}</div>
            <div class="label">性能评分</div>
            <el-progress 
              :percentage="siteInfo.performanceScore" 
              :color="scoreColor"
              :show-text="false"
            />
          </div>
        </el-col>
        <el-col :span="8">
          <div class="score-card">
            <div class="score">{{ siteInfo.securityScore }}</div>
            <div class="label">安全评分</div>
            <el-progress 
              :percentage="siteInfo.securityScore" 
              :color="scoreColor"
              :show-text="false"
            />
          </div>
        </el-col>
        <el-col :span="8">
          <div class="score-card">
            <div class="score">{{ siteInfo.stabilityScore }}</div>
            <div class="label">稳定性评分</div>
            <el-progress 
              :percentage="siteInfo.stabilityScore" 
              :color="scoreColor"
              :show-text="false"
            />
          </div>
        </el-col>
      </el-row>

      <!-- 性能指标卡片 -->
      <el-card class="metrics-card">
        <template #header>
          <div class="card-header">
            <h2>性能指标</h2>
          </div>
        </template>
        <el-row :gutter="24">
          <el-col :span="8" v-for="metric in performanceMetrics" :key="metric.label">
            <div class="metric-item">
              <div class="metric-value">{{ metric.value }}</div>
              <div class="metric-label">{{ metric.label }}</div>
              <div class="metric-trend" :class="metric.trend">
                {{ metric.trend === 'up' ? '↑' : '↓' }} {{ metric.change }}%
              </div>
            </div>
          </el-col>
        </el-row>
      </el-card>

      <!-- 历史测试记录 -->
      <el-card class="history-card">
        <template #header>
          <div class="card-header">
            <h2>测试历史</h2>
            <el-button type="primary" @click="startNewTest">开始新测试</el-button>
          </div>
        </template>
        
        <div class="chart-container">
          <el-tabs v-model="activeChart">
            <el-tab-pane label="响应时间趋势" name="responseTime">
              <!-- 这里可以使用 ECharts 等图表库 -->
              <div class="chart" ref="responseTimeChart"></div>
            </el-tab-pane>
            <el-tab-pane label="错误率趋势" name="errorRate">
              <div class="chart" ref="errorRateChart"></div>
            </el-tab-pane>
          </el-tabs>
        </div>

        <el-table :data="testHistory" style="width: 100%">
          <el-table-column prop="testTime" label="测试时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.testTime) }}
            </template>
          </el-table-column>
          <el-table-column prop="avgResponseTime" label="平均响应时间" width="150">
            <template #default="{ row }">
              {{ row.avgResponseTime }}ms
            </template>
          </el-table-column>
          <el-table-column prop="successRate" label="成功率" width="120">
            <template #default="{ row }">
              <el-tag :type="row.successRate > 95 ? 'success' : 'warning'">
                {{ row.successRate }}%
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="concurrency" label="并发数" width="120" />
          <el-table-column prop="totalRequests" label="总请求数" width="120" />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="viewDetail(row.id)">
                详情
              </el-button>
              <el-button link type="primary" @click="retest(row.id)">
                重测
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <!-- 问题分析 -->
      <el-card class="issues-card">
        <template #header>
          <div class="card-header">
            <h2>问题分析</h2>
          </div>
        </template>
        
        <el-timeline>
          <el-timeline-item
            v-for="issue in issues"
            :key="issue.id"
            :type="issue.severity"
            :timestamp="formatDate(issue.time)"
          >
            <h3>{{ issue.title }}</h3>
            <p>{{ issue.description }}</p>
            <el-tag :type="issue.severity">{{ issue.type }}</el-tag>
          </el-timeline-item>
        </el-timeline>
      </el-card>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

// 站点信息
const siteInfo = ref({
  domain: 'api.example.com',
  lastTestTime: new Date(),
  performanceScore: 85,
  securityScore: 92,
  stabilityScore: 88
})

// 性能指标
const performanceMetrics = ref([
  {
    label: '平均响应时间',
    value: '156ms',
    trend: 'down',
    change: 12
  },
  {
    label: '成功率',
    value: '99.8%',
    trend: 'up',
    change: 0.5
  },
  {
    label: '吞吐量',
    value: '1.2k/s',
    trend: 'up',
    change: 15
  }
])

// 测试历史
const testHistory = ref([])

// 问题列表
const issues = ref([
  {
    id: 1,
    title: '响应时间异常',
    description: '在高并发情况下，部分接口响应时间超过 1s',
    severity: 'warning',
    type: '性能问题',
    time: new Date()
  }
])

// emoji 列表
const emojis = [
  '📊', '📈', '🔍', '⚡️', '🎯', '📡',
  '💻', '🔐', '📱', '🌐', '⚙️', '🔄',
  '📶', '🔌', '💡', '🔋', '📡', '🎮'
]

// 格式化日期
const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

// 评分颜色
const scoreColor = (percentage) => {
  if (percentage > 80) return '#67C23A'
  if (percentage > 60) return '#E6A23C'
  return '#F56C6C'
}

// 方法
const startNewTest = () => {
  router.push('/home')
}

const viewDetail = (id) => {
  router.push(`/execute?id=${id}`)
}

const retest = async (id) => {
  // TODO: 实现重测逻辑
}

// 生命周期钩子
onMounted(() => {
  // TODO: 获取站点详情数据
})

// 通过 props 接收 domain 参数
const props = defineProps({
  domain: {
    type: String,
    required: true
  }
})

// 使用 props.domain 获取域名
console.log('当前域名:', props.domain)
</script>

<style lang="scss" scoped>
.site-detail-page {
  padding: 40px;
  max-width: 1280px;
  margin: 0 auto;
  min-height: 100vh;
  position: relative;
}

// 复用 Home.vue 的 emoji 背景样式
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

.site-header {
  margin-bottom: 40px;
  position: relative;
  z-index: 1;

  h1 {
    font-size: 44px;
    margin: 0 0 12px 0;
    font-weight: 600;
    background: linear-gradient(120deg, #24292f 30%, #57606a);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    letter-spacing: -0.5px;
  }

  .description {
    font-size: 16px;
    color: #57606a;
    margin: 0;
  }
}

.score-section {
  margin-bottom: 32px;

  .score-card {
    background: rgba(255, 255, 255, 0.7);
    backdrop-filter: blur(20px);
    border-radius: 16px;
    padding: 24px;
    text-align: center;
    border: 1px solid var(--color-border-default);
    transition: transform 0.3s ease;

    &:hover {
      transform: translateY(-2px);
    }

    .score {
      font-size: 48px;
      font-weight: 600;
      color: #24292f;
      margin-bottom: 8px;
    }

    .label {
      font-size: 14px;
      color: #57606a;
      margin-bottom: 16px;
    }
  }
}

.metrics-card, .history-card, .issues-card {
  margin-bottom: 32px;
  border-radius: 16px;
  border: 1px solid var(--color-border-default);
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  
  .metric-item {
    text-align: center;
    padding: 20px;

    .metric-value {
      font-size: 32px;
      font-weight: 600;
      color: #24292f;
    }

    .metric-label {
      font-size: 14px;
      color: #57606a;
      margin: 8px 0;
    }

    .metric-trend {
      font-size: 14px;
      
      &.up { color: var(--color-success-fg); }
      &.down { color: var(--color-danger-fg); }
    }
  }
}

.chart-container {
  margin: 20px 0;
  
  .chart {
    height: 300px;
    background: var(--color-canvas-subtle);
    border-radius: 12px;
    margin: 20px 0;
  }
}

// 复用其他样式
.card-header {
  h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #24292f;
    letter-spacing: -0.3px;
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