// File Storage UI - Main JavaScript

// API Base URL
const API_BASE = '/api';

// Show alert message
function showAlert(message, type = 'info') {
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type}`;
    alertDiv.textContent = message;

    const container = document.querySelector('.container');
    if (container) {
        container.insertBefore(alertDiv, container.firstChild);
        setTimeout(() => alertDiv.remove(), 5000);
    }
}

// Get token from cookie
function getToken() {
    const cookies = document.cookie.split(';');
    for (let cookie of cookies) {
        const [name, value] = cookie.trim().split('=');
        if (name === 'token') {
            return value;
        }
    }
    return null;
}

// Make authenticated API call
async function apiCall(endpoint, options = {}) {
    const token = getToken();

    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json',
            ...(token && { 'Authorization': `Bearer ${token}` })
        }
    };

    const response = await fetch(`${API_BASE}${endpoint}`, {
        ...defaultOptions,
        ...options,
        headers: { ...defaultOptions.headers, ...options.headers }
    });

    if (!response.ok) {
        const error = await response.json().catch(() => ({ error: 'Request failed' }));
        throw new Error(error.error || 'Request failed');
    }

    return response.json();
}

// Login function
async function login(email, password) {
    try {
        const response = await apiCall('/login', {
            method: 'POST',
            body: JSON.stringify({ email, password })
        });

        if (response.token) {
            showAlert('Login successful!', 'success');
            setTimeout(() => window.location.href = '/dashboard', 1000);
        }
    } catch (error) {
        showAlert(error.message, 'error');
    }
}

// Register function
async function register(username, email, password) {
    try {
        const response = await apiCall('/register', {
            method: 'POST',
            body: JSON.stringify({ username, email, password })
        });

        showAlert('Registration successful! Please login.', 'success');
        setTimeout(() => window.location.href = '/login', 2000);
    } catch (error) {
        showAlert(error.message, 'error');
    }
}

// Logout function
function logout() {
    document.cookie = 'token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    window.location.href = '/';
}

// Load users
async function loadUsers() {
    try {
        const response = await apiCall('/users');
        return response.users || [];
    } catch (error) {
        showAlert('Failed to load users', 'error');
        return [];
    }
}

// Load documents
async function loadDocuments(folderId = null) {
    try {
        const endpoint = folderId ? `/documents?folder_id=${folderId}` : '/documents';
        const documents = await apiCall(endpoint);
        return documents || [];
    } catch (error) {
        showAlert('Failed to load documents', 'error');
        return [];
    }
}

// Load folders
async function loadFolders(parentId = null) {
    try {
        const endpoint = parentId ? `/folders?parent_id=${parentId}` : '/folders';
        const folders = await apiCall(endpoint);
        return folders || [];
    } catch (error) {
        showAlert('Failed to load folders', 'error');
        return [];
    }
}

// Delete document
async function deleteDocument(documentId) {
    if (!confirm('Are you sure you want to delete this document?')) {
        return;
    }

    try {
        await apiCall(`/documents/${documentId}`, { method: 'DELETE' });
        showAlert('Document deleted successfully', 'success');
        return true;
    } catch (error) {
        showAlert('Failed to delete document', 'error');
        return false;
    }
}

// Delete folder
async function deleteFolder(folderId) {
    if (!confirm('Are you sure you want to delete this folder?')) {
        return;
    }

    try {
        await apiCall(`/folders/${folderId}`, { method: 'DELETE' });
        showAlert('Folder deleted successfully', 'success');
        return true;
    } catch (error) {
        showAlert('Failed to delete folder', 'error');
        return false;
    }
}

// Format file size
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

// Format date
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
}

// Initialize page
document.addEventListener('DOMContentLoaded', () => {
    // Attach login form handler
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            await login(email, password);
        });
    }

    // Attach register form handler
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = document.getElementById('username').value;
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            await register(username, email, password);
        });
    }

    // Attach logout button handler
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', (e) => {
            e.preventDefault();
            logout();
        });
    }
});

