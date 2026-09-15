CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- email_jobs
CREATE TABLE email_jobs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key TEXT UNIQUE,
    from_email      TEXT NOT NULL,
    to_email        TEXT NOT NULL,
    subject         TEXT NOT NULL,
    body            TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending',
    priority        SMALLINT NOT NULL DEFAULT 0,
    attempt_count   INT NOT NULL DEFAULT 0,
    max_attempts    INT NOT NULL DEFAULT 5,
    next_retry_at   TIMESTAMPTZ,
    locked_at       TIMESTAMPTZ,
    last_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    sent_at         TIMESTAMPTZ
);

CREATE INDEX idx_email_jobs_dequeue ON email_jobs (status, priority, next_retry_at);
CREATE INDEX idx_email_jobs_idempotency ON email_jobs (idempotency_key) WHERE idempotency_key IS NOT NULL;

-- email_job_attempts
CREATE TABLE email_job_attempts (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id             UUID NOT NULL REFERENCES email_jobs(id) ON DELETE CASCADE,
    attempt_number     INT NOT NULL,
    status             TEXT NOT NULL,
    error_message      TEXT,
    smtp_response_code TEXT,
    attempted_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_email_job_attempts_job_id ON email_job_attempts (job_id);
