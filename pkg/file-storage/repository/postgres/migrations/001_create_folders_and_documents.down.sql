-- Drop shares table
DROP INDEX IF EXISTS idx_shares_created_at;
DROP INDEX IF EXISTS idx_shares_share_link;
DROP INDEX IF EXISTS idx_shares_shared_with;
DROP INDEX IF EXISTS idx_shares_owner_id;
DROP INDEX IF EXISTS idx_shares_document_id;
DROP TABLE IF EXISTS shares;

-- Drop documents table
DROP INDEX IF EXISTS idx_documents_search;
DROP INDEX IF EXISTS idx_documents_is_deleted;
DROP INDEX IF EXISTS idx_documents_created_at;
DROP INDEX IF EXISTS idx_documents_name;
DROP INDEX IF EXISTS idx_documents_folder_id;
DROP INDEX IF EXISTS idx_documents_user_id;
DROP TABLE IF EXISTS documents;

-- Drop folders table
DROP INDEX IF EXISTS idx_folders_is_deleted;
DROP INDEX IF EXISTS idx_folders_parent_id;
DROP INDEX IF EXISTS idx_folders_user_id;
DROP TABLE IF EXISTS folders;

