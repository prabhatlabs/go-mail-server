-- name: EnqueueJob :one
INSERT INTO email_jobs (idempotency_key, from_email, to_email, subject, body, priority, max_attempts)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: GetJob :one
SELECT * FROM email_jobs WHERE id = $1;

-- name: CheckIdempotency :one
SELECT id FROM email_jobs WHERE idempotency_key = $1 LIMIT 1;

-- name: DequeueJob :one
UPDATE email_jobs
SET status = 'processing', locked_at = now(), updated_at = now()
WHERE id = (
    SELECT id FROM email_jobs
    WHERE status = 'pending'
      AND (next_retry_at IS NULL OR next_retry_at <= now())
    ORDER BY priority DESC, created_at ASC
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
RETURNING *;

-- name: MarkSent :exec
UPDATE email_jobs
SET status = 'sent', sent_at = now(), updated_at = now()
WHERE id = $1;

-- name: MarkFailed :exec
UPDATE email_jobs
SET status = 'failed',
    attempt_count = attempt_count + 1,
    last_error = $2,
    next_retry_at = $3,
    locked_at = NULL,
    updated_at = now()
WHERE id = $1;

-- name: MarkDeadLetter :exec
UPDATE email_jobs
SET status = 'dead_letter', updated_at = now()
WHERE id = $1;

-- name: FindStuckJobs :many
SELECT * FROM email_jobs
WHERE status = 'processing'
  AND locked_at < now() - INTERVAL '10 minutes'
ORDER BY locked_at ASC
LIMIT $1;

-- name: ListJobs :many
SELECT * FROM email_jobs
WHERE ($1::text IS NULL OR status = $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountJobs :one
SELECT COUNT(*) FROM email_jobs
WHERE ($1::text IS NULL OR status = $1);
