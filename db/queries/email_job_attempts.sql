-- name: RecordAttempt :one
INSERT INTO email_job_attempts (job_id, attempt_number, status, error_message, smtp_response_code)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetAttemptsForJob :many
SELECT * FROM email_job_attempts
WHERE job_id = $1
ORDER BY attempt_number ASC;
