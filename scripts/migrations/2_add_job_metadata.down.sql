ALTER TABLE jobs DROP COLUMN IF EXISTS priority;
ALTER TABLE jobs DROP COLUMN IF EXISTS partition_key;
ALTER TABLE jobs DROP COLUMN IF EXISTS tenant_id;
DROP INDEX IF EXISTS idx_jobs_partition_key;
DROP INDEX IF EXISTS idx_jobs_tenant_id;
