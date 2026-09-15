DROP INDEX IF EXISTS idx_email_jobs_dequeue;
DROP INDEX IF EXISTS idx_email_jobs_idempotency;
DROP INDEX IF EXISTS idx_email_job_attempts_job_id;

DROP TABLE IF EXISTS email_job_attempts;
DROP TABLE IF EXISTS email_jobs;
