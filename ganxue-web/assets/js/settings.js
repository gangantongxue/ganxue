import { authRequest, getCookie } from './auth.js';
import { showToast } from './ui.js';

// 页面加载时执行
window.addEventListener('DOMContentLoaded', async () => {
    // 检查用户登录状态
    if (!checkLogin()) {
        return;
    }
    
    // 加载用户信息
    await loadUserInfo();
    
    // 绑定事件监听器
    bindEventListeners();
});

// 检查用户是否已登录
function checkLogin() {
    const shortToken = localStorage.getItem('shortToken');
    const longToken = getCookie('longToken');
    
    if (!shortToken && !longToken) {
        window.location.href = 'index.html';
        return false;
    }
    return true;
}

// 加载用户信息
async function loadUserInfo() {
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
        
        // 更新UI显示当前用户名
        document.getElementById('currentUserName').textContent = userData.user_info.user_name || '未知用户';
    } catch (error) {
        console.error('获取用户信息失败:', error);
        showToast('获取用户信息失败，请刷新页面重试', 'error');
    }
}

// 绑定事件监听器
function bindEventListeners() {
    // 修改用户名按钮点击事件
    document.getElementById('updateUsernameBtn').addEventListener('click', handleUpdateUsername);
    
    // 确认注销账户按钮点击事件
    document.getElementById('deleteAccountBtn').addEventListener('click', handleDeleteAccount);
    
    // 取消注销按钮点击事件
    document.getElementById('cancelDeleteBtn').addEventListener('click', () => {
        document.getElementById('confirmPassword').value = '';
    });
    
    // 新用户名输入框事件监听
    const newUsernameInput = document.getElementById('newUsername');
    newUsernameInput.addEventListener('input', () => {
        // 实时验证用户名长度
        const username = newUsernameInput.value.trim();
        const updateButton = document.getElementById('updateUsernameBtn');
        
        if (username.length >= 3 && username.length <= 20) {
            updateButton.disabled = false;
        } else {
            updateButton.disabled = true;
        }
    });
}

// 处理修改用户名
async function handleUpdateUsername() {
    const newUsername = document.getElementById('newUsername').value.trim();
    
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
            document.getElementById('currentUserName').textContent = newUsername;
            document.getElementById('newUsername').value = '';
            
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
    const password = document.getElementById('confirmPassword').value;
    
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
            document.getElementById('confirmPassword').value = '';
        }
    } catch (error) {
        console.error('注销账户失败:', error);
        showToast('注销账户失败，请稍后重试', 'error');
        document.getElementById('confirmPassword').value = '';
    }
}

// 验证用户名
function validateUsername(username) {
    return username.length >= 3 && username.length <= 20;
}