// Dashboard functionality
document.addEventListener('DOMContentLoaded', async () => {
    await loadDashboardStats();
    await loadRecentFiles();
});

async function loadDashboardStats() {
    try {
        // Load documents count
        const docsResponse = await apiRequest('/api/v1/documents');
        if (docsResponse && docsResponse.ok) {
            const docs = await docsResponse.json();
            document.getElementById('totalDocs').textContent = docs?.length || 0;

            // Calculate total storage used
            let totalSize = 0;
            if (docs && docs.length > 0) {
                totalSize = docs.reduce((sum, doc) => sum + (doc.file_size || 0), 0);
            }
            document.getElementById('storageUsed').textContent = formatFileSize(totalSize);
        } else {
            document.getElementById('totalDocs').textContent = '0';
            document.getElementById('storageUsed').textContent = '0 B';
        }
        
        // Load folders count
        const foldersResponse = await apiRequest('/api/v1/folders');
        if (foldersResponse && foldersResponse.ok) {
            const folders = await foldersResponse.json();
            document.getElementById('totalFolders').textContent = folders?.length || 0;
        } else {
            document.getElementById('totalFolders').textContent = '0';
        }
        
        // Load shares count
        const sharesResponse = await apiRequest('/api/v1/shares');
        if (sharesResponse && sharesResponse.ok) {
            const shares = await sharesResponse.json();
            document.getElementById('totalShares').textContent = shares?.length || 0;
        } else {
            document.getElementById('totalShares').textContent = '0';
        }
    } catch (error) {
        console.error('Error loading dashboard stats:', error);
        document.getElementById('totalDocs').textContent = '-';
        document.getElementById('totalFolders').textContent = '-';
        document.getElementById('totalShares').textContent = '-';
        document.getElementById('storageUsed').textContent = '-';
    }
}

async function loadRecentFiles() {
    const recentFilesList = document.getElementById('recentFilesList');

    try {
        const response = await apiRequest('/api/v1/documents');
        if (!response || !response.ok) {
            recentFilesList.innerHTML = `
                <div class="empty-state">
                    <div class="empty-state-icon">📄</div>
                    <p>No recent files</p>
                </div>
            `;
            return;
        }

        const documents = await response.json();

        if (!documents || documents.length === 0) {
            recentFilesList.innerHTML = `
                <div class="empty-state">
                    <div class="empty-state-icon">📄</div>
                    <p>No recent files</p>
                </div>
            `;
            return;
        }
        
        // Sort by created_at and show only last 5 files
        const recent = documents
            .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
            .slice(0, 5);

        recentFilesList.innerHTML = recent.map(doc => `
            <div class="item-card">
                <div class="item-info">
                    <div class="item-icon">${getFileIcon(doc.mime_type)}</div>
                    <div class="item-details">
                        <div class="item-name">${doc.name || 'Untitled'}</div>
                        <div class="item-meta">
                            <span>${formatFileSize(doc.file_size || 0)}</span>
                            <span>${formatDate(doc.created_at)}</span>
                        </div>
                    </div>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Error loading recent files:', error);
        recentFilesList.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">📄</div>
                <p>Failed to load recent files</p>
            </div>
        `;
    }
}

