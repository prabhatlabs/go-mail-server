package jobs

import (
	"encoding/json/v2"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/prabhatlabs/go-mail-server/internal/db"
	"github.com/prabhatlabs/go-mail-server/internal/lib/response"
)

type enqueueRequest struct {
	IdempotencyKey *string `json:"idempotency_key,omitempty"`
	FromEmail      string  `json:"from_email"`
	ToEmail        string  `json:"to_email"`
	Subject        string  `json:"subject"`
	Body           string  `json:"body"`
	Priority       int16   `json:"priority"`
	MaxAttempts    int32   `json:"max_attempts"`
}

type jobResponse struct {
	ID             string  `json:"id"`
	IdempotencyKey *string `json:"idempotency_key,omitempty"`
	FromEmail      string  `json:"from_email"`
	ToEmail        string  `json:"to_email"`
	Subject        string  `json:"subject"`
	Body           string  `json:"body"`
	Status         string  `json:"status"`
	Priority       int16   `json:"priority"`
	AttemptCount   int32   `json:"attempt_count"`
	MaxAttempts    int32   `json:"max_attempts"`
	NextRetryAt    *string `json:"next_retry_at,omitempty"`
	LockedAt       *string `json:"locked_at,omitempty"`
	LastError      *string `json:"last_error,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	SentAt         *string `json:"sent_at,omitempty"`
	Attempts       []attemptResponse `json:"attempts,omitempty"`
}

type attemptResponse struct {
	ID               string  `json:"id"`
	JobID            string  `json:"job_id"`
	AttemptNumber    int32   `json:"attempt_number"`
	Status           string  `json:"status"`
	ErrorMessage     *string `json:"error_message,omitempty"`
	SmtpResponseCode *string `json:"smtp_response_code,omitempty"`
	AttemptedAt      string  `json:"attempted_at"`
}

type listResponse struct {
	Data   []jobResponse `json:"data"`
	Total  int64         `json:"total"`
	Limit  int32         `json:"limit"`
	Offset int32         `json:"offset"`
}

func toJobResponse(job db.EmailJob, attempts []db.EmailJobAttempt) jobResponse {
	resp := jobResponse{
		ID:           job.ID.String(),
		FromEmail:    job.FromEmail,
		ToEmail:      job.ToEmail,
		Subject:      job.Subject,
		Body:         job.Body,
		Status:       job.Status,
		Priority:     job.Priority,
		AttemptCount: job.AttemptCount,
		MaxAttempts:  job.MaxAttempts,
		CreatedAt:    job.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    job.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
	if job.IdempotencyKey.Valid {
		resp.IdempotencyKey = &job.IdempotencyKey.String
	}
	if job.NextRetryAt.Valid {
		s := job.NextRetryAt.Time.Format("2006-01-02T15:04:05Z07:00")
		resp.NextRetryAt = &s
	}
	if job.LockedAt.Valid {
		s := job.LockedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		resp.LockedAt = &s
	}
	if job.LastError.Valid {
		resp.LastError = &job.LastError.String
	}
	if job.SentAt.Valid {
		s := job.SentAt.Time.Format("2006-01-02T15:04:05Z07:00")
		resp.SentAt = &s
	}
	if len(attempts) > 0 {
		resp.Attempts = make([]attemptResponse, len(attempts))
		for i, a := range attempts {
			resp.Attempts[i] = attemptResponse{
				ID:            a.ID.String(),
				JobID:         a.JobID.String(),
				AttemptNumber: a.AttemptNumber,
				Status:        a.Status,
				AttemptedAt:   a.AttemptedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			}
			if a.ErrorMessage.Valid {
				resp.Attempts[i].ErrorMessage = &a.ErrorMessage.String
			}
			if a.SmtpResponseCode.Valid {
				resp.Attempts[i].SmtpResponseCode = &a.SmtpResponseCode.String
			}
		}
	}
	return resp
}

func parseUUIDParam(r *http.Request, key string) (pgtype.UUID, error) {
	str := chi.URLParam(r, key)
	var uuid pgtype.UUID
	err := uuid.Scan(str)
	return uuid, err
}

func (s *Service) listJobs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit := int32(20)
	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 100 {
			response.BadRequest(w, "limit must be between 1 and 100")
			return
		}
		limit = int32(n)
	}

	offset := int32(0)
	if v := q.Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			response.BadRequest(w, "offset must be >= 0")
			return
		}
		offset = int32(n)
	}

	status := q.Get("status")

	jobs, err := s.db.Q.ListJobs(r.Context(), db.ListJobsParams{
		Column1: status,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		response.InternalServerError(w, err, "failed to list jobs")
		return
	}

	total, err := s.db.Q.CountJobs(r.Context(), status)
	if err != nil {
		response.InternalServerError(w, err, "failed to count jobs")
		return
	}

	data := make([]jobResponse, len(jobs))
	for i, job := range jobs {
		data[i] = toJobResponse(job, nil)
	}

	response.OK(w, "jobs retrieved", listResponse{
		Data:   data,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

func (s *Service) enqueue(w http.ResponseWriter, r *http.Request) {
	var req enqueueRequest
	if err := json.UnmarshalRead(r.Body, &req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.FromEmail == "" || req.ToEmail == "" || req.Subject == "" || req.Body == "" {
		response.BadRequest(w, "from_email, to_email, subject, and body are required")
		return
	}

	if req.MaxAttempts == 0 {
		req.MaxAttempts = 5
	}

	var idempotencyKey pgtype.Text
	if req.IdempotencyKey != nil {
		idempotencyKey = pgtype.Text{String: *req.IdempotencyKey, Valid: true}
		existingID, err := s.db.Q.CheckIdempotency(r.Context(), idempotencyKey)
		if err == nil {
			response.Conflict(w, "job with this idempotency key already exists: "+existingID.String())
			return
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			response.InternalServerError(w, err, "failed to check idempotency")
			return
		}
	}

	id, err := s.db.Q.EnqueueJob(r.Context(), db.EnqueueJobParams{
		IdempotencyKey: idempotencyKey,
		FromEmail:      req.FromEmail,
		ToEmail:        req.ToEmail,
		Subject:        req.Subject,
		Body:           req.Body,
		Priority:       req.Priority,
		MaxAttempts:    req.MaxAttempts,
	})
	if err != nil {
		response.InternalServerError(w, err, "failed to enqueue job")
		return
	}

	response.Created(w, "job enqueued", map[string]string{"id": id.String()})
}

func (s *Service) getJob(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid job id")
		return
	}

	job, err := s.db.Q.GetJob(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.NotFound(w, "job not found")
			return
		}
		response.InternalServerError(w, err, "failed to get job")
		return
	}

	attempts, err := s.db.Q.GetAttemptsForJob(r.Context(), id)
	if err != nil {
		response.InternalServerError(w, err, "failed to get attempts")
		return
	}

	response.OK(w, "job retrieved", toJobResponse(job, attempts))
}

func (s *Service) getJobStatus(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid job id")
		return
	}

	job, err := s.db.Q.GetJob(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.NotFound(w, "job not found")
			return
		}
		response.InternalServerError(w, err, "failed to get job")
		return
	}

	attempts, err := s.db.Q.GetAttemptsForJob(r.Context(), id)
	if err != nil {
		response.InternalServerError(w, err, "failed to get attempts")
		return
	}

	resp := toJobResponse(job, attempts)
	response.OK(w, "job status retrieved", map[string]any{
		"status":   resp.Status,
		"attempts": resp.Attempts,
	})
}
