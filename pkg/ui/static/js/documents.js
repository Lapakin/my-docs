// Documents page functionality
document.addEventListener('DOMContentLoaded', () => {
    loadDocuments();
    
    const uploadBtn = document.getElementById('uploadBtn');
    const uploadForm = document.getElementById('uploadForm');
    const searchBtn = document.getElementById('searchBtn');
    const searchInput = document.getElementById('searchInput');
    const shareForm = document.getElementById('shareForm');

    if (uploadBtn) {
        uploadBtn.addEventListener('click', () => {
            document.getElementById('uploadModal').classList.add('show');
        });
    }

    if (uploadForm) {
        uploadForm.addEventListener('submit', handleUpload);
    }

    if (searchBtn) {
        searchBtn.addEventListener('click', handleSearch);
    }

    if (searchInput) {
        searchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') handleSearch();
        });
    }

    if (shareForm) {
        shareForm.addEventListener('submit', handleShare);
    }
});

function closeUploadModal() {
    document.getElementById('uploadModal').classList.remove('show');
    document.getElementById('uploadForm').reset();
}

function closeShareModal() {
    document.getElementById('shareModal').classList.remove('show');
    document.getElementById('shareForm').reset();
}

async function loadDocuments() {
    try {
        const response = await apiRequest('/api/v1/documents');
        if (!response || !response.ok) {
            displayEmptyState();
            return;
        }

        const documents = await response.json();
        displayDocuments(documents);
    } catch (error) {
        console.error('Error loading documents:', error);
        displayEmptyState();
    }
}

function displayEmptyState() {
    const documentsList = document.getElementById('documentsList');
    documentsList.innerHTML = `
        <div class="empty-state">
            <div class="empty-state-icon">📄</div>
            <p>No documents found</p>
        </div>
    `;
}

function displayDocuments(documents) {
    const documentsList = document.getElementById('documentsList');
    
    if (!documents || documents.length === 0) {
        displayEmptyState();
        return;
    }
    
    documentsList.innerHTML = documents.map(doc => `
        <div class="item-card">
            <div class="item-info">
                <div class="item-icon">${getFileIcon(doc.mime_type)}</div>
                <div class="item-details">
                    <div class="item-name">${doc.name || 'Untitled'}</div>
                    <div class="item-meta">
                        <span>${formatFileSize(doc.file_size || 0)}</span>
                        <span>${doc.mime_type || 'Unknown type'}</span>
                        <span>${formatDate(doc.created_at)}</span>
                    </div>
                </div>
            </div>
            <div class="item-actions">
                <button class="btn btn-sm btn-secondary" onclick="downloadDocument('${doc.id}')">Download</button>
                <button class="btn btn-sm btn-success" onclick="openShareModal('${doc.id}')">Share</button>
                <button class="btn btn-sm btn-danger" onclick="deleteDocument('${doc.id}')">Delete</button>
            </div>
        </div>
    `).join('');
}

async function handleUpload(e) {
    e.preventDefault();
    
    const fileInput = document.getElementById('file');
    if (!fileInput.files || fileInput.files.length === 0) {
        showToast('Please select a file', 'error');
        return;
    }

    const formData = new FormData();
    formData.append('file', fileInput.files[0]);

    const fileName = document.getElementById('fileName').value;
    const description = document.getElementById('description').value;

    if (fileName) formData.append('name', fileName);
    if (description) formData.append('description', description);

    try {
        const token = getCookie('token');
        // Use the dedicated upload endpoint that handles multipart form data
        const response = await fetch('/api/documents/upload', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
            body: formData,
        });
        
        if (response.ok) {
            closeUploadModal();
            loadDocuments();
            showToast('Document uploaded successfully', 'success');
        } else {
            const error = await response.json().catch(() => ({ error: 'Upload failed' }));
            showToast('Upload failed: ' + (error.message || error.error || 'Unknown error'), 'error');
        }
    } catch (error) {
        console.error('Error uploading document:', error);
        showToast('Upload failed: ' + error.message, 'error');
    }
}

async function handleSearch() {
    const query = document.getElementById('searchInput').value.trim();

    if (!query) {
        loadDocuments();
        return;
    }

    try {
        const response = await apiRequest('/api/v1/documents');
        if (!response || !response.ok) return;
        
        const documents = await response.json();
        // Client-side filtering
        const filtered = documents.filter(doc =>
            doc.name && doc.name.toLowerCase().includes(query.toLowerCase())
        );
        displayDocuments(filtered);
    } catch (error) {
        console.error('Error searching documents:', error);
    }
}

function downloadDocument(id) {
    // Use the dedicated download endpoint
    window.open(`/api/documents/${id}/download`, '_blank');
}

async function deleteDocument(id) {
    if (!confirm('Are you sure you want to delete this document?')) return;
    
    try {
        const response = await apiRequest(`/api/v1/documents/${id}`, {
            method: 'DELETE',
        });
        
        if (response && response.ok) {
            loadDocuments();
            showToast('Document deleted', 'success');
        } else {
            showToast('Delete failed', 'error');
        }
    } catch (error) {
        console.error('Error deleting document:', error);
        showToast('Delete failed', 'error');
    }
}

function openShareModal(documentId) {
    document.getElementById('shareDocumentId').value = documentId;
    document.getElementById('shareModal').classList.add('show');
}

async function handleShare(e) {
    e.preventDefault();

    const documentId = document.getElementById('shareDocumentId').value;
    const permission = document.getElementById('permission').value;
    const maxAccess = parseInt(document.getElementById('maxAccess').value) || 0;
    const expiresAt = document.getElementById('expiresAt').value;

    const shareData = {
        document_id: documentId,
        permission: permission,
        max_access: maxAccess,
    };

    if (expiresAt) {
        shareData.expires_at = new Date(expiresAt).toISOString();
    }

    try {
        const response = await apiRequest('/api/v1/shares', {
            method: 'POST',
            body: JSON.stringify(shareData),
        });

        if (response && response.ok) {
            const share = await response.json();
            closeShareModal();
            showToast('Share link created!', 'success');

            // Copy link to clipboard
            const shareLink = `${window.location.origin}/share/${share.share_link}`;
            navigator.clipboard.writeText(shareLink).then(() => {
                showToast('Link copied to clipboard', 'success');
            });
        } else {
            const error = await response.json();
            showToast('Failed to create share: ' + (error.message || error.error || 'Unknown error'), 'error');
        }
    } catch (error) {
        console.error('Error creating share:', error);
        showToast('Failed to create share', 'error');
    }
}

