package main

import (
	"context"
	"fmt"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"github.com/google/uuid"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"time"
)
import "github.com/urfave/cli/v2"

func main() {
	var command string
	var job string
	service := micro.NewService(
		micro.Name("rms-transcoder.client"),
		micro.Flags(
			&cli.StringFlag{
				Name:        "command",
				Usage:       "Must be one of: add, status, cancel",
				Required:    true,
				Destination: &command,
			},
			&cli.StringFlag{
				Name:        "job",
				Usage:       "job - job id",
				Required:    false,
				Destination: &job,
			},
		),
	)
	service.Init()

	client := rms_transcoder.NewTranscoderService("rms-transcoder", service.Client())

	switch command {
	case "add":
		if err := add(client); err != nil {
			panic(err)
		}
	case "status":
		if err := status(client, job); err != nil {
			panic(err)
		}
	case "cancel":
		if err := cancel(client, job); err != nil {
			panic(err)
		}
	default:
		panic("unknown command: " + command)
	}
}

func add(cli rms_transcoder.TranscoderService) error {
	req := rms_transcoder.AddJobRequest{
		Profile:      "telegram",
		Source:       "source1",
		Destination:  uuid.NewString() + ".mp4",
		AutoComplete: true,
	}

	resp, err := cli.AddJob(context.Background(), &req, client.WithRequestTimeout(40*time.Second))
	if err != nil {
		return err
	}

	fmt.Printf("Job: %s\n", resp.JobId)
	return nil
}

func status(cli rms_transcoder.TranscoderService, job string) error {
	result, err := cli.GetJob(context.Background(), &rms_transcoder.GetJobRequest{JobId: job})
	if err != nil {
		return err
	}
	fmt.Println(result.Status, result.Destination)
	return nil
}

func cancel(cli rms_transcoder.TranscoderService, job string) error {
	_, err := cli.CancelJob(context.Background(), &rms_transcoder.CancelJobRequest{JobId: job})
	return err
}
