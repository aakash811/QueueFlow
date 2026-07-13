ALTER TABLE jobs
ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS partition_key VARCHAR(255),
ADD COLUMN IF NOT EXISTS priority SMALLINT DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_jobs_tenant_id ON jobs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_jobs_partition_key ON jobs(partition_key);
