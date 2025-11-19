// 模板加载器
class TemplateLoader {
    constructor() {
        this.templates = {};
        this.templateBasePath = 'template/';
    }

    // 加载模板文件
    async loadTemplate(templateName) {
        if (this.templates[templateName]) {
            return this.templates[templateName];
        }

        try {
            const response = await fetch(`${this.templateBasePath}${templateName}.html`);
            if (!response.ok) {
                throw new Error(`模板文件 ${templateName} 加载失败: ${response.status}`);
            }
            
            const templateContent = await response.text();
            this.templates[templateName] = templateContent;
            return templateContent;
        } catch (error) {
            console.error('模板加载错误:', error);
            return `<div class="alert alert-danger">模板加载失败: ${templateName}</div>`;
        }
    }

    // 渲染模板到指定元素
    async renderTemplate(templateName, containerId) {
        const container = document.getElementById(containerId);
        if (!container) {
            console.error(`容器元素不存在: ${containerId}`);
            return;
        }

        const templateContent = await this.loadTemplate(templateName);
        container.innerHTML = templateContent;
    }

    // 批量加载模板
    async loadTemplates(templateNames) {
        const promises = templateNames.map(name => this.loadTemplate(name));
        return Promise.all(promises);
    }

    // 获取模板内容（不渲染）
    async getTemplate(templateName) {
        return await this.loadTemplate(templateName);
    }
}

// 全局模板加载器实例
const templateLoader = new TemplateLoader();

// 页面初始化时加载模板
async function initializeTemplates() {
    try {
        // 加载所有模板
        await templateLoader.loadTemplates([
            'header',
            'auth-tab', 
            'public-tab',
            'protected-tab'
        ]);

        // 渲染主内容
        await renderMainContent();
        
    } catch (error) {
        console.error('模板初始化失败:', error);
        // 显示错误信息
        document.body.innerHTML = `
            <div class="container mt-5">
                <div class="alert alert-danger">
                    <h4><i class="fas fa-exclamation-triangle me-2"></i>页面加载失败</h4>
                    <p>${error.message}</p>
                    <button class="btn btn-primary" onclick="location.reload()">重新加载</button>
                </div>
            </div>
        `;
    }
}

// 渲染主内容
async function renderMainContent() {
    const mainContainer = document.getElementById('main-container');
    if (!mainContainer) return;

    // 加载头部模板
    const headerContent = await templateLoader.getTemplate('header');
    
    // 构建主内容HTML
    const mainContent = `
        ${headerContent}
        
        <!-- 标签内容 -->
        <div class="tab-content" id="apiTabsContent">
            <div id="auth-container"></div>
            <div id="public-container"></div>
            <div id="protected-container"></div>
        </div>
        
        </div>
        </div>
        </div>
        </div>
    `;

    mainContainer.innerHTML = mainContent;

    // 渲染各个标签内容
    await templateLoader.renderTemplate('auth-tab', 'auth-container');
    await templateLoader.renderTemplate('public-tab', 'public-container');
    await templateLoader.renderTemplate('protected-tab', 'protected-container');

    // 初始化标签功能
    initializeTabs();
}

// 初始化标签功能
function initializeTabs() {
    // 标签切换事件
    const tabButtons = document.querySelectorAll('[data-bs-toggle="tab"]');
    tabButtons.forEach(button => {
        button.addEventListener('click', function() {
            // 移除所有激活状态
            tabButtons.forEach(btn => btn.classList.remove('active'));
            // 添加当前激活状态
            this.classList.add('active');
        });
    });

    // 服务器地址变化监听
    const serverUrlInput = document.getElementById('serverUrl');
    if (serverUrlInput) {
        serverUrlInput.addEventListener('change', function(e) {
            localStorage.setItem('server_url', e.target.value);
        });
    }
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    initializeTemplates();
});