package transcoder

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-packages/pkg/media"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"go-micro.dev/v4/logger"
	"os"
	"os/exec"
)

type transcodingTask struct {
	l           logger.Logger
	id          string
	settings    *media.TranscodingSettings
	source      string
	destination string
	dur         *uint32
}

func (t *transcodingTask) ID() string {
	return t.id
}

func normalizeResolution(r *uint32) int64 {
	if r == nil {
		return -1
	}
	return int64(*r)
}

func (t *transcodingTask) compileStream() *exec.Cmd {
	stream := ffmpeg_go.Input(t.source)
	p := t.settings
	outputArgs := ffmpeg_go.KwArgs{}

	if p.Video != nil {
		if p.Video.Height != nil || p.Video.Width != nil {
			stream = stream.Filter("scale", ffmpeg_go.Args{fmt.Sprintf("w=%d:h=%d", normalizeResolution(p.Video.Width), normalizeResolution(p.Video.Height))})
		}
		if p.Video.Codec != nil {
			outputArgs["c:v"] = *p.Video.Codec
		}
		if p.Video.Bitrate != nil {
			outputArgs["b:v"] = *p.Video.Bitrate
		}
	}

	if p.Audio != nil {
		if p.Audio.Codec != nil {
			outputArgs["c:a"] = *p.Audio.Codec
		}
		if p.Audio.Bitrate != nil {
			outputArgs["b:a"] = *p.Audio.Bitrate
		}
	}

	if t.dur != nil {
		outputArgs["t"] = *t.dur
	}

	return stream.Output(t.destination, outputArgs).Compile()
}

func (t *transcodingTask) cleanUp(success *bool) {
	if !*success {
		_ = os.Remove(t.destination)
	}
}

func (t *transcodingTask) Do(ctx context.Context) error {
	success := false
	defer t.cleanUp(&success)

	cmd := t.compileStream()
	ch := make(chan error)
	go func() {
		ch <- cmd.Run()
	}()

	var err error
	select {
	case err = <-ch:
	case <-ctx.Done():
		err = ctx.Err()
		_ = cmd.Cancel()
		<-ch
	}

	close(ch)
	success = err == nil
	return err
}
