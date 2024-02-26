package transcoder

import (
	"context"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"github.com/RacoonMediaServer/rms-transcoder/internal/model"
	"github.com/RacoonMediaServer/rms-transcoder/internal/worker"
	"go-micro.dev/v4/logger"
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

func (s *Service) getAndUpdateJob(id string) (*jobRecord, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.jobs[id]
	if !ok {
		return nil, false
	}

	isCancelled := false
	status, _ := record.receipt.Status()
	switch status {
	case worker.Done:
		record.job.Result = model.JobDone
	case worker.Failed:
		record.job.Result = model.JobFailed
	case worker.Cancelled:
		isCancelled = true
	default:
	}

	if record.job.AutoComplete {
		delete(s.jobs, id)
	}

	return record, isCancelled
}

func (s *Service) processReadyJob(id string) {
	record, isCancelled := s.getAndUpdateJob(id)
	if record == nil || isCancelled {
		return
	}

	if record.job.AutoComplete {
		if err := s.Database.RemoveJob(id); err != nil {
			s.l.Logf(logger.ErrorLevel, "Remove job %s failed: %s", id, err)
		}

	} else {
		if err := s.Database.UpdateJob(record.job); err != nil {
			s.l.Logf(logger.ErrorLevel, "Update job %s failed: %s", id, err)
		}
	}

}

func (s *Service) sendNotification(record *jobRecord) {
	kind := events.Notification_TranscodingDone
	if record.job.Result != model.JobDone {
		kind = events.Notification_TranscodingFailed
	}
	notification := events.Notification{
		Sender:        "rms-transcoder",
		Kind:          kind,
		TorrentID:     nil,
		MediaID:       nil,
		ItemTitle:     nil,
		VideoLocation: nil,
		SizeMB:        nil,
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.TODO(), 20*time.Second)
		defer cancel()
		if err := s.Publisher.Publish(ctx, &notification); err != nil {
			s.l.Logf(logger.ErrorLevel, "Send notification failed: %s", err)
		}
	}()
}
