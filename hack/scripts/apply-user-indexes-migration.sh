#!/bin/bash
# Script to apply user indexes migration to the user-management database

echo "Applying migration to add indexes to users tables..."

docker exec user-management-postgres psql -U postgres -d user_management << 'EOF'

-- Add indexes for users table to improve query performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_is_deleted ON users(is_deleted);
CREATE INDEX IF NOT EXISTS idx_users_active_not_deleted ON users(is_active, is_deleted) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Add indexes for user_passwords table
CREATE INDEX IF NOT EXISTS idx_user_passwords_user_id ON user_passwords(user_id);

-- Add indexes for refresh_tokens table
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked ON refresh_tokens(revoked);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_active ON refresh_tokens(user_id, revoked, expires_at) WHERE revoked = FALSE;

-- Show all indexes
SELECT
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
  AND tablename IN ('users', 'user_passwords', 'refresh_tokens')
ORDER BY tablename, indexname;

EOF

echo "Migration complete!"

