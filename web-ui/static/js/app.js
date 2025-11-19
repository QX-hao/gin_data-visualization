// Go-Web API 前端测试脚本
class ApiTester {
    constructor() {
        this.baseUrl = 'http://localhost:1234';
        this.apiBase = '/api/v1';
        this.token = localStorage.getItem('jwt_token') || '';
        this.init();
    }

    init() {
        // 从本地存储加载配置
        const savedUrl = localStorage.getItem('server_url');
        if (savedUrl) {
            this.baseUrl = savedUrl;
            document.getElementById('serverUrl').value = savedUrl;
        }
        
        // 监听服务器地址变化
        document.getElementById('serverUrl').addEventListener('change', (e) => {
            this.baseUrl = e.target.value;
            localStorage.setItem('server_url', this.baseUrl);
        });

        // 显示当前令牌状态
        this.updateTokenStatus();
    }

    // 通用 API 调用方法
    async callApi(endpoint, method = 'GET', data = null, requiresAuth = false) {
        const url = `${this.baseUrl}${this.apiBase}${endpoint}`;
        const options = {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            }
        };

        if (requiresAuth && this.token) {
            options.headers['Authorization'] = `Bearer ${this.token}`;
        }

        if (data && (method === 'POST' || method === 'PUT')) {
            options.body = JSON.stringify(data);
        }

        // 显示加载状态
        this.showLoading(endpoint);

        try {
            const response = await fetch(url, options);
            const result = await response.json();
            
            return {
                success: response.ok,
                status: response.status,
                data: result,
                headers: Object.fromEntries(response.headers.entries())
            };
        } catch (error) {
            return {
                success: false,
                error: error.message,
                status: 0
            };
        } finally {
            this.hideLoading(endpoint);
        }
    }

    // 显示加载状态
    showLoading(endpoint) {
        const button = document.querySelector(`[onclick*="${endpoint.split('/').pop()}"]`);
        if (button) {
            button.disabled = true;
            button.innerHTML = `<span class="loading"></span> 请求中...`;
        }
    }

    // 隐藏加载状态
    hideLoading(endpoint) {
        const button = document.querySelector(`[onclick*="${endpoint.split('/').pop()}"]`);
        if (button) {
            button.disabled = false;
            // 恢复原始文本（需要根据具体按钮调整）
            const originalText = this.getButtonOriginalText(endpoint);
            button.innerHTML = originalText;
        }
    }

    // 获取按钮原始文本
    getButtonOriginalText(endpoint) {
        const textMap = {
            'register': '<i class="fas fa-user-plus me-2"></i>注册',
            'login': '<i class="fas fa-sign-in-alt me-2"></i>登录',
            'logout': '<i class="fas fa-sign-out-alt me-2"></i>登出',
            'refreshToken': '<i class="fas fa-sync-alt me-2"></i>刷新令牌',
            'forgotPassword': '<i class="fas fa-key me-2"></i>忘记密码',
            'testHealth': '<i class="fas fa-heartbeat me-2"></i>测试健康检查',
            'getProfile': '<i class="fas fa-user me-2"></i>获取用户资料',
            'updateProfile': '<i class="fas fa-edit me-2"></i>更新用户资料',
            'changePassword': '<i class="fas fa-lock me-2"></i>修改密码'
        };
        return textMap[endpoint.split('/').pop()] || '操作';
    }

    // 显示结果
    displayResult(tabId, result) {
        const resultElement = document.getElementById(tabId + 'Result');
        const timestamp = new Date().toLocaleString();
        
        let html = `<div class="mb-3 p-3 rounded ${result.success ? 'bg-success bg-opacity-10' : 'bg-danger bg-opacity-10'}">
            <span class="${result.success ? 'status-success' : 'status-error'}">
                ${result.success ? '✅ 成功' : '❌ 失败'} - ${timestamp}
            </span>
        </div>`;

        if (result.status) {
            html += `<div class="mb-2"><strong>状态码:</strong> <code>${result.status}</code></div>`;
        }

        if (result.error) {
            html += `<div class="mb-2"><strong>错误:</strong> <span class="text-danger">${result.error}</span></div>`;
        }

        if (result.data) {
            html += `<div class="mb-2"><strong>响应数据:</strong></div>
                    <pre class="p-3 rounded bg-dark text-light">${JSON.stringify(result.data, null, 2)}</pre>`;
        }

        if (result.headers) {
            html += `<div class="mb-2"><strong>响应头:</strong></div>
                    <pre class="p-3 rounded bg-secondary text-light">${JSON.stringify(result.headers, null, 2)}</pre>`;
        }

        // 限制最多显示5条记录
        const currentContent = resultElement.innerHTML;
        const records = currentContent.split('<div class="mb-3 p-3 rounded');
        if (records.length > 5) {
            resultElement.innerHTML = html + currentContent.split('<div class="mb-3 p-3 rounded').slice(0, 4).join('<div class="mb-3 p-3 rounded');
        } else {
            resultElement.innerHTML = html + currentContent;
        }

        // 自动滚动到最新结果
        resultElement.scrollTop = 0;
    }

    // 更新令牌状态显示
    updateTokenStatus() {
        const tokenStatus = document.getElementById('tokenStatus');
        if (tokenStatus) {
            tokenStatus.textContent = this.token ? '✅ 已登录' : '❌ 未登录';
            tokenStatus.className = this.token ? 'badge bg-success' : 'badge bg-danger';
        }
    }

    // 认证接口
    async register() {
        const username = document.getElementById('registerUsername').value;
        const email = document.getElementById('registerEmail').value;
        const password = document.getElementById('registerPassword').value;

        if (!username || !email || !password) {
            this.showAlert('请填写所有必填字段', 'warning');
            return;
        }

        const result = await this.callApi('/auth/register', 'POST', {
            username: username,
            email: email,
            password: password
        });

        this.displayResult('auth', result);
    }

    async login() {
        const identifier = document.getElementById('loginIdentifier').value;
        const password = document.getElementById('loginPassword').value;

        if (!identifier || !password) {
            this.showAlert('请填写用户名/邮箱和密码', 'warning');
            return;
        }

        const result = await this.callApi('/auth/login', 'POST', {
            identifier: identifier,
            password: password
        });

        if (result.success && result.data && result.data.token) {
            this.token = result.data.token;
            localStorage.setItem('jwt_token', this.token);
            this.updateTokenStatus();
            this.showAlert('登录成功！', 'success');
        }

        this.displayResult('auth', result);
    }

    async logout() {
        const result = await this.callApi('/auth/logout', 'POST', null, true);
        
        if (result.success) {
            this.token = '';
            localStorage.removeItem('jwt_token');
            this.updateTokenStatus();
            this.showAlert('已登出', 'info');
        }

        this.displayResult('auth', result);
    }

    async refreshToken() {
        const result = await this.callApi('/auth/refresh', 'POST', null, true);
        
        if (result.success && result.data && result.data.token) {
            this.token = result.data.token;
            localStorage.setItem('jwt_token', this.token);
            this.updateTokenStatus();
            this.showAlert('令牌刷新成功！', 'success');
        }

        this.displayResult('auth', result);
    }

    async forgotPassword() {
        const email = prompt('请输入您的邮箱地址：');
        if (!email) return;

        const result = await this.callApi('/auth/forgot-password', 'POST', {
            email: email
        });

        this.displayResult('auth', result);
    }

    async resetPassword() {
        const token = prompt('请输入重置令牌：');
        const newPassword = prompt('请输入新密码：');
        
        if (!token || !newPassword) return;

        const result = await this.callApi('/auth/reset-password', 'POST', {
            token: token,
            new_password: newPassword
        });

        this.displayResult('auth', result);
    }

    // 公开接口
    async testHealth() {
        const result = await this.callApi('/public/health', 'GET');
        this.displayResult('public', result);
    }

    // 受保护接口
    async getProfile() {
        if (!this.token) {
            this.showAlert('请先登录获取令牌', 'warning');
            return;
        }

        const result = await this.callApi('/protected/users/profile', 'GET', null, true);
        this.displayResult('protected', result);
    }

    async updateProfile() {
        if (!this.token) {
            this.showAlert('请先登录获取令牌', 'warning');
            return;
        }

        const username = prompt('请输入新的用户名（留空保持不变）：');
        const email = prompt('请输入新的邮箱（留空保持不变）：');
        
        const updateData = {};
        if (username) updateData.username = username;
        if (email) updateData.email = email;

        if (Object.keys(updateData).length === 0) {
            this.showAlert('没有提供任何更新信息', 'info');
            return;
        }

        const result = await this.callApi('/protected/users/profile', 'PUT', updateData, true);
        this.displayResult('protected', result);
    }

    async changePassword() {
        if (!this.token) {
            this.showAlert('请先登录获取令牌', 'warning');
            return;
        }

        const oldPassword = prompt('请输入当前密码：');
        const newPassword = prompt('请输入新密码：');
        
        if (!oldPassword || !newPassword) {
            this.showAlert('请填写所有密码字段', 'warning');
            return;
        }

        const result = await this.callApi('/protected/users/password', 'PUT', {
            old_password: oldPassword,
            new_password: newPassword
        }, true);

        this.displayResult('protected', result);
    }

    // 手动设置令牌
    setToken() {
        const tokenInput = document.getElementById('jwtToken').value;
        if (tokenInput) {
            this.token = tokenInput;
            localStorage.setItem('jwt_token', this.token);
            this.updateTokenStatus();
            this.showAlert('令牌设置成功！', 'success');
        } else {
            this.showAlert('请输入有效的 JWT 令牌', 'warning');
        }
    }

    // 显示提示信息
    showAlert(message, type = 'info') {
        const alertDiv = document.createElement('div');
        alertDiv.className = `alert alert-${type} alert-dismissible fade show`;
        alertDiv.innerHTML = `
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        `;
        
        // 添加到页面顶部
        const container = document.querySelector('.container');
        container.insertBefore(alertDiv, container.firstChild);
        
        // 3秒后自动消失
        setTimeout(() => {
            if (alertDiv.parentNode) {
                alertDiv.remove();
            }
        }, 3000);
    }
}

// 全局函数供 HTML 调用
const apiTester = new ApiTester();

function register() { apiTester.register(); }
function login() { apiTester.login(); }
function logout() { apiTester.logout(); }
function refreshToken() { apiTester.refreshToken(); }
function forgotPassword() { apiTester.forgotPassword(); }
function resetPassword() { apiTester.resetPassword(); }
function testHealth() { apiTester.testHealth(); }
function getProfile() { apiTester.getProfile(); }
function updateProfile() { apiTester.updateProfile(); }
function changePassword() { apiTester.changePassword(); }
function setToken() { apiTester.setToken(); }

// 页面加载完成后的初始化
document.addEventListener('DOMContentLoaded', function() {
    // 添加令牌状态显示
    const protectedTab = document.getElementById('protected');
    if (protectedTab) {
        const statusHtml = `
            <div class="alert alert-info">
                <i class="fas fa-info-circle me-2"></i>
                当前令牌状态: <span id="tokenStatus" class="badge bg-secondary">${apiTester.token ? '✅ 已登录' : '❌ 未登录'}</span>
            </div>
        `;
        protectedTab.insertAdjacentHTML('afterbegin', statusHtml);
    }

    // 自动测试健康检查（可选）
    setTimeout(() => {
        // apiTester.testHealth();
    }, 1000);
});