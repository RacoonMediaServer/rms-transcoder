package transcoder

import (
	"context"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"github.com/RacoonMediaServer/rms-packages/pkg/worker"
	"go-micro.dev/v4/logger"
	"path/filepath"
	"time"
)

func (s *Service) processReadyJobs() {
	readyChan := s.Workers.DoneChannel()
	for {
		select {
		case taskId := <-readyChan:
			s.processReadyJob(taskId)
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Service) getAndUpdateJob(id string) *jobRecord {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.jobs[id]
	if !ok {
		return nil
	}

	status := record.receipt.Status()
	record.job.Done = status == worker.Done

	if record.job.AutoComplete {
		delete(s.jobs, id)
	}

	return record
}

func (s *Service) processReadyJob(id string) {
	record := s.getAndUpdateJob(id)
	if record == nil {
		return
	}

	if record.job.AutoComplete {
		if err := s.Database.RemoveJob(id); err != nil {
			s.l.Logf(logger.ErrorLevel, "Remove job %s failed: %s", id, err)
		}
		s.sendNotification(record)
	} else {
		if err := s.Database.UpdateJob(record.job); err != nil {
			s.l.Logf(logger.ErrorLevel, "Update job %s failed: %s", id, err)
		}
	}

}

func (s *Service) sendNotification(record *jobRecord) {
	kind := events.Notification_TranscodingDone
	if record.receipt.Status() != worker.Done {
		kind = events.Notification_TranscodingFailed
	}

	fileSize := uint32(getFileSize(filepath.Join(s.Config.Directory, record.job.Destination)) / 1024)
	notification := events.Notification{
		Sender:        "rms-transcoder",
		Kind:          kind,
		ItemTitle:     &record.job.Source,
		VideoLocation: &record.job.Destination,
		SizeMB:        &fileSize,
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ctx, cancel := context.WithTimeout(s.ctx, 20*time.Second)
		defer cancel()
		if err := s.Publisher.Publish(ctx, &notification); err != nil {
			s.l.Logf(logger.ErrorLevel, "Send notification failed: %s", err)
		}
	}()
}
