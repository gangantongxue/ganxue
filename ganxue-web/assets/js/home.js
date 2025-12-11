import { showToast } from './ui.js';
import { authRequest, getCookie, logout } from './auth.js';

// 切换用户设置悬浮窗口显示/隐藏
function toggleUserSettingsModal() {
    const modal = document.getElementById('userSettingsModal');
    if (modal) {
        modal.classList.toggle('show');
        
        // 点击遮罩层关闭模态框
        modal.addEventListener('click', function(event) {
            if (event.target === modal) {
                toggleUserSettingsModal();
            }
        });
        
        // 阻止模态框内部点击事件冒泡到遮罩层
        const modalContent = modal.querySelector('.modal-content');
        if (modalContent) {
            modalContent.addEventListener('click', function(event) {
                event.stopPropagation();
            });
        }
    }
}

// 初始化用户设置功能
function initUserSettings() {
    // 绑定修改用户名按钮事件
    const updateUsernameBtn = document.getElementById('modalUpdateUsernameBtn');
    if (updateUsernameBtn) {
        updateUsernameBtn.addEventListener('click', handleUpdateUsername);
    }
    
    // 绑定确认注销账户按钮事件
    const deleteAccountBtn = document.getElementById('modalDeleteAccountBtn');
    if (deleteAccountBtn) {
        deleteAccountBtn.addEventListener('click', handleDeleteAccount);
    }
    
    // 绑定取消注销按钮事件
    const cancelDeleteBtn = document.getElementById('modalCancelDeleteBtn');
    if (cancelDeleteBtn) {
        cancelDeleteBtn.addEventListener('click', () => {
            document.getElementById('modalConfirmPassword').value = '';
        });
    }
    
    // 新用户名输入框事件监听
    const newUsernameInput = document.getElementById('modalNewUsername');
    if (newUsernameInput) {
        newUsernameInput.addEventListener('input', () => {
            // 实时验证用户名长度
            const username = newUsernameInput.value.trim();
            const updateButton = document.getElementById('modalUpdateUsernameBtn');
            
            if (updateButton) {
                if (username.length >= 3 && username.length <= 20) {
                    updateButton.disabled = false;
                } else {
                    updateButton.disabled = true;
                }
            }
        });
    }
}

// 处理修改用户名
async function handleUpdateUsername() {
    const newUsername = document.getElementById('modalNewUsername').value.trim();
    
    // 验证用户名格式
    if (!validateUsername(newUsername)) {
        showToast('用户名长度必须在3-20个字符之间', 'warning');
        return;
    }
    
    try {
        // 调用修改用户名API
        const response = await authRequest('/api/auth/update-username', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({ new_username: newUsername })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            // 修改成功
            showToast(data.message || '用户名修改成功！', 'success');
            
            // 更新UI显示新用户名
            document.getElementById('userName').textContent = newUsername;
            document.getElementById('modalCurrentUserName').textContent = newUsername;
            document.getElementById('modalNewUsername').value = '';
            
            // 可以考虑更新localStorage中的用户信息
        } else {
            // 修改失败
            showToast(data.message || '用户名修改失败', 'error');
        }
    } catch (error) {
        console.error('修改用户名失败:', error);
        showToast('修改用户名失败，请稍后重试', 'error');
    }
}

// 处理注销账户
async function handleDeleteAccount() {
    const password = document.getElementById('modalConfirmPassword').value;
    
    if (!password) {
        showToast('请输入密码以确认注销', 'warning');
        return;
    }
    
    // 二次确认
    if (!confirm('确定要注销您的账户吗？此操作不可撤销，您的所有数据将被删除！')) {
        return;
    }
    
    try {
        // 调用注销账户API
        const response = await authRequest('/api/auth/delete-account', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({ password: password })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            // 注销成功
            showToast(data.message || '账户注销成功！', 'success');
            
            // 清除登录状态
            localStorage.removeItem('shortToken');
            document.cookie = 'longToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
            
            // 延迟跳转到登录页面
            setTimeout(() => {
                window.location.href = 'index.html';
            }, 2000);
        } else {
            // 注销失败
            showToast(data.message || '账户注销失败', 'error');
            document.getElementById('modalConfirmPassword').value = '';
        }
    } catch (error) {
        console.error('注销账户失败:', error);
        showToast('注销账户失败，请稍后重试', 'error');
        document.getElementById('modalConfirmPassword').value = '';
    }
}

// 验证用户名
function validateUsername(username) {
    return username.length >= 3 && username.length <= 20;
}

// 检查用户是否已登录
function checkLogin() {
    const shortToken = localStorage.getItem('shortToken');
    const longToken = getCookie('longToken');
    const autoLoginToken = getCookie('auto_login_token');
    
    console.log('检查登录状态：', {
        shortToken: shortToken ? '已存在' : 'undefined',
        longToken: longToken ? '已存在' : 'undefined',
        autoLoginToken: autoLoginToken ? '已存在' : 'undefined'
    });
    
    if (!shortToken && !longToken && !autoLoginToken) {
        console.log('未检测到登录凭证，重定向到登录页面');
        window.location.href = 'index.html';
        return false;
    }
    return true;
}

// 注意：getCookie函数已从auth.js导入，此处不再需要

// 获取用户信息
async function getUserInfo() {
    try {
        const response = await authRequest('api/auth/user/info', {
            method: 'GET',
            headers: {
                'Accept': 'application/json'
            }
        });

        if (!response.ok) {
            throw new Error('获取用户信息失败');
        }

        const data = await response.json();
        const userData = data.data;

        // 更新UI
        const userNameElement = document.getElementById('userName');
        userNameElement.textContent = userData.user_info.user_name || '未知用户';
        
        // 更新悬浮窗口中的当前用户名
        const modalCurrentUserName = document.getElementById('modalCurrentUserName');
        if (modalCurrentUserName) {
            modalCurrentUserName.textContent = userData.user_info.user_name || '未知用户';
        }
        // 使用固定头像，不再使用用户名首字母
        document.getElementById('streak_days').textContent = userData.study_stats.streak_days || 0;
        document.getElementById('total_days').textContent = userData.study_stats.total_days || 0;

        // 存储上次学习的章节信息
        if (userData.study_stats.last_time !== undefined) {
            const goLastChapter = userData.study_stats.go_last_chapter === "" ? "0000" : userData.study_stats.go_last_chapter;
            const cLastChapter = userData.study_stats.c_last_chapter === ""? "1000" : userData.study_stats.c_last_chapter;
            const cppLastChapter = userData.study_stats.cpp_last_chapter === ""? "2000" : userData.study_stats.cpp_last_chapter;
            const cloudcomputingLastChapter = userData.study_stats.cloudcomputing_last_chapter === ""? "3000" : userData.study_stats.cloudcomputing_last_chapter;
            localStorage.setItem('goLastChapter', goLastChapter);
            localStorage.setItem('cLastChapter', cLastChapter);
            localStorage.setItem('cppLastChapter', cppLastChapter);
            localStorage.setItem('cloudcomputingLastChapter', cloudcomputingLastChapter);
        }

        return userData;
    } catch (error) {
        console.error('获取用户信息失败:', error);
        // 处理特定的500服务器错误
        if (error.message === 'SERVER_ERROR_500') {
            showToast('服务器暂时不可用，请稍后重试', 'error');
        } else {
            showToast('获取用户信息失败，请重试', 'error');
        }
        return null;
    }
}

// 退出登录
async function handleLogout() {
    try {
        await logout();
    } catch (error) {
        console.error('登出失败:', error);
        showToast('登出失败，请稍后重试', 'error');
    }
}

// 跳转到设置页面
function goToSettings() {
    window.location.href = 'settings.html';
}

// 获取目录数据
async function getCatalogue() {
    try {
        const response = await authRequest('/api/auth/get-catalogue', {
            method: 'GET',
            headers: {
                'Accept': 'application/json'
            }
        });

        if (!response.ok) {
            throw new Error('获取目录失败');
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('获取目录失败:', error);
        showToast('获取目录失败，请重试', 'error');
        return null;
    }
}

// 渲染目录
function renderCatalogue(catalogueData) {
    if (!catalogueData) return;
    
    // 存储目录数据以便后续使用
    window.catalogueData = catalogueData;
    
    // 初始化目录区域
    const catalogueList = document.getElementById('catalogueList');
    catalogueList.innerHTML = '<li>欢迎使用敢学<br>让我们开启编程学习之旅吧</li>';
    
    // 为科目按钮添加鼠标悬停事件
    setupSubjectHoverEvents();
}

// 设置科目按钮的鼠标悬停事件
function setupSubjectHoverEvents() {
    const golangBtn = document.getElementById('golangBtn');
    const cBtn = document.getElementById('cBtn');
    const cppBtn = document.getElementById('cppBtn');
    const cloudcomputingBtn = document.getElementById('cloudcomputingBtn');
    const catalogueTitle = document.getElementById('catalogueTitle');
    const catalogueList = document.getElementById('catalogueList');
    
    // Golang按钮悬停事件
    golangBtn.addEventListener('mouseenter', () => {
        showCatalogue('golang', 'Golang 目录');
    });
    
    // C语言按钮悬停事件
    cBtn.addEventListener('mouseenter', () => {
        showCatalogue('c', 'C语言 目录');
    });
    
    // C++按钮悬停事件
    cppBtn.addEventListener('mouseenter', () => {
        showCatalogue('cpp', 'C++ 目录');
    });
    
    // CloudComputing按钮悬停事件
    cloudcomputingBtn.addEventListener('mouseenter', () => {
        showCatalogue('cloudcomputing', 'CloudComputing 目录');
    });
    
    // 显示指定科目的目录
    function showCatalogue(subject, title) {
        if (!window.catalogueData) return;
        
        catalogueTitle.textContent = title;
        catalogueList.innerHTML = '';
        
        Object.entries(window.catalogueData[subject]).forEach(([id, title], index) => {
            const li = document.createElement('li');
            li.textContent = title;
            li.dataset.id = id;
            // 添加动画索引属性
            li.style.setProperty('--item-index', index);
            li.addEventListener('click', () => {
                localStorage.setItem('nowGroup', subject);
                window.location.href = `learn.html?id=${id}&group=${subject}`;
            });
            catalogueList.appendChild(li);
        });
    }
}

// 页面加载完成后执行
window.addEventListener('DOMContentLoaded', async () => {
    // 检查登录状态
    if (!checkLogin()) {
        return;
    }
    
    // 获取用户信息
    await getUserInfo();
    
    // 为用户名绑定点击事件，显示悬浮窗口
    const userNameElement = document.getElementById('userName');
    if (userNameElement) {
        userNameElement.addEventListener('click', toggleUserSettingsModal);
    }
    
    // 为关闭按钮绑定点击事件
    const closeModalBtn = document.getElementById('closeModalBtn');
    if (closeModalBtn) {
        closeModalBtn.addEventListener('click', toggleUserSettingsModal);
    }
    
    // 初始化用户设置功能
    initUserSettings();
    
    // 获取目录数据
    const catalogueData = await getCatalogue();
    renderCatalogue(catalogueData);
    
    // 绑定退出登录按钮事件
    document.getElementById('logoutBtn').addEventListener('click', handleLogout);
    
    // 绑定Golang按钮事件
    document.getElementById('golangBtn').addEventListener('click', () => {
        showToast('即将开始Golang学习之旅！', 'success');
        // 获取上次学习的章节ID，如果没有则默认为'0000'
        const lastChapter = localStorage.getItem('goLastChapter') || '0000';

        localStorage.setItem('nowGroup','golang')
        // 跳转到学习页面，并传递章节ID
        window.location.href = `learn.html?id=${lastChapter}&group=golang`;
    });
    
    // 绑定C语言按钮事件
    document.getElementById('cBtn').addEventListener('click', () => {
        showToast('即将开始C语言学习之旅！', 'success');
        // 获取上次学习的章节ID，如果没有则默认为'1000'
        const lastChapter = localStorage.getItem('cLastChapter') || '1000';

        localStorage.setItem('nowGroup','c')
        // 跳转到学习页面，并传递章节ID
        window.location.href = `learn.html?id=${lastChapter}&group=c`;
    });
    
    // 绑定C++按钮事件
    document.getElementById('cppBtn').addEventListener('click', () => {
        showToast('即将开始C++学习之旅！', 'success');
        // 获取上次学习的章节ID，如果没有则默认为'2000'
        const lastChapter = localStorage.getItem('cppLastChapter') || '2000';

        localStorage.setItem('nowGroup','cpp')
        // 跳转到学习页面，并传递章节ID
        window.location.href = `learn.html?id=${lastChapter}&group=cpp`;
    });
    
    // 绑定CloudComputing按钮事件
    document.getElementById('cloudcomputingBtn').addEventListener('click', () => {
        showToast('即将开始CloudComputing学习之旅！', 'success');
        // 获取上次学习的章节ID，如果没有则默认为'3000'
        const lastChapter = localStorage.getItem('cloudcomputingLastChapter') || '3000';

        localStorage.setItem('nowGroup','cloudcomputing')
        // 跳转到学习页面，并传递章节ID
        window.location.href = `learn.html?id=${lastChapter}&group=cloudcomputing`;
    });
});