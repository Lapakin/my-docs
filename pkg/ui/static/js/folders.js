// Folders page functionality
document.addEventListener('DOMContentLoaded', () => {
    loadFolders();
    
    const createFolderBtn = document.getElementById('createFolderBtn');
    const createFolderForm = document.getElementById('createFolderForm');
    
    if (createFolderBtn) {
        createFolderBtn.addEventListener('click', () => {
            document.getElementById('createFolderModal').classList.add('show');
        });
    }

    if (createFolderForm) {
        createFolderForm.addEventListener('submit', handleCreateFolder);
    }
});

function closeFolderModal() {
    document.getElementById('createFolderModal').classList.remove('show');
    document.getElementById('createFolderForm').reset();
}

async function loadFolders() {
    try {
        const response = await apiRequest('/api/v1/folders');
        if (!response || !response.ok) {
            displayEmptyState();
            return;
        }

        const folders = await response.json();
        displayFolders(folders);
    } catch (error) {
        console.error('Error loading folders:', error);
        displayEmptyState();
    }
}

function displayEmptyState() {
    const foldersList = document.getElementById('foldersList');
    foldersList.innerHTML = `
        <div class="empty-state">
            <div class="empty-state-icon">📁</div>
            <p>No folders found</p>
        </div>
    `;
}

function displayFolders(folders) {
    const foldersList = document.getElementById('foldersList');
    
    if (!folders || folders.length === 0) {
        displayEmptyState();
        return;
    }
    
    foldersList.innerHTML = folders.map(folder => `
        <div class="item-card">
            <div class="item-info">
                <div class="item-icon" style="background: ${folder.color ? folder.color + '20' : '#eff6ff'};">
                    📁
                </div>
                <div class="item-details">
                    <div class="item-name">${folder.name}</div>
                    <div class="item-meta">
                        <span>${folder.is_public ? '🌐 Public' : '🔒 Private'}</span>
                        <span>Created: ${formatDate(folder.created_at)}</span>
                    </div>
                </div>
            </div>
            <div class="item-actions">
                <button class="btn btn-sm btn-danger" onclick="deleteFolder('${folder.id}')">Delete</button>
            </div>
        </div>
    `).join('');
}

async function handleCreateFolder(e) {
    e.preventDefault();
    
    const name = document.getElementById('folderName').value.trim();

    if (!name) {
        showToast('Please enter a folder name', 'error');
        return;
    }

    try {
        const response = await apiRequest('/api/v1/folders', {
            method: 'POST',
            body: JSON.stringify({ name }),
        });
        
        if (response && response.ok) {
            closeFolderModal();
            loadFolders();
            showToast('Folder created successfully', 'success');
        } else {
            const error = await response.json();
            showToast('Failed to create folder: ' + (error.message || error.error || 'Unknown error'), 'error');
        }
    } catch (error) {
        console.error('Error creating folder:', error);
        showToast('Failed to create folder', 'error');
    }
}

async function deleteFolder(id) {
    if (!confirm('Are you sure you want to delete this folder?')) return;
    
    try {
        const response = await apiRequest(`/api/v1/folders/${id}`, {
            method: 'DELETE',
        });
        
        if (response && response.ok) {
            loadFolders();
            showToast('Folder deleted', 'success');
        } else {
            showToast('Failed to delete folder', 'error');
        }
    } catch (error) {
        console.error('Error deleting folder:', error);
        showToast('Failed to delete folder', 'error');
    }
}

