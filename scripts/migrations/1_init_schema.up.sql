-- Create jobs table
CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    queue_name VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(50) NOT NULL,
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    scheduled_at TIMESTAMP,
    processed_at TIMESTAMP,
    failed_at TIMESTAMP,
    idempotency_key VARCHAR(255) UNIQUE
);

CREATE TABLE job_attempts (
    id UUID PRIMARY KEY,
    job_id UUID REFERENCES jobs(id) ON DELETE CASCADE,
    worker_id UUID,
    status VARCHAR(50) NOT NULL,
    error_message TEXT,
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

CREATE TABLE workers (
    id UUID PRIMARY KEY,
    hostname VARCHAR(255),
    status VARCHAR(50),
    last_heartbeat TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE dead_letter_jobs (
    id UUID PRIMARY KEY,
    original_job_id UUID,
    payload JSONB NOT NULL,
    failure_reason TEXT,
    failed_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'USER',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO users (id, username, password_hash, role)
VALUES (
    gen_random_uuid(),
    'admin',
    '$2a$10$zLBKopvun9cN3DlA7gwTaec1v/LyU8W4ARXkXoIlA.NjBe/mZOwy6',
    'ADMIN'
);
