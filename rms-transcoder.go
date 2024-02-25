package main

import (
	"fmt"
	rms_transcoder "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-transcoder"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/RacoonMediaServer/rms-transcoder/internal/config"
	"github.com/RacoonMediaServer/rms-transcoder/internal/db"
	"github.com/RacoonMediaServer/rms-transcoder/internal/service/profiles"
	"github.com/RacoonMediaServer/rms-transcoder/internal/service/transcoder"
	"github.com/RacoonMediaServer/rms-transcoder/internal/worker"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"

	// Plugins
	_ "github.com/go-micro/plugins/v4/registry/etcd"
)

var Version = "v0.0.0"

const serviceName = "rms-transcoder"

func main() {
	logger.Infof("%s %s", serviceName, Version)
	defer logger.Info("DONE.")

	useDebug := false

	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version(Version),
		micro.Flags(
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"debug"},
				Usage:       "debug log level",
				Value:       false,
				Destination: &useDebug,
			},
		),
	)

	service.Init(
		micro.Action(func(context *cli.Context) error {
			configFile := fmt.Sprintf("/etc/rms/%s.json", serviceName)
			if context.IsSet("config") {
				configFile = context.String("config")
			}
			return config.Load(configFile)
		}),
	)

	cfg := config.Config()

	if useDebug || cfg.Debug.Verbose {
		_ = logger.Init(logger.WithLevel(logger.DebugLevel))
	}

	_ = servicemgr.NewServiceFactory(service)

	database, err := db.Connect(config.Config().Database)
	if err != nil {
		logger.Fatalf("Connect to database failed: %s", err)
	}

	workers := worker.New(cfg.Transcoding.Workers)

	profilesService := &profiles.Service{
		Database: database,
	}
	transcoderService := &transcoder.Service{
		Database: database,
		Profiles: database,
		Workers:  workers,
	}

	if err = profilesService.Initialize(); err != nil {
		logger.Fatalf("Initialize profile service failed: %s", err)
	}

	if err = transcoderService.Initialize(); err != nil {
		logger.Fatalf("Initialize transcoder service failed: %s", err)
	}

	if err = rms_transcoder.RegisterProfilesHandler(service.Server(), profilesService); err != nil {
		logger.Fatalf("Register profile service failed: %s", err)
	}

	if err = rms_transcoder.RegisterTranscoderHandler(service.Server(), transcoderService); err != nil {
		logger.Fatalf("Register transcoder service failed: %s", err)
	}

	if err = service.Run(); err != nil {
		logger.Fatalf("Run service failed: %s", err)
	}

	workers.Stop()
}
