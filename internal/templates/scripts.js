// API 端点
const API_URL = '/api/v1';

// 提交测试
async function submitTest() {
    const curlInput = document.getElementById('curlInput').value;
    if (!curlInput) {
        alert('请输入 cURL 命令');
        return;
    }

    try {
        const response = await fetch(`${API_URL}/test-requests`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ curl: curlInput }),
        });

        const data = await response.json();
        if (response.ok) {
            alert('测试已开始');
            refreshResults();
        } else {
            alert(data.error || '请求失败');
        }
    } catch (error) {
        alert('网络错误');
        console.error(error);
    }
}

// 刷新结果
async function refreshResults() {
    const resultsContainer = document.getElementById('testResults');
    try {
        const response = await fetch(`${API_URL}/test-requests`);
        const data = await response.json();
        
        resultsContainer.innerHTML = data.map(result => `
            <div class="result-item">
                <div class="result-header">
                    <span>请求 ID: ${result.requestID}</span>
                    <span>测试数量: ${result.testCount}</span>
                </div>
                <div class="result-content">
                    ${formatTestObjects(result.testObjects)}
                </div>
            </div>
        `).join('');
    } catch (error) {
        resultsContainer.innerHTML = '<p>加载结果失败</p>';
        console.error(error);
    }
}

// 格式化测试对象
function formatTestObjects(testObjects) {
    if (!testObjects || !testObjects.length) return '无测试数据';
    
    return testObjects.map(obj => `
        <div class="test-object">
            <span class="status status-${getStatusClass(obj.status)}">${obj.status}</span>
            <span>${obj.reason}</span>
        </div>
    `).join('');
}

// 获取状态样式类
function getStatusClass(status) {
    switch (status) {
        case '待处理': return 'pending';
        case '成功': return 'success';
        case '失败': return 'error';
        default: return 'pending';
    }
}

// 清空输入
function clearInput() {
    document.getElementById('curlInput').value = '';
}

// 加载历史记录项
function loadHistoryItem(requestID) {
    // 实现加载历史记录的逻辑
}

// 页面加载完成后刷新结果
document.addEventListener('DOMContentLoaded', refreshResults); 