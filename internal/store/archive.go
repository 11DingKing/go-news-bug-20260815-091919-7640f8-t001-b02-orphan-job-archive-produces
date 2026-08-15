package store

import (
	"time"

	"example.com/ledgerd/internal/model"
)

// Archive atomically moves an active job into the archive area and attaches
// retention metadata. The migration must leave every archived job — whether or
// not it ever received a Finished event — carrying a complete RetentionInfo so
// the reaper and query paths can dereference it safely.
func (s *ledgerStore) Archive(jobID string, activeWindow, retentionWindow time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[jobID]
	if !ok {
		return ErrJobNotFound
	}
	if job.Status != model.StatusActive {
		return ErrNotActive
	}

	now := time.Now()
	job.Status = model.StatusArchived
	job.Retention = &model.RetentionInfo{
		ActiveWindow:    activeWindow,
		RetentionWindow: retentionWindow,
		ArchivedAt:      now,
		ExpiresAt:       now.Add(retentionWindow),
	}
	return nil
}
