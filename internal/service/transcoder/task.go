package transcoder

import (
	"context"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"go-micro.dev/v4/logger"
	"time"
)

type transcodingTask struct {
	l           logger.Logger
	id          string
	profile     *rms_transcoder.Profile
	source      string
	destination string
}

func (t *transcodingTask) ID() string {
	return t.id
}

func (t *transcodingTask) Do(ctx context.Context) error {
	for i := 0; i < 30; i++ {
		select {
		case <-time.After(1 * time.Second):
			t.l.Logf(logger.InfoLevel, "I am alive %d", i)
		case <-ctx.Done():
			t.l.Logf(logger.InfoLevel, "Cancelled")
			return ctx.Err()
		}
	}

	return nil
}
