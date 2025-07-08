-- Drop indexes
DROP INDEX IF EXISTS idx_audit_logs_timestamp;
DROP INDEX IF EXISTS idx_audit_logs_resource_id;
DROP INDEX IF EXISTS idx_audit_logs_resource_type;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_username;
DROP INDEX IF EXISTS idx_audit_logs_service_name;

-- Drop table
DROP TABLE IF EXISTS audit_logs;