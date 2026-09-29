-- Security initialization script for CodexGigantus
-- This script creates least-privilege database users

-- Read the passwords from Docker secrets
\set admin_password `cat /run/secrets/db_admin_password 2>/dev/null || echo 'changeme_admin'`
\set app_password `cat /run/secrets/db_password 2>/dev/null || echo 'changeme_app'`

-- Set admin password (if using default postgres user)
ALTER USER postgres WITH PASSWORD :'admin_password';

-- Create read-only application user with least privileges
CREATE USER codex_readonly WITH PASSWORD :'app_password';

-- Revoke all default privileges
REVOKE ALL ON DATABASE codex FROM codex_readonly;
REVOKE ALL ON SCHEMA public FROM codex_readonly;
REVOKE ALL ON ALL TABLES IN SCHEMA public FROM codex_readonly;
REVOKE ALL ON ALL SEQUENCES IN SCHEMA public FROM codex_readonly;
REVOKE ALL ON ALL FUNCTIONS IN SCHEMA public FROM codex_readonly;

-- Grant minimal required privileges for read-only access
GRANT CONNECT ON DATABASE codex TO codex_readonly;
GRANT USAGE ON SCHEMA public TO codex_readonly;
GRANT SELECT ON code_files TO codex_readonly;

-- Ensure future tables also have restricted access
ALTER DEFAULT PRIVILEGES IN SCHEMA public
    GRANT SELECT ON TABLES TO codex_readonly;

-- Revoke dangerous privileges from public
REVOKE CREATE ON SCHEMA public FROM PUBLIC;

-- Security: Disable connection from application user to other databases
REVOKE ALL ON DATABASE postgres FROM codex_readonly;
REVOKE ALL ON DATABASE template0 FROM codex_readonly;
REVOKE ALL ON DATABASE template1 FROM codex_readonly;

-- Create role for monitoring (optional, for DBAs)
CREATE ROLE codex_monitor WITH LOGIN PASSWORD :'admin_password';
GRANT pg_monitor TO codex_monitor;
GRANT CONNECT ON DATABASE codex TO codex_monitor;

-- Log privileges
DO $$
BEGIN
    RAISE NOTICE 'Security setup complete:';
    RAISE NOTICE '  - Admin user: postgres (full access)';
    RAISE NOTICE '  - Application user: codex_readonly (SELECT only on code_files)';
    RAISE NOTICE '  - Monitoring user: codex_monitor (monitoring queries only)';
    RAISE NOTICE 'IMPORTANT: Change default passwords in production!';
END $$;

-- Create audit table for security logging (optional)
CREATE TABLE IF NOT EXISTS audit_log (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    table_name VARCHAR(100),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    details JSONB
);

-- Grant audit log access to admin only
GRANT ALL ON audit_log TO postgres;
GRANT SELECT ON audit_log TO codex_monitor;

-- Prevent privilege escalation
REVOKE ALL ON audit_log FROM PUBLIC;
REVOKE ALL ON audit_log FROM codex_readonly;
