// Shares page functionality
document.addEventListener('DOMContentLoaded', () => {
    loadShares();
});

async function loadShares() {
    try {
        const response = await apiRequest('/api/v1/shares');
        if (!response || !response.ok) {
            displayEmptyState();
            return;
        }

        const shares = await response.json();
        displayShares(shares);
    } catch (error) {
        console.error('Error loading shares:', error);
        displayEmptyState();
    }
}

function displayEmptyState() {
    const sharesList = document.getElementById('sharesList');
    sharesList.innerHTML = `
        <div class="empty-state">
            <div class="empty-state-icon">🔗</div>
            <p>No shared files</p>
        </div>
    `;
}

function displayShares(shares) {
    const sharesList = document.getElementById('sharesList');

    if (!shares || shares.length === 0) {
        displayEmptyState();
        return;
    }

    sharesList.innerHTML = shares.map(share => {
        const shareLink = `${window.location.origin}/share/${share.share_link}`;
        const permissionBadge = getPermissionBadge(share.permission);

        return `
            <div class="item-card">
                <div class="item-info">
                    <div class="item-icon">🔗</div>
                    <div class="item-details">
                        <div class="item-name">Document: ${share.document_id}</div>
                        <div class="item-meta">
                            <span>${permissionBadge}</span>
                            <span>Views: ${share.access_count}${share.max_access > 0 ? '/' + share.max_access : ''}</span>
                            <span>${share.expires_at ? 'Expires: ' + formatDate(share.expires_at) : 'No expiry'}</span>
                        </div>
                    </div>
                </div>
                <div class="item-actions">
                    <button class="btn btn-sm btn-secondary" onclick="copyShareLink('${shareLink}')">Copy Link</button>
                    <button class="btn btn-sm btn-danger" onclick="deleteShare('${share.id}')">Delete</button>
                </div>
            </div>
        `;
    }).join('');
}

function getPermissionBadge(permission) {
    switch (permission) {
        case 'view': return '👁️ View';
        case 'download': return '⬇️ Download';
        case 'edit': return '✏️ Edit';
        default: return permission;
    }
}

function copyShareLink(link) {
    navigator.clipboard.writeText(link).then(() => {
        showToast('Link copied to clipboard!', 'success');
    }).catch(err => {
        console.error('Failed to copy:', err);
        showToast('Failed to copy link', 'error');
    });
}

async function deleteShare(id) {
    if (!confirm('Are you sure you want to delete this share?')) return;

    try {
        const response = await apiRequest(`/api/v1/shares/${id}`, {
            method: 'DELETE',
        });

        if (response && (response.ok || response.status === 204)) {
            loadShares();
            showToast('Share deleted', 'success');
        } else {
            showToast('Failed to delete share', 'error');
        }
    } catch (error) {
        console.error('Error deleting share:', error);
        showToast('Failed to delete share', 'error');
    }
}

