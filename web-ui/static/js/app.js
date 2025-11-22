// 数据可视化管理平台 - 前端逻辑

class AuthApp {
    constructor() {
        // API 基础地址 - 使用相对路径，通过nginx代理
        this.baseUrl = '';
        this.apiBase = '/api/v1';
        this.token = localStorage.getItem('jwt_token') || '';
        this.currentUser = JSON.parse(localStorage.getItem('user_info') || 'null');
        
        console.log('🔧 API 配置:', this.baseUrl + this.apiBase);
        
        this.init();
    }

    init() {
        // 检查是否已登录
        if (this.token && this.currentUser) {
            this.showMainPage();
        }
    }

    // 通用 API 调用方法
    async callApi(endpoint, method = 'GET', data = null, requiresAuth = false) {
        const url = `${this.baseUrl}${this.apiBase}${endpoint}`;
        console.log('🚀 发送请求:', method, url);
        console.log('📦 请求数据:', data);
        
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

        try {
            const response = await fetch(url, options);
            const result = await response.json();
            
            console.log('✅ 响应状态:', response.status);
            console.log('📥 响应数据:', result);
            
            return {
                success: response.ok,
                status: response.status,
                data: result
            };
        } catch (error) {
            console.error('❌ 请求失败:', error);
            return {
                success: false,
                error: error.message,
                status: 0
            };
        }
    }

    // 注册功能
    async register(username, email, password) {
        const result = await this.callApi('/auth/register', 'POST', {
            username: username,
            email: email,
            password: password
        });

        if (result.success) {
            this.showAlert('注册成功!请登录', 'success');
            // 切换到登录标签
            const loginTab = document.getElementById('login-tab');
            loginTab.click();
            // 自动填充用户名
            document.getElementById('loginIdentifier').value = username;
        } else {
            const errorMsg = result.data?.message || result.error || '注册失败,请重试';
            this.showAlert(errorMsg, 'danger');
        }

        return result;
    }

    // 登录功能
    async login(identifier, password) {
        const result = await this.callApi('/auth/login', 'POST', {
            username: identifier,
            password: password
        });

        if (result.success && result.data && result.data.token) {
            this.token = result.data.token;
            this.currentUser = result.data.user || { username: identifier };
            
            // 保存到本地存储
            localStorage.setItem('jwt_token', this.token);
            localStorage.setItem('user_info', JSON.stringify(this.currentUser));
            
            this.showAlert('登录成功!', 'success');
            
            // 延迟显示主页面
            setTimeout(() => {
                this.showMainPage();
            }, 500);
        } else {
            const errorMsg = result.data?.message || result.error || '登录失败,请检查用户名和密码';
            this.showAlert(errorMsg, 'danger');
        }

        return result;
    }

    // 登出功能
    async logout() {
        // 调用后端API
        await this.callApi('/auth/logout', 'POST', null, true);
        
        // 清除本地存储
        this.token = '';
        this.currentUser = null;
        localStorage.removeItem('jwt_token');
        localStorage.removeItem('user_info');
        
        this.showAlert('已退出登录', 'info');
        
        // 返回登录页面
        setTimeout(() => {
            this.showAuthPage();
        }, 500);
    }

    // 忘记密码
    async forgotPassword(email) {
        const result = await this.callApi('/auth/forgot-password', 'POST', {
            email: email
        });

        if (result.success) {
            this.showAlert('密码重置邮件已发送,请检查您的邮箱', 'success');
        } else {
            const errorMsg = result.data?.message || result.error || '发送失败,请重试';
            this.showAlert(errorMsg, 'danger');
        }

        return result;
    }

    // 显示主页面
    showMainPage() {
        document.getElementById('auth-page').classList.add('d-none');
        document.getElementById('main-page').classList.remove('d-none');
        
        // 更新用户名显示
        const usernameSpan = document.getElementById('currentUsername');
        if (usernameSpan && this.currentUser) {
            usernameSpan.textContent = this.currentUser.username || '用户';
        }
    }

    // 显示登录页面
    showAuthPage() {
        document.getElementById('auth-page').classList.remove('d-none');
        document.getElementById('main-page').classList.add('d-none');
        
        // 清空表单
        document.getElementById('loginForm').reset();
        document.getElementById('registerForm').reset();
    }

    // 显示提示信息
    showAlert(message, type = 'info') {
        // 删除旧的提示
        const oldAlerts = document.querySelectorAll('.custom-alert');
        oldAlerts.forEach(alert => alert.remove());

        const alertDiv = document.createElement('div');
        alertDiv.className = `alert alert-${type} alert-dismissible fade show custom-alert`;
        alertDiv.style.cssText = 'position: fixed; top: 20px; right: 20px; z-index: 9999; min-width: 300px;';
        alertDiv.innerHTML = `
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        `;
        
        document.body.appendChild(alertDiv);
        
        // 3秒后自动消失
        setTimeout(() => {
            if (alertDiv.parentNode) {
                alertDiv.remove();
            }
        }, 3000);
    }
}

// 创建全局实例
const app = new AuthApp();

// 处理注册表单提交
function handleRegister(event) {
    event.preventDefault();
    
    const username = document.getElementById('registerUsername').value.trim();
    const email = document.getElementById('registerEmail').value.trim();
    const password = document.getElementById('registerPassword').value;
    const confirmPassword = document.getElementById('registerConfirmPassword').value;
    
    // 验证密码
    if (password !== confirmPassword) {
        app.showAlert('两次输入的密码不一致', 'warning');
        return false;
    }
    
    if (password.length < 6) {
        app.showAlert('密码长度至少6位', 'warning');
        return false;
    }
    
    // 禁用按钮
    const btn = document.getElementById('registerBtn');
    btn.disabled = true;
    btn.innerHTML = '<i class="fas fa-spinner fa-spin me-2"></i>注册中...';
    
    // 调用注册接口
    app.register(username, email, password).finally(() => {
        btn.disabled = false;
        btn.innerHTML = '<i class="fas fa-user-plus me-2"></i>注册';
    });
    
    return false;
}

// 处理登录表单提交
function handleLogin(event) {
    event.preventDefault();
    
    const identifier = document.getElementById('loginIdentifier').value.trim();
    const password = document.getElementById('loginPassword').value;
    
    if (!identifier || !password) {
        app.showAlert('请输入用户名和密码', 'warning');
        return false;
    }
    
    // 禁用按钮
    const btn = document.getElementById('loginBtn');
    btn.disabled = true;
    btn.innerHTML = '<i class="fas fa-spinner fa-spin me-2"></i>登录中...';
    
    // 调用登录接口
    app.login(identifier, password).finally(() => {
        btn.disabled = false;
        btn.innerHTML = '<i class="fas fa-sign-in-alt me-2"></i>登录';
    });
    
    return false;
}

// 处理登出
function handleLogout() {
    if (confirm('确定要退出登录吗?')) {
        app.logout();
    }
}

// 处理忘记密码
function handleForgotPassword() {
    const email = prompt('请输入您的注册邮箱:');
    if (email && email.trim()) {
        app.forgotPassword(email.trim());
    }
}