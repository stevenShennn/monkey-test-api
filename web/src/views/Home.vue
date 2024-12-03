<template>
  <div class="home-page">
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

    <header class="home-header">
      <div class="title-wrapper">
        <h1>🐒 Monkey API Testor</h1>
        <p class="description">Monkey API Testor 是一款高效的 API 压力测试工具，支持并发测试、请求样本解析、结果分析等功能。</p>
      </div>
    </header>

    <main class="home-content">
      <el-card class="feature-card">
        <template #header>
          <div class="card-header">
            <h2>开始测试</h2>
          </div>
        </template>
        
        <div class="feature-content">
          <div class="input-section">
            <el-input
              v-model="curlCommands"
              type="textarea"
              :rows="6"
              placeholder="请输入多个 cURL 命令，每行一个..."
              @input="parseCurlCommands"
            />
            
            <el-card v-if="parsedCurls.length" class="curl-parse-result">
              <template #header>
                <div class="card-header">
                  <h3>CURL 解析结果</h3>
                </div>
              </template>
              
              <el-table :data="parsedCurls" style="width: 100%">
                <el-table-column prop="url" label="请求URL" />
                <el-table-column prop="method" label="请求方法" />
                <el-table-column label="请求头">
                  <template #default="{ row }">
                    <el-tag 
                      v-for="(value, key) in row.headers" 
                      :key="key"
                      class="mx-1"
                    >
                      {{ key }}: {{ value }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="data" label="请求体">
                  <template #default="{ row }">
                    <pre>{{ row.data }}</pre>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>
          </div>
          
          <div class="action-section">
            <el-button 
              type="primary" 
              @click="generateTestCase"
              :disabled="!parsedCurls.length"
            >
              生成测试案例
            </el-button>
          </div>
        </div>
      </el-card>

      <el-card class="history-card">
        <template #header>
          <div class="card-header">
            <h2>历史记录</h2>
          </div>
        </template>
        
        <el-table :data="historyList" style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="testTime" label="测试时间" width="160">
            <template #default="{ row }">
              {{ formatDate(row.testTime) }}
            </template>
          </el-table-column>
          <el-table-column prop="domain" label="域名" min-width="180" show-overflow-tooltip>
            <template #default="{ row }">
              <a @click="navigateToSiteDetail(row.domain)" class="domain-link">
                {{ row.domain }}
              </a>
            </template>
          </el-table-column>
          <el-table-column prop="vulnerabilityCount" label="漏洞数量" width="100" align="right">
            <template #default="{ row }">
              <el-tag :type="row.vulnerabilityCount > 0 ? 'danger' : 'success'">
                {{ row.vulnerabilityCount }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sampleCount" label="请求样本数量" width="120" align="right" />
          <el-table-column prop="avgResponseTime" label="平均响应时间" width="120" align="right">
            <template #default="{ row }">
              {{ row.avgResponseTime }}ms
            </template>
          </el-table-column>
          <el-table-column prop="maxResponseTime" label="最大响应时长" width="120" align="right">
            <template #default="{ row }">
              {{ row.maxResponseTime }}ms
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button 
                link 
                type="primary" 
                @click="viewResult(row.id)"
              >
                查看
              </el-button>
              <el-button 
                link 
                type="primary" 
                @click="retest(row.id)"
              >
                重新测试
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 表单数据
const curlCommands = ref('')

// 格式化日期的函数
const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

// 历史记录数据
const historyList = ref([
  {
    id: 1,
    testTime: '2024-01-15T10:30:00',
    domain: 'api.example.com',
    sampleCount: 1000,
    avgResponseTime: 156,
    maxResponseTime: 789,
    vulnerabilityCount: 2
  },
  {
    id: 2,
    testTime: '2024-01-15T11:45:00',
    domain: 'api2.example.com',
    sampleCount: 500,
    avgResponseTime: 89,
    maxResponseTime: 234,
    vulnerabilityCount: 0
  }
  // 可以添加更多测试数据
])

// 解析结果的响应式数组
const parsedCurls = ref([])

// 解析多个 cURL 命令
const parseCurlCommands = (value) => {
  // 先按换行分割，然后合并带有 \ 的行
  const lines = value.split('\n')
  const mergedCommands = []
  let currentCommand = ''

  for (let line of lines) {
    line = line.trim()
    if (line.endsWith('\\')) {
      // 如果行尾有 \，则去掉 \ 并与下一行合并
      currentCommand += line.slice(0, -1) + ' '
    } else {
      // 如果行尾没有 \，则这是完整的命令
      currentCommand += line
      if (currentCommand.trim().startsWith('curl')) {
        mergedCommands.push(currentCommand.trim())
      }
      currentCommand = ''
    }
  }

  // 解析每个合并后的命令
  parsedCurls.value = mergedCommands
    .map(parseCurlCommand)
    .filter(curl => curl.url)
}

// 解析单个 cURL 命令
const parseCurlCommand = (command) => {
  const parsed = {
    url: '',
    method: 'GET',
    headers: {},
    data: null
  }

  if (!command.trim()) return parsed

  try {
    // 解析 URL
    const urlMatch = command.match(/curl\s+['"]([^'"]+)['"]/) || 
                    command.match(/curl\s+'([^']+)'/) ||
                    command.match(/curl\s+"([^"]+)"/) ||
                    command.match(/curl\s+([^\s'"]+)/)
    
    if (urlMatch) {
      parsed.url = urlMatch[1]
    }

    // 解析请求头
    const headerMatches = command.matchAll(/-H\s+['"]([^'"]+)['"]|-H\s+'([^']+)'|-H\s+"([^"]+)"|-H\s+([^\s'"]+)/g)
    for (const match of headerMatches) {
      const headerStr = match[1] || match[2] || match[3] || match[4]
      const colonIndex = headerStr.indexOf(':')
      if (colonIndex > 0) {
        const key = headerStr.slice(0, colonIndex).trim()
        const value = headerStr.slice(colonIndex + 1).trim()
        parsed.headers[key] = value
      }
    }

    // 解析请求体（保持不变）
    const dataRegexes = [
      /--data\s+['"]([^'"]+)['"]/,       // --data "data"
      /-d\s+['"]([^'"]+)['"]/,           // -d "data"
      /--data\s+'([^']+)'/,              // --data 'data'
      /-d\s+'([^']+)'/,                  // -d 'data'
      /--data\s+([^\s'"'][^\s]+)/,       // --data data
      /-d\s+([^\s'"'][^\s]+)/            // -d data
    ]

    for (const regex of dataRegexes) {
      const match = command.match(regex)
      if (match) {
        try {
          parsed.data = JSON.parse(match[1])
        } catch {
          parsed.data = match[1]
        }
        break
      }
    }

    console.log('解析结果:', parsed) // 调试用

  } catch (error) {
    console.error('解析 CURL 失败:', error)
  }

  return parsed
}

// 方法
const generateTestCase = async () => {
  try {
    // TODO: 调用后端 API 生成测试案例
    router.push('/execute')
  } catch (error) {
    console.error('生成测试案例失败:', error)
  }
}

const viewResult = (id) => {
  router.push(`/execute?id=${id}`)
}

const retest = async (id) => {
  try {
    // TODO: 调用后端 API 重新测试
    router.push('/execute')
  } catch (error) {
    console.error('重新测试失败:', error)
  }
}

// emoji 列表
const emojis = [
  '🐒', '🍌', '🎯', '🚀', '🔍', '⚡️', 
  '🛠', '📊', '🎮', '🎨', '🔥', '⭐️',
  '🎪', '🎭', '🎪', '🎨', '🎯', '🎲',
  '🔮', '💫', '✨', '💡', '⚡️', '🌟',
  '🌈', '🎪', '🎭', '🎪', '🎨', '🎯',
  '🐒', '🍌', '🎯', '🚀', '🔍', '⚡️'
];

// 导航到网站详情页
const navigateToSiteDetail = (domain) => {
  router.push({ name: 'SiteDetail', params: { domain } })
}
</script>

<style lang="scss" scoped>
:root {
  --color-border-default: rgba(0, 0, 0, 0.1);
  --color-canvas-subtle: #f6f8fa;
  --color-accent-fg: #0969da;
  --color-danger-fg: #cf222e;
  --color-success-fg: #1a7f37;
}

.home-page {
  position: relative;
  padding: 40px;
  max-width: 1280px;
  margin: 0 auto;
  min-height: 100vh;
  background: #ffffff;

  // 方案一：半透明的 emoji 背景
  &::before {
    content: '🐒 🍌 🎯 🚀 🔍 ⚡️ 🛠 📊 🎮 🎨 🔥 ⭐️';
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 40px;
    padding: 40px;
    font-size: 24px;
    opacity: 0.05;
    pointer-events: none;
    z-index: 0;
    overflow: hidden;
    animation: floatingEmojis 60s linear infinite;
  }
}

// 方案二：动态的 emoji 背景
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

    &:nth-child(2n) {
      animation-duration: 12s;
    }
    &:nth-child(3n) {
      animation-duration: 10s;
    }
    &:nth-child(4n) {
      animation-duration: 14s;
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

@keyframes floatingEmojis {
  0% {
    transform: translateY(0);
  }
  100% {
    transform: translateY(-50%);
  }
}

.home-header {
  margin-bottom: 40px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--color-border-default);

  .title-wrapper {
    h1 {
      font-size: 44px;
      margin: 0 0 12px 0;
      display: flex;
      align-items: center;
      gap: 16px;
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
      line-height: 1.6;
      max-width: 600px;
    }
  }
}

.feature-card {
  margin-bottom: 32px;
  border-radius: 16px;
  border: 1px solid var(--color-border-default);
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
  }

  :deep(.el-card__header) {
    background: var(--color-canvas-subtle);
    border-bottom: 1px solid var(--color-border-default);
    padding: 20px;
    border-radius: 16px 16px 0 0;
  }

  :deep(.el-card__body) {
    padding: 24px;
  }
}

.feature-content {
  .input-section {
    margin-bottom: 24px;

    .el-input {
      :deep(.el-textarea__inner) {
        border-radius: 12px;
        border: 1px solid var(--color-border-default);
        font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, monospace;
        padding: 16px;
        transition: all 0.3s ease;
        
        &:focus {
          border-color: #0969da;
          box-shadow: 0 0 0 4px rgba(9, 105, 218, 0.15);
        }
      }
    }
  }

  .action-section {
    .el-button {
      padding: 12px 28px;
      font-size: 15px;
      font-weight: 500;
      border-radius: 12px;
      background: #2da44e;
      border: none;
      color: #ffffff;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-1px);
        background: #2c974b;
        box-shadow: 0 4px 12px rgba(45, 164, 78, 0.2);
      }

      &:disabled {
        background: #94d3a2;
        cursor: not-allowed;
        transform: none;
        box-shadow: none;
      }
    }
  }
}

.curl-parse-result {
  margin-top: 24px;
  border: 1px solid var(--color-border-default);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);

  .el-tag {
    margin: 0 6px 6px 0;
    border-radius: 8px;
    padding: 4px 10px;
    font-size: 12px;
    border: none;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  }

  pre {
    background: var(--color-canvas-subtle);
    border-radius: 12px;
    padding: 20px;
    margin: 0;
    font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, monospace;
    font-size: 13px;
    line-height: 1.5;
  }

  :deep(.el-table) {
    border-radius: 16px;
    overflow: hidden;
    
    th {
      background: var(--color-canvas-subtle);
      font-weight: 600;
      color: #24292f;
      height: 48px;
    }

    td {
      font-size: 14px;
      padding: 16px 0;
    }
  }
}

.history-card {
  border-radius: 16px;
  border: 1px solid var(--color-border-default);
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);

  :deep(.el-table) {
    border-radius: 16px;
    overflow: hidden;

    th {
      background: var(--color-canvas-subtle);
      font-weight: 600;
      color: #24292f;
      height: 48px;
    }

    td {
      padding: 16px 0;
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

  .el-tag {
    border-radius: 8px;
    padding: 4px 10px;
    font-size: 12px;
    border: none;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);

    &.el-tag--success {
      color: var(--color-success-fg);
      background: rgba(26, 127, 55, 0.08);
    }

    &.el-tag--danger {
      color: var(--color-danger-fg);
      background: rgba(207, 34, 46, 0.08);
    }
  }
}

.card-header {
  h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #24292f;
    letter-spacing: -0.3px;
  }
}

.domain-link {
  color: var(--color-accent-fg);
  cursor: pointer;
  text-decoration: underline;
}

.domain-link:hover {
  color: var(--color-accent-fg-hover);
}
</style> 