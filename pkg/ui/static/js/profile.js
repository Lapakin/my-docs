async function handleChangePassword(e) {
    console.log('Password change initiated...');
    if (e) {
        e.preventDefault();
        e.stopPropagation();
    }

    const currentPasswordInput = document.getElementById('currentPassword');
    const newPasswordInput = document.getElementById('newPassword');
    const confirmPasswordInput = document.getElementById('confirmPassword');

    if (!currentPasswordInput || !newPasswordInput || !confirmPasswordInput) {
        console.error('Core password inputs missing from DOM');
        return;
    }

    const currentPassword = currentPasswordInput.value;
    const newPassword = newPasswordInput.value;
    const confirmPassword = confirmPasswordInput.value;

    if (!currentPassword || !newPassword || !confirmPassword) {
        showToast('Please fill in all password fields', 'error');
        return;
    }

    if (newPassword !== confirmPassword) {
        showToast('The new passwords do not match', 'error');
        return;
    }

    const token = getCookie('token');
    const claims = token ? parseJwt(token) : null;
    if (!claims || !claims.user_id) {
        console.warn('Session invalid or expired during password change');
        window.location.href = '/login';
        return;
    }

    try {
        console.log('Sending secure password update request...');
        const response = await apiRequest(`/api/v1/users/${claims.user_id}/changePassword`, {
            method: 'POST',
            noRedirect: true,
            body: JSON.stringify({
                old_password: currentPassword,
                new_password: newPassword,
                confirm_password: confirmPassword
            })
        });

        if (!response) {
            console.error('Server connection failed');
            return;
        }

        console.log(`Update status: ${response.status}`);

        if (response.ok || response.status === 200 || response.status === 204) {
            showToast('Your password has been changed successfully', 'success');
            const form = document.getElementById('changePasswordForm');
            if (form) form.reset();
        } else {
            const error = await response.json().catch(() => ({}));
            console.log('Error payload:', error);
            showToast(error.message || error.error || 'Password update failed. Please check your current password.', 'error');
        }
    } catch (error) {
        console.error('Password change error:', error);
        showToast('A network error occurred. Please check your connection and try again.', 'error');
    }
}

document.addEventListener('DOMContentLoaded', () => {
    console.log('Profile page initialized');
    loadProfile();

    const changeBtn = document.getElementById('submitChangePassword');
    if (changeBtn) {
        console.log('Password change listener attached');
        changeBtn.addEventListener('click', handleChangePassword);
    } else {
        console.error('Password change button (#submitChangePassword) not found');
    }
});

function parseJwt(token) {
    try {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function (c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));

        return JSON.parse(jsonPayload);
    } catch (e) {
        return null;
    }
}

async function loadProfile() {
    console.log('Fetching user profile data...');
    const token = getCookie('token');
    if (!token) {
        window.location.href = '/login';
        return;
    }

    const claims = parseJwt(token);
    if (!claims || !claims.user_id) {
        showToast('Session error', 'error');
        return;
    }

    try {
        const response = await apiRequest(`/api/v1/users/${claims.user_id}`);
        if (!response || !response.ok) {
            throw new Error('Profile fetch failed');
        }

        const user = await response.json();

        const profileInfo = document.getElementById('profileInfo');
        if (profileInfo) {
            profileInfo.innerHTML = `
                <div class="info-group">
                    <label>Username:</label>
                    <span>${user.username}</span>
                </div>
                <div class="info-group">
                    <label>Email:</label>
                    <span>${user.email}</span>
                </div>
                <div class="info-group">
                    <label>Role:</label>
                    <span class="badge badge-${user.role === 'admin' ? 'primary' : 'secondary'}">${user.role}</span>
                </div>
                <div class="info-group">
                    <label>Joined:</label>
                    <span>${formatDate(user.created_at)}</span>
                </div>
            `;
        }
    } catch (error) {
        console.error('Profile loading error:', error);
        const profileInfo = document.getElementById('profileInfo');
        if (profileInfo) {
            profileInfo.innerHTML = '<p class="error">Failed to load profile details</p>';
        }
    }
}


