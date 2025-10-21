-- Create folders table
CREATE TABLE IF NOT EXISTS folders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    parent_id BIGINT,
    name VARCHAR(255) NOT NULL,
    path TEXT NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    color VARCHAR(50),
    icon VARCHAR(100),
    created_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP,
    CONSTRAINT fk_parent_folder FOREIGN KEY (parent_id) REFERENCES folders(id) ON DELETE CASCADE
);

-- Create indexes for folders
CREATE INDEX idx_folders_user_id ON folders(user_id);
CREATE INDEX idx_folders_parent_id ON folders(parent_id);
CREATE INDEX idx_folders_is_deleted ON folders(is_deleted);

-- Create documents table
CREATE TABLE IF NOT EXISTS documents (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    folder_id BIGINT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP,
    CONSTRAINT fk_folder FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE SET NULL
);

-- Create indexes for documents
CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_documents_folder_id ON documents(folder_id);
CREATE INDEX idx_documents_name ON documents(name);
CREATE INDEX idx_documents_created_at ON documents(created_at);
CREATE INDEX idx_documents_is_deleted ON documents(is_deleted);

-- Full-text search index for documents
CREATE INDEX idx_documents_search ON documents USING gin(to_tsvector('english', name || ' ' || COALESCE(description, '')));

-- Create shares table
CREATE TABLE IF NOT EXISTS shares (
    id BIGSERIAL PRIMARY KEY,
    document_id BIGINT NOT NULL,
    owner_id BIGINT NOT NULL,
    shared_with BIGINT,
    share_link VARCHAR(255) UNIQUE NOT NULL,
    permission VARCHAR(20) NOT NULL DEFAULT 'view',
    expires_at TIMESTAMP,
    access_count INT NOT NULL DEFAULT 0,
    max_access INT NOT NULL DEFAULT -1,
    password VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_document FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE,
    CHECK (permission IN ('view', 'download', 'edit'))
);

-- Create indexes for shares
CREATE INDEX idx_shares_document_id ON shares(document_id);
CREATE INDEX idx_shares_owner_id ON shares(owner_id);
CREATE INDEX idx_shares_shared_with ON shares(shared_with);
CREATE INDEX idx_shares_share_link ON shares(share_link);
CREATE INDEX idx_shares_created_at ON shares(created_at);

